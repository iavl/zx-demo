package utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os/exec"

	types "github.com/iavl/zx-demo"

	"github.com/iavl/zx-demo/paillier"
)

func GetRSAKeyPair() (*paillier.PublicKey, *paillier.PrivateKey) {
	pk, sk, _ := paillier.GenerateKeyPair(32)

	N, g := pk.ToDecimalString()
	fmt.Println(fmt.Sprintf("RSA公钥：n: %s g: %s", N, g))
	fmt.Println(fmt.Sprintf("RSA N2: %s", pk.N2.Text(10)))

	mu, lam := sk.ToDecimalString()
	fmt.Println(fmt.Sprintf("RSA私钥：λ: %s μ: %s", lam, mu))

	return pk, sk
}

func SetN2(n2 *big.Int) {
	command := "../cli/setN2.sh"

	fmt.Println(n2.String())
	cmd := exec.Command("/bin/bash", command, "alice", types.ContractBech32Addr, n2.String())
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
	command := "../cli/ClearResult.sh"

	cmd := exec.Command("/bin/bash", command, "alice", types.ContractBech32Addr, taskId)
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

func PaillerAdd(taskId string, value *big.Int) {

	command := "../cli/PaillerAdd.sh"

	cmd := exec.Command("/bin/bash", command, "alice", types.ContractBech32Addr, taskId, value.String())
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
}

func QueyrPaillierResult(taskId string) (result *big.Int) {
	type QueryResult struct {
		Gas    int64
		Result []*big.Int
	}

	command := "../cli/queyrPaillierResult.sh"

	cmd := exec.Command("/bin/bash", command, types.ContractBech32Addr, taskId)
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

func PaillerMain(pk *paillier.PublicKey, sk *paillier.PrivateKey, dataList []int64, idx int) {

	fmt.Println(fmt.Sprintf("===================== 第 %d 组测试数据 =====================", idx))
	// generate taskId
	taskId := fmt.Sprintf("00000000000000000000000000000000000000000000000000000000000000%02d", idx)
	fmt.Println(fmt.Sprintf("task id: %s", taskId))

	ClearResult(taskId)

	// 3. call contract to do paillier add
	for i, item := range dataList {
		//fmt.Println(fmt.Sprintf("%d", item))
		cipherText, _ := pk.Encrypt(item)
		fmt.Println(fmt.Sprintf("机构 %d, 明文贷款额：%d --> 加密密文：%d", i, item, cipherText))
		PaillerAdd(taskId, cipherText)

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
		fmt.Println(fmt.Errorf("decrypt failed: %v", err.Error()))
		return
	}

	fmt.Println(fmt.Sprintf("使用RSA私钥解密后的结果: [%d]", sum))
	fmt.Println(fmt.Sprintf("=========================================================="))

}
