package models

import (
	"fancykaraoke/db"
)

var conn = db.GetConnection()

type Song struct {
	Id        int    `json:"id"`
	Artist    string `json:"artist"`
	Title     string `json:"title"`
	YoutubeId string `json:"youtube_id"`
	CreatedAt string `json:"created_at"`
}

func GetSongById(id int) (Song, error) {
	q := `SELECT id, artist, title, youtube_id, created_at FROM songs WHERE id = ?`
	var song Song
	if err := conn.DB.QueryRow(q, id).Scan(&song.Id, &song.Artist, &song.Title, &song.YoutubeId, &song.CreatedAt); err != nil {
		return song, err
	}
	return song, nil
}

func SearchSongs(searchTerm string, limit int) ([]Song, error) {
	search := searchTerm + "*"
	q := `SELECT s.id, s.title, s.artist, s.youtube_id, s.created_at
				FROM songs_fts fts
				JOIN songs s on s.id = fts.rowid
				WHERE songs_fts MATCH ?
				ORDER BY rank DESC
				LIMIT ?`
	var songs []Song
	rows, err := conn.DB.Query(q, search, limit)
	if err != nil {
		return songs, err
	}
	for rows.Next() {
		var song Song
		err := rows.Scan(&song.Id, &song.Title, &song.Artist, &song.YoutubeId, &song.CreatedAt)
		if err != nil {
			return songs, err
		}
		songs = append(songs, song)

	}
	return songs, nil
}
