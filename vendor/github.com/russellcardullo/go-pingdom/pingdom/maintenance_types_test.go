package pingdom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaintenancePutParams(t *testing.T) {
	maintenance := MaintenanceWindow{
		Description:    "fake maintenance",
		From:           1,
		To:             1524040922,
		RecurrenceType: "none",
		RepeatEvery:    0,
		EffectiveTo:    1,
		UptimeIDs:      "12345,67890",
		TmsIDs:         "09876,54321",
	}
	params := maintenance.PutParams()
	want := map[string]string{
		"description":    "fake maintenance",
		"from":           "1",
		"to":             "1524040922",
		"recurrencetype": "none",
		"effectiveto":    "1",
		"uptimeids":      "12345,67890",
		"tmsids":         "09876,54321",
	}

	assert.Equal(t, want, params, "Maintenance.PutParams() should return correct map")
}

func TestMaintenancePostParams(t *testing.T) {
	maintenance := MaintenanceWindow{
		Description:    "fake maintenance",
		From:           1,
		To:             1524040922,
		RecurrenceType: "",
		UptimeIDs:      "",
		TmsIDs:         "",
	}
	params := maintenance.PostParams()
	want := map[string]string{
		"description": "fake maintenance",
		"from":        "1",
		"to":          "1524040922",
	}

	assert.Equal(t, want, params, "Maintenance.PostParams() should return correct map")
}

func TestMaintenanceValid(t *testing.T) {
	maintenance := MaintenanceWindow{
		Description: "fake maintenance",
		From:        1,
		To:          1524040922,
	}
	params := maintenance.Valid()

	assert.Equal(t, nil, params, "Maintenance.Valid() should return nil if valid")
}

func TestMaintenanceNotValid(t *testing.T) {
	maintenance := MaintenanceWindow{
		Description: "fake maintenance",
		From:        1,
	}
	params := maintenance.Valid()

	assert.NotEqual(t, nil, params, "Maintenance.Valid() should return not nil if not valid")
}
