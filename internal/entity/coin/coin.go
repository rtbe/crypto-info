package coin

import (
	"cryptoInfo/internal/entity/history"
	"cryptoInfo/internal/entity/price"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

//Info contains information about particular coin
type Info struct {
	CoinName string
	Price    price.Info
	History  history.Info
}

//SaveHistory saves 30 day price history in .csv file format into ../output directory
func (c *Info) SaveHistory() {
	if _, err := os.Stat("../output"); os.IsNotExist(err) {
		err := os.Mkdir("../output", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	file, err := os.Create(fmt.Sprintf("../output/%sHistory", c.CoinName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = csv.NewWriter(file).WriteAll(c.History.ParsedPrices)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Info) loadHistory() {
	file, err := os.Open(fmt.Sprintf("%sHistory.csv", c.CoinName))
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	for {
		v, err := r.Read()
		if err == io.EOF {
			return
		}
		c.History.ParsedPrices = append(c.History.ParsedPrices, []string{v[0], v[1]})
	}
}

//Print coloured information about requested coin into terminal
func (c *Info) Print() {
	if c.Price.Btc == 0.0 {
		fmt.Printf("| %s | %.1f$ | 24h_change_usd: %s |\n", c.CoinName, c.Price.Usd, colorizer(c.Price.Usd24hChange))
		return
	}
	fmt.Printf("| %s | %.1f$ | %f₿ | 24h_change_usd: %s | 24h_change_btc: %s |\n", c.CoinName, c.Price.Usd, c.Price.Btc, colorizer(c.Price.Usd24hChange), colorizer(c.Price.Btc24hChange))
}

func colorizer(i float64) string {
	var colorString string
	str := fmt.Sprintf("%.2f", i)
	if strings.HasPrefix(str, "-") {
		colorString = fmt.Sprint("\033[31m", str, "%", "\033[0m")
		return colorString
	}
	colorString = fmt.Sprint("\033[32m", str, "%", "\033[0m")
	return colorString
}
