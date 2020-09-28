package demo

import (
	"bufio"
	"fmt"
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
		fmt.Println(fmt.Sprintf("===================== 第 %d 组测试数据 =====================", idx))
		// generate taskId
		taskId := fmt.Sprintf("00000000000000000000000000000000000000000000000000000000000000%02d", idx)
		fmt.Println(fmt.Sprintf("task id: %s", taskId))

		utils.ClearResult(taskId)

		// 3. call contract to do paillier add
		for i, item := range items {
			//fmt.Println(fmt.Sprintf("%d", item))
			cipherText, _ := pk.Encrypt(item)
			fmt.Println(fmt.Sprintf("机构 %d, 明文贷款额：%d --> 加密密文：%d", i, item, cipherText))
			utils.PaillerAdd(taskId, cipherText)

			//break
		}

		// 4. query result from contract
		result := utils.QueyrPaillierResult(taskId)
		fmt.Println(fmt.Sprintf("===================== 第 %d 组测试结果 =====================", idx))
		fmt.Println(fmt.Sprintf("合约计算出的结果: %v", result))

		// 5. decrypt result
		// Test the homomorphic property
		sum, err := sk.Decrypt(result)
		if err != nil {
			t.Errorf("decrypt failed: %v", err.Error())
			return
		}

		fmt.Println(fmt.Sprintf("使用RSA私钥解密后的结果: [%d]", sum))
		fmt.Println(fmt.Sprintf("=========================================================="))

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
