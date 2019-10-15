package pingdom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutageSummaryServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/summary.outage/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"summary": {
                             "states": [
                                 {
                                     "status": "up",
                                     "timefrom": 1293143523,
                                     "timeto": 1294180263
                                 },
                                 {
                                     "status": "down",
                                     "timefrom": 1294180263,
                                     "timeto": 1294180323
                                 }
                             ]
                         }
		}`)
	})

	want := []OutageSummaryResponseState{
		{
			Status:   "up",
			FromTime: 1293143523,
			ToTime:   1294180263,
		},
		{
			Status:   "down",
			FromTime: 1294180263,
			ToTime:   1294180323,
		},
	}

	checks, err := client.OutageSummary.List(1, map[string]string{
		"from": "1293143523",
	})

	assert.NoError(t, err)
	assert.Equal(t, want, checks)
}
