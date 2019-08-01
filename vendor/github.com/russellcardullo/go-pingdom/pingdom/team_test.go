package pingdom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMTeamServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"teams": [
					{
							"id": "1",
							"name": "Team Rocket",
							"users": [
									{
											"id": "1",
											"email": "giovanni@team-rocket.org",
											"name": "Giovanni"
									}
							]
					},
					{
							"id": "2",
							"name": "The A-Team",
							"users": [
									{
											"id": "2",
											"email": "hannibal@ateam.org",
											"name": "John \"Hannibal\" Smith"
									},
									{
											"id": "3",
											"email": "faceman@ateam.org",
											"name": "Templeton \"Face(man)\" Peck"
									}
							]
					}
			]
	}`)
	})
	want := []TeamResponse{
		{
			ID:   "1",
			Name: "Team Rocket",
			Users: []TeamUserResponse{
				{
					ID:    "1",
					Email: "giovanni@team-rocket.org",
					Name:  "Giovanni",
				},
			},
		},
		{
			ID:   "2",
			Name: "The A-Team",
			Users: []TeamUserResponse{
				{
					ID:    "2",
					Email: "hannibal@ateam.org",
					Name:  "John \"Hannibal\" Smith",
				},
				{
					ID:    "3",
					Email: "faceman@ateam.org",
					Name:  "Templeton \"Face(man)\" Peck",
				},
			},
		},
	}

	teams, err := client.Teams.List()
	assert.NoError(t, err)
	assert.Equal(t, want, teams, "Teams.List() should return correct result")
}

func TestTeamServiceCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": "65",
			"name": "Operations",
			"users": [
					{
							"id": "10034512",
							"email": "hannibal@ateam.org",
							"name": "John \"Hannibal\" Smith"
					},
					{
							"id": "10043154",
							"email": "faceman@ateam.org",
							"name": "Templeton \"Face(man)\" Peck"
					}
			]
	}`)
	})

	team := TeamData{
		Name:    "Operations",
		UserIds: "10034512,10043154",
	}

	want := &TeamResponse{
		ID:   "65",
		Name: "Operations",
		Users: []TeamUserResponse{
			{
				ID:    "10034512",
				Email: "hannibal@ateam.org",
				Name:  "John \"Hannibal\" Smith",
			},
			{
				ID:    "10043154",
				Email: "faceman@ateam.org",
				Name:  "Templeton \"Face(man)\" Peck",
			},
		},
	}

	teams, err := client.Teams.Create(&team)
	assert.NoError(t, err)
	assert.Equal(t, want, teams, "Teams.Create() should return correct result")
}

func TestTeamServiceRead(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"team": {
					"id": "1",
					"name": "Team Rocket",
					"users": [
							{
									"id": "1",
									"email": "giovanni@team-rocket.org",
									"name": "Giovanni"
							}
					]
			}
	}`)
	})

	want := &TeamResponse{
		ID:   "1",
		Name: "Team Rocket",
		Users: []TeamUserResponse{
			{
				ID:    "1",
				Email: "giovanni@team-rocket.org",
				Name:  "Giovanni",
			},
		},
	}

	team, err := client.Teams.Read(1)
	assert.NoError(t, err)
	assert.Equal(t, want, team, "Teams.Read() should return correct result")
}

func TestTeamServiceUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/65", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": "65",
			"name": "Operations",
			"users": [
					{
							"id": "10034512",
							"email": "hannibal@ateam.org",
							"name": "John \"Hannibal\" Smith"
					},
					{
							"id": "10043154",
							"email": "faceman@ateam.org",
							"name": "Templeton \"Face(man)\" Peck"
					}
			]
	}`)
	})

	updateTeam := TeamData{
		Name:    "Operations",
		UserIds: "10034512,10043154",
	}

	want := &TeamResponse{
		ID:   "65",
		Name: "Operations",
		Users: []TeamUserResponse{
			{
				ID:    "10034512",
				Email: "hannibal@ateam.org",
				Name:  "John \"Hannibal\" Smith",
			},
			{
				ID:    "10043154",
				Email: "faceman@ateam.org",
				Name:  "Templeton \"Face(man)\" Peck",
			},
		},
	}

	team, err := client.Teams.Update(65, &updateTeam)
	assert.NoError(t, err)
	assert.Equal(t, want, team, "Teams.Update() should return correct result")
}

func TestTeamServiceDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"success": true
	}`)
	})
	want := &TeamDeleteResponse{Success: true}

	team, err := client.Teams.Delete(1234)
	assert.NoError(t, err)
	assert.Equal(t, want, team, "Teams.Delete() should return correct result")
}
