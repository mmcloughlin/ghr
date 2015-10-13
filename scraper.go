package main

import (
	"fmt"

	"github.com/google/go-github/github"
)

type Scraper struct {
	Client *github.Client
	Store  *Store
}

func (s *Scraper) Scrape(search *Search) error {
	for !search.Finished {
		opt := &github.SearchOptions{
			Sort:  "stars",
			Order: "asc",
			ListOptions: github.ListOptions{
				Page: search.CompletedPages + 1,
			},
		}

		repos, res, err := s.Client.Search.Repositories(search.Query, opt)
		if err != nil {
			return err
		}

		for _, r := range repos.Repositories {
			fmt.Println(*r.Name)
		}

		// Update
		search.CompletedPages++
		search.Finished = (res.NextPage == 0)
		q := s.Store.DB.Save(search)
		if q.Error != nil {
			return q.Error
		}
	}

	return nil
}

// FindEmail looks for the given user's email address by looking for commit
// events in their public feed.
func FindEmail(c *github.Client, user string) (string, error) {
	u, _, err := c.Users.Get(user)
	if err != nil {
		return "", err
	}
	if u.Email != nil {
		fmt.Println("email from user fetch:", *u.Email)
	}

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
