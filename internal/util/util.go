package util

import (
	"fmt"
	"os"
	"strings"
)

func Pause(show bool) { // {{{
	if show {
		fmt.Print("Presione Enter para continuar...")
	}
	_, _ = os.Stdin.Read(make([]byte, 1))
} // }}}

func IF[T any](cond bool, t, f T) T { // {{{
	if cond {
		return t
	}
	return f
} // }}}

// InSlice busca un elemento en una lista, devolviendo su indice o -1 si no lo encuentra,
// el criterio de comparación es la función `fn`
func InSlice(fn func(int) bool, slen int) int { // {{{
	for i := 0; i < slen; i++ {
		if fn(i) {
			return i
		}
	}

	return -1
} // }}}

// Wrap divide un texto en palabras de un número de caracteres determinado
func Wrap(texto string, anchoMaximo int) string { // {{{
	// Dividir el texto en palabras
	palabras := strings.Fields(texto)

	// Variables para el resultado y el ancho actual
	resultado := ""
	anchoActual := 0

	// Recorrer las palabras
	for _, palabra := range palabras {
		// Verificar si la palabra excede el ancho máximo
		if anchoActual+len(palabra)+1 > anchoMaximo {
			// Agregar nueva línea
			resultado += "\n"
			anchoActual = 0
		}

		// Agregar la palabra y espacio al resultado
		resultado += palabra + " "
		anchoActual += len(palabra) + 1
	}

	// Eliminar el espacio final y devolver el resultado
	return strings.TrimSpace(resultado)
} // }}}
