package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/markcheno/go-talib"
)

const (
	apiURL = "https://www.alphavantage.co/query"
	apiKey = "IM6C7CE0Q9TDNX71"
)

type HistoricalData struct {
	MetaData   map[string]string            `json:"Meta Data"`
	TimeSeries map[string]map[string]string `json:"Time Series (60min)"`
}

type ForexTimeSeries struct {
	MetaData   map[string]string
	TimeSeries map[string]struct {
		Open   string
		High   string
		Low    string
		Close  string
		Volume string
	}
}

func ForexHistoricalData(symbol string) error {
	data, err := forexFetchData(symbol)
	if err != nil {
		return err
	}
	err = saveToCSV(data, fmt.Sprintf("./Data/forex-%v.csv", symbol))
	return err
}

func forexFetchData(symbol string) (*ForexTimeSeries, error) {
	url := fmt.Sprintf("%s?function=TIME_SERIES_INTRADAY&symbol=%s&interval=60min&apikey=%s", apiURL, symbol, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()

	var data HistoricalData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	historicalData := &ForexTimeSeries{
		MetaData: data.MetaData,
		TimeSeries: make(map[string]struct {
			Open   string
			High   string
			Low    string
			Close  string
			Volume string
		}),
	}

	for timestamp, values := range data.TimeSeries {
		t, err := time.Parse("2006-01-02 15:04:05", timestamp)
		if err != nil {
			continue
		}
		if time.Since(t).Hours() <= 720 {
			historicalData.TimeSeries[timestamp] = struct {
				Open   string
				High   string
				Low    string
				Close  string
				Volume string
			}{
				Open:   values["1. open"],
				High:   values["2. high"],
				Low:    values["3. low"],
				Close:  values["4. close"],
				Volume: values["5. volume"],
			}
		}
	}

	return historicalData, nil
}

func saveToCSV(data *ForexTimeSeries, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{"Timestamp", "Open", "High", "Low", "Close", "Volume", "RSI", "SMA"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	var prices []float64
	for _, v := range data.TimeSeries {
		p, _ := strconv.ParseFloat(v.Close, 64)
		prices = append(prices, p)
	}

	rsi := talib.Rsi(prices, 14)
	sma := talib.Sma(prices, 20)

	// Write data rows
	for timestamp, values := range data.TimeSeries {
		row := []string{
			timestamp,
			values.Open,
			values.High,
			values.Low,
			values.Close,
			values.Volume,
			strconv.FormatFloat(rsi[0], 'f', -1, 64),
			strconv.FormatFloat(sma[0], 'f', -1, 64),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing CSV row: %v", err)
		}
	}

	return nil
}

func CheckForexExist(symbol string) (bool, error) {
	url := fmt.Sprintf("%s?function=TIME_SERIES_INTRADAY&symbol=%s&interval=60min&apikey=%s", apiURL, symbol, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making API request:", err)
		return false, err
	}
	defer resp.Body.Close()

	var data HistoricalData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding response:", err)
		return false, err
	}

	_, exists := data.MetaData["1. Information"]
	return exists, nil
}
