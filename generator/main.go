package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/RomainC75/crypto_socket/generator/trades"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("could not get the .env file ! ")
	}
	t := os.Getenv("CRYPTOLIST")
	cryptos := strings.Split(t, ",")
	for i, crypto := range cryptos {
		cryptos[i] = strings.Trim(strings.Trim(crypto, "\\"), "\"")
	}
	fmt.Print(cryptos)

	trades.SubscribeAndListen(cryptos)
}
