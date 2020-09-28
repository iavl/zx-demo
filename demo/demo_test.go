package demo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/iavl/zx-demo/paillier"

	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	abiFilePath = "/Users/zhuliting/go/zx-demo/contract/paillier.abi"
	yamlPath    = "/Users/zhuliting/go/zx-demo/config/sdk.yaml"

	contractBech32Addr = "nch1n5u896f2lrz9wf34z3xnquzqcprcys072u0mzq"
)

var (
	amount = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))
)

type CallResult struct {
	Height string `json:"height"`
	Txhash string `json:"txhash"`
}

func TestPaillierDemo(t *testing.T) {
	// 1. generate pk and sk
	pk, sk := getRSAKeyPair()
	//pk, _ := getRSAKeyPair()

	SetN2(t, pk.N2)

	// 2. read data from file
	lines := readDataFromFile("data.txt")

	for idx, items := range lines {
		fmt.Println(fmt.Sprintf("===================== 第 %d 组测试数据 =====================", idx))
		// generate taskId
		taskId := fmt.Sprintf("00000000000000000000000000000000000000000000000000000000000000%02d", idx+1)
		fmt.Println(fmt.Sprintf("task id: %s", taskId))

		clear(t, taskId)

		// 3. call contract to do paillier add
		for i, item := range items {
			//fmt.Println(fmt.Sprintf("%d", item))
			cipherText, _ := pk.Encrypt(item)
			fmt.Println(fmt.Sprintf("机构 %d, 明文贷款额：%d --> 加密密文：%d", i, item, cipherText))
			paillerAdd(t, taskId, cipherText)

			//break
		}

		// 4. query result from contract
		result := QueyrPaillierResult(taskId)
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

func clear(t *testing.T, taskId string) {
	command := "../cli/clear.sh"

	cmd := exec.Command("/bin/bash", command, "alice", contractBech32Addr, taskId)
	fmt.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Sprintf("Execute Command failed: %v", err))
		return
	}
	// result example: {"Gas":8516,"Result":[0]}
	//fmt.Println(fmt.Sprintf("clear response: %s", string(output)))
	var res CallResult
	err = json.Unmarshal(output, &res)
	if err != nil {
		fmt.Println(fmt.Sprintf("unmarshal result failed: %v", err))
		return
	}

	fmt.Println(fmt.Sprintf("clear txhash: %v", res.Txhash))
}

func SetN2(t *testing.T, n2 *big.Int) {
	command := "../cli/setN2.sh"

	fmt.Println(n2.String())
	cmd := exec.Command("/bin/bash", command, "alice", contractBech32Addr, n2.String())
	fmt.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Sprintf("Execute Command failed: %v", err))
		return
	}
	// result example: {"Gas":8516,"Result":[0]}
	//fmt.Println(fmt.Sprintf("setN2 response: %s", string(output)))
	var res CallResult
	err = json.Unmarshal(output, &res)
	if err != nil {
		fmt.Println(fmt.Sprintf("unmarshal result failed: %v", err))
		return
	}

	fmt.Println(fmt.Sprintf("SetN2 txhash: %v", res.Txhash))
}

func paillerAdd(t *testing.T, taskId string, value *big.Int) {

	command := "../cli/paillerAdd.sh"

	cmd := exec.Command("/bin/bash", command, "alice", contractBech32Addr, taskId, value.String())
	fmt.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Sprintf("Execute Command failed: %v", err))
		return
	}
	// result example: {"Gas":8516,"Result":[0]}
	var res CallResult
	err = json.Unmarshal(output, &res)
	if err != nil {
		fmt.Println(fmt.Sprintf("unmarshal result failed: %v", err))
		return
	}

	fmt.Println(fmt.Sprintf("paillerAdd txhash: %v", res.Txhash))
}

func getRSAKeyPair() (*paillier.PublicKey, *paillier.PrivateKey) {
	pk, sk, _ := paillier.GenerateKeyPair(32)

	N, g := pk.ToDecimalString()
	fmt.Println(fmt.Sprintf("RSA公钥：n: %s g: %s", N, g))
	fmt.Println(fmt.Sprintf("RSA N2: %s", pk.N2.Text(10)))

	mu, lam := sk.ToDecimalString()
	fmt.Println(fmt.Sprintf("RSA私钥：λ: %s μ: %s", lam, mu))

	return pk, sk
}

func QueyrPaillierResult(taskId string) (result *big.Int) {
	type QueryResult struct {
		Gas    int64
		Result []*big.Int
	}

	command := "../cli/queyrPaillierResult.sh"

	cmd := exec.Command("/bin/bash", command, contractBech32Addr, taskId)
	fmt.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Sprintf("Execute Command failed: %v", err))
		return
	}
	// result example: {"Gas":8516,"Result":[0]}
	//fmt.Println(fmt.Sprintf("%s", string(output)))

	var res QueryResult
	err = json.Unmarshal(output, &res)
	if err != nil {
		fmt.Println(fmt.Sprintf("unmarshal result failed: %v", err))
		return
	}

	result = res.Result[0]
	fmt.Println(fmt.Sprintf("result: %d", result))

	return result
}
