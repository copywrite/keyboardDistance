/*
  Thansk ANan, Siritinga, Nick Craig-Wood
  lots of room to improve
*/


package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type char byte

const max = 20 //modified here

var dist [80][80]int

var delCost = 6
var insertCost = 1

var dict = [][12]char{
	{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '_'},
	{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', '*'},
	{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', '*', '*'},
	{'z', 'x', 'c', 'v', 'b', 'n', 'm', '*', '*', '*', '*'}}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	calcWeight()

	f, err := os.Open("staff.txt")
	if err != nil {
		fmt.Printf("open error: %v\n", err)
	}

	o, err := os.Create("result-keyboard.txt")
	if err != nil {
		fmt.Printf("open error: %v\n", err)
	}

	var names []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		start := strings.LastIndex(line, ",") + 1
		end := strings.LastIndex(line, "(")
		name := line[start:end]
		names = append(names, name)
	}

	nameCount := len(names)
	procs := 2
	chunkSize := nameCount / procs
	// remain := nameCount % chunks

	// for t := 0; t < nameCount*2; t++ {
	// 	o.WriteString(<-c + "\r\n")
	// }

	ch := make(chan string, 2)
	var wg sync.WaitGroup

	t0 := time.Now()

	for i := 0; i < procs; i++ {
		//chunk bounds
		start := i * chunkSize
		end := (i+1)*chunkSize - 1

		wg.Add(1)
		go func(slices []string, allnames []string) {
			var dp = [max][max]int{} //modified here

			for _, slice := range slices {
				minDistance := 256
				distance := 0
				sum := 0
				closeNames := []string{}

				for _, name := range allnames {
					distance = calcEditDist(slice, name, &dp)
					sum += 1

					// fmt.Println(slice, name, distance, minDistance)

					if distance > 0 {
						if distance < minDistance {
							// fmt.Println(slice, name, distance, distance > 0)
							minDistance = distance
							closeNames = nil
							closeNames = append(closeNames, name)
						} else if distance == minDistance {
							// fmt.Println(slice, name, distance, distance > 0)
							closeNames = append(closeNames, name)
						}
					}
				}

				/*make array into a single line*/
				var singleLine = slice + " "
				for _, closeName := range closeNames {
					singleLine += closeName + ","
				}

				singleLine += strconv.Itoa(minDistance)

				ch <- singleLine
			}
			wg.Done()
		}(names[start:end], names)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := range ch {
		o.WriteString(i + "\r\n")
	}

	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

	defer f.Close()
	defer o.Close()
}

func iAbs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func calcWeight() {
	// fmt.Println(dict)
	for i := 0; i < 4; i++ {
		for j := 0; j < 11; j++ {
			if dict[i][j] == '*' {
				break
			} //遍历每个不是*的字符
			for x := 0; x < 4; x++ {
				for y := 0; y < 11; y++ {
					if dict[x][y] != '*' {
						// fmt.Println(dict[i][j], dict[i][j]-'0', dict[x][y], dict[x][y]-'0')
						dist[dict[i][j]-'0'][dict[x][y]-'0'] = iAbs(i-x) + iAbs(j-y)
					}
				}
			}
		}
	}
}

func getMin(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func toLower(a uint8) uint8 {
	if a >= 'A' && a <= 'Z' {
		return a + 32
	}
	return a
}

func calcEditDist(A string, B string, dp *[max][max]int) int {
	lenA := len(A)
	lenB := len(B)

	for i := 1; i <= lenA; i++ {
		dp[i][0] = i * insertCost
	}

	for j := 1; j <= lenB; j++ {
		dp[0][j] = j * insertCost
	}

	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			changeCost := 1
			if A[i-1] != B[j-1] {
				changeCost = dist[toLower(A[i-1])-'0'][toLower(B[j-1])-'0']
			} else {
				changeCost = 0
			}
			dp[i][j] = getMin(dp[i-1][j]+delCost, getMin(dp[i][j-1]+delCost, dp[i-1][j-1]+changeCost))
		}
	}

	return dp[lenA][lenB]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}
