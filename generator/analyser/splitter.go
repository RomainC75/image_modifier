package analyser

import (
	"errors"
	"strings"

	"github.com/RomainC75/crypto_socket/generator/trades"
)

func Splitter(done chan int, in <-chan trades.Ticker, cryptos []string) []chan trades.Ticker {
	tickersOut := make([]chan trades.Ticker, len(cryptos))

	for i := range cryptos {
		go func() {
			ch := make(chan trades.Ticker)
			tickersOut[i] = ch
			for {
				select {
				case <-done:
					return
				case message := <-in:
					if _, err := getIndexInCryptos(cryptos, message.Symbol); err != nil {
						ch <- message
					}
				}
			}
		}()
	}

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
