package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/markcheno/go-talib"
)

var (
	URl = "https://api.binance.com/api/v3"
)

type BinanceData struct {
	CloseTime time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	RSI       float32
	SMA       float32
}

type SymbolInfo struct {
	Symbol string `json:"symbol"`
}

type ExchangeInfo struct {
	Symbols []SymbolInfo `json:"symbols"`
}

type NowPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (bd *BinanceData) UnmarshalJSON(data []byte) error {
	var tmp []interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	closeTime := tmp[6].(float64)
	bd.CloseTime = time.Unix(int64(closeTime/1000), 0)

	open := tmp[1].(string)
	bd.Open, _ = strconv.ParseFloat(open, 64)

	high := tmp[2].(string)
	bd.High, _ = strconv.ParseFloat(high, 64)

	low := tmp[3].(string)
	bd.Low, _ = strconv.ParseFloat(low, 64)

	close := tmp[4].(string)
	bd.Close, _ = strconv.ParseFloat(close, 64)

	vol := tmp[5].(string)
	bd.Volume, _ = strconv.ParseFloat(vol, 64)

	return nil
}

func CheckSymbolExist(symbol string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%v/exchangeInfo", URl))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var exchangeInfo ExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&exchangeInfo); err != nil {
		return false, err
	}

	for _, s := range exchangeInfo.Symbols {
		if s.Symbol == symbol {
			return true, nil
		}
	}
	return false, nil
}

func GetNowPrice(symbol string) float64 {
	resp, err := http.Get(fmt.Sprintf("%v/ticker/price?symbol=%v", URl, symbol))
	if err != nil {
		return 0.0
	}
	defer resp.Body.Close()

	var price NowPrice
	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		return 0.0
	}

	result, _ := strconv.ParseFloat(price.Price, 64)
	return result
}

func GetSymbolHistory(symbol, interval string) error {
	resp, err := http.Get(fmt.Sprintf("%v/klines?symbol=%v&interval=%v&limit=%v", URl, symbol, interval, 720))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var binanceData []BinanceData

	if err := json.Unmarshal(body, &binanceData); err != nil {
		return err
	}

	if len(binanceData) > 0 {
		binanceData = binanceData[:len(binanceData)-1]
	}

	var prices []float64
	for _, price := range binanceData {
		prices = append(prices, price.Close)
	}

	var rsi []float64
	var sma []float64

	if interval == "1h" {
		rsi = talib.Rsi(prices, 14)
		sma = talib.Sma(prices, 21)
	} else if interval == "4h" {
		rsi = talib.Rsi(prices, 28)
		sma = talib.Sma(prices, 11)
	} else {
		rsi = talib.Rsi(prices, 14)
		sma = talib.Sma(prices, 21)
	}

	file, err := os.Create(fmt.Sprintf("./Data/%v.csv", symbol))
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{"Time", "Open", "High", "Low", "Close", "Volume", "RSI", "SMA"})
	for i, b := range binanceData {
		openTime := b.CloseTime.Add(time.Second).Format("2006-01-02 15:04:05")
		open := strconv.FormatFloat(b.Open, 'f', -1, 64)
		high := strconv.FormatFloat(b.High, 'f', -1, 64)
		low := strconv.FormatFloat(b.Low, 'f', -1, 64)
		close := strconv.FormatFloat(b.Close, 'f', -1, 64)
		volume := strconv.FormatFloat(b.Volume, 'f', -1, 64)
		rsi := strconv.FormatFloat(rsi[i], 'f', -1, 64)
		sma := strconv.FormatFloat(sma[i], 'f', -1, 64)
		w.Write([]string{openTime, open, high, low, close, volume, rsi, sma})
	}

	return nil
}
