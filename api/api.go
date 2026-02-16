package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"api-groupie-tracker/models"
	"sync"
)

const (
	ArtistsURL   = "https://groupietrackers.herokuapp.com/api/artists"
	LocationsURL = "https://groupietrackers.herokuapp.com/api/locations"
	DatesURL     = "https://groupietrackers.herokuapp.com/api/dates"
	RelationsURL = "https://groupietrackers.herokuapp.com/api/relation"
)

var (
	Artists   []models.Artist
	Locations models.LocationIndex
	Dates     models.DateIndex
	Relations models.RelationIndex
	mutex     sync.RWMutex
)

// FetchAllData récupère toutes les données de l'API en parallèle
func FetchAllData() error {
	var wg sync.WaitGroup
	errors := make(chan error, 4)

	wg.Add(4)

	// Récupérer les artistes
	go func() {
		defer wg.Done()
		if err := fetchArtists(); err != nil {
			errors <- fmt.Errorf("erreur artistes: %w", err)
		}
	}()

	// Récupérer les locations
	go func() {
		defer wg.Done()
		if err := fetchLocations(); err != nil {
			errors <- fmt.Errorf("erreur locations: %w", err)
		}
	}()

	// Récupérer les dates
	go func() {
		defer wg.Done()
		if err := fetchDates(); err != nil {
			errors <- fmt.Errorf("erreur dates: %w", err)
		}
	}()

	// Récupérer les relations
	go func() {
		defer wg.Done()
		if err := fetchRelations(); err != nil {
			errors <- fmt.Errorf("erreur relations: %w", err)
		}
	}()

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchArtists() error {
	resp, err := http.Get(ArtistsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	mutex.Lock()
	err = json.Unmarshal(body, &Artists)
	mutex.Unlock()

	return err
}

func fetchLocations() error {
	resp, err := http.Get(LocationsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	mutex.Lock()
	err = json.Unmarshal(body, &Locations)
	mutex.Unlock()

	return err
}

func fetchDates() error {
	resp, err := http.Get(DatesURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	mutex.Lock()
	err = json.Unmarshal(body, &Dates)
	mutex.Unlock()

	return err
}

func fetchRelations() error {
	resp, err := http.Get(RelationsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	mutex.Lock()
	err = json.Unmarshal(body, &Relations)
	mutex.Unlock()

	return err
}

// GetArtistByID retourne un artiste par son ID
func GetArtistByID(id int) (*models.Artist, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, artist := range Artists {
		if artist.ID == id {
			return &artist, nil
		}
	}
	return nil, fmt.Errorf("artiste non trouvé")
}

// GetFullArtistByID retourne toutes les infos d'un artiste
func GetFullArtistByID(id int) (*models.FullArtist, error) {
	artist, err := GetArtistByID(id)
	if err != nil {
		return nil, err
	}

	mutex.RLock()
	defer mutex.RUnlock()

	fullArtist := &models.FullArtist{
		Artist: *artist,
	}

	// Ajouter les locations
	if id <= len(Locations.Index) {
		fullArtist.LocationsList = Locations.Index[id-1].Locations
	}

	// Ajouter les dates
	if id <= len(Dates.Index) {
		fullArtist.DatesList = Dates.Index[id-1].Dates
	}

	// Ajouter les relations
	if id <= len(Relations.Index) {
		fullArtist.DatesLocations = Relations.Index[id-1].DatesLocations
	}

	return fullArtist, nil
}

// GetAllArtists retourne tous les artistes
func GetAllArtists() []models.Artist {
	mutex.RLock()
	defer mutex.RUnlock()
	return Artists
}