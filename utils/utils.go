package utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os/exec"
	"testing"

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

func SetN2(t *testing.T, n2 *big.Int) {
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

func ClearResult(t *testing.T, taskId string) {
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

func PaillerAdd(t *testing.T, taskId string, value *big.Int) {

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
