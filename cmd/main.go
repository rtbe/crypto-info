package main

import (
	"cryptoInfo/internal/coin"
	"cryptoInfo/internal/history"
	"cryptoInfo/internal/job"
	"cryptoInfo/internal/price"
	"flag"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	timeoutFlag := flag.Int("t", 5, "Timeout time for requests in seconds, 5 second by default")
	coinNameFlag := flag.String("c", "bitcoin", "Coin name, bitcoin by default")
	priceHistoryFlag := flag.Bool("h", false, "Price history: true or false, by default false")
	flag.Parse()
	cl := http.Client{
		Timeout: time.Duration(time.Duration(*timeoutFlag) * time.Second),
	}
	coinNames := strings.Split(*coinNameFlag, ",")

	for i, v := range coinNames {
		coinNames[i] = strings.Trim(v, " ")
	}

	for _, coinName := range coinNames {
		coin := coin.Info{
			CoinName: coinName,
		}

		jobs := job.NewQueue()
		//if coinName=="bitcoin" you should check price only vs USD
		if *priceHistoryFlag && coinName != "bitcoin" {
			jobs.Push(job.New("getHistory", coinName, "btc"))
		}
		if *priceHistoryFlag && coinName == "bitcoin" {
			jobs.Push(job.New("getHistory", coinName, "usd"))
		}
		if coinName != "bitcoin" {
			jobs.Push(job.New("getPrice", coinName, "btc"))
			jobs.Push(job.New("getPrice", coinName, "usd"))
		}
		if coinName == "bitcoin" {
			jobs.Push(job.New("getPrice", coinName, "usd"))
		}

		wg := sync.WaitGroup{}
		for _, job := range *jobs {
			priceCh := make(chan map[string]interface{})

			switch job[0] {
			case "getPrice":
				go func(p *price.Info) {
					p.Request(job, &cl, priceCh)
				}(&coin.Price)
				//block untill price will be recieved
				priceMap := <-priceCh
				wg.Add(1)

				go func(p *price.Info) {
					p.Parse(coin.CoinName, priceMap, &wg)
				}(&coin.Price)
			case "getHistory":
				go func(h *history.Info) {
					h.Request(job, 30, &cl, priceCh)
				}(&coin.History)
				//block untill history will be recieved
				<-priceCh

				wg.Add(1)
				go func(h *history.Info) {
					h.Parse(&wg)
				}(&coin.History)
			}
		}
		//wait for parsing values
		wg.Wait()

		coin.Print()
		if *priceHistoryFlag == true {
			coin.SaveHistory()
		}
	}
}
