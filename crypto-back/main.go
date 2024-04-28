package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

type CoinInfo struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Price  float64 `json:"current_price"`
}

var c = cache.New(10*time.Minute, 10*time.Minute) 

func getCoinPrices() ([]CoinInfo, error) {
	cached, found := c.Get("coinPrices")
	if found {
		fmt.Println("Using cached data.")
		return cached.([]CoinInfo), nil
	}

	resp, err := http.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var coins []CoinInfo
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return nil, err
	}

	c.Set("coinPrices", coins, cache.DefaultExpiration)
	return coins, nil
}

func getSpecificCoin(coinID string) (*CoinInfo, error) {
	coins, err := getCoinPrices()
	if err != nil {
		return nil, err
	}

	for _, coin := range coins {
		if coin.ID == coinID {
			return &coin, nil
		}
	}

	return nil, fmt.Errorf("coin with ID '%s' not found", coinID)
}

func main() {
	coins, err := getCoinPrices()
	if err != nil {
		fmt.Println("Error fetching coin prices:", err)
		return
	}
	fmt.Println("All coins:", coins)

	specificCoin, err := getSpecificCoin("bitcoin")
	if err != nil {
		fmt.Println("Error fetching specific coin:", err)
		return
	}
	fmt.Println("Specific coin:", specificCoin)
}
