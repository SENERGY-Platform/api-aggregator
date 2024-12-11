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
	"github.com/SENERGY-Platform/device-repository/lib/client"
	importRepo "github.com/SENERGY-Platform/import-repository/lib/client"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
)

type Interface interface {
	Config() Config
	ListGateways(token auth.Token, limit int64, offset int64) (result []map[string]interface{}, err error)
	GetExtendedProcessList(token auth.Token, query url.Values) (result []map[string]interface{}, err error)
	CompleteDeviceHistory(token auth.Token, duration string, devices []map[string]interface{}) (result []map[string]interface{}, err error)
	CompleteGatewayHistory(token auth.Token, duration string, devices []map[string]interface{}) (result []map[string]interface{}, err error)
	ListAllGateways(token auth.Token) (result []map[string]interface{}, err error)
	FindDevices(token auth.Token, limit int, offset int) ([]map[string]interface{}, error)
	GetMeasuringFunctionsForAspect(token auth.Token, aspectId string) (functions []Function, err error, code int)
	GetMeasuringFunctions(token auth.Token, functionIds []string) (functions []Function, err error, code int)
	GetImportTypesWithAspect(token auth.Token, aspectIds []string) (importTypes []ImportTypeWithCriteria, err error, code int)
	GetAspectNodes(ids []string, token auth.Token) ([]model.AspectNode, error)
	GetAspectNodesWithMeasuringFunction(token auth.Token) ([]model.AspectNode, error)
	GetImportTypes(token auth.Token) (importTypes []ImportTypeWithCriteria, err error, code int)
	GetDeviceClassUses(token auth.Token) (result interface{}, err error)
}

type Lib struct {
	config     Config
	deviceRepo client.Interface
	importRepo importRepo.Interface
}

func (this *Lib) Config() Config {
	return this.config
}

func New(config Config) *Lib {
	return &Lib{config: config, deviceRepo: client.NewClient(config.IotUrl), importRepo: importRepo.NewClient(config.ImportRepoUrl)}
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
