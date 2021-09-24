package converters

import (
    "fmt"
)

func ConvertImageAddr(file string) string {
	return fmt.Sprintf("/qccoo/my/images/%s", file)
}
