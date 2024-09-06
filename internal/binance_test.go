package internal

import "testing"

type CheckExist struct {
	symbol   string
	expected bool
}

var CheckExistTest = []CheckExist{
	CheckExist{symbol: "BTCUSDT", expected: true},
	CheckExist{symbol: "ETHUSDT", expected: true},
	CheckExist{symbol: "BTCUSD", expected: false},
	CheckExist{symbol: "ETH", expected: false},
}

func CheckSymbolExistTest(t *testing.T) {
	for _, test := range CheckExistTest {
		if output, _ := CheckSymbolExist(test.symbol); output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}
