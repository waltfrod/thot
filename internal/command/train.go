package command

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wrodriguez/thot/internal/db"
	"github.com/wrodriguez/thot/internal/kbd"
	"github.com/wrodriguez/thot/internal/ui"
	"github.com/wrodriguez/thot/internal/util"
	"gitlab.com/tozd/go/errors"
)

const LenWords = 2

var defStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("36")).Italic(true)

func validateLang(lang string) bool {
	return lang == "spa" || lang == "eng"
}

func validateRows(rows []string) bool {
	for _, row := range rows {
		if row != "row1" && row != "row2" && row != "row3" && row != "row4" && row != "all" {
			return false
		}
	}

	return true
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func Train(layoutName, lang string, rows []string) errors.E { // {{{
	rows = unique(rows)
	if layout := kbd.FindLayout(layoutName); layout != nil {
		if validateLang(lang) {
			lng := util.IF(lang == "spa", db.Spanish, db.English)
			if validateRows(rows) {
				checkAll := util.InSlice(func(i int) bool {
					return rows[i] == "all"
				}, len(rows))
				if checkAll != -1 {
					rows = []string{"row1", "row2", "row3", "row4"}
				}

				kdb, err := db.NewDatabase()
				if err != nil {
					return errors.WithMessage(err, "No se pudo conectar a la Base de Datos")
				}
				words, err := kdb.Words(LenWords, lng, layout.GetKeys(rows...))
				if err != nil {
					return errors.WithMessage(err, "No se pudo obtener las palabras")
				}
				fmt.Println(defStyle.Render("󰌓  Layout:"), layoutName)
				fmt.Println(defStyle.Render("  Idioma:"), lang)
				fmt.Println(defStyle.Render("󰠷  Filas:"), strings.Join(rows, ", "))
				fmt.Println(defStyle.Render("󱀍  Cantidad de palabras: "), len(words))
				fmt.Println(defStyle.Render("󰘝  Letras a practicar: "), layout.GetKeys(rows...))
				fmt.Println(defStyle.Render("  Para comenzar pulse Enter ..."))
				util.Pause(false)

				model := ui.NewModel(words)
				// fmt.Printf("model: %#v\n", model)
				p := tea.NewProgram(model)

				// Pause()
				model.Start()
				if _, err := p.Run(); err != nil {
					fmt.Printf("Alas, there's been an error: %v", err)
					os.Exit(1)
				}
			} else {
				fmt.Println(errStyle.Render("Los valores válidos para las filas son `row1`, `row2`, `row3` o `row4`"))
				os.Exit(2)
			}
		} else {
			fmt.Println(errStyle.Render(fmt.Sprintf("El idioma %q no es valido", lang)))
			os.Exit(2)
		}
	} else {
		fmt.Println(errStyle.Render(fmt.Sprintf("Layout %q no encontrado", layoutName)))
		os.Exit(2)
	}

	return nil

} // }}}
