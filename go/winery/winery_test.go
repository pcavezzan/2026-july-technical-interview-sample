package winery

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.EqualError(t, err, "wine must have a given price")

	var w2 = Object{
		"name":   "Château Angelus",
		"region": "Bordeaux",
		"year":   2017,
		"price":  -1928.,
	}
	_, err = NewWine(w2)
	assert.EqualError(t, err, "wine price must be a positive floating value, got (price: -1928.00)")

	var wOK = Object{
		"name":   "Château Angelus",
		"region": "Bordeaux",
		"year":   2017,
		"price":  1928.,
	}
	w, err := NewWine(wOK)
	assert.NoError(t, err)
	assert.Equal(t, w.Year, 2017, "wine properly created")
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

/***********************************************************************************************/
/****************************** BENCHMARKS TESTS FUNCTIONS *************************************/
/***********************************************************************************************/

func sampleCellar(n int) Cellar {
	colors := []string{COLOR_RED, COLOR_WHITE, COLOR_ROSE}
	c := make(Cellar, 0, n)
	for i := 0; i < n; i++ {
		c = append(c, Wine{
			Id:     int64(i),
			Name:   "Château Something Long Enough To Matter",
			Price:  float64(i) * 12.5,
			Year:   1990 + i%30,
			Color:  colors[i%len(colors)],
			Region: "Bordeaux",
		})
	}
	return c
}

// A - idiomatic range, value bound to a name
func classifyIdiomatic(c Cellar) map[string]Cellar {
	res := make(map[string]Cellar)
	for _, wine := range c {
		res[wine.Color] = append(res[wine.Color], wine)
	}
	return res
}

// B - what was written live during the interview
func classifyIndexThenCopy(c Cellar) map[string]Cellar {
	res := make(map[string]Cellar)
	for i, _ := range c {
		wine := c[i]
		res[wine.Color] = append(res[wine.Color], wine)
	}
	return res
}

// C - index only, never bind the value to a local name
func classifyIndexOnly(c Cellar) map[string]Cellar {
	res := make(map[string]Cellar)
	for i := range c {
		res[c[i].Color] = append(res[c[i].Color], c[i])
	}
	return res
}

func BenchmarkClassifyIdiomatic(b *testing.B) {
	c := sampleCellar(1000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = classifyIdiomatic(c)
	}
}

func BenchmarkClassifyIndexThenCopy(b *testing.B) {
	c := sampleCellar(1000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = classifyIndexThenCopy(c)
	}
}

func BenchmarkClassifyIndexOnly(b *testing.B) {
	c := sampleCellar(1000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = classifyIndexOnly(c)
	}
}
