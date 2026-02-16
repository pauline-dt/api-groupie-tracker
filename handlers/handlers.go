package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"api-groupie-tracker/api"
	"api-groupie-tracker/models"
	"api-groupie-tracker/utils"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

/*
	PageData
	➡️ STRUCT UNIQUE pour index.html
	➡️ TOUS les champs utilisés dans le template existent TOUJOURS
*/
type PageData struct {
	Artists      []models.FullArtist
	MinCreation  int
	MaxCreation  int
	MinAlbum     int
	MaxAlbum     int
	MinMembers   int
	MaxMembers   int
	AllLocations []string

	Query    string
	Filtered bool
}

// =======================
// HOME
// =======================
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	artists := api.GetAllArtists()
	fullArtists := make([]models.FullArtist, 0, len(artists))

	for _, artist := range artists {
		fullArtists = append(fullArtists, models.FullArtist{
			Artist:         artist,
			FirstAlbumYear: utils.ExtractYear(artist.FirstAlbum),
		})
	}

	minCreation, maxCreation, minAlbum, maxAlbum := utils.GetYearRange(artists)
	minMembers, maxMembers := utils.GetMembersRange(artists)

	data := PageData{
		Artists:      fullArtists,
		MinCreation:  minCreation,
		MaxCreation:  maxCreation,
		MinAlbum:     minAlbum,
		MaxAlbum:     maxAlbum,
		MinMembers:   minMembers,
		MaxMembers:   maxMembers,
		AllLocations: utils.GetUniqueLocations(artists, api.Relations.Index),

		Query:    "",
		Filtered: false,
	}

	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// =======================
// ARTIST DETAILS
// =======================
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	fullArtist, err := api.GetFullArtistByID(id)
	if err != nil {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	fullArtist.FirstAlbumYear = utils.ExtractYear(fullArtist.FirstAlbum)

	if err := templates.ExecuteTemplate(w, "artist.html", fullArtist); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// =======================
// FILTER
// =======================
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	criteria := models.FilterCriteria{}
	r.ParseForm()

	if v := r.FormValue("creation_min"); v != "" {
		criteria.CreationDateMin, _ = strconv.Atoi(v)
	}
	if v := r.FormValue("creation_max"); v != "" {
		criteria.CreationDateMax, _ = strconv.Atoi(v)
	}

	if v := r.FormValue("album_min"); v != "" {
		criteria.FirstAlbumMin, _ = strconv.Atoi(v)
	}
	if v := r.FormValue("album_max"); v != "" {
		criteria.FirstAlbumMax, _ = strconv.Atoi(v)
	}

	if v := r.FormValue("members_min"); v != "" {
		criteria.MembersMin, _ = strconv.Atoi(v)
	}
	if v := r.FormValue("members_max"); v != "" {
		criteria.MembersMax, _ = strconv.Atoi(v)
	}

	criteria.Locations = r.Form["locations"]

	artists := api.GetAllArtists()
	filtered := utils.FilterArtists(artists, criteria)

	if len(criteria.Locations) > 0 {
		var final []models.FullArtist
		for _, artist := range filtered {
			full, err := api.GetFullArtistByID(artist.ID)
			if err != nil {
				continue
			}
			for _, loc := range criteria.Locations {
				if utils.SearchInLocations(full.LocationsList, loc) {
					final = append(final, *full)
					break
				}
			}
		}
		filtered = final
	}

	minCreation, maxCreation, minAlbum, maxAlbum := utils.GetYearRange(artists)
	minMembers, maxMembers := utils.GetMembersRange(artists)

	data := PageData{
		Artists:      filtered,
		MinCreation:  minCreation,
		MaxCreation:  maxCreation,
		MinAlbum:     minAlbum,
		MaxAlbum:     maxAlbum,
		MinMembers:   minMembers,
		MaxMembers:   maxMembers,
		AllLocations: utils.GetUniqueLocations(artists, api.Relations.Index),

		Query:    "",
		Filtered: true,
	}

	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// =======================
// SEARCH
// =======================
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	artists := api.GetAllArtists()
	var results []models.FullArtist

	query = strings.ToLower(query)

	for _, artist := range artists {
		match := false

		if strings.Contains(strings.ToLower(artist.Name), query) {
			match = true
		}

		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				match = true
				break
			}
		}

		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) ||
			strings.Contains(strconv.Itoa(artist.CreationDate), query) {
			match = true
		}

		full, err := api.GetFullArtistByID(artist.ID)
		if err == nil && utils.SearchInLocations(full.LocationsList, query) {
			match = true
		}

		if match && full != nil {
			full.FirstAlbumYear = utils.ExtractYear(full.FirstAlbum)
			results = append(results, *full)
		}
	}

	minCreation, maxCreation, minAlbum, maxAlbum := utils.GetYearRange(artists)
	minMembers, maxMembers := utils.GetMembersRange(artists)

	data := PageData{
		Artists:      results,
		MinCreation:  minCreation,
		MaxCreation:  maxCreation,
		MinAlbum:     minAlbum,
		MaxAlbum:     maxAlbum,
		MinMembers:   minMembers,
		MaxMembers:   maxMembers,
		AllLocations: utils.GetUniqueLocations(artists, api.Relations.Index),

		Query:    query,
		Filtered: true,
	}

	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// =======================
// AUTOCOMPLETE
// =======================
func SuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	artists := api.GetAllArtists()
	suggestions := utils.SearchArtists(artists, query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

// =======================
// ERROR
// =======================
func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	data := struct {
		Code    int
		Message string
	}{
		Code:    status,
		Message: http.StatusText(status),
	}

	if status == http.StatusNotFound {
		data.Message = "La page que vous recherchez n'existe pas."
	} else if status == http.StatusInternalServerError {
		data.Message = "Une erreur interne du serveur s'est produite."
	}

	if err := templates.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Error %d", status), status)
	}
}
