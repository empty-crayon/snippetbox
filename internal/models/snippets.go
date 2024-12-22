package models

import (
	"database/sql"
	"errors"
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

// This will insert a new snippet in the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	// this placeholder way of doing it isntead of string formatting to prevent sql injections
	stmt := `INSERT INTO snippets (title, content, created, expires)
			VALUES (?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`

	result, err := m.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// This will return a specific snippet for the id
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `
		SELECT id, title, content, created FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?
	`

	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return s, nil
}

// This will return 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	// sql statement we want to execute
	stmt := `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10;
	`

	// using query statement on the connection pool to execute the  query
	rows, err := m.DB.Query(stmt)


	if err != nil {
		return nil, err
	}

	// this to make sure sql.rows always closes properly before latest returns
	// critical: as long as the result set is open, underlyring db connection remains open
	// so if something goes wrong, the resultset isnt closed. 
	defer rows.Close()

	var snippets []Snippet

	// rows.next to iterate over rows this prepares each row for .scan

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// successful iteration over rows doesnt imply there was no error in iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
