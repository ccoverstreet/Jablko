// Cookbook implementation/
// Cale Overstreet
// Jun 27, 2021

/* Handles recipe database and JMOD routing
 */

package cookbook

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/ccoverstreet/Jarmuz-Cookbook/jablkodev"
)

const defaultConfig = `
{
	"instances": [
		{
			"title": "Cookbook"
		}
	]
}
`

type Instance struct {
	Title string `json:"title"`
}

type Recipe struct {
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
}

type Cookbook struct {
	sync.RWMutex
	Instances      []Instance `json:"instances"`
	jablkoCorePort string
	jmodPort       string
	jmodKey        string
	jmodDataDir    string

	mux *http.ServeMux

	recipes map[string]Recipe
}

func CreateCookbook(jablkoCorePort, jmodPort, jmodKey, jmodDataDir, jmodConfig string) *Cookbook {
	book := &Cookbook{
		sync.RWMutex{},
		nil,
		jablkoCorePort,
		jmodPort,
		jmodKey,
		jmodDataDir,
		http.NewServeMux(),
		make(map[string]Recipe),
	}

	fShouldSave := len(jmodConfig) < 4

	if fShouldSave {
		jmodConfig = defaultConfig
	}

	err := json.Unmarshal([]byte(jmodConfig), &book)
	if err != nil {
		panic(err)
	}

	if fShouldSave {
		book.SaveConfig()
	}

	// Try to read recipe database if it exists
	b, err := ioutil.ReadFile(jmodDataDir + "/jarmuzrecipes.json")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Println(err)
			panic(err)
		}

		b = []byte("{}")
	}

	err = json.Unmarshal(b, &book.recipes)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// Routes
	book.mux.HandleFunc("/webComponent", book.WebComponentHandler)
	book.mux.HandleFunc("/instanceData", book.InstanceHandler)

	book.mux.HandleFunc("/jmod/getRecipeList", book.GetRecipeListHandler)
	book.mux.HandleFunc("/jmod/addRecipe", book.AddRecipeHandler)
	book.mux.HandleFunc("/jmod/removeRecipe", book.RemoveRecipeHandler)
	book.mux.HandleFunc("/jmod/getRecipe", book.GetRecipeHandler)
	book.mux.HandleFunc("/jmod/updateRecipe", book.UpdateRecipeHandler)

	return book
}

func (book *Cookbook) GetRouter() *http.ServeMux {
	return book.mux
}

func (book *Cookbook) SaveConfig() error {
	b, err := json.Marshal(book)
	if err != nil {
		return err
	}

	err = jablkodev.JablkoSaveConfig(book.jablkoCorePort,
		book.jmodPort,
		book.jmodKey,
		b)

	return nil
}

func (book *Cookbook) SaveRecipeDatabase() error {
	b, err := json.MarshalIndent(book.recipes, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(book.jmodDataDir+"/jarmuzrecipes.json", b, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (book *Cookbook) GetRecipeNames() []string {
	book.RLock()
	defer book.RUnlock()
	names := []string{}
	for name, _ := range book.recipes {
		names = append(names, name)
	}

	// From StackOverflow
	sort.Slice(names, func(i, j int) bool { return strings.ToLower(names[i]) < strings.ToLower(names[j]) })

	return names
}

func (book *Cookbook) AddRecipe(name string, ingredients string, instructions string) error {
	book.Lock()
	defer book.Unlock()

	// Check if recipe already exists
	if _, ok := book.recipes[name]; ok {
		return fmt.Errorf("Recipe already exists")
	}

	book.recipes[name] = Recipe{ingredients, instructions}

	return book.SaveRecipeDatabase()
}

func (book *Cookbook) RemoveRecipe(name string) error {
	book.Lock()
	defer book.Unlock()

	if _, ok := book.recipes[name]; !ok {
		return fmt.Errorf("Recipe does not exist")
	}

	delete(book.recipes, name)

	return book.SaveRecipeDatabase()
}

func (book *Cookbook) GetRecipe(name string) (Recipe, error) {
	book.RLock()
	defer book.RUnlock()

	recipe, ok := book.recipes[name]
	if !ok {
		return Recipe{}, fmt.Errorf("Recipe does not exist")
	}

	return recipe, nil
}

func (book *Cookbook) UpdateRecipe(name string, ingredients string, instructions string) error {
	if _, ok := book.recipes[name]; !ok {
		return fmt.Errorf("Recipe does not exist")
	}

	book.recipes[name] = Recipe{ingredients, instructions}

	return book.SaveRecipeDatabase()
}
