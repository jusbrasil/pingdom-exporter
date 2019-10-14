package pingdom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/checks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"checks": [
				{
					"hostname": "example.com",
					"id": 85975,
					"lasterrortime": 1297446423,
					"lastresponsetime": 355,
					"lasttesttime": 1300977363,
					"name": "My check 1",
					"resolution": 1,
					"status": "up",
					"type": "http",
					"tags": [
						{
							"name": "apache",
							"type": "a",
							"count": 2
						}
					],
					"responsetime_threshold": 2300
				},
				{
					"hostname": "mydomain.com",
					"id": 161748,
					"lasterrortime": 1299194968,
					"lastresponsetime": 1141,
					"lasttesttime": 1300977268,
					"name": "My check 2",
					"resolution": 5,
					"status": "up",
					"type": "ping",
					"tags": [
						{
							"name": "nginx",
							"type": "u",
							"count": 1
						}
					]
				},
				{
					"hostname": "example.net",
					"id": 208655,
					"lasterrortime": 1300527997,
					"lastresponsetime": 800,
					"lasttesttime": 1300977337,
					"name": "My check 3",
					"resolution": 1,
					"status": "down",
					"type": "http",
					"tags": [
						{
							"name": "apache",
							"type": "a",
							"count": 2
						}
					]
				}
			]
		}`)
	})

	var countA, countB float64 = 1, 2

	want := []CheckResponse{
		{
			ID:                    85975,
			Name:                  "My check 1",
			LastErrorTime:         1297446423,
			LastResponseTime:      355,
			LastTestTime:          1300977363,
			Hostname:              "example.com",
			Resolution:            1,
			Status:                "up",
			ResponseTimeThreshold: 2300,
			Type: CheckResponseType{
				Name: "http",
			},
			Tags: []CheckResponseTag{
				{
					Name:  "apache",
					Type:  "a",
					Count: countB,
				},
			},
		},
		{
			ID:               161748,
			Name:             "My check 2",
			LastErrorTime:    1299194968,
			LastResponseTime: 1141,
			LastTestTime:     1300977268,
			Hostname:         "mydomain.com",
			Resolution:       5,
			Status:           "up",
			Type: CheckResponseType{
				Name: "ping",
			},
			Tags: []CheckResponseTag{
				{
					Name:  "nginx",
					Type:  "u",
					Count: countA,
				},
			},
		},
		{
			ID:               208655,
			Name:             "My check 3",
			LastErrorTime:    1300527997,
			LastResponseTime: 800,
			LastTestTime:     1300977337,
			Hostname:         "example.net",
			Resolution:       1,
			Status:           "down",
			Type: CheckResponseType{
				Name: "http",
			},
			Tags: []CheckResponseTag{
				{
					Name:  "apache",
					Type:  "a",
					Count: countB,
				},
			},
		},
	}

	checks, err := client.Checks.List()
	assert.NoError(t, err)
	assert.Equal(t, want, checks)
}
