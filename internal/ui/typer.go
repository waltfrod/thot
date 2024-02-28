package ui

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wrodriguez/thot/internal/util"
)

type Status int

const (
	Ok Status = iota
	Err
	None
)

var (
	cursorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("240"))
	okStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("40"))
	errStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	itemStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("32"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("151"))
)

var (
	wordStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	timeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("29"))
	mistakeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("197"))
	wpmStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("73"))
	precStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	boxStyle     = lipgloss.NewStyle().Padding(1).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("241"))
)

type KeyMap struct {
	Salir key.Binding
}

type Stats struct {
	tchar int
	tempo float64
	cerr  int
}

type Character struct {
	char   string
	style  lipgloss.Style
	active bool
	status Status
}

type Size struct {
	Width  int
	Height int
}

type Model struct {
	current []Character
	cursor  int
	line    int
	lines   []string
	wsize   Size
	start   time.Time
	help    help.Model
	cerr    int
	end     bool
	stats   Stats
	first   bool
}

var keyMap = KeyMap{
	Salir: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Salir"),
	),
}

func NewCharacter(char string) Character { // {{{
	return Character{
		char:   char,
		style:  defaultStyle,
		active: false,
		status: None,
	}
} // }}}

func countMistakes(c []Character) int { // {{{
	cont := 0
	lc := len(c)
	for i := 0; i < lc; i++ {
		if c[i].status == Err {
			cont++
		}
	}

	return cont
} // }}}

func NewModel(lines []string) *Model { // {{{
	m := &Model{
		current: []Character{},
		cursor:  0,
		line:    0,
		lines:   lines,
		wsize:   Size{},
		help:    help.New(),
		cerr:    0,
		end:     false,
		first:   true,
	}

	return m
} // }}}

func (s Stats) String() string { // {{{
	sb := strings.Builder{}

	sb.WriteString(boxStyle.Render(
		wordStyle.Render(fmt.Sprintf(" 󰀬 Total caracteres: %d", s.tchar)) + "\n " +
			timeStyle.Render(fmt.Sprintf(" Tiempo: %.2fs(%.2fm)", s.tempo, s.tempo/60)) + "\n " +
			mistakeStyle.Render(fmt.Sprintf("󰚌 errores: %d", s.cerr)) + "\n " +
			wpmStyle.Render(
				fmt.Sprintf("󰌓 WPM: %.2f", s.WPM(s.tchar, s.cerr, s.tempo/60)),
			) + "\n " +
			precStyle.Render(fmt.Sprintf("󰓾 Precisión: %.2f%%", s.Accuracy(s.tchar, s.cerr))),
	))

	return sb.String()
} // }}}

// WPM calcula los caracteres por minuto, basada en la ecuación:
// https://www.speedtypingonline.com/typing-equations
func (s Stats) WPM(all, uncorrect int, minutes float64) float64 { // {{{
	return (float64(all/5) - float64(uncorrect)) / minutes
} // }}}

// Accuracy calcula el porcentaje de caracteres correctos
func (s Stats) Accuracy(all, uncorrect int) float64 { // {{{
	return (float64(all-uncorrect) / float64(all)) * 100
} // }}}

// ShortHelp devuelve las combinaciones de teclas que se mostrarán en la vista de mini ayuda.
// Forma parte de la interfaz key.Map.
func (k KeyMap) ShortHelp() []key.Binding { // {{{
	return []key.Binding{k.Salir}
} // }}}

// FullHelp devuelve las combinaciones de teclas para la vista de ayuda expandida.
// Forma parte de la interfaz interfaz key.Map.
func (k KeyMap) FullHelp() [][]key.Binding { // {{{
	return [][]key.Binding{{k.Salir}}
} // }}}

func (c Character) Rune() rune { // {{{
	return []rune(c.char)[0]
} // }}}

func (c Character) Char() string { // {{{
	return c.char
} // }}}

func (c Character) String() string { // {{{
	return util.IF(c.active, cursorStyle, c.style).Render(c.char)
} // }}}

func (c *Character) Style(s lipgloss.Style) { // {{{
	c.style = s
} // }}}

func (c *Character) Active() { // {{{
	c.active = true
} // }}}

func (c *Character) Inactive() { // {{{
	c.active = false
} // }}}

func (c *Character) Ok() { // {{{
	c.style = okStyle
	c.status = Ok
} // }}}

func (c *Character) Err() { // {{{
	c.style = errStyle
	c.status = Err
} // }}}

func (m Model) Init() tea.Cmd { // {{{
	// Simplemente devuelve `nil`, que significa "nada de E/S ahora mismo, por favor".
	return nil
} // }}}

func (m *Model) Start() { // {{{
	m.start = time.Now()
} // }}}

func (m *Model) Stop() { // {{{
	m.end = true
	seg := time.Since(m.start).Seconds()
	txtlen := utf8.RuneCountInString(strings.Join(m.lines, " "))

	m.stats = Stats{
		tchar: txtlen,
		cerr:  m.cerr,
		tempo: seg,
	}
} // }}}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { // {{{
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.wsize.Width = msg.Width
		m.wsize.Height = msg.Height
		return m, nil
	case tea.KeyMsg:
		lmsg := len([]rune(msg.String()))
		ms := msg.String()
		_ = lmsg
		switch ms {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.cursor >= len(m.current) {
				m.cerr += countMistakes(m.current)
				m.line++
				if m.line < len(m.lines) {
					m.current = m.ToChars()
					m.cursor = 0
				}
			}
			if m.line >= len(m.lines) {
				m.Stop()
				return m, tea.Quit
			}
		case "backspace":
			if m.cursor > 0 && m.cursor < len(m.current) {
				m.current[m.cursor].Inactive()
				m.current[m.cursor].Style(defaultStyle)
				m.cursor -= 1
				m.current[m.cursor].Active()
			}
		default:
			if lmsg == 1 {
				if m.cursor < len(m.current) {
					if m.current[m.cursor].Char() == ms {
						m.current[m.cursor].Ok()
					} else {
						m.current[m.cursor].Err()
					}

					m.current[m.cursor].Inactive()
					m.cursor += util.IF(m.cursor < len(m.current), 1, 0)
					if m.cursor < len(m.current) {
						m.current[m.cursor].Active()
					}
				}
			}
		}
	}

	return m, nil
} // }}}

func (m *Model) View() string { // {{{
	if len(m.current) == 0 && m.cursor == 0 {
		m.current = m.ToChars()
	}
	sb := strings.Builder{}
	sb.WriteString(infoStyle.Render(" Inicio: "+m.start.Format("03:04:05 PM")) + "\n\n")
	if !m.end {
		sb.WriteString(itemStyle.Render("  "))
		for _, char := range m.current {
			sb.WriteString(char.String())
		}

		sb.WriteString("\n\n\n" + m.help.View(keyMap))
	} else {
		sb.WriteString(m.stats.String() + "\n\n")
	}

	return sb.String()
} // }}}

func (m *Model) ToChars() []Character { // {{{
	var chars []Character = []Character{}
	word := strings.Split(m.lines[m.line], "")
	for _, char := range word {
		chars = append(chars, NewCharacter(char))
	}
	if len(chars) > 0 {
		chars[0].Active()
	}

	return chars
} // }}}
