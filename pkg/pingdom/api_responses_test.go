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

func TestHasIgnoredTag(t *testing.T) {
	testCases := []struct {
		tag      CheckResponseTag
		expected bool
	}{
		{
			tag: CheckResponseTag{
				Name:  "pingdom_exporter_ignored",
				Type:  "a",
				Count: 2,
			},
			expected: true,
		},
		{
			tag: CheckResponseTag{
				Name:  "pingdom_exporter_not_ignored",
				Type:  "a",
				Count: 2,
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		response := CheckResponse{
			Tags: []CheckResponseTag{testCase.tag},
		}

		actual := response.HasIgnoreTag()
		assert.Equal(t, actual, testCase.expected)
	}
}

func TestCheckResponseUptimeSLOFromTags(t *testing.T) {
	testCases := []struct {
		tag               CheckResponseTag
		expectedUptimeSLO float64
	}{
		{
			tag: CheckResponseTag{
				Name:  "uptime_slo_99999",
				Type:  "a",
				Count: 2,
			},
			expectedUptimeSLO: 99.999,
		},
		{
			tag: CheckResponseTag{
				Name:  "uptime_slo_99995",
				Type:  "a",
				Count: 2,
			},
			expectedUptimeSLO: 99.995,
		},
		{
			tag: CheckResponseTag{
				Name:  "uptime_slo_9999",
				Type:  "a",
				Count: 2,
			},
			expectedUptimeSLO: 99.99,
		},
		{
			tag: CheckResponseTag{
				Name:  "uptime_slo_9995",
				Type:  "a",
				Count: 2,
			},
			expectedUptimeSLO: 99.95,
		},
		{
			tag: CheckResponseTag{
				Name:  "uptime_slo_999",
				Type:  "a",
				Count: 2,
			},
			expectedUptimeSLO: 99.9,
		},
		{
			tag: CheckResponseTag{
				Name:  "uptime_slo_995",
				Type:  "a",
				Count: 2,
			},
			expectedUptimeSLO: 99.5,
		},
		{
			tag: CheckResponseTag{
				Name:  "uptime_slo_99",
				Type:  "a",
				Count: 2,
			},
			expectedUptimeSLO: 99,
		},
		{
			tag: CheckResponseTag{
				Name:  "uptime_slo_95",
				Type:  "a",
				Count: 2,
			},
			expectedUptimeSLO: 95,
		},
		{
			tag: CheckResponseTag{
				Name:  "uptime_slo_9",
				Type:  "a",
				Count: 2,
			},
			expectedUptimeSLO: 9,
		},
		{
			tag: CheckResponseTag{
				Name:  "no_uptime_slo_tag",
				Type:  "a",
				Count: 2,
			},
			expectedUptimeSLO: 91,
		},
	}

	for _, testCase := range testCases {
		response := CheckResponse{
			Tags: []CheckResponseTag{testCase.tag},
		}

		uptimeSLO := response.UptimeSLOFromTags(91)
		assert.Equal(t, uptimeSLO, testCase.expectedUptimeSLO)
	}
}
