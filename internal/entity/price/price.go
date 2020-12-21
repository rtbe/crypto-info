package price

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

//Info contains price in different formats
type Info struct {
	sync.Mutex
	Btc          float64
	Usd          float64
	Usd24hChange float64
	Btc24hChange float64
}

//Request makes a request to CoinGecko API to get information about current price and 24 hours price changw
func (p *Info) Request(query []string, cl *http.Client, priceCh chan<- map[string]interface{}) {
	coinName, vsCurrency := query[1], query[2]
	queryString := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s&include_24hr_change=true", coinName, vsCurrency)

	req, err := http.NewRequest("GET", queryString, nil)
	req.Header.Add("accept", "application/json")
	resp, err := cl.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	res := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&res)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Error reading from body")
	}

	priceCh <- res
	close(priceCh)
}

//Parse parses requested information into convinient format in concurrent-safe manner
func (p *Info) Parse(coinName string, coinMap map[string]interface{}, wg *sync.WaitGroup) {
	if coinMap[coinName] == nil {
		log.Fatal("Invalid coin name")
	}
	p.Lock()
	defer p.Unlock()
	for k, v := range coinMap[coinName].(map[string]interface{}) {
		switch k {
		case "btc":
			p.Btc = v.(float64)
		case "usd":
			p.Usd = v.(float64)
		case "btc_24h_change":
			p.Btc24hChange = v.(float64)
		case "usd_24h_change":
			p.Usd24hChange = v.(float64)
		}
	}
	wg.Done()
}
