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
const lineLenVariance = 10 // this is added to te line length
const lineSpeedMin = 0.25
const lineSpeedVariance = 2
const maxLines = 100

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
}

func NewLine() line {
	return line{y: 999999}
}

type model struct {
	width, height int
	rain          [][]string //yx
	hello         string
	lines         []line
	lastTick      time.Time
	speed         float64
	multi         bool // multicolour mode
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

var ideographicSpace = '　'

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

	case tea.KeyPressMsg:
		// "press 'any' key to continue" or quit in this case...
		return m, tea.Quit
	case time.Time:
		currentTime := msg
		elapsedTime := currentTime.Sub(m.lastTick)

		// clear rain
		for row := range m.rain {
			for column := range m.rain[row] {
				m.rain[row][column] = string(ideographicSpace)
			}
		}
		// lines
		for line, _ := range m.lines {
			line := &m.lines[line]
			// reset line off screen
			if line.y > len(m.rain) || line.x >= m.width {
				line.speed = float64(float64(lineSpeedMin+(rand.Float64()*lineSpeedVariance)) / float64(m.height)) ///40 is a nicer range to work around
				line.len = lineLenMin + rand.Intn(lineLenVariance)
				line.y = 0 - line.len
				line.yF = float64(0 - line.len)
				line.x = rand.Intn(m.width)
				if m.multi {
					line.colour = colours[rand.Int31n(int32(len(colours)))]
				} else {
					line.colour = white
				}
			}
			// update line
			line.yF = line.yF + (float64(elapsedTime.Milliseconds()) * line.speed)
			line.y = int(line.yF)

			// draw line
			for yOffset := range line.len {
				y := line.y + yOffset
				if y >= len(m.rain) {
					break
				}
				if y >= 0 {
					m.rain[y][line.x] = line.colour + string(randKat()) + reset
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
const black = "\033[30;1m"
const red = "\033[31;1m"
const green = "\033[32;1m"
const yellow = "\033[33;1m"
const blue = "\033[34;1m"
const violet = "\033[35;1m"
const lightblue = "\033[36;1m"
const white = "\033[36;1m"

var colours = []string{red, green, yellow, blue, violet, lightblue}

func (m model) View() tea.View {
	var sb strings.Builder

	for row := range m.rain {
		sb.WriteString(strings.Join(m.rain[row], ""))
		// no blank line at the bottom
		if row != len(m.rain)-1 {
			sb.WriteRune('\n')
		}
	}
	return tea.NewView(sb.String())
}

type args struct {
	multi bool
}

func main() {
	args := args{}
	flag.BoolVar(&args.multi, "multi", false, "adds colour")
	flag.Parse()

	m := model{speed: 1, multi: args.multi}
	m.lines = make([]line, maxLines)
	for range maxLines {
		m.lines = append(m.lines, NewLine())
	}

	if err := booba.Run(m); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
