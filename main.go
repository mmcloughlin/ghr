package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-github/github"
)

// FindEmail looks for the given user's email address by looking for commit
// events in their public feed.
func FindEmail(c *github.Client, user string) (string, error) {
	events, _, err := c.Activity.ListEventsPerformedByUser(user, true, nil)
	if err != nil {
		return "", err
	}

	for _, ev := range events {
		if *ev.Type != "PushEvent" {
			continue
		}
		push := ev.Payload().(*github.PushEvent)
		for _, commit := range push.Commits {
			author := commit.Author
			if author.Email != nil {
				return *author.Email, nil
			}
		}
	}

	return "", nil
}

func main() {
	u, _ := url.Parse("http://localhost:11112")
	trans := http.Transport{
		Proxy: http.ProxyURL(u),
	}
	client := &http.Client{Transport: &trans}

	c := github.NewClient(client)
	opt := &github.SearchOptions{
		Sort:  "stars",
		Order: "asc",
	}

	repos, _, err := c.Search.Repositories("stars:>10 language:go", opt)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range repos.Repositories {
		user := *r.Owner.Login
		email, err := FindEmail(c, user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(*r.Name, *r.StargazersCount, user, email)
	}

}
