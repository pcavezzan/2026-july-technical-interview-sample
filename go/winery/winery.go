package winery

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

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
func NewWine(o Object) (*Wine, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	var w Wine
	err = json.Unmarshal(data, &w)
	if err != nil {
		return nil, err
	}

	if w.Price == 0.0 {
		return nil, errors.New("wine must have a given price")
	}

	if w.Price < 0.0 {
		return nil, fmt.Errorf("wine price must be a positive floating value, got (price: %0.02f)", w.Price)
	}

	return &w, err
}

// ClassifyByColor Classifies all wines in a wine dictionary by color
func (c Cellar) ClassifyByColor() map[string]Cellar {
	var res map[string]Cellar

	// TODO: Candidate Codes

	return res
}

// SortByPrice sorts all wines by their price
func (c Cellar) SortByPrice(desc bool) Cellar {
	res := slices.Clone(c)
	var sort func(a, b Wine) int
	if desc {
		sort = func(a, b Wine) int {
			return cmp.Compare(b.Price, a.Price)
		}
	} else {
		sort = func(a, b Wine) int {
			return cmp.Compare(a.Price, b.Price)
		}
	}
	slices.SortStableFunc(res, sort)
	return res
}

// Search returns wine that matches the search terms on below wine properties
// - Name
// - Region
// - Year
// Note that the search implementation must be case insensitive
func (c Cellar) Search(str string) Cellar {
	var res Cellar
	searchTerm := fold(str)
	for _, wine := range c {
		if strings.Contains(fold(wine.Name), searchTerm) ||
			strings.Contains(fold(wine.Region), searchTerm) ||
			strings.Contains(fold(strconv.Itoa(wine.Year)), searchTerm) {
			res = append(res, wine)
		}
	}
	return res
}

var folder = cases.Fold()

// fold abaisse la casse ET retire les accents : "La tâche" -> "la tache"
func fold(s string) string {
	t := transform.Chain(
		norm.NFD,                           // é -> e + ´
		runes.Remove(runes.In(unicode.Mn)), // supprime les marques (Mn)
		norm.NFC,                           // recompose
	)
	out, _, _ := transform.String(t, s)
	return folder.String(out)
}
