package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Conf Config

type Config struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
	CryptoApi string `json:"CryptoApi"`
}

func ConfigInit() {

	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Error reading config file")
		panic(err)
	}

	err = json.Unmarshal(file, &Conf)
	if err != nil {
		return
	}

}
