package utils

import (
	"strings"
)

func EscapeLikeWildcards(src string) (res string) {
	var m = map[string]string{
		`_`: `\_`,
		`%`: `\%`,
		`'`: `\'`,
		`"`: `\"`,
		`<`: `\<`,
		`>`: `\>`,
		`&`: `\&`,
		`*`: `\*`,
	}
	for oldS, newS := range m {
		src = strings.ReplaceAll(src, oldS, newS)
	}
	return src
}
