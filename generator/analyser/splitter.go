package analyser

import (
	"errors"
	"strings"

	"github.com/RomainC75/crypto_socket/generator/trades"
)

func Splitter(done chan int, in <-chan trades.Ticker, cryptos []string) []chan trades.Ticker {
	tickersOut := make([]chan trades.Ticker, len(cryptos))

	for i := range cryptos {
		tickersOut[i] = make(chan trades.Ticker)
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case message := <-in:
				if index, err := getIndexInCryptos(cryptos, message.Symbol); err != nil {
					tickersOut[index] <- message
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
