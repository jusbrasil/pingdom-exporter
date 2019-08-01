package pingdom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamPutParams(t *testing.T) {
	team := TeamData{
		Name:    "fake team",
		UserIds: "1,2",
	}
	params := team.PutParams()
	want := map[string]string{
		"name":    "fake team",
		"userids": "1,2",
	}

	assert.Equal(t, want, params, "Team.PutParams() should return correct map")
}

func TestTeamPostParams(t *testing.T) {
	team := TeamData{
		Name:    "fake team",
		UserIds: "",
	}
	params := team.PostParams()
	want := map[string]string{
		"name": "fake team",
	}

	assert.Equal(t, want, params, "Team.PostParams() should return correct map")
}

func TestTeamValid(t *testing.T) {
	team := TeamData{
		Name:    "fake team",
		UserIds: "",
	}
	params := team.Valid()

	assert.Equal(t, nil, params, "Team.Valid() should return nil if valid")
}

func TestTeamNotValid(t *testing.T) {
	team := TeamData{
		Name:    "",
		UserIds: "1,3",
	}
	params := team.Valid()

	assert.NotEqual(t, nil, params, "Team.Valid() should return not nil if not valid")
}
