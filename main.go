package main

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

	// Common network options
	var proxy, token, userAgent string
	searchFlags := pflag.NewFlagSet("connection", pflag.ExitOnError)
	searchFlags.StringVarP(&proxy, "proxy", "", "",
		"HTTP proxy")
	searchFlags.StringVarP(&token, "token", "", "",
		"API token")
	searchFlags.StringVarP(&userAgent, "useragent", "", "github.com/mmcloughlin/ghr",
		"User agent. Please use your github username.")

	// Start a new search
	var searchQuery string
	searchCmd := &cobra.Command{
		Use:   "search --query=QUERY --proxy=PROXY --token=TOKEN --useragent=AGENT",
		Short: "Start a new query",
		Run: func(cmd *cobra.Command, args []string) {
			client := BuildHTTPClient(proxy, token)
			c := github.NewClient(client)
			c.UserAgent = userAgent

			s, err := store.NewSearch(searchQuery)
			if err != nil {
				log.Fatal(err)
			}

			Query(c, store, s)
		},
	}

	searchCmd.PersistentFlags().AddFlagSet(searchFlags)
	searchCmd.PersistentFlags().StringVarP(&searchQuery, "query", "", "",
		"Search keywords and qualifiers (https://developer.github.com/v3/search/#search-repositories)")

	mainCmd.AddCommand(searchCmd)

	// Resume a search
	var searchID uint
	resumeCmd := &cobra.Command{
		Use:   "resume --id=ID --proxy=PROXY --token=TOKEN --useragent=AGENT",
		Short: "Resume another query",
		Run: func(cmd *cobra.Command, args []string) {
			client := BuildHTTPClient(proxy, token)
			c := github.NewClient(client)
			c.UserAgent = userAgent

			s := &Search{}
			q := store.DB.First(s, searchID)
			if q.Error != nil {
				log.Fatal(q.Error)
			}

			Query(c, store, s)
		},
	}

	resumeCmd.PersistentFlags().AddFlagSet(searchFlags)
	resumeCmd.PersistentFlags().UintVarP(&searchID, "id", "", 0,
		"Search ID")

	mainCmd.AddCommand(resumeCmd)

	// Execute
	_ = mainCmd.Execute()
}

func BuildHTTPClient(proxy, token string) *http.Client {
	var base http.RoundTripper
	base = http.DefaultTransport

	// Proxy layer
	if len(proxy) > 0 {
		u, _ := url.Parse(proxy)
		base = &http.Transport{
			Proxy: http.ProxyURL(u),
		}
	}

	// Authentication layer
	if len(token) > 0 {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		base = &oauth2.Transport{
			Source: ts,
			Base:   base,
		}
	}

	// Rate limiting
	transport := &RateLimitedTransport{
		Base: base,
	}

	return &http.Client{
		Transport: transport,
	}
}

func Query(c *github.Client, store *Store, s *Search) {
	// scraper
	scraper := &Scraper{
		Client: c,
		Store:  store,
	}

	err := scraper.Scrape(s)
	if err != nil {
		log.Fatal(err)
	}
}
