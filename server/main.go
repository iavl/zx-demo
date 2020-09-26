package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
)

type PrivKey struct {
	Lamb int `json:"lambda"`
	Miu  int `json:"miu"`
}

type PubKey struct {
	N int `json:"n"`
	G int `json:"g"`
}

type RSAKeyPair struct {
	PrivK PrivKey `json:"pri"`
	PubK  PubKey  `json:"pub"`
}

type RetGetRSAKeyPair struct {
	Data RSAKeyPair `json:"data"`
}

type RetEncryptCompute struct {
	Result []int  `json:"result"`
	Log    string `json:"log"`
}

func GetRSAKeyPair(w http.ResponseWriter, req *http.Request) {
	ret := new(RetGetRSAKeyPair)
	//id := req.FormValue("id")
	//id := req.PostFormValue('id')
	//fmt.Println(id)

	pri := PrivKey{3312616468, 1437516892}
	pub := PubKey{3312731593, 3312731594}
	ret.Data = RSAKeyPair{pri, pub}

	retJson, _ := json.Marshal(ret)

	io.WriteString(w, string(retJson))
}
func EncryptCompute(w http.ResponseWriter, req *http.Request) {
	ret := new(RetEncryptCompute)
	//id := req.FormValue("id")
	//id := req.PostFormValue('id')
	//fmt.Println(id)

	randArr := make([]int, 30, 30)
	for i, _ := range randArr {
		randArr[i] = rand.Intn(1000000)
	}
	ret.Result = randArr
	ret.Log = "success"
	retJson, _ := json.Marshal(ret)

	io.WriteString(w, string(retJson))
}

func main() {
	http.HandleFunc("/api/compute", EncryptCompute)
	http.HandleFunc("/api/get_rsa", GetRSAKeyPair)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
