package command

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/wrodriguez/thot/internal/kbd"
	"github.com/wrodriguez/thot/internal/util"
)

var errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))

func PrintLayout(layoutName string) {
	if layout := kbd.FindLayout(layoutName); layout != nil {
		var width, _ = util.GetConsoleSize()
		if width < 105 {
			fmt.Println(
				wStyle.Render(
					"El ancho de la consola es demasiado corto para mostrar el diagrama. Se recomienda al menos 105 caracteres de ancho.",
				),
			)
			kbd.PrintMiniKeyboard(layoutName, layout)
			return
		}
		kbd.PrintKeyboard(layoutName, layout)
	} else {
		fmt.Println(errStyle.Render(fmt.Sprintf("Layout %q no encontrado", layoutName)))
	}
}
