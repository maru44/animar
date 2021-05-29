package tools

import (
	"os"
)

func GetenvOrDefault(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ConvertUintsToInts(lst []uint8) []int {
	var ret []int
	for _, u := range lst {
		i := int(u)
		ret = append(ret, i)
	}
	return ret
}
