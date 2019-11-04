package pingdom

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/common/log"
)

// Uptime SLO tag format.
var uptimeSLORegexp = regexp.MustCompile(`^uptime_slo_(?P<SLO>\d+)$`)

// Response represents a general response from the Pingdom API.
type Response struct {
	Message string `json:"message"`
}

// Error represents an error response from the Pingdom API.
type Error struct {
	StatusCode int    `json:"statuscode"`
	StatusDesc string `json:"statusdesc"`
	Message    string `json:"errormessage"`
}

// OutageSummaryResponse represents the JSON response for a outage summary list from the Pingdom API.
type OutageSummaryResponse struct {
	States []OutageSummaryResponseState `json:"states"`
}

// OutageSummaryResponseState represents the JSON response for each outage summary.
type OutageSummaryResponseState struct {
	Status   string `json:"status"`
	FromTime int64  `json:"timefrom"`
	ToTime   int64  `json:"timeto"`
}

// CheckResponse represents the JSON response for a check from the Pingdom API.
type CheckResponse struct {
	ID                       int                 `json:"id"`
	Name                     string              `json:"name"`
	Resolution               int                 `json:"resolution,omitempty"`
	SendNotificationWhenDown int                 `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int                 `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool                `json:"notifywhenbackup,omitempty"`
	Created                  int64               `json:"created,omitempty"`
	Hostname                 string              `json:"hostname,omitempty"`
	Status                   string              `json:"status,omitempty"`
	LastErrorTime            int64               `json:"lasterrortime,omitempty"`
	LastTestTime             int64               `json:"lasttesttime,omitempty"`
	LastResponseTime         int64               `json:"lastresponsetime,omitempty"`
	IntegrationIds           []int               `json:"integrationids,omitempty"`
	SeverityLevel            string              `json:"severity_level,omitempty"`
	Type                     CheckResponseType   `json:"type,omitempty"`
	Tags                     []CheckResponseTag  `json:"tags,omitempty"`
	UserIds                  []int               `json:"userids,omitempty"`
	Teams                    []CheckTeamResponse `json:"teams,omitempty"`
	ResponseTimeThreshold    int                 `json:"responsetime_threshold,omitempty"`
	ProbeFilters             []string            `json:"probe_filters,omitempty"`
}

// CheckTeamResponse holds the team names for each check.
type CheckTeamResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CheckResponseType is the type of the Pingdom check.
type CheckResponseType struct {
	Name string                    `json:"-"`
	HTTP *CheckResponseHTTPDetails `json:"http,omitempty"`
	TCP  *CheckResponseTCPDetails  `json:"tcp,omitempty"`
}

// CheckResponseTag is an optional tag that can be added to checks.
type CheckResponseTag struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Count interface{} `json:"count"`
}

// SummaryPerformanceResponse represents the JSON response for a summary performance from the Pingdom API.
type SummaryPerformanceResponse struct {
	Summary SummaryPerformanceMap `json:"summary"`
}

// SummaryPerformanceMap is the performance broken down over different time intervals.
type SummaryPerformanceMap struct {
	Hours []SummaryPerformanceSummary `json:"hours,omitempty"`
	Days  []SummaryPerformanceSummary `json:"days,omitempty"`
	Weeks []SummaryPerformanceSummary `json:"weeks,omitempty"`
}

// SummaryPerformanceSummary is the metrics for a performance summary.
type SummaryPerformanceSummary struct {
	AvgResponse int `json:"avgresponse"`
	Downtime    int `json:"downtime"`
	StartTime   int `json:"starttime"`
	Unmonitored int `json:"unmonitored"`
	Uptime      int `json:"uptime"`
}

// UnmarshalJSON converts a byte array into a CheckResponseType.
func (c *CheckResponseType) UnmarshalJSON(b []byte) error {
	var raw interface{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		c.Name = v
	case map[string]interface{}:
		if len(v) != 1 {
			return fmt.Errorf("Check detailed response `check.type` contains more than one object: %+v", v)
		}
		for k := range v {
			c.Name = k
		}

		// Allow continue use json.Unmarshall using a type != Unmarshaller
		// This avoid enter in a infinite loop
		type t CheckResponseType
		var rawCheckDetails t

		err := json.Unmarshal(b, &rawCheckDetails)
		if err != nil {
			return err
		}
		c.HTTP = rawCheckDetails.HTTP
		c.TCP = rawCheckDetails.TCP
	}
	return nil
}

// CheckResponseHTTPDetails represents the details specific to HTTP checks.
type CheckResponseHTTPDetails struct {
	URL              string            `json:"url,omitempty"`
	Encryption       bool              `json:"encryption,omitempty"`
	Port             int               `json:"port,omitempty"`
	Username         string            `json:"username,omitempty"`
	Password         string            `json:"password,omitempty"`
	ShouldContain    string            `json:"shouldcontain,omitempty"`
	ShouldNotContain string            `json:"shouldnotcontain,omitempty"`
	PostData         string            `json:"postdata,omitempty"`
	RequestHeaders   map[string]string `json:"requestheaders,omitempty"`
}

// CheckResponseTCPDetails represents the details specific to TCP checks.
type CheckResponseTCPDetails struct {
	Port           int    `json:"port,omitempty"`
	StringToSend   string `json:"stringtosend,omitempty"`
	StringToExpect string `json:"stringtoexpect,omitempty"`
}

// Return string representation of  Error.
func (r *Error) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

// TagsString returns the check tags as a comma-separated string.
func (cr *CheckResponse) TagsString() string {
	var tagsRaw []string
	for _, tag := range cr.Tags {
		tagsRaw = append(tagsRaw, tag.Name)
	}
	return strings.Join(tagsRaw, ",")
}

// UptimeSLOFromTag returns the uptime SLO configured to this check via a tag,
// i.e. "uptime_slo_999" for 99.9 uptime SLO. Returns the argument as the
// default uptime SLO in case no uptime SLO tag exists for this check.
func (cr *CheckResponse) UptimeSLOFromTags(defaultUptimeSLO float64) float64 {
	for _, tag := range cr.Tags {
		matches := uptimeSLORegexp.FindStringSubmatch(tag.Name)

		if len(matches) > 0 {
			n, err := strconv.ParseFloat(matches[1], 64)

			if err != nil {
				log.Errorf("Error parsing uptime SLO tag %s: %v", matches[1], err)
				break
			}

			return n / math.Pow(10, math.Max(0, float64(len(matches[1])-2)))
		}
	}

	return defaultUptimeSLO
}

// private types used to unmarshall JSON responses from Pingdom.

type listChecksJSONResponse struct {
	Checks []CheckResponse `json:"checks"`
}

type listOutageSummaryJSONResponse struct {
	Summary OutageSummaryResponse `json:"summary"`
}

type errorJSONResponse struct {
	Error *Error `json:"error"`
}
