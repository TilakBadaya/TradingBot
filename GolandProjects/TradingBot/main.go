package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

const (
	apiKey    string = "mxpu0k2av02w9egx"
	apiSecret string = "u00yc5qhxshwakvei5pcpuxr0ub69zw7"
)

var (
	kc           *kiteconnect.Client
	positionOpen bool
	long         bool
	instruments  []kiteconnect.Instrument
	quotes       []kiteconnect.Quote
)

type StockData struct {
	Price float64   `json:"Price"`
	Time  time.Time `json:"Time"`
}
type StockDataTemp struct {
	Price float64 `json:"Price"`
	Time  string  `json:"Time"`
}

func main() {
	//	var stockDataTemp []StockDataTemp
	//	var stockData []StockDataTemp
	//	//startSession()
	//	data := `[
	//  {
	//    "Price": 100.0,
	//    "Timestamp": "2022-01-01T00:00:00Z"
	//  },
	//  {
	//    "Price": 102.0,
	//    "Timestamp": "2022-01-02T00:00:00Z"
	//  },
	//  {
	//    "Price": 104.0,
	//    "Timestamp": "2022-01-03T00:00:00Z"
	//  },
	//  {
	//    "Price": 106.0,
	//    "Timestamp": "2022-01-04T00:00:00Z"
	//  },
	//  {
	//    "Price": 108.0,
	//    "Timestamp": "2022-01-05T00:00:00Z"
	//  },
	//  {
	//    "Price": 110.0,
	//    "Timestamp": "2022-01-06T00:00:00Z"
	//  },
	//  {
	//    "Price": 112.0,
	//    "Timestamp": "2022-01-07T00:00:00Z"
	//  },
	//  {
	//    "Price": 115.0,
	//    "Timestamp": "2022-01-08T00:00:00Z"
	//  },
	//  {
	//    "Price": 116.0,
	//    "Timestamp": "2022-01-09T00:00:00Z"
	//  },
	//  {
	//    "Price": 117.0,
	//    "Timestamp": "2022-01-10T00:00:00Z"
	//  },
	//  {
	//    "Price": 118.0,
	//    "Timestamp": "2022-01-11T00:00:00Z"
	//  },
	//  {
	//    "Price": 116.182,
	//    "Timestamp": "2022-01-12T00:00:00Z"
	//  },
	//  {
	//    "Price": 114.364,
	//    "Timestamp": "2022-01-13T00:00:00Z"
	//  },
	//  {
	//    "Price": 112.546,
	//    "Timestamp": "2022-01-14T00:00:00Z"
	//  },
	//  {
	//    "Price": 110.728,
	//    "Timestamp": "2022-01-15T00:00:00Z"
	//  },
	//  {
	//    "Price": 108.91,
	//    "Timestamp": "2022-01-16T00:00:00Z"
	//  },
	//  {
	//    "Price": 107.092,
	//    "Timestamp": "2022-01-17T00:00:00Z"
	//  },
	//  {
	//    "Price": 105.274,
	//    "Timestamp": "2022-01-18T00:00:00Z"
	//  },
	//  {
	//    "Price": 103.456,
	//    "Timestamp": "2022-01-19T00:00:00Z"
	//  },
	//  {
	//    "Price": 101.638,
	//    "Timestamp": "2022-01-20T00:00:00Z"
	//  }
	//]
	//`
	//	err := json.Unmarshal([]byte(data), &stockDataTemp)
	//	if err != nil {
	//		log.Println(err.Error())
	//		return
	//	}
	//	log.Println(stockDataTemp)
	//	//dataset := //[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	//	// Generate all combinations of A, B, C, and D values
	//	combinations := generateCombinationsThreeSwings(stockData)
	//	// Check each combination for Butterfly pattern
	//	for _, combination := range combinations {
	//		if isCrabPattern(combination) {
	//			log.Printf("Valid Crab Pattern detected! \n %+v", combination)
	//			return
	//		}
	//		log.Println("nothin ")
	//	}

	//fmt.Println("No Crab Pattern detected.")
	// Main loop
	for {
		updateInstruments()
		updateQuotes()

		checkOpenSignals()
		checkCloseSignals()

		// Sleep for a specified duration before the next check
		time.Sleep(10 * time.Second) // Change the duration as per your requirements
	}
}

func startSession() {
	// Create a new Kite connect instance
	kc = kiteconnect.New(apiKey)

	// Login URL from which request token can be obtained
	fmt.Println(kc.GetLoginURL())

	// Obtained request token after Kite Connect login flow
	requestToken := "SItKhM33r5TA5ZqZ54cZc9UskPPD75vU"

	// Get user details and access token
	data, err := kc.GenerateSession(requestToken, apiSecret)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Set access token
	kc.SetAccessToken(data.AccessToken)
}

//func startBinanceSession() {
//	// Create a new Binance client
//	client := binance.NewClient(apiKey, secretKey)
//
//	// Get server time
//	serverTime, err := client.NewServerTimeService().Do()
//	if err != nil {
//		log.Fatalf("Error: %v", err)
//	}
//
//	fmt.Println("Server Time:", serverTime)
//
//	// You may need to adjust the logic to get the request token or authenticate with Binance API if required
//
//	// Set the symbol for trading
//	client = client.NewWithSymbol(symbol)
//}
func updateInstruments() {
	// Get instruments
	getInstruments, err := kc.GetInstruments()
	if err != nil {
		log.Print(err)
		return
	}

	instruments = nil

	for _, instrument := range getInstruments {
		if instrument.Name == "NIFTY" && instrument.InstrumentType == "FUT" {
			instruments = append(instruments, instrument)
		}
	}
}

func updateQuotes() {
	quotes = nil

	for _, instrument := range instruments {
		quote, err := kc.GetQuote(fmt.Sprintf(instrument.Exchange + ":" + instrument.Tradingsymbol))
		if err != nil {
			log.Print(err)
			return
		}
		quotes = append(quotes, quote)
	}
}

//func updateQuoteBinance() {
//	// Get the latest ticker price
//	ticker, err := client.NewTickerPriceService().Symbol(symbol).Do()
//	if err != nil {
//		log.Print(err)
//		return
//	}
//
//	// Print the current price
//	fmt.Println("Current Price:", ticker.Price)
//}

func checkOpenSignals() {
	if quotes[0][""].BuyQuantity > quotes[0][""].SellQuantity &&
		quotes[1][""].BuyQuantity > quotes[1][""].SellQuantity &&
		quotes[2][""].BuyQuantity > quotes[2][""].SellQuantity {
		log.Print("LONG SIGNAL") //Signal to buy
		positionOpen = true
		long = true
	} else if quotes[0][""].BuyQuantity < quotes[0][""].SellQuantity &&
		quotes[1][""].BuyQuantity < quotes[1][""].SellQuantity &&
		quotes[2][""].BuyQuantity < quotes[2][""].SellQuantity {
		log.Print("SHORT SIGNAL")
		positionOpen = true
		long = false
	} else {
		log.Print("NO SIGNAL")
	}
}

func checkCloseSignals() {
	if positionOpen && long {
		if quotes[0][""].BuyQuantity < quotes[0][""].SellQuantity ||
			quotes[1][""].BuyQuantity < quotes[1][""].SellQuantity ||
			quotes[2][""].BuyQuantity < quotes[2][""].SellQuantity {
			log.Print("CLOSE SIGNAL FOR OPEN LONG POSITION")
			positionOpen = false
		}
	} else if positionOpen && !long {
		if quotes[0][""].BuyQuantity > quotes[0][""].SellQuantity ||
			quotes[1][""].BuyQuantity > quotes[1][""].SellQuantity ||
			quotes[2][""].BuyQuantity > quotes[2][""].SellQuantity {
			log.Print("CLOSE SIGNAL FOR OPEN SHORT POSITION")
			positionOpen = false
		}
	}
}

const Tolerance = 0.05 // Tolerance level for ratio comparisons

// Generate all combinations of A, B, C, and D values from the dataset
func generateCombinationsThreeSwings(stockData []StockData) [][]StockData {
	var combinations [][]StockData
	n := len(stockData)

	for a := 0; a < n-3; a++ {
		for b := a + 1; b < n-2; b++ {
			for c := b + 1; c < n-1; c++ {
				for d := c + 1; d < n; d++ {
					if stockData[a].Time.Before(stockData[b].Time) &&
						stockData[b].Time.Before(stockData[c].Time) &&
						stockData[c].Time.Before(stockData[d].Time) {
						combination := []StockData{stockData[a], stockData[b], stockData[c], stockData[d]}
						combinations = append(combinations, combination)
					}
				}
			}
		}
	}

	return combinations
}

func generateCombinationsFourSwings(stockData []StockData) [][]StockData {
	combinations := [][]StockData{}
	n := len(stockData)

	for a := 0; a < n-4; a++ {
		for b := a + 1; b < n-3; b++ {
			for c := b + 1; c < n-2; c++ {
				for d := c + 1; d < n-1; d++ {
					for e := d + 1; e < n; e++ {
						if stockData[a].Time.Before(stockData[b].Time) &&
							stockData[b].Time.Before(stockData[c].Time) &&
							stockData[c].Time.Before(stockData[d].Time) &&
							stockData[d].Time.Before(stockData[e].Time) {
							combination := []StockData{stockData[a], stockData[b], stockData[c], stockData[d], stockData[e]}
							combinations = append(combinations, combination)
						}
					}
				}
			}
		}
	}

	return combinations
}

// Check if the given combination forms a Butterfly pattern
func isButterflyPattern(combinations []StockData) bool {
	// Calculate the ratios for the Butterfly pattern

	ABRatio := combinations[1].Price / combinations[0].Price
	BCRatio := combinations[2].Price / combinations[1].Price
	CDRatio := combinations[3].Price / combinations[0].Price

	// Define the expected Fibonacci ratios for a Butterfly pattern
	expectedABRatio := 0.786 // Expected ratio for AB
	expectedBCRatio := 0.382 // Expected ratio for BC
	expectedCDRatio := 1.272 // Expected ratio for CD

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid Butterfly pattern
	return abMatch && bcMatch && cdMatch
}
func isCrabPattern(combination []StockData) bool {
	fmt.Println("crab check")
	// Calculate the ratios for the Crab pattern
	ABRatio := combination[1].Price / combination[0].Price
	BCRatio := combination[2].Price / combination[1].Price
	CDRatio := combination[3].Price / combination[2].Price
	EDRatio := combination[4].Price / combination[3].Price

	// Define the expected Fibonacci ratios for a Crab pattern
	expectedABRatio := 0.382 // Expected ratio for AB
	expectedBCRatio := 0.382 // Expected ratio for BC
	expectedCDRatio := 2.618 // Expected ratio for CD
	expectedEDRatio := 0.786 // Expected ratio for ED

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance
	edMatch := math.Abs(EDRatio-expectedEDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid Crab pattern
	return abMatch && bcMatch && cdMatch && edMatch
}
func isGartleyPattern(data []float64) bool {
	// Calculate the ratios for the Gartley pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]
	EDRatio := data[4] / data[3]

	// Define the expected Fibonacci ratios for a Gartley pattern
	expectedABRatio := 0.618 // Expected ratio for AB
	expectedBCRatio := 0.382 // Expected ratio for BC
	expectedCDRatio := 1.272 // Expected ratio for CD
	expectedEDRatio := 0.786 // Expected ratio for ED

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance
	edMatch := math.Abs(EDRatio-expectedEDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid Gartley pattern
	return abMatch && bcMatch && cdMatch && edMatch
}
func isBatPattern(data []float64) bool {
	// Calculate the ratios for the Bat pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]
	EDRatio := data[4] / data[3]

	// Define the expected Fibonacci ratios for a Bat pattern
	expectedABRatio := 0.382 // Expected ratio for AB
	expectedBCRatio := 0.382 // Expected ratio for BC
	expectedCDRatio := 0.886 // Expected ratio for CD
	expectedEDRatio := 1.618 // Expected ratio for ED

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance
	edMatch := math.Abs(EDRatio-expectedEDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid Bat pattern
	return abMatch && bcMatch && cdMatch && edMatch
}
func isCypherPattern(data []float64) bool {
	// Calculate the ratios for the Cypher pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]
	EDRatio := data[4] / data[3]

	// Define the expected Fibonacci ratios for a Cypher pattern
	expectedABRatio := 0.382 // Expected ratio for AB
	expectedBCRatio := 0.382 // Expected ratio for BC
	expectedCDRatio := 0.786 // Expected ratio for CD
	expectedEDRatio := 1.272 // Expected ratio for ED

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance
	edMatch := math.Abs(EDRatio-expectedEDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid Cypher pattern
	return abMatch && bcMatch && cdMatch && edMatch
}
func isSharkPattern(data []float64) bool {
	// Calculate the ratios for the Shark pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]
	EDRatio := data[4] / data[3]

	// Define the expected Fibonacci ratios for a Shark pattern
	expectedABRatio := 0.886 // Expected ratio for AB
	expectedBCRatio := 1.13  // Expected ratio for BC
	expectedCDRatio := 1.618 // Expected ratio for CD
	expectedEDRatio := 2.24  // Expected ratio for ED

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance
	edMatch := math.Abs(EDRatio-expectedEDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid Shark pattern
	return abMatch && bcMatch && cdMatch && edMatch
}
func isFiveZeroPattern(data []float64) bool {
	// Calculate the ratios for the 5-0 pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]
	EDRatio := data[4] / data[3]

	// Define the expected Fibonacci ratios for a 5-0 pattern
	expectedABRatio := 0.5   // Expected ratio for AB
	expectedBCRatio := 0.382 // Expected ratio for BC
	expectedCDRatio := 2.618 // Expected ratio for CD
	expectedEDRatio := 0.786 // Expected ratio for ED

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance
	edMatch := math.Abs(EDRatio-expectedEDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid 5-0 pattern
	return abMatch && bcMatch && cdMatch && edMatch
}
func isThreeDrivesPattern(data []float64) bool {
	// Calculate the ratios for the Three Drives pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]

	// Define the expected Fibonacci ratios for a Three Drives pattern
	expectedABRatio := 0.618 // Expected ratio for AB
	expectedBCRatio := 0.382 // Expected ratio for BC
	expectedCDRatio := 1.272 // Expected ratio for CD

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid Three Drives pattern
	return abMatch && bcMatch && cdMatch
}
func isABCDPattern(data []float64) bool {
	// Calculate the ratios for the AB=CD pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]

	// Define the expected Fibonacci ratios for an AB=CD pattern
	expectedABRatio := 0.618 // Expected ratio for AB
	expectedBCRatio := 0.382 // Expected ratio for BC
	expectedCDRatio := 1.0   // Expected ratio for CD

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid AB=CD pattern
	return abMatch && bcMatch && cdMatch
}
func isAlternateBatPattern(data []float64) bool {
	// Calculate the ratios for the Alternate Bat pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]
	XDRatio := data[3] / data[0]

	// Define the expected Fibonacci ratios for the Alternate Bat pattern
	expectedABRatio := 0.382 // Expected ratio for AB
	expectedBCRatio := 0.382 // Expected ratio for BC
	expectedCDRatio := 2.618 // Expected ratio for CD
	expectedXDRatio := 0.886 // Expected ratio for XD

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance
	xdMatch := math.Abs(XDRatio-expectedXDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid Alternate Bat pattern
	return abMatch && bcMatch && cdMatch && xdMatch
}
func isDeepCrabPattern(data []float64) bool {
	// Calculate the ratios for the Deep Crab pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]
	XDRatio := data[3] / data[0]

	// Define the expected Fibonacci ratios for the Deep Crab pattern
	expectedABRatio := 0.382 // Expected ratio for AB
	expectedBCRatio := 2.24  // Expected ratio for BC
	expectedCDRatio := 0.886 // Expected ratio for CD
	expectedXDRatio := 2.618 // Expected ratio for XD

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance
	xdMatch := math.Abs(XDRatio-expectedXDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid Deep Crab pattern
	return abMatch && bcMatch && cdMatch && xdMatch
}
func isWhiteSwanPattern(data []float64) bool {
	// Calculate the ratios for the White Swan pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]

	// Define the expected Fibonacci ratios for the White Swan pattern
	expectedABRatio := 0.382 // Expected ratio for AB
	expectedBCRatio := 1.13  // Expected ratio for BC
	expectedCDRatio := 1.618 // Expected ratio for CD

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid White Swan pattern
	return abMatch && bcMatch && cdMatch
}
func isBlackSwanPattern(data []float64) bool {
	// Calculate the ratios for the Black Swan pattern
	ABRatio := data[1] / data[0]
	BCRatio := data[2] / data[1]
	CDRatio := data[3] / data[2]

	// Define the expected Fibonacci ratios for the Black Swan pattern
	expectedABRatio := 0.886 // Expected ratio for AB
	expectedBCRatio := 2.618 // Expected ratio for BC
	expectedCDRatio := 0.382 // Expected ratio for CD

	// Check if the calculated ratios are within the tolerance level of the expected ratios
	abMatch := math.Abs(ABRatio-expectedABRatio) <= Tolerance
	bcMatch := math.Abs(BCRatio-expectedBCRatio) <= Tolerance
	cdMatch := math.Abs(CDRatio-expectedCDRatio) <= Tolerance

	// Return true if all ratios match the expected ratios, indicating a valid Black Swan pattern
	return abMatch && bcMatch && cdMatch
}
