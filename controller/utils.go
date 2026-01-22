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

var favorisStore = make(map[string]map[string]interface{})

func AddFavoris(gameID string, gameData map[string]interface{}) {
	favorisStore[gameID] = gameData
	fmt.Println("Ajouté aux favoris :", gameID)
}

func RemoveFavoris(gameID string) {
	delete(favorisStore, gameID)
	fmt.Println("Supprimé des favoris :", gameID)
}

func GetFavoris() []interface{} {
	favoris := make([]interface{}, 0, len(favorisStore))
	for _, game := range favorisStore {
		favoris = append(favoris, game)
	}
	return favoris
}

func IsFavoris(gameID string) bool {
	_, exists := favorisStore[gameID]
	return exists
}

// GetCollectionJS retourne le code JavaScript pour la page Collection
func GetCollectionJS() template.JS {
	jsCode := `
		// ============ GESTION DES FILTRES ============
		const filterCheckboxes = document.querySelectorAll('.filter-checkbox');
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
			params.set('page', '1');
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

		// ============ GESTION DES FAVORIS ============
		
		// Récupérer les favoris depuis localStorage
		function getFavoris() {
			const favoris = localStorage.getItem('favoris');
			return favoris ? JSON.parse(favoris) : {};
		}

		// Sauvegarder les favoris dans localStorage
		function saveFavoris(favoris) {
			localStorage.setItem('favoris', JSON.stringify(favoris));
		}

		// Ajouter un jeu aux favoris
		async function addToFavoris(gameId, gameData) {
			try {
				const response = await fetch('/api/favoris/add', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify(gameData)
				});

				if (response.ok) {
					const favoris = getFavoris();
					favoris[gameId] = gameData;
					saveFavoris(favoris);
					console.log('Ajouté aux favoris:', gameId);
					return true;
				}
			} catch (error) {
				console.error('Erreur lors de l\'ajout aux favoris:', error);
			}
			return false;
		}

		// Retirer un jeu des favoris
		async function removeFromFavoris(gameId) {
			try {
				const response = await fetch('/api/favoris/remove', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({ id: gameId })
				});

				if (response.ok) {
					const favoris = getFavoris();
					delete favoris[gameId];
					saveFavoris(favoris);
					console.log('Retiré des favoris:', gameId);
					return true;
				}
			} catch (error) {
				console.error('Erreur lors du retrait des favoris:', error);
			}
			return false;
		}

		// Vérifier si un jeu est dans les favoris
		function isFavoris(gameId) {
			const favoris = getFavoris();
			return favoris.hasOwnProperty(gameId.toString());
		}

		// Mettre à jour l'état visuel du bouton favori
		function updateFavoriteButton(button, isFav) {
			const heartPath = button.querySelector('.heart-icon path');
			if (isFav) {
				button.classList.add('active');
				heartPath.setAttribute('fill', 'currentColor');
			} else {
				button.classList.remove('active');
				heartPath.setAttribute('fill', 'none');
			}
		}

		// Initialiser les boutons favoris
		function initFavoriteButtons() {
			const favoriteButtons = document.querySelectorAll('.favorite-btn');
			
			favoriteButtons.forEach(button => {
				const gameId = button.getAttribute('data-game-id');
				
				// Restaurer l'état depuis localStorage
				if (isFavoris(gameId)) {
					updateFavoriteButton(button, true);
				}

				// Gérer le clic sur le bouton favori
				button.addEventListener('click', async (e) => {
					e.stopPropagation(); // Empêcher l'ouverture de la modale
					
					const gameData = JSON.parse(button.getAttribute('data-game'));
					const isCurrentlyFav = isFavoris(gameId);

					if (isCurrentlyFav) {
						const success = await removeFromFavoris(gameId);
						if (success) {
							updateFavoriteButton(button, false);
							
							// Si on est sur la page favoris, retirer la carte
							if (window.location.pathname === '/favoris') {
								button.closest('.game-card').style.animation = 'fadeOut 0.3s ease';
								setTimeout(() => {
									button.closest('.game-card').remove();
									
									// Vérifier s'il reste des jeux
									const remainingCards = document.querySelectorAll('.game-card');
									if (remainingCards.length === 0) {
										location.reload(); // Recharger pour afficher le message vide
									}
								}, 300);
							}
						}
					} else {
						const success = await addToFavoris(gameId, gameData);
						if (success) {
							updateFavoriteButton(button, true);
						}
					}
				});
			});
		}

		// ============ GESTION DE LA MODALE ============
		
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

		// ============ EVENT LISTENERS ============
		
		// Filtres (uniquement sur la page collection)
		if (filterCheckboxes.length > 0) {
			filterCheckboxes.forEach(checkbox => {
				checkbox.addEventListener('change', applyFilters);
			});
		}

		if (resetBtn) {
			resetBtn.addEventListener('click', resetFilters);
		}

		// Cartes de jeux cliquables
		gameCards.forEach(card => {
			card.addEventListener('click', (e) => {
				// Ne pas ouvrir la modale si on clique sur le bouton favori
				if (e.target.closest('.favorite-btn')) {
					return;
				}
				const gameData = JSON.parse(card.getAttribute('data-game'));
				openModal(gameData);
			});
		});

		// Modale
		if (closeBtn) {
			closeBtn.addEventListener('click', () => {
				modal.style.display = 'none';
			});
		}

		if (modal) {
			window.addEventListener('click', (event) => {
				if (event.target == modal) {
					modal.style.display = 'none';
				}
			});
		}

		// ============ INITIALISATION ============
		
		// Restaurer les filtres au chargement (uniquement sur la page collection)
		if (filterCheckboxes.length > 0) {
			restoreFilters();
		}

		// Initialiser les boutons favoris
		initFavoriteButtons();
	`
	return template.JS(jsCode)
}
