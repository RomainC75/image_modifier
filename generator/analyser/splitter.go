package analyser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/RomainC75/crypto_socket/generator/tools"
	"github.com/RomainC75/crypto_socket/generator/trades"
)

func Splitter(done chan int, in <-chan trades.Ticker, cryptos []string) []chan trades.Ticker {
	tickersOut := make([]chan trades.Ticker, len(cryptos))
	fmt.Println("pre affectation", tickersOut)

	for i := 0; i < len(cryptos); i++ {
		tickersOut[i] = make(chan trades.Ticker)
	}

	go func() {
		for {
			fmt.Println("waiting")
			select {
			case <-done:
				fmt.Println("DONE !")
				// close channels
				return
			case m := <-in:
				fmt.Println("<-in : ", m)
				index, err := getIndexInCryptos(cryptos, m.Symbol)
				fmt.Println("=", m, index, err)
				if i, err := getIndexInCryptos(cryptos, m.Symbol); err == nil {
					fmt.Println("==>", i, err)
					tickersOut[i] <- m
				}
			}
		}
	}()

	return tickersOut
}

func getIndexInCryptos(cryptos []string, name string) (int, error) {
	for i, crypto := range cryptos {
		if strings.ToLower(crypto) == strings.ToLower(name) {
			return i, nil
		}
	}
	return -1, errors.New("not found")
}

func SingleCryptoListener(i int, done <-chan int, wg *sync.WaitGroup, ch <-chan trades.Ticker) {
	go func() {
		maxLength := 5
		values := make([]float64, maxLength)
		sl := values[:0]
		for {
			select {
			case <-done:
				wg.Done()
				return
			case m := <-ch:
				value, err := strconv.ParseFloat(m.Price, 64)
				if err == nil {
					tools.InsertValue[float64](&sl, value, maxLength)
				}
				if i == 0 {
					fmt.Println("SLICE : ", sl)
					fmt.Printf("index: %d // crypto : %s // value : %s\n", i, m.Symbol, m.Price)
				}
			}
		}
	}()
}
