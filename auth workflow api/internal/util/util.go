package util

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	n := rand.Intn(10)
	fmt.Println(n)
}
