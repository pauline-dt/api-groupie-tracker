package utils

import (
	"groupie-tracker/models"
	"testing"
)

func TestExtractYear(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"14-12-2000", 2000},
		{"01-01-1990", 1990},
		{"25-06-2015", 2015},
		{"invalid", 0},
		{"", 0},
	}

	for _, test := range tests {
		result := ExtractYear(test.input)
		if result != test.expected {
			t.Errorf("ExtractYear(%s) = %d; expected %d", test.input, result, test.expected)
		}
	}
}

func TestNormalizeLocation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"seattle-washington-usa", "seattle washington usa"},
		{"New_York_City", "new york city"},
		{"Paris-France", "paris france"},
		{"  Tokyo  ", "tokyo"},
	}

	for _, test := range tests {
		result := NormalizeLocation(test.input)
		if result != test.expected {
			t.Errorf("NormalizeLocation(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestLocationContains(t *testing.T) {
	tests := []struct {
		location string
		search   string
		expected bool
	}{
		{"seattle-washington-usa", "seattle", true},
		{"seattle-washington-usa", "washington", true},
		{"seattle-washington-usa", "usa", true},
		{"seattle-washington-usa", "california", false},
		{"New_York_City", "york", true},
		{"Paris-France", "paris", true},
	}

	for _, test := range tests {
		result := LocationContains(test.location, test.search)
		if result != test.expected {
			t.Errorf("LocationContains(%s, %s) = %v; expected %v", 
				test.location, test.search, result, test.expected)
		}
	}
}

func TestFilterArtists(t *testing.T) {
	artists := []models.Artist{
		{
			ID:           1,
			Name:         "Queen",
			Members:      []string{"Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"},
			CreationDate: 1970,
			FirstAlbum:   "14-12-1973",
		},
		{
			ID:           2,
			Name:         "Pink Floyd",
			Members:      []string{"David Gilmour", "Roger Waters", "Nick Mason"},
			CreationDate: 1965,
			FirstAlbum:   "05-08-1967",
		},
	}

	// Test filtre par date de création
	criteria := models.FilterCriteria{
		CreationDateMin: 1968,
		CreationDateMax: 1975,
	}
	filtered := FilterArtists(artists, criteria)
	if len(filtered) != 1 || filtered[0].Name != "Queen" {
		t.Errorf("Filter by creation date failed")
	}

	// Test filtre par nombre de membres
	criteria = models.FilterCriteria{
		MembersMin: 4,
		MembersMax: 4,
	}
	filtered = FilterArtists(artists, criteria)
	if len(filtered) != 1 || filtered[0].Name != "Queen" {
		t.Errorf("Filter by member count failed")
	}

	// Test filtre par année du premier album
	criteria = models.FilterCriteria{
		FirstAlbumMin: 1970,
		FirstAlbumMax: 1975,
	}
	filtered = FilterArtists(artists, criteria)
	if len(filtered) != 1 || filtered[0].Name != "Queen" {
		t.Errorf("Filter by first album year failed")
	}
}

func TestSearchArtists(t *testing.T) {
	artists := []models.Artist{
		{
			ID:           1,
			Name:         "Queen",
			Members:      []string{"Freddie Mercury", "Brian May"},
			CreationDate: 1970,
			FirstAlbum:   "14-12-1973",
		},
		{
			ID:           2,
			Name:         "Pink Floyd",
			Members:      []string{"David Gilmour"},
			CreationDate: 1965,
			FirstAlbum:   "05-08-1967",
		},
	}

	// Test recherche par nom d'artiste
	results := SearchArtists(artists, "queen")
	if len(results) == 0 {
		t.Errorf("Search by artist name failed")
	}

	// Test recherche par membre
	results = SearchArtists(artists, "freddie")
	if len(results) == 0 {
		t.Errorf("Search by member name failed")
	}

	// Test recherche vide
	results = SearchArtists(artists, "")
	if len(results) != 0 {
		t.Errorf("Empty search should return no results")
	}
}

func TestGetYearRange(t *testing.T) {
	artists := []models.Artist{
		{
			CreationDate: 1970,
			FirstAlbum:   "14-12-1973",
		},
		{
			CreationDate: 1965,
			FirstAlbum:   "05-08-1967",
		},
		{
			CreationDate: 1980,
			FirstAlbum:   "10-03-1982",
		},
	}

	minCreation, maxCreation, minAlbum, maxAlbum := GetYearRange(artists)

	if minCreation != 1965 {
		t.Errorf("Min creation year = %d; expected 1965", minCreation)
	}
	if maxCreation != 1980 {
		t.Errorf("Max creation year = %d; expected 1980", maxCreation)
	}
	if minAlbum != 1967 {
		t.Errorf("Min album year = %d; expected 1967", minAlbum)
	}
	if maxAlbum != 1982 {
		t.Errorf("Max album year = %d; expected 1982", maxAlbum)
	}
}

func TestGetMembersRange(t *testing.T) {
	artists := []models.Artist{
		{Members: []string{"A", "B"}},
		{Members: []string{"C", "D", "E", "F"}},
		{Members: []string{"G"}},
	}

	min, max := GetMembersRange(artists)

	if min != 1 {
		t.Errorf("Min members = %d; expected 1", min)
	}
	if max != 4 {
		t.Errorf("Max members = %d; expected 4", max)
	}
}