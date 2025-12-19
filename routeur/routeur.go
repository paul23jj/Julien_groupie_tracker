package routeur

import (
	"net/http"
)

func InitRoutes() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	tpl := http.FileServer(http.Dir("template"))
	http.Handle("/template/", http.StripPrefix("/template/", tpl))

	http.HandleFunc("/aPropos", aProposHandler)
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/collection", collectionHandler)
	http.HandleFunc("/favoris", favorisHandler)
	http.HandleFunc("/recherche", rechercheHandler)
	http.HandleFunc("/ressources", ressourcesHandler)

	http.HandleFunc("/", indexHandler)
}

//faire les fonction pour les handlers
