package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	//URL API
	urlAPI := "https://api.rawg.io/api/"

	//init du client HTTP qui va émettre les reequêtes
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	//Création de la requête HTTP vers l'api avec init de la méthode HTTP, la route et le corps de la requête
	req, errReq := http.NewRequest(http.MethodGet, urlAPI, nil)
	if errReq != nil {
		fmt.Println("une erreur est survenue : ", errReq.Error())
	}

	//Ajout d'une métadonnée dans le header, User_Agent permet d'identifier l'application, système ...
	req.Header.Add("User-Agent", "Ynov campus cours")

	//execution de la requête HTTP vers l'API
	res, errResp := httpClient.Do(req)
	if errResp != nil {
		fmt.Println("Une erreur est survenue : ", errResp.Error())
		return
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	//lecture et récup du corps de la requête HTTP
	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		fmt.Println("Une erreur est survenue : ", errResp.Error())
	}

	//décla de la variable qui va contenir les données
	var decodeData ApiData //structure a faire apres l'erreur s'enlevera

	//decodage des données en format JSON et ajout des données à la variable: decodeData
	json.Unmarshal(body, &decodeData)

	//affichage des données
	fmt.Println(decodeData.Results[0])
}
