package utils

import (
	"database/sql"
	"encoding/json"
	"fancykaraoke/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type SongImport struct {
	Artist    string `json:"artist"`
	Title     string `json:"title"`
	YoutubeId string `json:"youtube_id"`
}

type ImportLog struct {
	ID                 int       `json:"id"`
	StartedAt          time.Time `json:"started_at"`
	EndedAt            time.Time `json:"ended_at"`
	ImportCount        int       `json:"import_count"`
	TotalRemoteRecords int       `json:"total_remote_records"`
	Success            bool      `json:"success"`
	ErrorMessage       error     `json:"error_message"`
}

type KaraokeResponse struct {
	RecordsTotal int        `json:"recordsTotal"`
	Data         [][]string `json:"data"`
}

var (
	importLog = ImportLog{
		StartedAt:          time.Now(),
		ImportCount:        0,
		TotalRemoteRecords: 0,
		Success:            false,
	}
	conn      = db.GetConnection()
	re        = regexp.MustCompile(`<a[^>]*>([^<]+)</a>`)
	startedAt = time.Now()
	youtubeRe = regexp.MustCompile(`v=([^&|]+)`)
)

func fetchData(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if (err) != nil {
		return nil, err
	}
	return resp, nil
}

func selectStartCount() (int, error) {
	q := `SELECT total_remote_records FROM import_log ORDER BY id DESC LIMIT 1`
	var count int
	if err := conn.DB.QueryRow(q).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func parseResp(resp *http.Response) ([]SongImport, int, error) {
	defer resp.Body.Close()
	var songs []SongImport
	var karaokeResp KaraokeResponse

	if err := json.NewDecoder(resp.Body).Decode(&karaokeResp); err != nil {
		return nil, 0, err
	}

	total := karaokeResp.RecordsTotal
	for _, htmlStrings := range karaokeResp.Data {
		artistMatch := re.FindStringSubmatch(htmlStrings[1])
		if len(artistMatch) < 2 {
			continue
		}
		artist := artistMatch[1]

		titleMatch := re.FindStringSubmatch(htmlStrings[2])
		if len(titleMatch) < 2 {
			continue
		}
		title := titleMatch[1]

		youtubeMatch := youtubeRe.FindStringSubmatch(htmlStrings[5])
		if len(youtubeMatch) < 2 {
			continue
		}
		youtubeID := youtubeMatch[1]

		songs = append(songs, SongImport{
			Artist:    artist,
			Title:     title,
			YoutubeId: youtubeID,
		})
	}

	return songs, total, nil
}

func saveSongs(songs []SongImport) (int, error) {
	successCount := 0
	for _, song := range songs {
		sql := "INSERT INTO songs (artist, title, youtube_id) VALUES (?, ?, ?)"
		_, err := conn.DB.Exec(sql, song.Artist, song.Title, song.YoutubeId)
		if err != nil {
			log.Println("Error importing", song.Artist, song.Title, song.YoutubeId, err)
			continue // Skip to next song
		}
		successCount++
	}
	return successCount, nil
}

func saveImportLog() error {
	sql := "INSERT INTO import_log (started_at, ended_at, import_count, total_remote_records, success) VALUES (?, ?, ?, ?, ?)"
	_, err := conn.DB.Exec(sql, importLog.StartedAt, importLog.EndedAt, importLog.ImportCount, importLog.TotalRemoteRecords, importLog.Success)
	return err
}

func saveSuccessLog(importCount int, totalRecords int) error {
	importLog.Success = true
	importLog.ImportCount = importCount
	importLog.TotalRemoteRecords = totalRecords
	importLog.EndedAt = time.Now()
	err := saveImportLog()
	return err
}

func saveErrorLog(err error) {
	importLog.Success = false
	importLog.EndedAt = time.Now()
	importLog.ErrorMessage = err
	_ = saveImportLog()
	log.Fatal(err)
}

func ImportSongs() {
	length := 50
	if len(os.Args) > 2 {
		num, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Invalid number", err)
		}
		length = num
	}

	importLog.StartedAt = time.Now()
	startCount, err := selectStartCount()
	if err != nil {
		saveErrorLog(err)
	}
	url := fmt.Sprintf(`https://www.karaokenerds.com/Community/BrowseJson/?length=%d&start=%d&order[0][column]=3&order[0][dir]=desc`, length, startCount)
	resp, err := fetchData(url)
	if err != nil {
		saveErrorLog(err)
	}

	songs, remoteCount, err := parseResp(resp)
	if err != nil {
		saveErrorLog(err)
	}

	successCount, err := saveSongs(songs)
	if err != nil {
		saveErrorLog(err)
		log.Fatal(err)
	}

	if err = saveSuccessLog(remoteCount, successCount); err != nil {
		log.Fatal(err)
	}

	log.Println(fmt.Sprintf("Import complete. Imported %d songs. %d total remote records.", len(songs), remoteCount))
}
