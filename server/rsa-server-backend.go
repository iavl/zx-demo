package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/iavl/zx-demo/paillier"

	"github.com/iavl/zx-demo/utils"
)

type PrivKey struct {
	Lamb string `json:"lambda"`
	Mu   string `json:"mu"`
}

type PubKey struct {
	N string `json:"n"`
	G string `json:"g"`
}

type RSAKeyPair struct {
	PrivK PrivKey `json:"pri_key"`
	PubK  PubKey  `json:"pub_key"`
}

type RetGetRSAKeyPair struct {
	Data RSAKeyPair `json:"data"`
}

type RetEncryptCompute struct {
	Result []int64 `json:"result"`
	Log    string  `json:"log"`
}

type ReqBody struct {
	DataList []int64 `json:"data_list"`
	PrivK    PrivKey `json:"pri_key"`
	PubK     PubKey  `json:"pub_key"`
}

func GetRSAKeyPair(w http.ResponseWriter, req *http.Request) {
	// 1. generate pk and sk
	pk, sk := utils.GetRSAKeyPair()
	N, g := pk.ToDecimalString()
	mu, lam := sk.ToDecimalString()
	//N2 := pk.N2.Text(10)

	pri := PrivKey{lam, mu}
	pub := PubKey{N, g}

	ret := new(RetGetRSAKeyPair)
	ret.Data = RSAKeyPair{pri, pub}

	retJson, _ := json.Marshal(ret)
	io.WriteString(w, string(retJson))
}

func EncryptCompute(w http.ResponseWriter, req *http.Request) {
	/*
		post data:
			{
				"data_list": [112, 333, 444, 555],
				"pri_key": {
					"lambda": "2774103120",
					"mu": "882170834"
				},
				"pub_key": {
					"n": "2774208617",
					"g": "2774208618"
				}
			}
	*/

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		io.WriteString(w, "request post data format error")
		return
	}

	var body ReqBody
	if err = json.Unmarshal(buf, &body); err != nil {
		io.WriteString(w, "unmarshal post data error")
		return
	}

	fmt.Println(fmt.Sprintf("data list: %v", body.DataList))
	fmt.Println(fmt.Sprintf("pri key: %v", body.PrivK))
	fmt.Println(fmt.Sprintf("pub key: %v", body.PubK))

	pk, err := paillier.NewPublicKey(utils.ToHexString(body.PubK.N), utils.ToHexString(body.PubK.G))
	if err != nil {
		io.WriteString(w, "parse pubkey error")
		return
	}

	sk, err := paillier.NewPrivateKey(utils.ToHexString(body.PrivK.Mu), utils.ToHexString(body.PrivK.Lamb), pk)
	if err != nil {
		io.WriteString(w, "parse priv key error")
		return
	}

	N, g := pk.ToDecimalString()
	fmt.Println(fmt.Sprintf("RSA公钥：n: %s g: %s", N, g))
	fmt.Println(fmt.Sprintf("RSA N2: %s", pk.N2.Text(10)))

	mu, lam := sk.ToDecimalString()
	fmt.Println(fmt.Sprintf("RSA私钥：λ: %s μ: %s", lam, mu))

	// generate taskId
	rand.Seed(time.Now().Unix())
	taskId := fmt.Sprintf("0000000000000000000000000000000000000000000000000000000000000%03d", rand.Intn(100))
	utils.PaillerMain(pk, sk, body.DataList, taskId)

	ret := new(RetEncryptCompute)
	ret.Result = body.DataList
	ret.Log = "success"
	retJson, _ := json.Marshal(ret)

	io.WriteString(w, string(retJson))
}

func main() {
	http.HandleFunc("/api/compute", EncryptCompute)
	http.HandleFunc("/api/get_rsa", GetRSAKeyPair)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
