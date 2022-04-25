package main

import (
	"encoding/json"
	"io/ioutil"
)

type LoginCred struct {
	Username string
	Password string
	IPs      []string
}

func SaveCred(cred *LoginCred) error {
	js, _ := json.Marshal(cred)
	return ioutil.WriteFile("ac.bin", js, 0755)
}

func LoadCred() (*LoginCred, error) {
	content, err := ioutil.ReadFile("ac.bin")
	if err != nil {
		return nil, err
	}

	thisCred := LoginCred{}
	err = json.Unmarshal(content, &thisCred)
	return &thisCred, err
}
