package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Search struct {
	ID             uint
	Query          string
	CompletedPages int
	Finished       bool
}

type Prospect struct {
	ID       uint
	User     string
	Repo     string
	Name     string
	Email    string
	Location string
}

type Store struct {
	Path string
	DB   gorm.DB
}

func NewStore(path string) (*Store, error) {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &Store{
		Path: path,
		DB:   db,
	}, nil
}

func (s *Store) Init() error {
	s.DB.CreateTable(&Search{})
	s.DB.CreateTable(&Prospect{})
	return nil
}

func (s *Store) NewSearch(query string) (*Search, error) {
	search := &Search{Query: query}
	q := s.DB.Create(search)
	if q.Error != nil {
		return nil, q.Error
	}
	return search, nil
}
