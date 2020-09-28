package demo

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/iavl/zx-demo/utils"
)

func TestPaillierDemo(t *testing.T) {
	// 1. generate pk and sk
	pk, sk := utils.GetRSAKeyPair()

	utils.SetN2(pk.N2)

	// 2. read data from file
	lines := readDataFromFile("data.txt")
	for idx, items := range lines {

		utils.PaillerMain(pk, sk, items, idx)
		break
		time.Sleep(2)
	}
}

func readDataFromFile(filePath string) (results [][]int64) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	for {
		// read data by line
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		//fmt.Println(fmt.Sprintf("new line: %s", line))

		// trim '[' / ']'
		line = strings.Trim(line, "[")
		line = strings.Trim(line, "]")

		stringSlice := strings.Split(line, ",")
		var result []int64
		for _, s := range stringSlice {
			s = strings.TrimSpace(s)

			i64, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic(err)
			}
			result = append(result, i64)
		}
		results = append(results, result)

		if err != nil && err == io.EOF {
			break
		}

	}
	return results
}
