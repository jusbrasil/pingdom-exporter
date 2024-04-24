package pingdom

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
)

var (
	reqLimitHeaderKeys = []string{
		"req-limit-short",
		"req-limit-long",
	}
	reqLimitRe = regexp.MustCompile(`Remaining: (\d+) Time until reset: (\d+)`)
)

// CheckService provides an interface to Pingdom checks.
type CheckService struct {
	client *Client
}

// List returns a list of checks from Pingdom.
// This returns type CheckResponse rather than Check since the
// Pingdom API does not return a complete representation of a check.
func (cs *CheckService) List(params ...map[string]string) ([]CheckResponse, float64, error) {
	param := map[string]string{}
	if len(params) == 1 {
		param = params[0]
	}
	req, err := cs.client.NewRequest("GET", "/checks", param)
	if err != nil {
		return nil, 0, err
	}

	resp, err := cs.client.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	minRequestLimit := minRequestLimitFromHeader(resp.Header)

	if err := validateResponse(resp); err != nil {
		return nil, minRequestLimit, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	m := &listChecksJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)

	return m.Checks, minRequestLimit, err
}

func minRequestLimitFromHeader(header http.Header) float64 {
	minRequestLimit := math.MaxFloat64

	for _, key := range reqLimitHeaderKeys {
		matches := reqLimitRe.FindStringSubmatch(header.Get(key))
		if len(matches) > 0 {
			limit, err := strconv.ParseFloat(matches[1], 64)
			if err == nil && limit < minRequestLimit {
				minRequestLimit = limit
			}
		}
	}

	return minRequestLimit
}
