package models

import (
	"fancykaraoke/db"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var conn = db.GetConnection()

type Party struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

func NewParty(name string) (Party, error) {
	// Strip everything except letters, numbers, spaces, and hyphens
	nameRe := regexp.MustCompile(`[^a-zA-Z 0-9-]+`)
	safeName := nameRe.ReplaceAllString(name, "")
	// For URL-friendly slug: replace spaces with hyphens, lowercase
	slug := strings.ToLower(strings.ReplaceAll(safeName, " ", "-"))
	// Remove any multiple consecutive hyphens
	slugRe := regexp.MustCompile(`-+`)
	slug = slugRe.ReplaceAllString(slug, "-")
	// Trim hyphens from start/end
	slug = strings.Trim(slug, "-")
	// Add number to already used slug
	slugCount, err := GetSlugCount(slug)
	if err != nil {
		return Party{}, err
	}
	if slugCount > 0 {
		slug = slug + "-" + strconv.Itoa(slugCount)
	}

	return Party{
		Id:        uuid.New().String(),
		Name:      safeName,
		Slug:      slug,
		CreatedAt: time.Now(),
	}
}

func GetPartyBySlug(slug string) (Party, error) {
	q := `SELECT id, name, slug, created_at FROM parties WHERE slug = ?`
	var party Party
	if err := conn.DB.QueryRow(q, slug).Scan(&party.Id, &party.Name, &party.Slug, &party.CreatedAt); err != nil {
		return party, err
	}
	return party, nil
}

func CreateParty(name string) (Party, error) {
	party, err := NewParty(name)
	if err != nil {
		return Party{}, err
	}

	sql := `INSERT INTO parties (id, name, slug) VALUES (?, ?, ?)`
	result, err := conn.DB.Exec(sql, party.Id, party.Name, party.Slug)
	if err != nil {
		return Party{}, err
	}
	return party, nil
}

func GetSlugCount(slug string) (int, error) {
	param := slug + "%"
	sql := `SELECT COUNT(*) FROM parties WHERE slug LIKE ?`
	count := 0
	err := conn.DB.QueryRow(sql, param).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
