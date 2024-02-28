package command

import (
	"fmt"
	"math"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/wrodriguez/thot/internal/kbd"
	"github.com/wrodriguez/thot/internal/util"
)

const banner = `LAYOUTS SOPORTADOS POR THOT
===========================`

var wStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
var tStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("38"))
var iStyle = lipgloss.NewStyle().Italic(true)

func ListLayouts() {
	var width, _ = util.GetConsoleSize()
	cols := width / 30
	fmt.Println(tStyle.Render(banner))
	layouts := kbd.ListLayouts()
	sort.Strings(layouts)
	tbl := sliceToMatrix(layouts, cols)
	t := table.New().
		Border(lipgloss.DoubleBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(tbl...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("194")).Padding(0, 1)
		})
	fmt.Println(t.Render())
}

func sliceToMatrix(s []string, cols int) [][]string {
	l := len(s)
	rows := int(math.Ceil(float64(l) / float64(cols)))
	matrix := make([][]string, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]string, cols)
		for j := 0; j < cols; j++ {
			lpos := cols*i + j
			if lpos < l {
				matrix[i][j] = "  " + s[lpos]
			} else {
				matrix[i][j] = ""
			}
		}
	}

	return matrix
}
