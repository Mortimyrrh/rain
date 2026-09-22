package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/NimbleMarkets/go-booba"
)

const fps = 24

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type line struct {
	x, y, len int
}

type model struct {
	width, height int
	rain          [][]rune //yx faster to draw I hope
	hello         string
	line          line
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
		m.rain = make([][]rune, m.height)
		for row := range m.rain {
			m.rain[row] = make([]rune, m.width)
		}

	case tea.KeyPressMsg:
		// "press 'any' key to continue" or quit in this case...
		return m, tea.Quit
	case tickMsg:
		// clear rain
		for row := range m.rain {
			for column := range m.rain[row] {
				m.rain[row][column] = ideographicSpace
			}
		}
		// delete lines off screen
		if m.line.y > len(m.rain) || m.line.x >= m.width {
			// reset line
			m.line.y = 0 - m.line.len
			m.line.x = rand.Intn(m.width)
		}
		// update lines
		m.line.y++
		// draw lines
		for yOffset := range m.line.len {
			y := m.line.y + yOffset
			if y >= len(m.rain) {
				break
			}
			if y >= 0 {
				m.rain[y][m.line.x] = randKat()
			}
		}
		return m, tick()
	}
	return m, nil
}

func (m model) View() tea.View {
	var sb strings.Builder
	for row := range m.rain {
		sb.WriteString(string(m.rain[row]))
		// no blank line at the bottom
		if row != len(m.rain)-1 {
			sb.WriteRune('\n')
		}
	}
	return tea.NewView(sb.String())
}

func main() {
	m := model{line: line{x: 10, y: 0, len: 10}}
	if err := booba.Run(m); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
