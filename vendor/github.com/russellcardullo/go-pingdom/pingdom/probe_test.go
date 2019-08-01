package pingdom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProbesServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/probes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"probes": [
				{
					"id": 32,
					"country": "United States",
					"city": "Los Angeles",
					"name": "Los Angeles, CA",
					"active": true,
					"hostname": "s410.pingdom.com",
					"ip": "204.152.200.42",
					"countryiso": "US",
					"ipv6": "2607:fcd0:100:8d00::410",
					"region": "NA"
				},
				{
					"id": 184,
					"country": "Brazil",
					"city": "São Paulo",
					"name": "Sao Paulo 2, Brazil",
					"active": true,
					"hostname": "s4028.pingdom.com",
					"ip": "52.67.148.55",
					"countryiso": "BR",
					"ipv6": "2600:1f1e:d7c:fd05::4028",
					"region": "LATAM"
		  	}
			]
		}`)
	})
	want := []ProbeResponse{
		{
			ID:         32,
			Country:    "United States",
			City:       "Los Angeles",
			Name:       "Los Angeles, CA",
			Active:     true,
			Hostname:   "s410.pingdom.com",
			IP:         "204.152.200.42",
			IPv6:       "2607:fcd0:100:8d00::410",
			CountryISO: "US",
			Region:     "NA",
		},
		{
			ID:         184,
			Country:    "Brazil",
			City:       "São Paulo",
			Name:       "Sao Paulo 2, Brazil",
			Active:     true,
			Hostname:   "s4028.pingdom.com",
			IP:         "52.67.148.55",
			IPv6:       "2600:1f1e:d7c:fd05::4028",
			CountryISO: "BR",
			Region:     "LATAM",
		},
	}

	params := make(map[string]string)

	probes, err := client.Probes.List(params)
	assert.NoError(t, err)
	assert.Equal(t, want, probes, "Probes.List() should return correct result")
}
