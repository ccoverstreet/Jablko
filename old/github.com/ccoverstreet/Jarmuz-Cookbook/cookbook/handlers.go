package cookbook

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Handles called by Jablko

//go:embed webcomponent.js
var webcomponentFile []byte

func (book *Cookbook) WebComponentHandler(w http.ResponseWriter, r *http.Request) {
	book.RLock()
	defer book.RUnlock()

	//fmt.Fprintf(w, "%s", webcomponentFile)
	b, err := ioutil.ReadFile("./cookbook/webcomponent.js")
	if err != nil {
		log.Printf("ERROR: Unable to read webcomponent file - %v", err)
	}

	fmt.Fprintf(w, "%s", b)
}

func (book *Cookbook) InstanceHandler(w http.ResponseWriter, r *http.Request) {
	book.RLock()
	defer book.RUnlock()

	b, err := json.Marshal(book.Instances)
	if err != nil {
		log.Printf("ERROR: Unable to marshal data - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", b)
}

// App Routes

func handleError(message string, err error, w http.ResponseWriter) {
	log.Printf("ERROR: %s - %v", message, err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%v", err)
}

func (book *Cookbook) GetRecipeListHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(book.GetRecipeNames())
	if err != nil {
		handleError("Unable to marshal", err, w)
		return
	}

	fmt.Fprintf(w, "%s", b)
}

func (book *Cookbook) AddRecipeHandler(w http.ResponseWriter, r *http.Request) {
	type NewRecipe struct {
		Name         string `json:"name"`
		Ingredients  string `json:"ingredients"`
		Instructions string `json:"instructions"`
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError("Unable to read body", err, w)
		return
	}

	recipe := NewRecipe{}
	err = json.Unmarshal(b, &recipe)
	if err != nil {
		handleError("Unable to marshal body", err, w)
		return
	}

	err = book.AddRecipe(recipe.Name, recipe.Ingredients, recipe.Instructions)
	if err != nil {
		handleError("Unable to add recipe", err, w)
		return
	}
}

func (book *Cookbook) RemoveRecipeHandler(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Name string `json:"name"`
	}

	req := ReqBody{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError("Unable to read body", err, w)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		handleError("Unable to marshal body", err, w)
		return
	}

	err = book.RemoveRecipe(req.Name)
	if err != nil {
		handleError("Unable to remove recipe", err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Removed recipe")
}

func (book *Cookbook) GetRecipeHandler(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Name string `json:"name"`
	}

	req := ReqBody{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError("Unable to read body", err, w)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		handleError("Unable to marshal body", err, w)
		return
	}

	recipe, err := book.GetRecipe(req.Name)
	if err != nil {
		handleError("Unable to get recipe", err, w)
		return
	}

	b, err := json.Marshal(recipe)
	if err != nil {
		handleError("Unable to marshal recipe", err, w)
		return
	}

	fmt.Fprintf(w, "%s", b)
}

func (book *Cookbook) UpdateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Name         string `json:"name"`
		Ingredients  string `json:"ingredients"`
		Instructions string `json:"instructions"`
	}

	req := ReqBody{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError("Unable to read body", err, w)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		handleError("Unable to marshal body", err, w)
		return
	}

	err = book.UpdateRecipe(req.Name, req.Ingredients, req.Instructions)
	if err != nil {
		handleError("Unable to update recipe", err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Updated recipe")
}
