//生成测试数据

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	o, err := os.Create("sampleData.txt")
	if err != nil {
		fmt.Printf("open error: %v\n", err)
	}

	for i := 0; i < 40000; i++ {
		// fmt.Println(randomString(10))
		o.WriteString(strings.ToLower(randomString(8)) + "\n")
	}

	o.Close()

}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
