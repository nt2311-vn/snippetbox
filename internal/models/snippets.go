package models

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

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
	query := `INSERT INTO snippets (title, content, created, expires)
	          VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := sm.DB.Exec(query, title, content, expires)
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
	query := `select id, title, content, created, expires from snippets
	where expires > utc_timestamp() and id = ?
	`
	s := &Snippet{}
	if err := sm.DB.QueryRow(query, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (sm *SnippetModel) Latest() ([]*Snippet, error) {
	query := `select id, title, content, created, expires from snippets
	where expires > utc_timestamp() order by id desc limit 10
	`

	rows, err := sm.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, err
}
