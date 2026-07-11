package winery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var noPriceErr = errors.New("wine must have a given price")

const (
	COLOR_RED   = ""
	COLOR_WHITE = "white"
	COLOR_ROSE  = "rose"
)

type Object = map[string]any

// Define what is a wine for our winery
type Wine struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Year   int     `json:"year"`
	Color  string  `json:"color"`
	Region string  `json:"region"`
}

// Catalog of wines stored in the winery
type Cellar []Wine

// Length return the number of wine in the cellar
func (c Cellar) Length() int {
	return len(c)
}

// Dump returns a concatenation of all wine Ids in the cellar (useful for test)
func (c Cellar) dump() string {
	var ids = make([]string, 0, len(c))
	for _, wine := range c {
		ids = append(ids, strconv.Itoa(int(wine.Id)))
	}
	return strings.Join(ids, ",")
}

/******************************************************************************************/
/******************************************************************************************/
/*********************************** FUNCTIONS TO CODE ************************************/
/******************************************************************************************/
/******************************************************************************************/

// NewWine handles all error handling when creating a wine
func NewWine(w Object) (*Wine, error) {
	data, err := json.Marshal(w)
	if err != nil {
		return nil, err
	}
	var newWine Wine
	err = json.Unmarshal(data, &newWine)
	if err != nil {
		return nil, err
	}

	if newWine.Price == 0.0 {
		return nil, noPriceErr
	}

	if newWine.Price < 0.0 {
		return nil, fmt.Errorf("wine price must be a positive floating value, got (price: %0.02f)", newWine.Price)
	}

	return &newWine, nil
}

// ClassifyByColor Classifies all wines in a wine dictionary by color
func (c Cellar) ClassifyByColor() map[string]Cellar {
	var res map[string]Cellar

	// TODO: Candidate Codes

	return res
}

// SortByPrice sorts all wines by their price
func (c Cellar) SortByPrice(desc bool) Cellar {
	var res Cellar

	// TODO: Candidate Codes

	return res
}

// Search returns wine that matches the search terms on below wine properties
// - Name
// - Region
// - Year
// Note that the search implementation must be case insensitive
func (c Cellar) Search(str string) Cellar {
	var res Cellar

	// TODO: Candidate Codes

	return res
}
