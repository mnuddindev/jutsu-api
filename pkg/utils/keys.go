package utils

import (
	"errors"
	"regexp"
	"strconv"
)

type ErrorLoadingException struct{
	error
}

func NewErrorLoadingException(message string) error {
	return &ErrorLoadingException{error: errors.New(message)}
}

func MatchingKey(value, script string) (string, error) {
	re := regexp.MustCompile("," + regexp.QuoteMeta(value) + "=((?:0x)?([0-9a-fA-F]+))")
	match := re.FindStringSubmatch(script)
	if len(match) >= 3 {
		return match[2], nil
	}
	return "", NewErrorLoadingException("Failed to match the key")
}

func GetKeys(script string) [][]int {
	re := regexp.MustCompile(`case\s*0x[0-9a-f]+:(?![^;]*=partKey)\s*\w+\s*=\s*(\w+)\s*,\s*\w+\s*=\s*(\w+);`)
	matches := re.FindAllStringSubmatch(script, -1)
	var out [][]int
	for _, m := range matches {
		k1, err1 := MatchingKey(m[1], script)
		k2, err2 := MatchingKey(m[2], script)
		if err1 != nil || err2 != nil {
			continue
		}
		p1, errp1 := strconv.ParseInt(k1, 16, 32)
		p2, errp2 := strconv.ParseInt(k2, 16, 32)
		if errp1 != nil || errp2 != nil {
			continue
		}
		out = append(out, []int{int(p1), int(p2)})
	}
	return out
}
