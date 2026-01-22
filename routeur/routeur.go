package routeur

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"Steam-API/controller"
)

// New crée les routes
func New() http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("icons"))))

	mux.HandleFunc("/aPropos", aProposHandler)
	mux.HandleFunc("/collection", collectionHandler)
	mux.HandleFunc("/favoris", favorisHandler)
	mux.HandleFunc("/recherche", rechercheHandler)
	mux.HandleFunc("/traitment/search", traitmentSearchHandler)
	mux.HandleFunc("/traitement/search", traitmentSearchHandler)

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

func jsonMarshal(v interface{}) template.JS {
	b, _ := json.Marshal(v)
	return template.JS(string(b))
}

func collectionHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if pi, err := strconv.Atoi(p); err == nil && pi > 0 {
			page = pi
		}
	}

	// Récupérer les genres depuis les paramètres GET
	genres := r.URL.Query()["genre"]

	// Utiliser CollectionWithGenres si des genres sont spécifiés, sinon Collection
	var data controller.PageData
	if len(genres) > 0 {
		data = controller.CollectionWithGenres(page, genres)
	} else {
		data = controller.Collection(page)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// create template with funcs so formatDate is available inside templates
	tmpl := template.New("collection.html").Funcs(template.FuncMap{
		"formatDate":    controller.FormatDate,
		"convertToJSON": jsonMarshal,
	})
	tmpl, err := tmpl.ParseFiles("template/collection.html")
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
	// Récupérer le paramètre de recherche depuis l'URL (?search=...)
	query := r.URL.Query().Get("search")
	data := controller.Recherche(query)
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

func traitmentSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse error", http.StatusBadRequest)
		return
	}
	// Récupérer la valeur du champ "search" du formulaire
	query := r.FormValue("search")
	// CORRIGÉ: Passer query à la fonction Recherche
	data := controller.Recherche(query)
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := controller.Index()
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
