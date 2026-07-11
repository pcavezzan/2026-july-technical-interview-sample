package winery

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/******************************************************************************************/
/****************************** WINERY TEST (TODO: 3/4) ***********************************/
/******************************************************************************************/

func TestWineCreation(t *testing.T) {
	var w1 = Object{
		"name":   "Château Angelus",
		"region": "Bordeaux",
		"year":   2017,
		"color":  COLOR_RED,
	}
	_, err := NewWine(w1)
	require.EqualError(t, err, "wine must have a given price")

	var w2 = Object{
		"name":   "Château Angelus",
		"region": "Bordeaux",
		"year":   2017,
		"price":  -1928.,
	}
	_, err = NewWine(w2)
	require.EqualError(t, err, "wine price must be a positive floating value, got (price: -1928.00)")

	var wOK = Object{
		"name":   "Château Angelus",
		"region": "Bordeaux",
		"year":   2017,
		"price":  1928.,
	}
	w, err := NewWine(wOK)
	require.NoError(t, err)
	require.Equal(t, w.Year, 2017, "wine properly created")
}

func TestClassifyByColor(t *testing.T) {
	cellar := openWineCatalog(t)
	assert.Equal(t, cellar.Length(), 6)

	byColor := cellar.ClassifyByColor()
	redWines := byColor[COLOR_RED]
	whiteWines := byColor[COLOR_WHITE]
	roseWines := byColor[COLOR_ROSE]

	assert.Equal(t, 4, redWines.Length())
	assert.Equal(t, 2, whiteWines.Length())
	assert.Nil(t, roseWines)

	cellar = append(cellar, Wine{
		Name:   "Château Puech Haut",
		Region: "Languedoc",
		Year:   2018,
		Price:  18.95,
		Color:  COLOR_ROSE,
	})

	byColor = cellar.ClassifyByColor()
	redWines = byColor[COLOR_RED]
	whiteWines = byColor[COLOR_WHITE]
	roseWines = byColor[COLOR_ROSE]

	assert.Equal(t, 4, redWines.Length())
	assert.Equal(t, 2, whiteWines.Length())
	assert.Equal(t, 1, roseWines.Length())

	assert.Equal(t, "0,2,3,4", redWines.dump())
	assert.Equal(t, "1,5", whiteWines.dump())
}

func TestSortPrice(t *testing.T) {
	cellar := openWineCatalog(t)
	assert.Equal(t, 6, cellar.Length())

	sortAsc := cellar.SortByPrice(false)
	assert.Equal(t, "5,3,1,2,0,4", sortAsc.dump())
	assert.Equal(t, "0,1,2,3,4,5", cellar.dump())

	sortDesc := cellar.SortByPrice(true)
	assert.Equal(t, "4,0,1,2,3,5", sortDesc.dump())
	assert.Equal(t, "0,1,2,3,4,5", cellar.dump())
}

func TestSearchWines(t *testing.T) {
	cellar := openWineCatalog(t)
	assert.Equal(t, 6, cellar.Length())

	res1 := cellar.Search("")
	assert.Equal(t, "0,1,2,3,4,5", res1.dump())

	res2 := cellar.Search("cru")
	assert.Equal(t, "1,4,5", res2.dump())

	res3 := cellar.Search("bord")
	assert.Equal(t, "0,2,3", res3.dump())

	res4 := cellar.Search("88")
	assert.Equal(t, "0,3", res4.dump())

	res5 := cellar.Search("AUX")
	assert.Equal(t, "0,2,3", res5.dump())
}

/******************************************************************************************/
/****************************** TEST HELPER FUNCTIONS *************************************/
/******************************************************************************************/

func openWineCatalog(t *testing.T) Cellar {
	jsonFile, err := os.Open("wines.json")
	if err != nil {
		assert.Fail(t, "cannot open wines catalog json file")
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		assert.Fail(t, "cannot read wines catalog json file")
	}

	var wineCatalog map[string][]Object
	err = json.Unmarshal([]byte(byteValue), &wineCatalog)
	if err != nil {
		assert.Fail(t, "cannot unmarshal wines catalog json value")
	}

	jsonWines := wineCatalog["wines"]
	if len(jsonWines) == 0 {
		assert.Fail(t, "wines catalog is empty")
	}
	cellar := make(Cellar, 0, len(jsonWines))
	for _, w := range jsonWines {
		wine, err := NewWine(w)
		if err != nil {
			assert.Fail(t, "cannot create a Wine instance", err)
		}
		cellar = append(cellar, *wine)
	}
	return cellar
}
