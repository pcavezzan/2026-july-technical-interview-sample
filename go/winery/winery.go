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
	Color  string  `json:"color,omitempty"`
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

	return &newWine, validate(&newWine)
}

func validate(newWine *Wine) error {
	if newWine.Price == 0.0 {
		return noPriceErr
	}

	if newWine.Price < 0.0 {
		return fmt.Errorf("wine price must be a positive floating value, got (price: %0.02f)", newWine.Price)
	}
	return nil
}

// FromObject creates a Wine instance from the provided Object map by extracting and converting its fields.
func FromObject(o Object) (*Wine, error) {
	var price float64
	if priceValue, ok := o["price"]; ok {
		price = parseFloat64(priceValue)
	}
	var color string
	if colorValue, ok := o["color"]; ok {
		color = colorValue.(string)
	}

	w := &Wine{
		Name:   o["name"].(string),
		Region: o["region"].(string),
		Year:   o["year"].(int),
		Price:  price,
		Color:  color,
	}
	return w, validate(w)
}

func parseFloat64(priceValue any) float64 {
	var price float64
	switch priceValue.(type) {
	case float64:
		price = priceValue.(float64)
	case int:
		price = float64(priceValue.(int))
	}
	return price
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
