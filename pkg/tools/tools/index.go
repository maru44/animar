package tools

import (
	"animar/v1/configs"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const slugLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

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

func ErrorLog(err error) {
	if IsProductionEnv() {
		log.Println(err)
	} else {
		fmt.Println(err)
	}
}

func IsProductionEnv() bool {
	// 本番環境IPリスト
	hosts := strings.Split(configs.ProductionIpList, ", ")
	host, _ := os.Hostname()

	// if runtime.GOOS != "linux" {
	// 	return false
	// }
	for _, v := range hosts {
		if v == host {
			return true
		}
	}
	return false
}

// insert & update用
func NewNullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// insert & update用
func NewNullInt(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func GenRandSlug(n int) string {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = slugLetters[rand.Intn(len(slugLetters))]
	}
	return string(b)
}
