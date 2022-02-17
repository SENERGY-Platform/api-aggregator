package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/auth"
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/model"
	"net/http"
	"runtime/debug"
	"strconv"
)

type AspectNodeQuery struct {
	Ids []string `json:"ids"`
}

func (this *Lib) GetAspectNodes(ids []string, token auth.Token) ([]model.AspectNode, error) {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(AspectNodeQuery{Ids: ids})
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	req, err := http.NewRequest("POST", this.config.IotUrl+"/query/aspect-nodes", requestBody)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	req.Header.Set("Authorization", token.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, errors.New("unexpected status code " + strconv.Itoa(resp.StatusCode))
	}

	nodes := []model.AspectNode{}
	err = json.NewDecoder(resp.Body).Decode(&nodes)
	return nodes, err
}

func (this *Lib) GetAspectNodesWithMeasuringFunction(token auth.Token) ([]model.AspectNode, error) {
	req, err := http.NewRequest("GET", this.config.IotUrl+"/aspect-nodes?function=measuring-function", nil)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	req.Header.Set("Authorization", token.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, errors.New("unexpected status code " + strconv.Itoa(resp.StatusCode))
	}

	nodes := []model.AspectNode{}
	err = json.NewDecoder(resp.Body).Decode(&nodes)
	return nodes, err
}
