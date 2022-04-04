/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/auth"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
)

func (this *Lib) PermListGateways(token auth.Token, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return this.PermList(token, "hubs", right, limit, offset)
}

func (this *Lib) PermListGatewaysOrdered(token auth.Token, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermListOrdered(token, "hubs", right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermSearchGateways(token auth.Token, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return this.PermSearch(token, "hubs", query, right, limit, offset)
}

func (this *Lib) PermSearchGatewaysOrdered(token auth.Token, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermSearchOrdered(token, "hubs", query, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermDeviceIdList(token auth.Token, ids []string, right string) (result []map[string]interface{}, err error) {
	return this.PermSelectIds(token, "devices", right, ids)
}

func (this *Lib) PermDeviceIdListOrdered(token auth.Token, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermIdListOrdered(token, "devices", ids, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermIdListOrdered(token auth.Token, kind string, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		ListIds: &QueryListIds{
			QueryListCommons: QueryListCommons{
				Limit:    len(ids),
				Rights:   right,
				SortBy:   orderfeature,
				SortDesc: direction == "desc",
			},
			Ids: ids,
		},
	}, &result)
	return
}

func (this *Lib) PermListAllGateways(token auth.Token, right string) (result []map[string]interface{}, err error) {
	return this.PermListAll(token, "hubs", right)
}

func (this *Lib) PermListAll(token auth.Token, kind string, right string) (result []map[string]interface{}, err error) {
	resp, err := get(token.Token, this.config.PermissionsUrl+"/v3/resources/:"+url.PathEscape(kind)+"?limit=9999&rights="+right)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("access denied")
		return result, err
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func (this *Lib) PermList(token auth.Token, kind string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:  l,
				Offset: o,
				Rights: right,
			},
		},
	}, &result)
	return
}

func (this *Lib) PermListOrdered(token auth.Token, kind string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:    l,
				Offset:   o,
				Rights:   right,
				SortBy:   orderfeature,
				SortDesc: direction == "desc",
			},
		},
	}, &result)
	return
}

func (this *Lib) PermSearch(token auth.Token, kind string, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:  l,
				Offset: o,
				Rights: right,
			},
			Search: query,
		},
	}, &result)
	return
}

func (this *Lib) PermSearchOrdered(token auth.Token, kind string, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:    l,
				Offset:   o,
				Rights:   right,
				SortBy:   orderfeature,
				SortDesc: direction == "desc",
			},
			Search: query,
		},
	}, &result)
	return
}

func (this *Lib) PermSelectIds(token auth.Token, kind string, right string, ids []string) (result []map[string]interface{}, err error) {
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		ListIds: &QueryListIds{
			QueryListCommons: QueryListCommons{
				Limit:  len(ids),
				Rights: right,
			},
			Ids: ids,
		},
	}, &result)
	return
}

func (this *Lib) PermSelectDeviceTypesByIdRead(token auth.Token, ids []string) (result []map[string]interface{}, err error) {
	return this.PermSelectIds(token, "device-types", "r", ids)
}

func (this *Lib) PermCheckDeviceAdmin(token auth.Token, ids []string) (result map[string]bool, err error) {
	return this.PermCheck(token, "devices", ids, "a")
}

func (this *Lib) PermCheckDeviceRead(token auth.Token, ids []string) (result map[string]bool, err error) {
	return this.PermCheck(token, "devices", ids, "r")
}

func (this *Lib) PermCheck(token auth.Token, kind string, ids []string, right string) (result map[string]bool, err error) {
	result = map[string]bool{}
	err = postJson(token.Token, this.config.PermissionsUrl+"/ids/check/"+url.PathEscape(kind)+"/"+right, ids, &result)
	return
}

func QueryPermissionSearchFindAll[T any](lib *Lib, token string, query QueryMessage, sortFieldValueGetter func(e T) interface{}, idGetter func(e T) string) (result []T, err error, code int) {
	var after *ListAfter
	temp := []T{}
	limit := 9999
	for {
		if query.Find != nil {
			query.Find.QueryListCommons.Limit = limit
			query.Find.QueryListCommons.Offset = 0
			query.Find.QueryListCommons.After = after
		}
		err, code = lib.QueryPermissionsSearch(token, query, &temp)
		if err != nil {
			return result, err, code
		}
		result = append(result, temp...)
		if len(temp) < limit {
			return result, nil, http.StatusOK
		}
		if len(temp) > 0 {
			after = &ListAfter{
				SortFieldValue: sortFieldValueGetter(temp[len(temp)-1]),
				Id:             idGetter(temp[len(temp)-1]),
			}
		}
		temp = []T{}
	}
}

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
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
	After    *ListAfter `json:"after"`
	Rights   string     `json:"rights"`
	SortBy   string     `json:"sort_by"`
	SortDesc bool       `json:"sort_desc"`
}

type ListAfter struct {
	SortFieldValue interface{} `json:"sort_field_value"`
	Id             string      `json:"id"`
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
