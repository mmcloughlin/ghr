package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

func main() {
	// Main command
	var storePath string
	var store *Store
	mainCmd := &cobra.Command{
		Use:   "ghr",
		Short: "Scrape GitHub for tech talent",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			s, err := NewStore("./store.db")
			if err != nil {
				log.Fatal(err)
			}
			store = s
		},
	}

	mainCmd.PersistentFlags().StringVarP(&storePath, "store", "", "",
		"Path to datastore (SQLite3 database).")

	// Initialize datastore
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize datastore",
		Run: func(cmd *cobra.Command, args []string) {
			store.Init()
		},
	}

	mainCmd.AddCommand(initCmd)

	// Start a query
	var searchQuery string
	searchCmd := &cobra.Command{
		Use:   "search --query=QUERY",
		Short: "Start a new query",
		Run: func(cmd *cobra.Command, args []string) {
			Query(store, searchQuery)
		},
	}

	searchCmd.PersistentFlags().StringVarP(&searchQuery, "query", "", "",
		"Search keywords and qualifiers (https://developer.github.com/v3/search/#search-repositories)")

	mainCmd.AddCommand(searchCmd)

	// Execute
	_ = mainCmd.Execute()
}

func Query(store *Store, q string) {
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
	c.UserAgent = "github.com/mmcloughlin/ghr"

	// scraper
	scraper := &Scraper{
		Client: c,
		Store:  store,
	}

	// Start scrape
	s, err := store.NewSearch(q)
	if err != nil {
		log.Fatal(err)
	}

	err = scraper.Scrape(s)
	if err != nil {
		log.Fatal(err)
	}
}
