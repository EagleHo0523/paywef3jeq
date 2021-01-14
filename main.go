package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	r "./router"
	verify "./verification"
)

func main() {
	paymentService("8080")
	// verification()
}

func paymentService(port string) {
	fmt.Println("service starting on port:", port)
	port = ":" + port
	router := r.NewRouter()
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal("Payment service listen and serve: ", err)
	}
}

func verification() {
	plainText := "{\"uid\":\"d6MZfLOvm6SGudhABJx8fGn6QRG3onCw\",\"txid\":\"wPnahCR8SRNRtjrH1UcGTt4WN75uJT5R\",\"amount\":\"0.00000001\"}"

	account := "0xa83c7f79f5BF4C0Cc883E44ca9c0987BdD94B7a1"
	pubKey := "046c75524035e4e4c4d232965e86ccb5daa060ecff0c450aa61223a31ee804daf5c368cdcf1c0774d9e2a991d047f6d760583365ef1a895531a41d8158c41765ca"
	fmt.Println("pubkey len:", len(pubKey))
	privKey := "1f2317482fff7ec9b4137479289babe80dd0162a3fea3494cfe6951fcc6cfbd5"
	fmt.Println("privkey len:", len(privKey))

	start := time.Now()
	// for i := 0; i < 1000; i++ {
	pubVerify, err := verify.ImportPubKey(pubKey)
	if err != nil {
		fmt.Println("#1", err)
	}

	ciphertext, err := pubVerify.EncryptWithCheck(account, plainText)
	if err != nil {
		fmt.Println("#2", err)
	}
	fmt.Println("ciphertext len:", len(ciphertext))
	loginVerify, err := verify.ImportPrivKey(privKey)
	if err != nil {
		fmt.Println("#3", err)
	}

	pt, err := loginVerify.DecryptWithCheck(account, ciphertext) // 4aaa3bd9cb189c668735fd65995bc117f1361bc82fbeb685dcbb47bafb60bbf3
	if err != nil {
		fmt.Println("#4", err)
	}
	fmt.Println("before plainText is", plainText, "after decrypt is", pt)
	// }
	during := time.Since(start)
	fmt.Println("during:", during)
}
