package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/RomainC75/crypto_socket/generator/analyser"
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

	done := make(chan int)

	dataOut, err := trades.SubscribeAndListen(cryptos)
	if err != nil {
		log.Fatal("could not connect")
	}

	var wg sync.WaitGroup
	for i, ch := range analyser.Splitter(done, dataOut, cryptos) {
		fmt.Println("inside llop : ", i)
		wg.Add(1)
		analyser.SingleCryptoListener(i, done, &wg, ch)

	}

	wg.Wait()
	done <- 1

}
