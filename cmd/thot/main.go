package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/integrii/flaggy"
	"github.com/wrodriguez/thot/internal/command"
)

var Version = "devel"
var CommitHash = "devel"
var BuildTimestamp = "--now--"

const Logo = `┏┳┓┓
 ┃ ┣┓┏┓╋
 ┻ ┛┗┗┛┗`

// const Logo = "╔╦╗┬ ┬┌─┐┌┬┐\n" +
// 	" ║ ├─┤│ │ │ \n" +
// 	" ╩ ┴ ┴└─┘ ┴"

var listCommand *flaggy.Subcommand
var printCommand *flaggy.Subcommand
var trainCommand *flaggy.Subcommand
var dbCommand *flaggy.Subcommand

var layoutName string = "qwerty"
var layout string = "qwerty"
var lang string = "spa"
var rows []string = []string{"row3"}

var errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Bold(true)

func main() {
	configArgs()
	flaggy.Parse()

	if listCommand != nil && listCommand.Used {
		command.ListLayouts()
	} else if printCommand != nil && printCommand.Used {
		command.PrintLayout(layoutName)
	} else if trainCommand != nil && trainCommand.Used {
		command.Train(layout, lang, rows)
	} else if dbCommand != nil && dbCommand.Used {
		err := command.CopyDB()
		if err != nil {
			fmt.Println(errStyle.Render(err.Error()))
			os.Exit(1)
		}
	} else {
		flaggy.ShowHelp("")
	}

}

func configArgs() { // {{{
	flaggy.DefaultParser.AdditionalHelpPrepend = Logo
	flaggy.SetName("thot")
	flaggy.SetDescription("Un pequeño entrenador de teclado")
	flaggy.SetVersion(fmt.Sprintf("%s (%s, %s)", Version, CommitHash, BuildTimestamp))

	listCommand = flaggy.NewSubcommand("list")
	listCommand.Description = "Lista los layouts que pueden ser utilizados por Thot"
	printCommand = flaggy.NewSubcommand("print")
	printCommand.Description = "Imprime el layout de teclado seleccionado"
	printCommand.AddPositionalValue(
		&layoutName,
		"layout-name",
		1,
		true,
		"El nombre del layout, se puede consultar la lista de layouts disponibles a través del comando `thot list`",
	)
	trainCommand = flaggy.NewSubcommand("train")
	trainCommand.Description = "Practica con Thot para mejorar el método de mecanografía"
	trainCommand.AddPositionalValue(
		&layout,
		"layout",
		1,
		true,
		"El nombre del layout, se puede consultar la lista de layouts disponibles a través del comando `thot list`",
	)
	trainCommand.AddPositionalValue(
		&lang,
		"lang",
		2,
		true,
		"El idioma a mostrar las palabras, acepta solo los valores `spa` o `eng`",
	)
	trainCommand.StringSlice(
		&rows,
		"r",
		"row",
		"Las filas a mostrar, acepta solo los valores `row1`, `row2`, `row3`, `row4` o `all` para mostrar todas",
	)

	dbCommand = flaggy.NewSubcommand("db")
	dbCommand.Description = "Crea la base de datos de palabras de Thot"

	flaggy.AttachSubcommand(listCommand, 1)
	flaggy.AttachSubcommand(printCommand, 1)
	flaggy.AttachSubcommand(trainCommand, 1)
	flaggy.AttachSubcommand(dbCommand, 1)

} // }}}
