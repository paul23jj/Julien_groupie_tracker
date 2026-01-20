package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var RAWGAPIKey = "41b359451b3046c6b86785db848f7ded"

type PageData map[string]interface{}

func APropos() PageData {
	return PageData{
		"Title":   "À propos",
		"Content": "Contenu de la page À propos",
	}
}

func Categories() PageData {
	return PageData{
		"Title": "Catégories",
	}
}

func Collection(page int) PageData {
	return CollectionWithGenres(page, []string{})
}

// CollectionWithGenres récupère les jeux avec possibilité de filtrer par genres
func CollectionWithGenres(page int, genres []string) PageData {
	pd := PageData{
		"Title":        "Collection",
		"Page":         page,
		"Games":        []interface{}{},
		"PrevPage":     0,
		"NextPage":     0,
		"CollectionJS": GetCollectionJS(),
	}

	apiKey := getAPIKey()
	if apiKey == "" {
		pd["Error"] = "RAWG API key not set; returning empty collection"
		return pd
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.rawg.io/api/games", nil)
	if err != nil {
		pd["Error"] = err.Error()
		return pd
	}
	q := req.URL.Query()
	q.Add("key", apiKey)
	q.Add("page_size", "12")
	q.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		pd["Error"] = err.Error()
		return pd
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		pd["Error"] = "RAWG API returned status: " + resp.Status
		return pd
	}

	var parsed map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		pd["Error"] = err.Error()
		return pd
	}

	if res, ok := parsed["results"]; ok {
		if s, ok := res.([]interface{}); ok {
			// Appliquer le filtre par genres si des genres sont spécifiés
			filtered := FilterGamesByGenre(s, genres)
			pd["Games"] = filtered
			fmt.Println("Collection: fetched", len(s), "games, filtered to", len(filtered), "for page", page)
		} else {
			pd["Games"] = []interface{}{}
			fmt.Println("Collection: results present but not a slice, type:", fmt.Sprintf("%T", res))
		}
	} else {
		pd["Games"] = []interface{}{}
		fmt.Println("Collection: no results key in parsed response")
	}
	if parsed["next"] != nil {
		pd["NextPage"] = page + 1
	}
	if page > 1 {
		pd["PrevPage"] = page - 1
	}

	return pd
}

func Favoris() PageData {
	return PageData{
		"Title": "Favoris",
	}
}

func Recherche() PageData {
	return PageData{
		"Title": "Recherche",
	}
}

func Ressources() PageData {
	return PageData{
		"Title": "Ressources",
	}
}

func getAPIKey() string {
	k := os.Getenv("RAWG_API_KEY")
	if k != "" {
		return k
	}
	return RAWGAPIKey
}

func Search(query string) PageData {
	pd := PageData{
		"Title":   "Recherche",
		"Query":   query,
		"Results": []interface{}{},
	}
	if query == "" {
		return pd
	}

	apiKey := getAPIKey()
	if apiKey == "" {
		pd["Error"] = "RAWG API key not set; returning empty results"
		return pd
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", "https://api.rawg.io/api/games", nil)
	if err != nil {
		pd["Error"] = err.Error()
		return pd
	}
	q := req.URL.Query()
	q.Add("key", apiKey)
	q.Add("search", query)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		pd["Error"] = err.Error()
		return pd
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		pd["Error"] = "RAWG API returned status: " + resp.Status
		return pd
	}

	var parsed map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		pd["Error"] = err.Error()
		return pd
	}

	if res, ok := parsed["results"]; ok {
		pd["Results"] = res
	} else {
		pd["Results"] = []interface{}{}
	}
	return pd
}

func Index() PageData {
	pd := PageData{
		"Title":         "RAWR API Explorer",
		"FeaturedGames": []interface{}{},
		"Endpoints": []string{
			"/api/games",
			"/api/games/{id}",
			"/api/categories",
			"/api/search",
		},
	}

	apiKey := getAPIKey()
	if apiKey == "" {
		pd["Error"] = "RAWG_API_KEY not set; featured games disabled"
		return pd
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.rawg.io/api/games", nil)
	if err != nil {
		pd["Error"] = err.Error()
		return pd
	}
	q := req.URL.Query()
	q.Add("key", apiKey)
	q.Add("page_size", "4")
	q.Add("ordering", "-rating")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		pd["Error"] = err.Error()
		return pd
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		pd["Error"] = "RAWG API returned status: " + resp.Status
		return pd
	}

	var parsed map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		pd["Error"] = err.Error()
		return pd
	}
	if res, ok := parsed["results"]; ok {
		pd["FeaturedGames"] = res
	}
	return pd
}
