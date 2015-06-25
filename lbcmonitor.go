package main

import (
	"encoding/json"
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
)

const (
	sellPage = "https://localbitcoins.com/sell-bitcoins-online/US/united-states/cash-deposit/.json"
	buyPage  = "https://localbitcoins.com/buy-bitcoins-online/US/united-states/cash-deposit/.json"
)

type Options struct {
	Buying       bool    `long:"buy" description:"Switch to buying mode instead of selling."`
	XBTAmount    float64 `long:"xbt" description:"Amount of XBT that requires conversion to fiat."`
	FiatAmount   float64 `long:"fiat" description:"Amount of fiat needed from the conversion of given XBT."`
	XBTPrice     float64 `long:"xbtprice" description:"Price of 1 XBT that the software should alert at."`
	Pages        uint    `long:"pages" description:"Number of pages to follow. Default 5."`
	priceMonitor bool
}

var opts = new(Options)
var w = new(tabwriter.Writer)

func main() {
	_, err := flags.Parse(opts)
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			return
		}
		os.Exit(1)
	}

	if opts.XBTAmount == 0 && opts.FiatAmount == 0 { // invalid/price monitoring
		if opts.XBTPrice == 0 { // invalid
			fmt.Fprintln(os.Stderr, "Invalid arguments")
			os.Exit(1)
		} else {
			opts.priceMonitor = true
		}
	} else if opts.XBTAmount == 0 || opts.FiatAmount == 0 { // insufficient args
		fmt.Fprintln(os.Stderr, "Invalid arguments")
		os.Exit(1)
	}
	if opts.Pages == 0 {
		opts.Pages = 5
	}

	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Rate\tMin\tMax\tBanks\tLink")
	fmt.Fprintln(w, "----\t---\t---\t-----\t----")
	if opts.Buying {
		process(buyPage)
	} else {
		process(sellPage)
	}
	w.Flush()
}

func process(jsonURL string) {
	res, err := http.Get(jsonURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "got HTTP error %v\n", err)
		os.Exit(1)
	}
	if res.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "got status code %d\n", res.StatusCode)
		os.Exit(1)
	}
	m := new(MainJson)
	err = json.NewDecoder(res.Body).Decode(m)
	res.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "json decode failed: %v\n", err)
		os.Exit(1)

	}

	ads := m.Data.AdList
	for i, ad := range ads {
		var maxAmount, minAmount float64
		if ad.Data.MaxAmount != "" {
			maxAmount, err = strconv.ParseFloat(ad.Data.MaxAmount, 64)
			if err != nil {
				log.Printf("maxAmount decode for ad #%d failed: %v\n", i, err)
			}
		}
		if ad.Data.MinAmount != "" {
			minAmount, err = strconv.ParseFloat(ad.Data.MinAmount, 64)
			if err != nil {
				log.Printf("minAmount decode for ad #%d failed: %v\n", i, err)
			}
		}
		price, err := strconv.ParseFloat(ad.Data.Price, 64)
		if err != nil {
			log.Printf("price decode for ad #%d failed: %v\n", i, err)
		}

		bankName := ad.Data.BankName
		url := ad.Actions.PublicView

		if opts.priceMonitor {
			var cond bool
			if !opts.Buying { // want highest selling price
				cond = price >= opts.XBTPrice
			} else { // want lowest selling price
				cond = price <= opts.XBTPrice
			}
			if cond {
				fmt.Fprintf(w, "%.2f\t%.0f\t%.0f\t%s\t%s\n", price, minAmount,
					maxAmount, bankName, url)
			}
		} else {
			if maxAmount < opts.FiatAmount {
				continue
			}
			if minAmount > opts.FiatAmount {
				continue
			}
			calcFiat := opts.XBTAmount * price

			var cond bool
			if !opts.Buying { // want max fiat
				cond = calcFiat >= opts.FiatAmount
			} else { // want lowest fiat
				cond = calcFiat <= opts.FiatAmount
			}

			if cond {
				fmt.Fprintf(w, "%.2f\t%.0f\t%.0f\t%s\t%s\n", price, minAmount,
					maxAmount, bankName, url)
			}
		}
	}
	opts.Pages -= 1
	if opts.Pages > 0 && m.Pagination.Next != "" { // we can still go over more
		process(m.Pagination.Next)
	}
}

type MainJson struct {
	Pagination PaginationJson `json:"pagination"`
	Data       DataJson       `json:"data"`
}

type PaginationJson struct {
	Next string `json:"next"`
}

type DataJson struct {
	AdList []AdJson `json:"ad_list"`
}

type AdJson struct {
	Data    AdDataJson    `json:"data"`
	Actions AdActionsJson `json:"actions"`
}

type AdActionsJson struct {
	PublicView string `json:"public_view"`
}

type AdDataJson struct {
	Price     string `json:"temp_price"`
	BankName  string `json:"bank_name"`
	MinAmount string `json:"min_amount"`
	MaxAmount string `json:"max_amount"`
	Message   string `json:"msg"`
}

type result struct {
	price float64
	url   string
}
