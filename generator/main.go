package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

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

	multiCryptoChan := analyser.Splitter(done, dataOut, cryptos)

	var wg sync.WaitGroup
	fmt.Println("len", len(multiCryptoChan))
	for i, ch := range multiCryptoChan {
		fmt.Println("inside llop : ", i)
		wg.Add(1)
		go func() {
			for {
				select {
				case <-done:
					wg.Done()
					return
				case m := <-ch:
					fmt.Printf("index: %d // crypto : %s // value : %s\n", i, m.Symbol, m.Price)
				}
			}
		}()
	}
	fmt.Printf("after loop\n")
	for {
		time.Sleep(time.Second * 10)
		wg.Wait()
		done <- 1

	}

	fmt.Println("suite")

}
