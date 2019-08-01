package pingdom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpCheckPutParams(t *testing.T) {
	check := HttpCheck{
		Name:     "fake check",
		Hostname: "example.com",
		Url:      "/foo",
		RequestHeaders: map[string]string{
			"User-Agent": "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
			"Pragma":     "no-cache",
		},
		Username:              "user",
		Password:              "pass",
		IntegrationIds:        []int{33333333, 44444444},
		UserIds:               []int{123, 456},
		TeamIds:               []int{789},
		ResponseTimeThreshold: 2300,
	}
	want := map[string]string{
		"name":                   "fake check",
		"host":                   "example.com",
		"paused":                 "false",
		"resolution":             "0",
		"notifyagainevery":       "0",
		"notifywhenbackup":       "false",
		"url":                    "/foo",
		"requestheader0":         "Pragma:no-cache",
		"requestheader1":         "User-Agent:Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
		"auth":                   "user:pass",
		"encryption":             "false",
		"shouldnotcontain":       "",
		"postdata":               "",
		"integrationids":         "33333333,44444444",
		"tags":                   "",
		"probe_filters":          "",
		"userids":                "123,456",
		"teamids":                "789",
		"responsetime_threshold": "2300",
	}

	params := check.PutParams()
	assert.Equal(t, want, params)
}

func TestHttpCheckPostParams(t *testing.T) {
	check := HttpCheck{
		Name:     "fake check",
		Hostname: "example.com",
		Url:      "/foo",
		RequestHeaders: map[string]string{
			"User-Agent": "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
			"Pragma":     "no-cache",
		},
		Username:              "user",
		Password:              "pass",
		IntegrationIds:        []int{33333333, 44444444},
		UserIds:               []int{123, 456},
		TeamIds:               []int{789},
		ResponseTimeThreshold: 2300,
	}
	want := map[string]string{
		"name":                   "fake check",
		"host":                   "example.com",
		"paused":                 "false",
		"resolution":             "0",
		"notifyagainevery":       "0",
		"notifywhenbackup":       "false",
		"type":                   "http",
		"url":                    "/foo",
		"requestheader0":         "Pragma:no-cache",
		"requestheader1":         "User-Agent:Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
		"auth":                   "user:pass",
		"encryption":             "false",
		"integrationids":         "33333333,44444444",
		"userids":                "123,456",
		"teamids":                "789",
		"responsetime_threshold": "2300",
	}

	params := check.PostParams()
	assert.Equal(t, want, params)
}

func TestHttpCheckValid(t *testing.T) {
	check := HttpCheck{Name: "fake check", Hostname: "example.com", Resolution: 15}
	assert.NoError(t, check.Valid())

	badCheck := HttpCheck{Name: "fake check", Hostname: "example.com"}
	assert.Error(t, badCheck.Valid())

	badContainsCheck := HttpCheck{
		Name:             "fake check",
		Hostname:         "example.com",
		Resolution:       15,
		ShouldContain:    "foo",
		ShouldNotContain: "bar",
	}
	assert.Error(t, badContainsCheck.Valid())
}

func TestPingCheckPostParams(t *testing.T) {
	check := PingCheck{
		Name:                  "fake check",
		Hostname:              "example.com",
		IntegrationIds:        []int{33333333, 44444444},
		UserIds:               []int{123, 456},
		TeamIds:               []int{789},
		ResponseTimeThreshold: 2300,
	}
	want := map[string]string{
		"name":                   "fake check",
		"host":                   "example.com",
		"paused":                 "false",
		"resolution":             "0",
		"notifyagainevery":       "0",
		"notifywhenbackup":       "false",
		"type":                   "ping",
		"integrationids":         "33333333,44444444",
		"userids":                "123,456",
		"teamids":                "789",
		"responsetime_threshold": "2300",
	}

	params := check.PostParams()
	assert.Equal(t, want, params)
}

func TestPingCheckPutParams(t *testing.T) {
	check := PingCheck{
		Name:           "fake check",
		Hostname:       "example.com",
		IntegrationIds: []int{33333333, 44444444},
		UserIds:        []int{123, 456},
		TeamIds:        []int{789},
	}
	want := map[string]string{
		"name":             "fake check",
		"host":             "example.com",
		"resolution":       "0",
		"paused":           "false",
		"notifyagainevery": "0",
		"notifywhenbackup": "false",
		"integrationids":   "33333333,44444444",
		"probe_filters":    "",
		"userids":          "123,456",
		"teamids":          "789",
	}

	params := check.PutParams()
	assert.Equal(t, want, params)
}

func TestPingCheckValid(t *testing.T) {
	check := PingCheck{Name: "fake check", Hostname: "example.com", Resolution: 15}
	assert.NoError(t, check.Valid())

	badCheck := PingCheck{Name: "fake check", Hostname: "example.com"}
	assert.Error(t, badCheck.Valid())
}

func TestTCPCheckPostParams(t *testing.T) {
	check := TCPCheck{
		Name:           "fake check",
		Hostname:       "example.com",
		IntegrationIds: []int{33333333, 44444444},
		UserIds:        []int{123, 456},
		TeamIds:        []int{789},
		Port:           8080,
		StringToSend:   "Hello World",
		StringToExpect: "Hi there",
	}
	want := map[string]string{
		"name":             "fake check",
		"host":             "example.com",
		"paused":           "false",
		"resolution":       "0",
		"notifyagainevery": "0",
		"notifywhenbackup": "false",
		"type":             "tcp",
		"integrationids":   "33333333,44444444",
		"userids":          "123,456",
		"teamids":          "789",
		"port":             "8080",
		"stringtosend":     "Hello World",
		"stringtoexpect":   "Hi there",
	}

	params := check.PostParams()
	assert.Equal(t, want, params)
}

func TestTCPCheckValid(t *testing.T) {
	check := TCPCheck{Name: "fake check", Hostname: "example.com", Resolution: 15, Port: 8080}
	assert.NoError(t, check.Valid())

	badCheck := TCPCheck{Name: "fake check", Hostname: "example.com", Resolution: 15}
	assert.Error(t, badCheck.Valid())
}

func TestSummaryPerformanceRequestValid(t *testing.T) {
	t.Run("missing field 'id'", func(t *testing.T) {
		assert.Equal(t, ErrMissingId, SummaryPerformanceRequest{}.Valid())
	})

	t.Run("resolution", func(t *testing.T) {
		assert.Nil(t, SummaryPerformanceRequest{
			Id:         123,
			Resolution: "hour",
		}.Valid())
		assert.Nil(t, SummaryPerformanceRequest{
			Id:         123,
			Resolution: "day",
		}.Valid())
		assert.Nil(t, SummaryPerformanceRequest{
			Id:         123,
			Resolution: "week",
		}.Valid())
		assert.Equal(t, ErrBadResolution, SummaryPerformanceRequest{
			Id:         123,
			Resolution: "month",
		}.Valid())

	})
}

func TestSummaryPerformanceRequestGetParams(t *testing.T) {
	id := 1337
	t.Run("empty request", func(t *testing.T) {
		want := map[string]string{}

		params := SummaryPerformanceRequest{
			Id: id,
		}.GetParams()

		assert.Equal(t, want, params)
	})

	t.Run("with some params", func(t *testing.T) {
		want := map[string]string{
			"resolution":    "week",
			"includeuptime": "true",
		}

		params := SummaryPerformanceRequest{
			Id:            id,
			IncludeUptime: true,
			Resolution:    "week",
		}.GetParams()

		assert.Equal(t, want, params)
	})
}
