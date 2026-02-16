package utils

import (
	"api-groupie-tracker/models"
	"strconv"
	"strings"
)

// ExtractYear extrait l'année d'une date au format "DD-MM-YYYY"
func ExtractYear(dateStr string) int {
	parts := strings.Split(dateStr, "-")
	if len(parts) == 3 {
		year, _ := strconv.Atoi(parts[2])
		return year
	}
	return 0
}

// FilterArtists filtre les artistes selon les critères
func FilterArtists(artists []models.Artist, criteria models.FilterCriteria) []models.FullArtist {
	var filtered []models.FullArtist

	for _, artist := range artists {
		// Filtre par date de création
		if criteria.CreationDateMin > 0 && artist.CreationDate < criteria.CreationDateMin {
			continue
		}
		if criteria.CreationDateMax > 0 && artist.CreationDate > criteria.CreationDateMax {
			continue
		}

		// Filtre par année du premier album
		firstAlbumYear := ExtractYear(artist.FirstAlbum)
		if criteria.FirstAlbumMin > 0 && firstAlbumYear < criteria.FirstAlbumMin {
			continue
		}
		if criteria.FirstAlbumMax > 0 && firstAlbumYear > criteria.FirstAlbumMax {
			continue
		}

		// Filtre par nombre de membres
		memberCount := len(artist.Members)
		if criteria.MembersMin > 0 && memberCount < criteria.MembersMin {
			continue
		}
		if criteria.MembersMax > 0 && memberCount > criteria.MembersMax {
			continue
		}

		// Créer FullArtist
		fullArtist := models.FullArtist{
			Artist:         artist,
			FirstAlbumYear: firstAlbumYear,
		}

		filtered = append(filtered, fullArtist)
	}

	return filtered
}

// NormalizeLocation normalise une location pour la comparaison
func NormalizeLocation(location string) string {
	location = strings.ToLower(location)
	location = strings.ReplaceAll(location, "_", " ")
	location = strings.ReplaceAll(location, "-", " ")
	return strings.TrimSpace(location)
}

// LocationContains vérifie si une location contient une sous-chaîne
func LocationContains(location, search string) bool {
	normLoc := NormalizeLocation(location)
	normSearch := NormalizeLocation(search)
	return strings.Contains(normLoc, normSearch)
}

// SearchArtists recherche dans les artistes
func SearchArtists(artists []models.Artist, query string) []models.SearchSuggestion {
	var suggestions []models.SearchSuggestion
	query = strings.ToLower(strings.TrimSpace(query))

	if query == "" {
		return suggestions
	}

	seen := make(map[string]bool)

	for _, artist := range artists {
		// Recherche dans le nom de l'artiste/groupe
		if strings.Contains(strings.ToLower(artist.Name), query) {
			key := "artist-" + artist.Name
			if !seen[key] {
				suggestions = append(suggestions, models.SearchSuggestion{
					Value: artist.Name,
					Type:  "artist/band",
					ID:    artist.ID,
				})
				seen[key] = true
			}
		}

		// Recherche dans les membres
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				key := "member-" + member
				if !seen[key] {
					suggestions = append(suggestions, models.SearchSuggestion{
						Value: member,
						Type:  "member",
						ID:    artist.ID,
					})
					seen[key] = true
				}
			}
		}

		// Recherche dans la date de création
		creationStr := strconv.Itoa(artist.CreationDate)
		if strings.Contains(creationStr, query) {
			key := "creation-" + creationStr
			if !seen[key] {
				suggestions = append(suggestions, models.SearchSuggestion{
					Value: creationStr,
					Type:  "creation date",
					ID:    artist.ID,
				})
				seen[key] = true
			}
		}

		// Recherche dans la date du premier album
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			key := "firstalbum-" + artist.FirstAlbum
			if !seen[key] {
				suggestions = append(suggestions, models.SearchSuggestion{
					Value: artist.FirstAlbum,
					Type:  "first album date",
					ID:    artist.ID,
				})
				seen[key] = true
			}
		}
	}

	// Limiter à 10 suggestions
	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}

	return suggestions
}

// SearchInLocations recherche dans les locations avec support des relations géographiques
func SearchInLocations(locations []string, query string) bool {
	query = NormalizeLocation(query)

	for _, location := range locations {
		if LocationContains(location, query) {
			return true
		}
	}

	return false
}

// GetUniqueLocations retourne toutes les locations uniques
func GetUniqueLocations(artists []models.Artist, relations []models.Relation) []string {
	locationSet := make(map[string]bool)
	var locations []string

	for _, relation := range relations {
		for location := range relation.DatesLocations {
			normalized := NormalizeLocation(location)
			if !locationSet[normalized] {
				locationSet[normalized] = true
				locations = append(locations, location)
			}
		}
	}

	return locations
}

// GetYearRange retourne la plage d'années pour les filtres
func GetYearRange(artists []models.Artist) (minCreation, maxCreation, minAlbum, maxAlbum int) {
	if len(artists) == 0 {
		return
	}

	minCreation = artists[0].CreationDate
	maxCreation = artists[0].CreationDate
	minAlbum = ExtractYear(artists[0].FirstAlbum)
	maxAlbum = minAlbum

	for _, artist := range artists {
		if artist.CreationDate < minCreation {
			minCreation = artist.CreationDate
		}
		if artist.CreationDate > maxCreation {
			maxCreation = artist.CreationDate
		}

		year := ExtractYear(artist.FirstAlbum)
		if year > 0 {
			if year < minAlbum {
				minAlbum = year
			}
			if year > maxAlbum {
				maxAlbum = year
			}
		}
	}

	return
}

// GetMembersRange retourne la plage de nombre de membres
func GetMembersRange(artists []models.Artist) (min, max int) {
	if len(artists) == 0 {
		return
	}

	min = len(artists[0].Members)
	max = min

	for _, artist := range artists {
		count := len(artist.Members)
		if count < min {
			min = count
		}
		if count > max {
			max = count
		}
	}

	return
}