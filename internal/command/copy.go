package command

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/tozd/go/errors"
)

//go:embed data/*.db
var fs embed.FS

func CopyDB() errors.E {
	fmt.Println(defStyle.Render("Copiando Base de Datos..."))
	dbBytes, err := fs.ReadFile("data/words.db")

	if err != nil {
		return errors.WithMessage(err, "No se pudo crear la Base de Datos")
	}

	h, err := os.UserHomeDir()
	if err != nil {
		h = os.Getenv("HOME")
		if h == "" {
			return errors.WithMessage(err, "No se pudo obtener el directorio HOME")
		}
	}

	baseDir := filepath.Join(h, ".config", "thot")
	_ = os.MkdirAll(baseDir, 0755)
	dbPath := filepath.Join(baseDir, "words.db")

	err = os.WriteFile(dbPath, dbBytes, 0644)
	if err != nil {
		return errors.WithMessage(err, "No se pudo copiar la Base de Datos")
	}

	fmt.Println(defStyle.Render("Base de datos copiada en " + dbPath))

	return nil
}
