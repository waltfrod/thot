//go:build (unix && ignore) || !windows
// +build unix,ignore !windows

package util

import (
	"syscall"
	"unsafe"
)

// GetConsoleSize devuelve el número actual de columnas y filas de la ventana de consola activa.
// El valor de retorno de esta función está en el orden de `columnas, filas`.
func GetConsoleSize() (int, int) { // {{{
	var sz struct {
		rows    uint16
		cols    uint16
		xpixels uint16
		ypixels uint16
	}
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&sz)))
	return int(sz.cols), int(sz.rows)
} // }}}
