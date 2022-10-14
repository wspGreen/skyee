package sfmt

import "fmt"

func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}
