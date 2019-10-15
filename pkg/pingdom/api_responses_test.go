package pingdom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorError(t *testing.T) {
	errorResponse := Error{200, "OK", "Message"}
	assert.Equal(t, "200 OK: Message", errorResponse.Error())
}

func TestCheckResponseTagsString(t *testing.T) {
	checkResponse := CheckResponse{
		Tags: []CheckResponseTag{
			{
				Name:  "apache",
				Type:  "a",
				Count: 2,
			},
			{
				Name:  "server",
				Type:  "a",
				Count: 2,
			},
		},
	}
	assert.Equal(t, "apache,server", checkResponse.TagsString())
}
