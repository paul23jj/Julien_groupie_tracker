package controller

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

// FormatDate formate une date en "JJ mois AAAA" (mois en français).
func FormatDate(v interface{}) string {
	if v == nil {
		return ""
	}

	switch t := v.(type) {
	case time.Time:
		return fmt.Sprintf("%02d %s %d", t.Day(), frenchMonth(t.Month()), t.Year())
	case *time.Time:
		if t == nil {
			return ""
		}
		return fmt.Sprintf("%02d %s %d", t.Day(), frenchMonth(t.Month()), t.Year())
	}

	s, ok := v.(string)
	if !ok || s == "" {
		return ""
	}

	layouts := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02 15:04:05",
	}
	var parsed time.Time
	var err error
	for _, l := range layouts {
		parsed, err = time.Parse(l, s)
		if err == nil {
			return fmt.Sprintf("%02d %s %d", parsed.Day(), frenchMonth(parsed.Month()), parsed.Year())
		}
	}
	// si on n'a pas réussi à parser, renvoyer la chaîne brute
	return s
}

func frenchMonth(m time.Month) string {
	months := []string{
		"",
		"janvier",
		"février",
		"mars",
		"avril",
		"mai",
		"juin",
		"juillet",
		"août",
		"septembre",
		"octobre",
		"novembre",
		"décembre",
	}
	if int(m) >= 1 && int(m) <= 12 {
		return months[int(m)]
	}
	return ""
}

// FilterGamesByGenre filtre les jeux en fonction des genres sélectionnés
func FilterGamesByGenre(games []interface{}, selectedGenres []string) []interface{} {
	if len(selectedGenres) == 0 {
		return games
	}

	var filtered []interface{}

	for _, game := range games {
		gameMap, ok := game.(map[string]interface{})
		if !ok {
			continue
		}

		if hasMatchingGenre(gameMap, selectedGenres) {
			filtered = append(filtered, game)
		}
	}

	return filtered
}

// hasMatchingGenre vérifie si un jeu a au moins un des genres sélectionnés
func hasMatchingGenre(game map[string]interface{}, selectedGenres []string) bool {
	genres, ok := game["genres"].([]interface{})
	if !ok {
		return false
	}

	for _, genre := range genres {
		genreMap, ok := genre.(map[string]interface{})
		if !ok {
			continue
		}

		genreName, ok := genreMap["name"].(string)
		if !ok {
			continue
		}

		for _, selected := range selectedGenres {
			if strings.EqualFold(genreName, selected) {
				return true
			}
		}
	}

	return false
}

// GetCollectionJS retourne le code JavaScript pour la page Collection
func GetCollectionJS() template.JS {
	jsCode := `const filterCheckboxes = document.querySelectorAll('.filter-checkbox');
		const resetBtn = document.getElementById('resetFilters');
		const gameCards = document.querySelectorAll('.game-card.clickable');
		const modal = document.getElementById('gameModal');
		const modalContent = document.getElementById('gameDetails');
		const closeBtn = document.querySelector('.close');

		// Récupérer les genres depuis les paramètres d'URL
		function getGenresFromURL() {
			const params = new URLSearchParams(window.location.search);
			return params.getAll('genre') || [];
		}

		// Restaurer les checkboxes cochées à partir de l'URL
		function restoreFilters() {
			const selectedGenres = getGenresFromURL();
			filterCheckboxes.forEach(checkbox => {
				if (selectedGenres.includes(checkbox.value)) {
					checkbox.checked = true;
				}
			});
		}

		// Construire l'URL avec les genres sélectionnés
		function buildFilterURL() {
			const selectedGenres = Array.from(filterCheckboxes)
				.filter(cb => cb.checked)
				.map(cb => cb.value);

			const params = new URLSearchParams();
			params.set('page', '1'); // Retour à la page 1 quand on filtre
			selectedGenres.forEach(genre => {
				params.append('genre', genre);
			});

			return '/collection?' + params.toString();
		}

		// Appliquer le filtre
		function applyFilters() {
			window.location.href = buildFilterURL();
		}

		// Réinitialiser les filtres
		function resetFilters() {
			window.location.href = '/collection';
		}

		// Afficher la modale avec les détails du jeu
		function openModal(game) {
			let genres = game.genres ? game.genres.map(g => g.name).join(', ') : 'N/A';
			let platforms = game.platforms ? game.platforms.map(p => p.platform.name).join(', ') : 'N/A';
			
			modalContent.innerHTML = ` + "`" + `
				<div class="modal-body">
					<img src="${game.background_image}" alt="${game.name}" style="width: 100%; border-radius: 8px; margin-bottom: 15px;">
					<h2>${game.name}</h2>
					<p><strong>Date de sortie :</strong> ${game.released || 'Inconnue'}</p>
					<p><strong>Note :</strong> ${game.rating || 'N/A'} / 5</p>
					<p><strong>Genres :</strong> ${genres}</p>
					<p><strong>Plateformes :</strong> ${platforms}</p>
					${game.description_raw ? ` + "`" + `<p><strong>Description :</strong><br>${game.description_raw}</p>` + "`" + ` : ''}
				</div>
			` + "`" + `;
			modal.style.display = 'block';
		}

		filterCheckboxes.forEach(checkbox => {
			checkbox.addEventListener('change', applyFilters);
		});

		resetBtn.addEventListener('click', resetFilters);

		gameCards.forEach(card => {
			card.addEventListener('click', () => {
				const gameData = JSON.parse(card.getAttribute('data-game'));
				openModal(gameData);
			});
		});

		closeBtn.addEventListener('click', () => {
			modal.style.display = 'none';
		});

		window.addEventListener('click', (event) => {
			if (event.target == modal) {
				modal.style.display = 'none';
			}
		});

		// Restaurer les filtres au chargement
		restoreFilters();`
	return template.JS(jsCode)
}
