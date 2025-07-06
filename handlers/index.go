package handlers

import (
	"fancykaraoke/templates/pages"
	"net/http"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	component := pages.IndexPage("World")
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
