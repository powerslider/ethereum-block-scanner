package numbers

import (
	"fmt"
	"strconv"
	"strings"
)

// HexToInt parse hex string value to int
func HexToInt(value string) (int, error) {
	i, err := strconv.ParseInt(strings.TrimPrefix(value, "0x"), 16, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

// IntToHex convert int to hexadecimal representation
func IntToHex(i int) string {
	return fmt.Sprintf("0x%x", i)
}
