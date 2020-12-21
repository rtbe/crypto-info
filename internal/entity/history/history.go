package history

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Info struct {
	Prices       [][]float64 `json:"prices"`
	ParsedPrices [][]string
}

func (h *Info) Request(query []string, days int, cl *http.Client, priceCh chan<- map[string]interface{}) {

	coinName, vsCurrency := query[1], query[2]
	requestQueryString := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=%s&days=%s", coinName, vsCurrency, strconv.Itoa(days))

	req, err := http.NewRequest("GET", requestQueryString, nil)
	if err != nil {
		log.Fatal("Error creating new request")
	}
	req.Header.Add("accept", "application/json")

	resp, err := cl.Do(req)
	if err != nil {
		log.Fatal("Error sending request to server:", err.Error())
	}
	err = json.NewDecoder(resp.Body).Decode(h)

	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Error reading from request body:", err.Error())
	}
	var res map[string]interface{}
	priceCh <- res
	close(priceCh)
}

func (h *Info) Parse(wg *sync.WaitGroup) {
	//Time zone output format
	loc := time.FixedZone("UTC-0", 0)
	for _, v := range h.Prices {
		t := time.Unix(int64(v[0]/1000), 0)
		time := strings.TrimSuffix(t.In(loc).Format(time.Stamp), " UTC-0")
		price := fmt.Sprintf("%.7f", v[1])
		h.ParsedPrices = append(h.ParsedPrices, []string{time, price})
	}
	wg.Done()
}
