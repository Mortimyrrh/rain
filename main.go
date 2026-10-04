package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/NimbleMarkets/go-booba"
)

const fps = 120
const lineLenMin = 3
const lineLenVariance = 15 // this is added to the line length
const lineSpeedMin = 0.01
const lineSpeedVariance = .03

func tick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return time.Time(t)
	})
}

type line struct {
	x, y, len int
	speed     float64
	yF        float64
	colour    string
	history   []string
}

func NewLine() line {
	return line{y: 999999} // will get reset immediately
}

type model struct {
	width, height int
	rain          [][]string //yx
	hello         string
	lines         []line
	lastTick      time.Time
	speed         float64
	multi         bool // multicolour mode
	maxLines      int
	density       float64 // does update dynamically
}

func (m model) Init() tea.Cmd {
	return tick()
}

var katakana = []rune{
	'１', '２', '３', '４', '５', '６', '７', '８', '９', '０',
	'゠', 'ァ', 'ア', 'ィ', 'イ', 'ゥ', 'ウ', 'ェ', 'エ', 'ォ', 'オ', 'カ', 'ガ', 'キ', 'ギ', 'ク',
	'グ', 'ケ', 'ゲ', 'コ', 'ゴ', 'サ', 'ザ', 'シ', 'ジ', 'ス', 'ズ', 'セ', 'ゼ', 'ソ', 'ゾ', 'タ',
	'ダ', 'チ', 'ヂ', 'ッ', 'ツ', 'ヅ', 'テ', 'デ', 'ト', 'ド', 'ナ', 'ニ', 'ヌ', 'ネ', 'ノ', 'ハ',
	'バ', 'パ', 'ヒ', 'ビ', 'ピ', 'フ', 'ブ', 'プ', 'ヘ', 'ベ', 'ペ', 'ホ', 'ボ', 'ポ', 'マ', 'ミ',
	'ム', 'メ', 'モ', 'ャ', 'ヤ', 'ュ', 'ユ', 'ョ', 'ヨ', 'ラ', 'リ', 'ル', 'レ', 'ロ', 'ヮ', 'ワ',
	'ヰ', 'ヱ', 'ヲ', 'ン', 'ヴ', 'ヵ', 'ヶ', 'ヷ', 'ヸ', 'ヹ', 'ヺ', '・', 'ー', 'ヽ', 'ヾ', 'ヿ',
}

const ideographicSpace = string('　')

func randKat() rune {
	return katakana[rand.Intn(len(katakana))]
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width / 2 // using double width unicode characters
		m.height = msg.Height
		// recreate rain
		m.rain = make([][]string, m.height)
		for row := range m.rain {
			m.rain[row] = make([]string, m.width)
		}

		m.maxLines = int(float64(m.width*m.height) * m.density)
		m.lines = make([]line, m.maxLines)
		for range m.maxLines {
			m.lines = append(m.lines, NewLine())
		}

	case tea.KeyPressMsg:
		// "press 'any' key to continue" or quit in this case...
		return m, tea.Quit
	case time.Time:
		currentTime := msg
		elapsedTime := currentTime.Sub(m.lastTick)

		// clear rain
		for row := range m.rain {
			for column := range m.rain[row] {
				m.rain[row][column] = ideographicSpace
			}
		}
		// lines
		for line, _ := range m.lines {
			line := &m.lines[line]
			// reset line off screen
			if line.y > len(m.rain) || line.x >= m.width {
				line.speed = float64(float64(lineSpeedMin*float64(m.height)/1000) + (rand.Float64() * lineSpeedVariance)) ///1000 is a nicer range to work around
				line.len = lineLenMin + rand.Intn(lineLenVariance)
				line.y = 0
				line.yF = -float64(line.len + rand.Intn(m.height))
				line.x = rand.Intn(m.width)
				line.history = make([]string, m.height)
				for i := range m.height {
					line.history[i] = string(randKat())
				}
				if m.multi {
					line.colour = colours[rand.Int31n(int32(len(colours)))]
				} else {
					line.colour = white
				}
			}
			// update line
			line.yF = line.yF + (float64(elapsedTime.Milliseconds()) * m.speed * line.speed)
			line.y = int(line.yF)

			// draw line
			for yOffset := range line.len {
				y := line.y + yOffset
				if y >= len(m.rain) {
					break
				}
				if y >= 0 {
					rune_ := line.history[y]
					if yOffset == line.len-1 {
						rune_ = string(randKat())
					}
					if m.multi {
						m.rain[y][line.x] = line.colour + rune_ // + reset (no need to mid way reset do it once last)
					} else {
						m.rain[y][line.x] = rune_ // skip setting colour
					}
				}
			}
		}
		m.lastTick = currentTime
		return m, tick()
	}
	return m, nil
}

// https://en.wikipedia.org/wiki/ANSI_escape_code#/media/File:ANSI_sample_program_output.png
// const ESC = "\033"
const reset = "\033[0m"
const black = "\033[30m"
const red = "\033[31m"
const green = "\033[32m"
const yellow = "\033[33m"
const blue = "\033[34m"
const violet = "\033[35m"
const lightblue = "\033[36m"
const white = "\033[36m"

var colours = []string{red, green, yellow, blue, violet, lightblue}

func (m model) View() tea.View {
	var sb strings.Builder

	for row := range m.rain {
		for _, ru := range m.rain[row] {
			sb.WriteString(ru)
		}
		// no blank line at the bottom
		if row != len(m.rain)-1 {
			sb.WriteRune('\n')
		}
	}
	sb.WriteString(reset)
	return tea.NewView(sb.String())
}

func main() {
	m := model{maxLines: 1}
	flag.Float64Var(&m.speed, "speed", 1, "adds speed (1 is default)")
	flag.BoolVar(&m.multi, "multi", false, "adds colour (use a gpu accelerated terminal for less lag)")
	flag.Float64Var(&m.density, "density", .05, "adds lines (.1 is default)")
	flag.Parse()

	if err := booba.Run(m); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
