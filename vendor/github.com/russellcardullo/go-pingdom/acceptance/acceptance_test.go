package acceptance

import (
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/russellcardullo/go-pingdom/pingdom"
	"github.com/stretchr/testify/assert"
)

var client *pingdom.Client

var runAcceptance bool

func init() {
	if os.Getenv("PINGDOM_ACCEPTANCE") == "1" {
		runAcceptance = true

		config := pingdom.ClientConfig{
			User:     os.Getenv("PINGDOM_USER"),
			Password: os.Getenv("PINGDOM_PASSWORD"),
			APIKey:   os.Getenv("PINGDOM_API_KEY"),
			HTTPClient: &http.Client{
				Timeout: time.Second * 10,
			},
		}
		client, _ = pingdom.NewClientWithConfig(config)
	}
}

func TestListChecks(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}
	checks, err := client.Checks.List()
	assert.NoError(t, err)
	assert.NotNil(t, checks)
}

func TestHTTPCheck(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}
	newCheck := pingdom.HttpCheck{
		Name:                     "Test Check",
		Hostname:                 "example.com",
		Resolution:               5,
		SendNotificationWhenDown: 100,
		Tags:                     "tag",
	}
	check, err := client.Checks.Create(&newCheck)
	assert.NoError(t, err)
	assert.NotNil(t, check)

	newCheck.Name = "Test Check 2"
	up, err := client.Checks.Update(check.ID, &newCheck)
	assert.NoError(t, err)
	assert.NotNil(t, up)

	resp, err := client.Checks.Read(check.ID)
	assert.NoError(t, err)
	assert.Equal(t, newCheck.Name, resp.Name)

	delMsg, err := client.Checks.Delete(check.ID)
	assert.NoError(t, err)
	assert.NotNil(t, delMsg)
}

func TestTagSupport(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}
	newCheck := pingdom.HttpCheck{
		Name:                     "Test Check",
		Hostname:                 "example.com",
		Resolution:               5,
		SendNotificationWhenDown: 100,
		Tags:                     "tag",
	}
	check, err := client.Checks.Create(&newCheck)
	assert.NoError(t, err)
	assert.NotNil(t, check)

	params := make(map[string]string)
	params["include_tags"] = "true"
	params["tags"] = "tag"

	checks, err := client.Checks.List(params)
	assert.NoError(t, err)
	assert.NotNil(t, checks)
	assert.Equal(t, 1, len(checks))

	delMsg, err := client.Checks.Delete(check.ID)
	assert.NoError(t, err)
	assert.NotNil(t, delMsg)
}

func TestProbes(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}
	params := make(map[string]string)

	probes, err := client.Probes.List(params)
	assert.NoError(t, err)
	assert.NotNil(t, probes)
	assert.NotEmpty(t, probes)
}

func TestTeams(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}
	teams, err := client.Teams.List()
	assert.NoError(t, err)
	assert.NotNil(t, teams)

	td := pingdom.TeamData{
		Name: "Team",
	}
	team, err := client.Teams.Create(&td)
	assert.NoError(t, err)
	assert.NotNil(t, team)

	id, _ := strconv.Atoi(team.ID)
	teamr, err := client.Teams.Read(id)
	assert.NoError(t, err)
	assert.NotNil(t, teamr)

	success, err := client.Teams.Delete(id)
	assert.NoError(t, err)
	assert.NotNil(t, success)
}
