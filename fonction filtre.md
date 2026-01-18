package main



import (

&nbsp;	"encoding/json"

&nbsp;	"log"

&nbsp;	"net/http"

&nbsp;	"strconv"

&nbsp;	"strings"

)



// Structure représentant un produit

type Product struct {

&nbsp;	ID    int     `json:"id"`

&nbsp;	Name  string  `json:"name"`

&nbsp;	Price float64 `json:"price"`

}



// Données simulées (base en mémoire)

var products = \[]Product{

&nbsp;	{1, "Chaussures", 49.99},

&nbsp;	{2, "T-shirt", 19.99},

&nbsp;	{3, "Pantalon", 39.99},

&nbsp;	{4, "Veste", 89.99},

&nbsp;	{5, "Casquette", 14.99},

}



// Handler API avec filtres

func getProducts(w http.ResponseWriter, r \*http.Request) {

&nbsp;	w.Header().Set("Content-Type", "application/json")



&nbsp;	// Récupération des paramètres de filtre

&nbsp;	search := strings.ToLower(r.URL.Query().Get("search"))

&nbsp;	minPriceStr := r.URL.Query().Get("minPrice")

&nbsp;	maxPriceStr := r.URL.Query().Get("maxPrice")



&nbsp;	var minPrice, maxPrice float64

&nbsp;	var err error



&nbsp;	if minPriceStr != "" {

&nbsp;		minPrice, err = strconv.ParseFloat(minPriceStr, 64)

&nbsp;		if err != nil {

&nbsp;			http.Error(w, "minPrice invalide", http.StatusBadRequest)

&nbsp;			return

&nbsp;		}

&nbsp;	}

&nbsp;	if maxPriceStr != "" {

&nbsp;		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)

&nbsp;		if err != nil {

&nbsp;			http.Error(w, "maxPrice invalide", http.StatusBadRequest)

&nbsp;			return

&nbsp;		}

&nbsp;	}



&nbsp;	// Filtrage

&nbsp;	var filtered \[]Product

&nbsp;	for \_, p := range products {

&nbsp;		if search != "" \&\& !strings.Contains(strings.ToLower(p.Name), search) {

&nbsp;			continue

&nbsp;		}

&nbsp;		if minPriceStr != "" \&\& p.Price < minPrice {

&nbsp;			continue

&nbsp;		}

&nbsp;		if maxPriceStr != "" \&\& p.Price > maxPrice {

&nbsp;			continue

&nbsp;		}

&nbsp;		filtered = append(filtered, p)

&nbsp;	}



&nbsp;	json.NewEncoder(w).Encode(filtered)

}



func main() {

&nbsp;	http.HandleFunc("/api/products", getProducts)



&nbsp;	log.Println("Serveur démarré sur http://localhost:8080")

&nbsp;	log.Fatal(http.ListenAndServe(":8080", nil))

}




<!DOCTYPE html>

<html lang="fr">

<head>

&nbsp;   <meta charset="UTF-8">

&nbsp;   <title>Filtre Produits</title>

</head>

<body>

&nbsp;   <h1>Filtrer les produits</h1>

&nbsp;   <form id="filterForm">

&nbsp;       <input type="text" id="search" placeholder="Recherche par nom">

&nbsp;       <input type="number" step="0.01" id="minPrice" placeholder="Prix min">

&nbsp;       <input type="number" step="0.01" id="maxPrice" placeholder="Prix max">

&nbsp;       <button type="submit">Filtrer</button>

&nbsp;   </form>



&nbsp;   <h2>Résultats</h2>

&nbsp;   <ul id="results"></ul>



&nbsp;   <script>

&nbsp;       document.getElementById("filterForm").addEventListener("submit", function(e) {

&nbsp;           e.preventDefault();



&nbsp;           const search = document.getElementById("search").value;

&nbsp;           const minPrice = document.getElementById("minPrice").value;

&nbsp;           const maxPrice = document.getElementById("maxPrice").value;



&nbsp;           let url = "http://localhost:8080/api/products?";

&nbsp;           if (search) url += `search=${encodeURIComponent(search)}\&`;

&nbsp;           if (minPrice) url += `minPrice=${encodeURIComponent(minPrice)}\&`;

&nbsp;           if (maxPrice) url += `maxPrice=${encodeURIComponent(maxPrice)}`;



&nbsp;           fetch(url)

&nbsp;               .then(res => res.json())

&nbsp;               .then(data => {

&nbsp;                   const results = document.getElementById("results");

&nbsp;                   results.innerHTML = "";

&nbsp;                   if (data.length === 0) {

&nbsp;                       results.innerHTML = "<li>Aucun produit trouvé</li>";

&nbsp;                       return;

&nbsp;                   }

&nbsp;                   data.forEach(p => {

&nbsp;                       results.innerHTML += `<li>${p.name} - ${p.price} €</li>`;

&nbsp;                   });

&nbsp;               })

&nbsp;               .catch(err => console.error(err));

&nbsp;       });

&nbsp;   </script>

</body>

</html>



