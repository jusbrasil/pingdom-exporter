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

func TestCheckServiceCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/checks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"check":{
				"id":138631,
				"name":"My new HTTP check"
			}
		}`)
	})

	newCheck := HttpCheck{
		Name:           "My new HTTP check",
		Hostname:       "example.com",
		Resolution:     5,
		IntegrationIds: []int{33333333, 44444444},
	}
	want := &CheckResponse{ID: 138631, Name: "My new HTTP check"}

	check, err := client.Checks.Create(&newCheck)
	assert.NoError(t, err)
	assert.Equal(t, want, check)
}

func TestCheckServiceRead(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/checks/85975", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"check" : {
        "created" : 1240394682,
        "hostname" : "s7.mydomain.com",
        "id" : 85975,
        "integrationids": [
            33333333,
            44444444
        ],
        "ipv6": false,
        "lasterrortime" : 1293143467,
        "lasttesttime" : 1294064823,
        "name" : "My check 7",
        "notifyagainevery" : 0,
        "notifywhenbackup" : false,
        "probe_filters": [],
        "resolution" : 1,
        "sendnotificationwhendown" : 0,
        "responsetime_threshold": 2300,
        "status" : "up",
        "tags": [],
        "teams": [
            {
                "id": 123456,
                "name": "Oncall"
            }
        ],
        "type" : {
          "http" : {
            "encryption": false,
            "port" : 80,
            "requestheaders" : {
              "User-Agent" : "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)"
            },
            "url" : "/"
          }
        }
			}
		}`)
	})

	want := &CheckResponse{
		ID:                       85975,
		Name:                     "My check 7",
		Resolution:               1,
		SendNotificationWhenDown: 0,
		NotifyAgainEvery:         0,
		NotifyWhenBackup:         false,
		Created:                  1240394682,
		Hostname:                 "s7.mydomain.com",
		Status:                   "up",
		LastErrorTime:            1293143467,
		LastTestTime:             1294064823,
		ResponseTimeThreshold:    2300,
		Teams: []CheckTeamResponse{
			{
				Name: "Oncall",
				ID:   123456,
			},
		},
		TeamIds: []int{123456},
		Type: CheckResponseType{
			Name: "http",
			HTTP: &CheckResponseHTTPDetails{
				Url:              "/",
				Encryption:       false,
				Port:             80,
				Username:         "",
				Password:         "",
				ShouldContain:    "",
				ShouldNotContain: "",
				PostData:         "",
				RequestHeaders: map[string]string{
					"User-Agent": "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
				},
			},
		},
		IntegrationIds: []int{33333333, 44444444},
		Tags:           []CheckResponseTag{},
		ProbeFilters:   []string{},
	}

	check, err := client.Checks.Read(85975)
	assert.NoError(t, err)
	assert.Equal(t, want, check)
}

func TestCheckServiceUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/checks/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"message":"Modification of check was successful!"}`)
	})

	updateCheck := HttpCheck{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
	want := &PingdomResponse{Message: "Modification of check was successful!"}

	msg, err := client.Checks.Update(12345, &updateCheck)
	assert.NoError(t, err)
	assert.Equal(t, want, msg)
}

func TestCheckServiceDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/checks/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"message":"Deletion of check was successful!"}`)
	})

	want := &PingdomResponse{Message: "Deletion of check was successful!"}

	msg, err := client.Checks.Delete(12345)
	assert.NoError(t, err)
	assert.Equal(t, want, msg)
}

func TestCheckServiceSummaryPerformance(t *testing.T) {
	id := 1337
	t.Run("passes on error from API", func(t *testing.T) {
		setup()
		defer teardown()

		errorMsg := `{"error":{"statuscode":401,"statusdesc":"Unauthorized","errormessage":"Invalid email and\/or password"}}`
		request := SummaryPerformanceRequest{
			Id: id,
		}

		mux.HandleFunc(fmt.Sprintf("/summary.performance/%v", id), func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(401)
			fmt.Fprint(w, errorMsg)
		})

		_, err := client.Checks.SummaryPerformance(request)

		assert.Equal(t, &PingdomError{
			StatusCode: 401,
			StatusDesc: "Unauthorized",
			Message:    "Invalid email and/or password",
		}, err)
	})

	t.Run("passes on response as datastructure", func(t *testing.T) {
		setup()
		defer teardown()

		request := SummaryPerformanceRequest{
			Id: id,
		}

		expectedResponse := SummaryPerformanceResponse{
			Summary: SummaryPerformanceMap{
				Hours: []SummaryPerformanceSummary{
					{
						AvgResponse: 222,
						Downtime:    0,
						StartTime:   1536926400,
						Unmonitored: 0,
						Uptime:      3600,
					},
					{
						AvgResponse: 225,
						Downtime:    0,
						StartTime:   1536930000,
						Unmonitored: 0,
						Uptime:      3442,
					},
				},
			},
		}

		mux.HandleFunc(fmt.Sprintf("/summary.performance/%v", id), func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{
	"summary": {
		"hours": [
			{
				"avgresponse": 222,
				"downtime": 0,
        		"starttime": 1536926400,
        		"unmonitored": 0,
        		"uptime": 3600
			},
      		{
	        	"avgresponse": 225,
	        	"downtime": 0,
	        	"starttime": 1536930000,
	        	"unmonitored": 0,
	        	"uptime": 3442
	      	}
		]
	}
}`)
		})

		resp, err := client.Checks.SummaryPerformance(request)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, *resp)
	})
}

func TestCheckServiceResults(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/results/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
    "activeprobes": [
        259,
        255,
        93,
        94,
        87
    ],
    "results": [
        {
            "probeid": 259,
            "time": 1563370611,
            "status": "up",
            "responsetime": 145,
            "statusdesc": "OK",
            "statusdesclong": "OK"
        },
        {
            "probeid": 87,
            "time": 1563370551,
            "status": "up",
            "responsetime": 56,
            "statusdesc": "OK",
            "statusdesclong": "OK"
        },
        {
            "probeid": 93,
            "time": 1563370491,
            "status": "up",
            "responsetime": 962,
            "statusdesc": "OK",
            "statusdesclong": "OK"
        },
        {
            "probeid": 255,
            "time": 1563370431,
            "status": "up",
            "responsetime": 395,
            "statusdesc": "OK",
            "statusdesclong": "OK"
        },
        {
            "probeid": 94,
            "time": 1563370371,
            "status": "up",
            "responsetime": 1084,
            "statusdesc": "OK",
            "statusdesclong": "OK"
        }
    ]
}`)
	})

	want := &ResultsResponse{
		ActiveProbes: []int{259, 255, 93, 94, 87},
		Results: []Result{
			{ProbeID: 259, Time: 1563370611, Status: "up", ResponseTime: 145, StatusDesc: "OK", StatusDescLong: "OK"},
			{ProbeID: 87, Time: 1563370551, Status: "up", ResponseTime: 56, StatusDesc: "OK", StatusDescLong: "OK"},
			{ProbeID: 93, Time: 1563370491, Status: "up", ResponseTime: 962, StatusDesc: "OK", StatusDescLong: "OK"},
			{ProbeID: 255, Time: 1563370431, Status: "up", ResponseTime: 395, StatusDesc: "OK", StatusDescLong: "OK"},
			{ProbeID: 94, Time: 1563370371, Status: "up", ResponseTime: 1084, StatusDesc: "OK", StatusDescLong: "OK"},
		},
	}

	results, err := client.Checks.Results(12345)
	assert.NoError(t, err)
	assert.Equal(t, want, results)
}
