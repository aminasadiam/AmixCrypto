package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
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
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func ForexHistoricalData(symbol string) error {
	data, err := forexFetchData(symbol)
	if err != nil {
		return err
	}
	err = saveToCSV(data, fmt.Sprintf("./Data/forex-%v.csv", symbol))
	return err
}

func forexFetchData(symbol string) ([]ForexTimeSeries, error) {
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

	var result []ForexTimeSeries

	for timestamp, values := range data.TimeSeries {
		t, err := time.Parse("2006-01-02 15:04:05", timestamp)
		if err != nil {
			continue
		}
		open, _ := strconv.ParseFloat(values["1. open"], 64)
		low, _ := strconv.ParseFloat(values["3. low"], 64)
		high, _ := strconv.ParseFloat(values["2. high"], 64)
		close, _ := strconv.ParseFloat(values["4. close"], 64)
		volume, _ := strconv.ParseFloat(values["5. volum"], 64)
		result = append(result, ForexTimeSeries{
			Time:   t,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func saveToCSV(data []ForexTimeSeries, filename string) error {
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
	for _, v := range data {
		prices = append(prices, v.Close)
	}

	rsi := talib.Rsi(prices, 14)
	sma := talib.Sma(prices, 20)

	// Write data rows
	for i, values := range data {
		t := values.Time.Add(time.Hour * 3).Add(time.Minute * 30)
		row := []string{
			t.Format("2006-01-02 15:04:05"),
			strconv.FormatFloat(values.Open, 'f', -1, 64),
			strconv.FormatFloat(values.High, 'f', -1, 64),
			strconv.FormatFloat(values.Low, 'f', -1, 64),
			strconv.FormatFloat(values.Close, 'f', -1, 64),
			strconv.FormatFloat(values.Volume, 'f', -1, 64),
			strconv.FormatFloat(rsi[i], 'f', -1, 64),
			strconv.FormatFloat(sma[i], 'f', -1, 64),
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
