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
	colsInt, _ := strconv.ParseInt(sn[0], 10, 64)
	rowsInt, _ := strconv.ParseInt(sn[1], 10, 64)

	return int(colsInt), int(rowsInt)
}
