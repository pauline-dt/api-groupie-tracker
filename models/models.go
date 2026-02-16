package models

// Artist représente un artiste/groupe
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// Location représente les lieux de concerts
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

// LocationIndex contient tous les lieux
type LocationIndex struct {
	Index []Location `json:"index"`
}

// Date représente les dates de concerts
type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// DateIndex contient toutes les dates
type DateIndex struct {
	Index []Date `json:"index"`
}

// Relation lie les artistes, dates et lieux
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// RelationIndex contient toutes les relations
type RelationIndex struct {
	Index []Relation `json:"index"`
}

// FullArtist combine toutes les informations d'un artiste
type FullArtist struct {
	Artist
	LocationsList  []string
	DatesList      []string
	DatesLocations map[string][]string
	FirstAlbumYear int
}

// FilterCriteria représente les critères de filtrage
type FilterCriteria struct {
	CreationDateMin int
	CreationDateMax int
	FirstAlbumMin   int
	FirstAlbumMax   int
	MembersMin      int
	MembersMax      int
	Locations       []string
}

// SearchSuggestion représente une suggestion de recherche
type SearchSuggestion struct {
	Value string `json:"value"`
	Type  string `json:"type"`
	ID    int    `json:"id"`
}