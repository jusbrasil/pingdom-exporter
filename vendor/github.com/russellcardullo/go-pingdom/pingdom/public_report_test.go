package pingdom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicReportList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/reports.public", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		    "public": [
		        {
		            "checkid": 3276510,
		            "checkname": "Foo Check 1",
		            "reporturl": "http://stats.pingdom.com/a1sdfrg20clt/3276510"
		        },
		        {
		            "checkid": 3276511,
		            "checkname": "Foo Check 2",
		            "reporturl": "http://stats.pingdom.com/a1sdfrg20clt/3276511"
		        },
		        {
		            "checkid": 3276514,
		            "checkname": "Foo Check 3",
		            "reporturl": "http://stats.pingdom.com/a1sdfrg20clt/3276514"
		        },
		        {
		            "checkid": 3276515,
		            "checkname": "Foo Check 4",
		            "reporturl": "http://stats.pingdom.com/a1sdfrg20clt/3276515"
		        }
			]
		}`)
	})
	want := []PublicReportResponse{
		{
			ID:        3276510,
			Name:      "Foo Check 1",
			ReportURL: "http://stats.pingdom.com/a1sdfrg20clt/3276510",
		},
		{
			ID:        3276511,
			Name:      "Foo Check 2",
			ReportURL: "http://stats.pingdom.com/a1sdfrg20clt/3276511",
		},
		{
			ID:        3276514,
			Name:      "Foo Check 3",
			ReportURL: "http://stats.pingdom.com/a1sdfrg20clt/3276514",
		},
		{
			ID:        3276515,
			Name:      "Foo Check 4",
			ReportURL: "http://stats.pingdom.com/a1sdfrg20clt/3276515",
		},
	}

	checks, err := client.PublicReport.List()
	assert.NoError(t, err)
	assert.Equal(t, want, checks, "PublicReport.List() should return correct result")
}

func TestPublicReportPublishCheck(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/reports.public/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"message": "Check published"}`)
	})
	want := &PingdomResponse{Message: "Check published"}

	msg, err := client.PublicReport.PublishCheck(12345)
	assert.NoError(t, err)
	assert.Equal(t, want, msg, "PublicReport.PublishCheck() should return correct result")
}

func TestPublicReportWithdrawlCheck(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/reports.public/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"message": "Check published"}`)
	})
	want := &PingdomResponse{Message: "Check published"}

	msg, err := client.PublicReport.WithdrawlCheck(12345)
	assert.NoError(t, err)
	assert.Equal(t, want, msg, "PublicReport.WithdrawlCheck() should return correct result")
}
