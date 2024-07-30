package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

func WASender(data map[string]string) {

	url := viper.GetString("WA_GATEWAY_URL")
	token := viper.GetString("TOKEN_WA_GATEWAY")

	fmt.Println("token>>", token)
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url+"/chat/send/text", bytes.NewBuffer(jsonData))
	PanicIfError(err)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Token", token)
	req.Header.Set("Content-Type", "application/json")

	//Mengirim permintaan
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()
	// Membaca dan menampilkan respons
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	fmt.Println("Status Code:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
