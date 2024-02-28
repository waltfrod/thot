package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	sqlite3 "github.com/mattn/go-sqlite3"
	"gitlab.com/tozd/go/errors"
)

type Lang string

const (
	English Lang = "dic_en"
	Spanish Lang = "dic_es"
)

type Database struct {
	db *sql.DB
}

// fnRegexp es una función auxiliar para SQLite3 ya que no soporta su implementación nativa
func fnRegexp(re, search string) bool { // {{{
	match, e := regexp.MatchString(re, search)
	if e != nil {
		return false
	}

	return match
} // }}}

func NewDatabase() (*Database, errors.E) { // {{{
	h, err := os.UserHomeDir()
	if err != nil {
		h = os.Getenv("HOME")
		if h == "" {
			return nil, errors.WithMessage(err, "No se pudo obtener el directorio HOME")
		}
	}

	dbPath := filepath.Join(h, ".config", "thot", "words.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, errors.WithDetails(
			errors.WithMessage(err, "No se pudo encontrar la base de datos"),
			"path",
			dbPath,
		)
	}

	sql.Register("sqlite3_with_regexp", &sqlite3.SQLiteDriver{
		ConnectHook: func(cnx *sqlite3.SQLiteConn) error {
			if err := cnx.RegisterFunc("regexg", fnRegexp, true); err != nil {
				return errors.WithMessage(err, "No se pudo registrar la función fnRegexp")
			}
			return nil
		},
	})

	db, err := sql.Open("sqlite3_with_regexp", dbPath)
	if err != nil {
		errors.WithMessage(err, "No se pudo conectar a la base de datos")
	}

	return &Database{db}, nil
} // }}}

func (d Database) Words(limit int, l Lang, re string) ([]string, errors.E) { //{{{
	var words []string = make([]string, 0)
	qry := fmt.Sprintf(
		"SELECT palabra FROM %s WHERE regexg('^[%s]+$', palabra) ORDER BY RANDOM() LIMIT ?",
		string(l),
		re,
	)
	rows, err := d.db.Query(
		qry,
		limit,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "No se pudo obtener las palabras")
	}
	defer rows.Close()

	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, errors.WithMessage(err, "No se pudo obtener la palabra")
		}
		words = append(words, word)
	}

	return words, nil
} //}}}
