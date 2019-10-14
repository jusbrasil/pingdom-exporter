package pingdom

import (
	"encoding/json"
	"io/ioutil"
)

// CheckService provides an interface to Pingdom checks.
type CheckService struct {
	client *Client
}

// List returns a list of checks from Pingdom.
// This returns type CheckResponse rather than Check since the
// Pingdom API does not return a complete representation of a check.
func (cs *CheckService) List(params ...map[string]string) ([]CheckResponse, error) {
	param := map[string]string{}
	if len(params) == 1 {
		param = params[0]
	}
	req, err := cs.client.NewRequest("GET", "/checks", param)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	m := &listChecksJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)

	return m.Checks, err
}
