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
	"github.com/SENERGY-Platform/api-aggregator/pkg/auth"
	"github.com/SENERGY-Platform/api-aggregator/pkg/model"
	"github.com/SENERGY-Platform/permission-search/lib/client"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
)

type Interface interface {
	Config() Config
	SortByName(input []map[string]interface{}, sortAsc bool) (output []map[string]interface{})
	FilterDevicesByState(token auth.Token, devices []map[string]interface{}, state string) (result []map[string]interface{}, err error)
	CompleteDevices(token auth.Token, ids []string) (result []map[string]interface{}, err error)
	CompleteDevicesOrdered(token auth.Token, ids []string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error)
	GetGatewaysHistory(token auth.Token, duration string) (result []map[string]interface{}, err error)
	ListGateways(token auth.Token, limit string, offset string) (result []map[string]interface{}, err error)
	SearchGateways(token auth.Token, query string, limit string, offset string) (result []map[string]interface{}, err error)
	ListGatewaysOrdered(token auth.Token, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error)
	SearchGatewaysOrdered(token auth.Token, query string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error)
	GetExtendedProcessList(token auth.Token, query url.Values) (result []map[string]interface{}, err error)
	CompleteDeviceHistory(token auth.Token, duration string, devices []map[string]interface{}) (result []map[string]interface{}, err error)
	CompleteGatewayHistory(token auth.Token, duration string, devices []map[string]interface{}) (result []map[string]interface{}, err error)
	ListAllGateways(token auth.Token) (result []map[string]interface{}, err error)
	GetGatewayDevices(token auth.Token, id string) (ids []string, err error)
	GetDeviceTypeDevices(token auth.Token, id string, limit string, offset string, orderFeature string, direction string) (ids []string, err error)
	FindDevices(token auth.Token, search string, list []string, limit int, offset int, orderfeature string, direction string, location string, state string) ([]map[string]interface{}, error)
	FindDevicesAfter(token auth.Token, search string, list []string, limit int, afterId string, orderfeature string, direction string, location string, state string) ([]map[string]interface{}, error)
	GetMeasuringFunctionsForAspect(token auth.Token, aspectId string) (functions []Function, err error, code int)
	GetMeasuringFunctions(token auth.Token, functionIds []string) (functions []Function, err error, code int)
	GetImportTypesWithAspect(token auth.Token, aspectIds []string) (importTypes []ImportTypePermissionSearch, err error, code int)
	GetNestedFunctionInfos(token auth.Token) (result []model.FunctionInfo, err error)
	GetAspectNodes(ids []string, token auth.Token) ([]model.AspectNode, error)
	GetAspectNodesWithMeasuringFunction(token auth.Token) ([]model.AspectNode, error)
	GetImportTypes(token auth.Token) (importTypes []ImportTypePermissionSearch, err error, code int)
	GetDeviceClassUses(token auth.Token) (result interface{}, err error)
}

type Lib struct {
	config           Config
	permissionsearch client.Client
}

func (this *Lib) Config() Config {
	return this.config
}

func New(config Config) *Lib {
	return &Lib{config: config, permissionsearch: client.NewClient(config.PermissionsUrl)}
}

func post(token string, url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", token)
	return http.DefaultClient.Do(req)
}

func postJson(token string, url string, in interface{}, out interface{}) (err error) {
	requestBody := new(bytes.Buffer)
	err = json.NewEncoder(requestBody).Encode(in)
	if err != nil {
		return err
	}
	resp, err := post(token, url, "application/json", requestBody)
	if err != nil {
		debug.PrintStack()
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		responseMsg, _ := io.ReadAll(resp.Body)
		debug.PrintStack()
		return errors.New(string(responseMsg))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func get(token string, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	req.Header.Set("Authorization", token)
	return http.DefaultClient.Do(req)
}

func GetJson(token string, url string, out interface{}) (err error) {
	resp, err := get(token, url)
	if err != nil {
		debug.PrintStack()
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		responseMsg, _ := io.ReadAll(resp.Body)
		debug.PrintStack()
		return errors.New(string(responseMsg))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
