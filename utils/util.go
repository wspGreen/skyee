package utils

import (
	"reflect"
	"strings"
)

func GetClassName(refType reflect.Type) string {
	sType := refType.String()
	index := strings.Index(sType, ".")
	if index != -1 {
		sType = sType[index+1:]
	}
	return sType
}
