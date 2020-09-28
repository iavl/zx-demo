package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

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
	Result []int  `json:"result"`
	Log    string `json:"log"`
}

type ReqBody struct {
	DataList []int   `json:"data_list"`
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
	}

	fmt.Println(fmt.Sprintf("data list: %v", body.DataList))
	fmt.Println(fmt.Sprintf("pri key: %v", body.PrivK))
	fmt.Println(fmt.Sprintf("pub key: %v", body.PubK))

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
