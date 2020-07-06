package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	Charset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	RandLength = 5
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateCode(prefix string) string {
	t := time.Now()
	y := ""
	if t.Year()%100 < 10 {
		y = fmt.Sprintf("0%d", t.Year()%100)
	} else {
		y = fmt.Sprintf("%d", t.Year()%100)
	}
	m := ""
	if t.Month() < 10 {
		m = fmt.Sprintf("0%d", t.Month())
	} else {
		m = fmt.Sprintf("%d", t.Month())
	}
	d := ""
	if t.Day() < 10 {
		d = fmt.Sprintf("0%d", t.Day())
	} else {
		d = fmt.Sprintf("%d", t.Day())
	}
	code := fmt.Sprintf("%s%s%s%s%s", prefix, y, m, d, genStringWithLength(RandLength))
	return strings.ToUpper(code)
}

func stringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = Charset[seededRand.Intn(len(Charset))]
	}
	return string(b)
}

func genStringWithLength(length int) string {
	return stringWithCharset(length)
}
