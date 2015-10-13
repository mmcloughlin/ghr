package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-github/github"
)

func main() {
	// Create DB
	store, err := NewStore("./store.db")
	if err != nil {
		log.Fatal(err)
	}
	store.Init()

	// proxy layer
	u, _ := url.Parse("http://localhost:11112")
	proxyTransport := &http.Transport{
		Proxy: http.ProxyURL(u),
	}

	// http client
	client := &http.Client{
		Transport: &RateLimitedTransport{
			Base: proxyTransport,
		},
	}

	// github client
	c := github.NewClient(client)

	// scraper
	scraper := &Scraper{
		Client: c,
		Store:  store,
	}

	// Start scrape
	s, err := store.NewSearch("stars:>=1000 language:Go")
	if err != nil {
		log.Fatal(err)
	}

	err = scraper.Scrape(s)
	if err != nil {
		log.Fatal(err)
	}
}
