package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	queryFmt := `insert into snippets (title, content, created, expires)
	values (?, ?, utc_timestamp(), date_add(utc_timestamp(), interval ? day))`

	result, err := sm.DB.Exec(queryFmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

func (sm *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
