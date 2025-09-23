package asciicast

import (
	"strconv"
	"strings"
)

// ParseResize formatted as "{COLS}x{ROWS}"
func ParseResize(str string) (cols, rows int) {
	sn := strings.Split(str, "x")
	if len(sn) != 2 {
		return
	}
	cols, _ = strconv.Atoi(sn[0])
	rows, _ = strconv.Atoi(sn[1])

	return
}
