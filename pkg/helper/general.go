package helper

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	max_price float64 = 5.05
	min_price float64 = 111.83
)

func RandomIntegerCreator(min, max int) int {
	DelaySystem()
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func DelaySystem() {
	time.Sleep(500 * time.Nanosecond)
}

func StringUpper(data string) string {
	return strings.ToUpper(data)
}

func RandomFloatCreator() float64 {
	DelaySystem()
	rand.Seed(time.Now().UnixNano())
	return math.Round((rand.Float64()*(max_price-min_price)+min_price)*100) / 100
}

func DateDiffDay(fd, sd time.Time) int {
	return int((fd.Sub(sd)).Hours() / 24)
}
