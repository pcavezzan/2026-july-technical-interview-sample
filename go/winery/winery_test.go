package winery

import (
	"encoding/json"
	"errors"
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
	t.Parallel()

	type args struct {
		wine Object
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
		want    Wine
	}{
		{
			name: "wine must have a given price",
			args: args{
				wine: Object{
					"name":   "Château Angelus",
					"region": "Bordeaux",
					"year":   2017,
					"color":  COLOR_RED,
				},
			},
			wantErr: errors.New("wine must have a given price"),
		},
		{
			name: "wine price must be a positive floating value",
			args: args{
				wine: Object{
					"name":   "Château Angelus",
					"region": "Bordeaux",
					"year":   2017,
					"price":  -1928.,
				},
			},
			wantErr: errors.New("wine price must be a positive floating value, got (price: -1928.00)"),
		},
		{
			name: "wine properly created",
			args: args{
				wine: Object{
					"name":   "Château Angelus",
					"region": "Bordeaux",
					"year":   2017,
					"price":  1928.,
				},
			},
			want: Wine{
				Name:   "Château Angelus",
				Region: "Bordeaux",
				Year:   2017,
				Price:  1928.,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewWine(tt.args.wine)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, &tt.want, got)
			}
		})
	}
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
	t.Parallel()

	tests := []struct {
		name string
		desc bool
		want string
	}{
		{name: "ascending", desc: false, want: "5,3,1,2,0,4"},
		{name: "descending", desc: true, want: "4,0,1,2,3,5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cellar := openWineCatalog(t)
			require.Equal(t, 6, cellar.Length())

			got := cellar.SortByPrice(tt.desc)

			assert.Equal(t, tt.want, got.dump())
			// l'invariant clé, vérifié pour chaque cas
			assert.Equal(t, "0,1,2,3,4,5", cellar.dump(), "original should stay the same")
		})
	}
}

func TestSearchWines(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "should return all wines",
			query: "",
			want:  "0,1,2,3,4,5",
		},
		{
			name:  "should return cru wines",
			query: "cru",
			want:  "1,4,5",
		},
		{
			name:  "should return 88's wines",
			query: "88",
			want:  "0,3",
		},
		{
			name:  "should return AUX wines",
			query: "AUX",
			want:  "0,2,3",
		},
		{
			name:  "should return romanee result",
			query: "romanee",
			want:  "4",
		},
	}

	// We intentionnaly share cella because search should not modify the cellar.
	// So we can reuse it for all tests without having any effect on any tests.
	cellar := openWineCatalog(t)
	require.Equal(t, 6, cellar.Length())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := cellar.Search(tt.query)

			assert.Equal(t, tt.want, res.dump())
		})
	}
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
