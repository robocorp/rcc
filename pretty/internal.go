package pretty

import "fmt"

const (
	Escape = 0x1b
)

func csi(value string) string {
	return fmt.Sprintf("%c[%s", Escape, value)
}
