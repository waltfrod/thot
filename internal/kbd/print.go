package kbd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wrodriguez/thot/internal/util"
)

const (
	Menique = "\x1b[38;5;227m"
	Anular  = "\x1b[38;5;69m"
	Corazon = "\x1b[38;5;70m"
	Indicei = "\x1b[38;5;161m"
	Indiced = "\x1b[38;5;172m"
)

var defStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("192"))

var (
	meniqueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("227"))
	anularStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("69"))
	corazonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("70"))
	indiceiStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("161"))
	indicedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("172"))
)

const (
	Row1 = "row1"
	Row2 = "row2"
	Row3 = "row3"
	Row4 = "row4"
)

const (
	Sup = "\x1e"
	Inf = "\x1f"
)

var ansi map[string]string = make(map[string]string)
var iso map[string]string = make(map[string]string)
var miniAnsi map[string]string = make(map[string]string)
var miniIso map[string]string = make(map[string]string)

//go:embed layout.json
var blayouts []byte
var layouts map[string]Keyboard

const rowBottom = `
╔═══════╗╔═══════╗╔════════╗╔════════════════════════════════════════════╗╔═══════╗╔═══════╗╔═══════╗
║░░░░░░░║║░░░░░░░║║░░░░░░░░║║░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░║║░░░░░░░║║░░░░░░░║║░░░░░░░║
║░░░░░░░║║░░░░░░░║║░░░░░░░░║║░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░║║░░░░░░░║║░░░░░░░║║░░░░░░░║
╚═══════╝╚═══════╝╚════════╝╚════════════════════════════════════════════╝╚═══════╝╚═══════╝╚═══════╝`

const miniRowBottom = `╔════╗╔════╗╔════╗╔═══════════════════╗╔════╗╔════╗╔════╗
║░░░░║║░░░░║║░░░░║║░░░░░░░░░░░░░░░░░░░║║░░░░║║░░░░║║░░░░║
╚════╝╚════╝╚════╝╚═══════════════════╝╚════╝╚════╝╚════╝`

var reLetter = regexp.MustCompile(`[a-zA-ZñÑ]`)

type Keyboard struct {
	Type string              `json:"type"`
	Keys map[string][]string `json:"keys"`
}

func (k *Keyboard) String() string { // {{{
	return fmt.Sprintf(
		"Type: %s\n, Keys:\n\t%q\n\t%q\n\t%q\n\t%q",
		k.Type,
		k.Keys[Row1],
		k.Keys[Row2],
		k.Keys[Row3],
		k.Keys[Row4],
	)
} // }}}

func init() { // {{{
	ansi["row1"] = fmt.Sprintf(
		"%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮╭─────╮\x1b[0m╔════════╗\n",
		Menique,
		Anular,
		Corazon,
		Indicei,
		Indiced,
		Corazon,
		Anular,
		Menique,
	) +
		fmt.Sprintf(
			"%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   ││ \x1e   │\x1b[0m║░░░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   ││ \x1f   │\x1b[0m║░░░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"╰%s─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯%s╰─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯╰─────╯\x1b[0m╚════════╝\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		)
	ansi["row2"] = fmt.Sprintf(
		"╔════════╗%s╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮╭─────╮╭─────╮\x1b[0m\n",
		Menique,
		Anular,
		Corazon,
		Indicei,
		Indiced,
		Corazon,
		Anular,
		Menique,
	) +
		fmt.Sprintf(
			"║░░░░░░░░║%s│ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   ││ \x1e   ││ \x1e   │\x1b[0m\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"║░░░░░░░░║%s│ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   ││ \x1f   ││ \x1f   │\x1b[0m\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"╚════════╝%s╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯%s╰─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯╰─────╯╰─────╯\x1b[0m\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		)
	ansi["row3"] = fmt.Sprintf(
		"╔══════════╗%s╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮\x1b[0m╔══════════╗\n",
		Menique,
		Anular,
		Corazon,
		Indicei,
		Indiced,
		Corazon,
		Anular,
		Menique,
	) +
		fmt.Sprintf(
			"║░░░░░░░░░░║%s│ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   │\x1b[0m║░░░░░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"║░░░░░░░░░░║%s│ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   │\x1b[0m║░░░░░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"╚══════════╝%s╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯%s╰─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯\x1b[0m╚══════════╝\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		)
	ansi["row4"] = fmt.Sprintf(
		"╔═════════════╗%s╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮\x1b[0m╔══════════════╗\n",
		Menique,
		Anular,
		Corazon,
		Indicei,
		Indiced,
		Corazon,
		Anular,
		Menique,
	) +
		fmt.Sprintf(
			"║░░░░░░░░░░░░░║%s│ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   │\x1b[0m║░░░░░░░░░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"║░░░░░░░░░░░░░║%s│ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   │\x1b[0m║░░░░░░░░░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"╚═════════════╝%s╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯%s╰─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯\x1b[0m╚══════════════╝",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		)

	iso["row1"] = fmt.Sprintf(
		"%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮╭─────╮\x1b[0m╔════════╗\n",
		Menique,
		Anular,
		Corazon,
		Indicei,
		Indiced,
		Corazon,
		Anular,
		Menique,
	) +
		fmt.Sprintf(
			"%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   ││ \x1e   │\x1b[0m║░░░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   ││ \x1f   │\x1b[0m║░░░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"%s╰─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯%s╰─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯╰─────╯\x1b[0m╚════════╝\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		)
	iso["row2"] = fmt.Sprintf(
		"╔═══════╗%s╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮╭─────╮\x1b[0m╔══════╗\n",
		Menique,
		Anular,
		Corazon,
		Indicei,
		Indiced,
		Corazon,
		Anular,
		Menique,
	) +
		fmt.Sprintf(
			"║░░░░░░░║%s│ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   ││ \x1e   │\x1b[0m║░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"║░░░░░░░║%s│ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   ││ \x1f   │\x1b[0m║░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"╚═══════╝%s╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯%s╰─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯╰─────╯\x1b[0m╚═╗░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		)
	iso["row3"] = fmt.Sprintf(
		"╔═════════╗%s╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮╭─────╮\x1b[0m║░░░░║\n",
		Menique,
		Anular,
		Corazon,
		Indicei,
		Indiced,
		Corazon,
		Anular,
		Menique,
	) +
		fmt.Sprintf(
			"║░░░░░░░░░║%s│ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   ││ \x1e   │\x1b[0m║░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"║░░░░░░░░░║%s│ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   ││ \x1f   │\x1b[0m║░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"╚═════════╝%s╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯%s╰─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯╰─────╯\x1b[0m╚════╝\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		)
	iso["row4"] = fmt.Sprintf(
		"╔═══════╗%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮╭─────╮%s╭─────╮╭─────╮%s╭─────╮%s╭─────╮%s╭─────╮\x1b[0m╔═════════════╗\n",
		Menique,
		Anular,
		Corazon,
		Indicei,
		Indiced,
		Corazon,
		Anular,
		Menique,
	) +
		fmt.Sprintf(
			"║░░░░░░░║%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   ││ \x1e   │%s│ \x1e   │%s│ \x1e   │%s│ \x1e   │\x1b[0m║░░░░░░░░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"║░░░░░░░║%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   ││ \x1f   │%s│ \x1f   │%s│ \x1f   │%s│ \x1f   │\x1b[0m║░░░░░░░░░░░░░║\n",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		) +
		fmt.Sprintf(
			"╚═══════╝%s╰─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯╰─────╯%s╰─────╯╰─────╯%s╰─────╯%s╰─────╯%s╰─────╯\x1b[0m╚═════════════╝",
			Menique,
			Anular,
			Corazon,
			Indicei,
			Indiced,
			Corazon,
			Anular,
			Menique,
		)

		// Mini Keyboard

	miniAnsi[Row1] = "╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╔═══╗\n" +
		"│\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f│║░░░║\n" +
		"╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╚═══╝\n"
	miniAnsi[Row2] = "╔═══╗╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮\n" +
		"║░░░║│\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f│\n" +
		"╚═══╝╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯\n"
	miniAnsi[Row3] = "╔═════╗╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╔════╗\n" +
		"║░░░░░║│\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f│║░░░░║\n" +
		"╚═════╝╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╚════╝\n"
	miniAnsi[Row4] = "╔═══════╗╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╔══════╗\n" +
		"║░░░░░░░║│\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f│║░░░░░░║\n" +
		"╚═══════╝╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╚══════╝\n"

	miniIso[Row1] = "╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╔═══╗\n" +
		"│\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f│║░░░║\n" +
		"╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╚═══╝\n"
	miniIso[Row2] = "╔══╗╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╔═══╗\n" +
		"║░░║│\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f│║░░░║\n" +
		"╚══╝╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯║░░░║\n"
	miniIso[Row3] = "╔═══════╗╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮║░░║\n" +
		"║░░░░░░░║│\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f│║░░║\n" +
		"╚═══════╝╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╚══╝\n"
	miniIso[Row4] = "╔═══════╗╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╭──╮╔══════╗\n" +
		"║░░░░░░░║│\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f││\x1e\x1f│║░░░░░░║\n" +
		"╚═══════╝╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╰──╯╚══════╝\n"

	if err := json.Unmarshal(blayouts, &layouts); err != nil {
		panic(err)
	}
} // }}}

func FindLayout(name string) *Keyboard { // {{{
	if layout, ok := layouts[name]; ok {
		return &layout
	}

	return nil
} // }}}

func ListLayouts() []string { // {{{
	var names []string = make([]string, 0)

	for name := range layouts {
		names = append(
			names,
			name,
		) // TODO: Revisar el formato de name)
	}
	return names
} // }}}

func (k *Keyboard) GetKeys(rows ...string) string { // {{{
	sb := strings.Builder{}
	for _, row := range rows {
		if r, ok := k.Keys[row]; ok {
			as := strings.Split(strings.Join(r, ""), "")
			for _, s := range as {
				if reLetter.MatchString(s) { // TODO: Revisar el regex
					sb.WriteString(s)
				}
			}
		}
	}

	return sb.String()
} // }}}

func PrintKeyboard(name string, k *Keyboard) { // {{{
	box := lipgloss.NewStyle().
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("27"))
	sbi := strings.Builder{}
	_ = box

	sbi.WriteString(
		defStyle.Render(
			"󰌓 Nombre: ",
		) + name + " | " + defStyle.Render(
			"󰌓 Tipo: ",
		) + strings.ToUpper(
			k.Type,
		),
	)
	// fmt.Println(box.Render(sbi.String()))

	sbk := strings.Builder{}
	template := util.IF(k.Type == "ansi", ansi, iso)
	for _, row := range []string{"row1", "row2", "row3", "row4"} {
		sbk.WriteString(replace(template[row], k.Keys[row]))
	}
	sbk.WriteString(rowBottom + "\033[0m\n")

	sbk.WriteString(fmt.Sprintf(
		"%s%s%s%s 󰹆                            󰹇  %s%s%s%s\n",
		meniqueStyle.Render(" Meñique "),
		anularStyle.Render(" Anular "),
		corazonStyle.Render(" Corazon "),
		indiceiStyle.Render(" Indice "),
		indicedStyle.Render(" Indice "),
		corazonStyle.Render(" Corazon "),
		anularStyle.Render(" Anular "),
		meniqueStyle.Render(" Meñique "),
	))

	fmt.Println(box.Render(lipgloss.JoinVertical(lipgloss.Left, sbi.String(), sbk.String())))
} // }}}

func PrintMiniKeyboard(name string, k *Keyboard) { // {{{
	box := lipgloss.NewStyle().
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("27"))
	sbi := strings.Builder{}
	sbi.WriteString(
		defStyle.Render(
			"󰌓 Nombre: ",
		) + name + " | " + defStyle.Render(
			"󰌓 Tipo: ",
		) + strings.ToUpper(
			k.Type,
		),
	)
	sbk := strings.Builder{}
	template := util.IF(k.Type == "ansi", miniAnsi, miniIso)
	for _, row := range []string{"row1", "row2", "row3", "row4"} {
		sbk.WriteString(replace(template[row], k.Keys[row]))
	}
	sbk.WriteString(miniRowBottom + "\033[0m\n")

	fmt.Println(box.Render(lipgloss.JoinVertical(lipgloss.Left, sbi.String(), sbk.String())))

} // }}}

func chars(s string) (string, string) { // {{{
	sp := strings.Split(s, "")
	if len(sp) == 2 {
		return sp[0], sp[1]
	}

	return "", ""
} // }}}

func replace(tpl string, data []string) string { // {{{
	str := tpl
	for _, v := range data {
		i, s := chars(v)
		str = strings.Replace(str, Sup, s, 1)
		str = strings.Replace(str, Inf, i, 1)
	}

	return str
} // }}}
