package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

type char byte

const max = 1024

var dist [80][80]int
var dp [max][max]int
var delCost = 6
var insertCost = 1

var dict = [][12]char{
	{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '_'},
	{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', '*'},
	{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', '*', '*'},
	{'z', 'x', 'c', 'v', 'b', 'n', 'm', '*', '*', '*', '*'}}

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

func CalcEditDist(A string, B string) int {
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
			changeCost := 0
			changeCost = 1
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

func parallel(slices []string, total []string, resultChan chan []string) {
	results := []string{} //存放chunk的处理结果

	for _, slice := range slices { //chunk循环
		minDistance := 50
		var result string              //存放chunk中一条记录的处理结果
		closeNames := []string{}       //相似记录
		for _, single := range total { //遍历所有的记录
			distance := CalcEditDist(slice, single)      //计算距离
			if distance != 0 && distance < minDistance { //存在更小的距离，重置相似记录
				minDistance = distance
				closeNames = nil
				closeNames = append(closeNames, single)
			} else if distance == minDistance { //存放等距离的记录
				closeNames = append(closeNames, single)
			}
		}

		/*格式化一条记录*/
		result = result + slice + ":"
		for _, closeName := range closeNames {
			result = result + closeName + ","
		}
		result = result + strconv.Itoa(minDistance)
		results = append(results, result)
	}

	resultChan <- results //发送至channel
}

func main() {

	numCPU := runtime.NumCPU()

	fmt.Printf("We have %d CPUs\n", numCPU)

	runtime.GOMAXPROCS(numCPU)

	calcWeight()

	//将记录读取至names数组

	f, err := os.Open("sampleData.txt")
	if err != nil {
		fmt.Printf("open error: %v\n", err)
	}

	names := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		names = append(names, line)
	}

	nameCount := len(names)
	fmt.Println(nameCount)

	chunkSize := 100
	//对任务进行分解，并行处理。分解成每个单元100个需要2.85，那么需要分解成N/100个单元
	taskCount := nameCount / chunkSize
	taskRemain := nameCount % chunkSize

	fmt.Printf("%dTasks, %dRemain\n", taskCount, taskRemain)

	resultChan := make(chan []string, 2)

	for i := 0; i < taskCount; i++ {
		fmt.Println(i*chunkSize, (i+1)*chunkSize-1)
		go parallel(names[i*chunkSize:(i+1)*chunkSize-1], names, resultChan)

		results := <-resultChan

		for _, result := range results {
			fmt.Println(result)
		}
	}
}
