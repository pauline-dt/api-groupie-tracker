package main

import (
	"fmt"
	"log"
	"net/http"
	"api-groupie-tracker/api"
	"api-groupie-tracker/handlers"
)

func main() {
	// Charger les données de l'API au démarrage
	if err := api.FetchAllData(); err != nil {
		log.Fatal("Erreur lors du chargement des données de l'API:", err)
	}

	// Configuration des routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist/", handlers.ArtistHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/filter", handlers.FilterHandler)
	http.HandleFunc("/api/suggestions", handlers.SuggestionsHandler)
	
	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Démarrer le serveur
	port := ":8080"
	fmt.Printf("🎵 Groupie Tracker démarré sur http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}