package main

import (
	"fmt"

	"github.com/google/go-github/github"
)

type Scraper struct {
	Client *github.Client
	Store  *Store
}

func (s *Scraper) ProspectFromRepository(repo github.Repository) (*Prospect, error) {
	user := *repo.Owner.Login

	// Initialize Prospect object
	p := &Prospect{
		User: user,
		Repo: *repo.FullName,
	}

	// Try to fetch the owner's information directly
	u, _, err := s.Client.Users.Get(user)
	if err != nil {
		return nil, err
	}
	if u.Email != nil {
		p.Email = *u.Email
		if u.Name != nil {
			p.Name = *u.Name
		}
		if u.Location != nil {
			p.Location = *u.Location
		}
		return p, nil
	}

	// Otherwise try to fetch it from the activity list
	events, _, err := s.Client.Activity.ListEventsPerformedByUser(user, true, nil)
	if err != nil {
		return nil, err
	}

	for _, ev := range events {
		if *ev.Type != "PushEvent" {
			continue
		}
		push := ev.Payload().(*github.PushEvent)
		for _, commit := range push.Commits {
			author := commit.Author
			if author.Email != nil {
				p.Email = *author.Email
				if author.Name != nil {
					p.Name = *author.Name
				}
				return p, nil
			}
		}
	}

	return nil, nil
}

func (s *Scraper) Scrape(search *Search) error {
	for !search.Finished {
		// Fetch page of results
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

		// Iterate through them and find information about the owner, if we can.
		for _, r := range repos.Repositories {
			fmt.Println(*r.Name)
			p, err := s.ProspectFromRepository(r)
			if err != nil {
				return err
			}
			if p != nil {
				s.Store.DB.Create(p)
			}
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
