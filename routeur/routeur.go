package routeur

import (
	"html/template"
	"net/http"

	"Steam-API/controller"
)

// New construit et retourne un http.Handler (ServeMux) avec toutes les routes enregistrées.
func New() http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/aPropos", aProposHandler)
	mux.HandleFunc("/categories", categoriesHandler)
	mux.HandleFunc("/collection", collectionHandler)
	mux.HandleFunc("/favoris", favorisHandler)
	mux.HandleFunc("/recherche", rechercheHandler)
	mux.HandleFunc("/ressources", ressourcesHandler)

	mux.HandleFunc("/", indexHandler)

	return mux
}

func aProposHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.APropos()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("template/aPropos.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.Categories()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("template/categories.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func collectionHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.Collection()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("template/collection.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func favorisHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.Favoris()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("template/favoris.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rechercheHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.Recherche()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("template/recherche.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ressourcesHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.Ressources()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("template/ressources.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Accueil"}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
