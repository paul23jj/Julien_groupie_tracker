package controller

// PageData représente les données passées aux templates.
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

func Collection() PageData {
	return PageData{
		"Title": "Collection",
	}
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
