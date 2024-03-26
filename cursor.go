package climenu

import (
	"fmt"
	"io"
)

// hide cursor
func cursorHide(w io.Writer) {
	w.Write([]byte("\033[?25l"))
}

// Display cursor
func cursorShow(w io.Writer) {
	w.Write([]byte("\033[?25h"))
}

// Save cursor position
func cursorSave(w io.Writer) {
	w.Write([]byte("\033[s"))
}

// restore cursor position
func cursorRestore(w io.Writer) {
	w.Write([]byte("\033[u"))
}

// Move cursor up n lines
func cursorMoveup(w io.Writer, n int) {
	fmt.Fprintf(w, "\033[%dA", n)
}

// Move cursor down n lines
func cursorMovedown(w io.Writer, n int) {
	fmt.Fprintf(w, "\033[%dB", n)
}

// Cursor moves n lines to the right
func cursorMoveright(w io.Writer, n int) {
	fmt.Fprintf(w, "\033[%dC", n)
}

// Move cursor to the left by n lines
func cursorMoveleft(w io.Writer, n int) {
	fmt.Fprintf(w, "\033[%dD", n)
}
