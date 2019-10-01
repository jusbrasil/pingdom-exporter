package pingdom

import (
	"encoding/json"
	"io/ioutil"
)

// TsmService provides an interface to Pingdom transactions.
type TsmService struct {
	client *Client
}

// TODO
func (cs *TsmService) List() ([]TsmResponse, error) {
	req, err := cs.client.NewRequest("GET", "/tms.recipes", nil)
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

	p := &listTsmJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &p)

	return p.Tsm, err
}
