package models

import (
	"database/sql"
	"fmt"
)

const selectionLimit = 50

type Track struct {
	Id     int
	Name   string
	Artist string
	Album  string
}

func GetSelectionForHomePage(db *sql.DB) ([]Track, error) {
	tracks := make([]Track, 0)
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM tracks LIMIT %d`, selectionLimit))
	if err != nil {
		return nil, err
	}

	var track Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Name, &track.Artist, &track.Album); err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}
