package pingdom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// OutageService provides an interface to Pingdom outage summary.
type OutageSummaryService struct {
	client *Client
}

// List returns a list of outage summaries from Pingdom.
func (os *OutageSummaryService) List(checkId int, params ...map[string]string) ([]OutageSummaryResponseState, error) {
	param := map[string]string{}
	if len(params) == 1 {
		param = params[0]
	}

	req, err := os.client.NewRequest("GET", fmt.Sprintf("/summary.outage/%d", checkId), param)
	if err != nil {
		return nil, err
	}

	resp, err := os.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	m := &listOutageSummaryJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)

	return m.Summary.States, err
}
