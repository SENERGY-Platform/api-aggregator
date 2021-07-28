package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
)

func (this *Lib) QueryPermissionsSearch(token string, query QueryMessage, result interface{}) (err error, code int) {
	requestBody := new(bytes.Buffer)
	err = json.NewEncoder(requestBody).Encode(query)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	req, err := http.NewRequest("POST", this.config.PermissionsUrl+"/v3/query", requestBody)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		err = errors.New(buf.String())
		log.Println("ERROR: ", resp.StatusCode, err)
		debug.PrintStack()
		return err, resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

type QueryMessage struct {
	Resource string         `json:"resource"`
	Find     *QueryFind     `json:"find"`
	ListIds  *QueryListIds  `json:"list_ids"`
	CheckIds *QueryCheckIds `json:"check_ids"`
}
type QueryFind struct {
	QueryListCommons
	Search string     `json:"search"`
	Filter *Selection `json:"filter"`
}

type QueryListIds struct {
	QueryListCommons
	Ids []string `json:"ids"`
}

type QueryCheckIds struct {
	Ids    []string `json:"ids"`
	Rights string   `json:"rights"`
}

type QueryListCommons struct {
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Rights   string `json:"rights"`
	SortBy   string `json:"sort_by"`
	SortDesc bool   `json:"sort_desc"`
}

type QueryOperationType string

const (
	QueryEqualOperation             QueryOperationType = "=="
	QueryUnequalOperation           QueryOperationType = "!="
	QueryAnyValueInFeatureOperation QueryOperationType = "any_value_in_feature"
)

type ConditionConfig struct {
	Feature   string             `json:"feature"`
	Operation QueryOperationType `json:"operation"`
	Value     interface{}        `json:"value"`
	Ref       string             `json:"ref"`
}

type Selection struct {
	And       []Selection     `json:"and"`
	Or        []Selection     `json:"or"`
	Not       *Selection      `json:"not"`
	Condition ConditionConfig `json:"condition"`
}
