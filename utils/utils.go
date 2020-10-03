package utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"

	types "github.com/iavl/zx-demo"

	"github.com/iavl/zx-demo/paillier"
)

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func ToHexString(decimalString string) string {
	i, _ := strconv.ParseInt(decimalString, 10, 64)
	return strconv.FormatInt(i, 16)
}

func ParseBigInt(s string) (*big.Int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return big.NewInt(0), err
	}
	n := big.NewInt(int64(num))
	return n, nil
}

func GetRSAKeyPair() (*paillier.PublicKey, *paillier.PrivateKey) {
	pk, sk, _ := paillier.GenerateKeyPair(32)

	N, g := pk.ToDecimalString()
	fmt.Println(fmt.Sprintf("RSA公钥：n: %s g: %s", N, g))
	fmt.Println(fmt.Sprintf("RSA N2: %d", pk.N2))

	mu, lam := sk.ToDecimalString()
	fmt.Println(fmt.Sprintf("RSA私钥：λ: %s μ: %s", lam, mu))

	return pk, sk
}

func SetN2(n2 *big.Int) {
	dir, _ := os.Getwd()
	command := dir + "/cli/setN2.sh"

	cmd := exec.Command("/bin/bash", command, "alice", types.ContractBech32Addr, dir, n2.String())
	fmt.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Sprintf("Execute Command failed: %v", err))
		return
	}
	// result example: {"Gas":8516,"Result":[0]}
	//fmt.Println(fmt.Sprintf("setN2 response: %s", string(output)))
	var res types.CallResult
	err = json.Unmarshal(output, &res)
	if err != nil {
		fmt.Println(fmt.Sprintf("unmarshal result failed: %v", err))
		return
	}

	fmt.Println(fmt.Sprintf("SetN2 txhash: %v", res.Txhash))
}

func ClearResult(taskId string) {
	dir, _ := os.Getwd()
	command := dir + "/cli/clearResult.sh"

	cmd := exec.Command("/bin/bash", command, "alice", types.ContractBech32Addr, dir, taskId)
	fmt.Println(cmd.String())
	output, err := cmd.Output()

	if err != nil {
		fmt.Println(fmt.Sprintf("Execute Command failed: %v", err))
		return
	}
	// result example: {"Gas":8516,"Result":[0]}
	//fmt.Println(fmt.Sprintf("ClearResult response: %s", string(output)))
	var res types.CallResult
	err = json.Unmarshal(output, &res)
	if err != nil {
		fmt.Println(fmt.Sprintf("unmarshal result failed: %v", err))
		return
	}

	fmt.Println(fmt.Sprintf("ClearResult txhash: %v", res.Txhash))
}

func PaillerAdd(taskId string, value *big.Int, name string) (txHash string) {
	dir, _ := os.Getwd()
	command := dir + "/cli/paillerAdd.sh"

	cmd := exec.Command("/bin/bash", command, name, types.ContractBech32Addr, dir, taskId, value.String())
	fmt.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Sprintf("Execute Command failed: %v", err))
		return
	}
	// result example: {"Gas":8516,"Result":[0]}
	var res types.CallResult
	err = json.Unmarshal(output, &res)
	if err != nil {
		fmt.Println(fmt.Sprintf("unmarshal result failed: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("PaillerAdd txhash: %v", res.Txhash))

	txHash = res.Txhash
	return txHash
}

func QueryPaillierResult(taskId string) (result int64) {
	type QueryResult struct {
		Gas    int64
		Result []int64
	}

	dir, _ := os.Getwd()
	command := dir + "/cli/queryPaillierResult.sh"

	cmd := exec.Command("/bin/bash", command, types.ContractBech32Addr, dir, taskId)
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

func PaillerMain(pk *paillier.PublicKey, sk *paillier.PrivateKey, dataList []int64, taskId string) (cipherTextList []string, txHashList []string, encryptResult, decryptResult int64, err error) {
	fmt.Println(fmt.Sprintf("task id: %s", taskId))

	fmt.Println(fmt.Sprintf("N2: %x", pk.N2))
	ClearResult(taskId)

	SetN2(pk.N2)

	// 3. call contract to do paillier add
	cipherTextList = make([]string, 0)
	txHashList = make([]string, 0)
	for i, item := range dataList {
		//n, g := pk.ToString()
		//fmt.Println(fmt.Sprintf("data: %d, pub key, n: %s, g: %s", item, n, g))
		cipherText, _ := pk.Encrypt(item)
		fmt.Println(fmt.Sprintf("机构 %d, 明文贷款额：%d --> 加密密文：%d", i, item, cipherText))
		name := fmt.Sprintf("alice%d", i)
		fmt.Println(fmt.Sprintf("name: %s", name))
		txHash := PaillerAdd(taskId, cipherText, name)

		cipherTextList = append(cipherTextList, cipherText.String())
		txHashList = append(txHashList, txHash)
		//break
	}

	// 4. query result from contract
	encryptResult = QueryPaillierResult(taskId)
	fmt.Println(fmt.Sprintf("===================== taskID: %s 测试结果 =====================", taskId))
	fmt.Println(fmt.Sprintf("合约计算出的结果: %v", encryptResult))

	// 5. decrypt result
	// Test the homomorphic property

	decryptResult, err = sk.Decrypt(big.NewInt(encryptResult))
	if err != nil {
		fmt.Println(fmt.Errorf("decrypt failed: %v", err.Error()))
		return cipherTextList, txHashList, 0, 0, err
	}

	fmt.Println(fmt.Sprintf("使用RSA私钥解密后的结果: [%d]", decryptResult))
	fmt.Println(fmt.Sprintf("=========================================================="))
	return cipherTextList, txHashList, encryptResult, decryptResult, nil
}
