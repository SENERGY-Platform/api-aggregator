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
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
)

func (this *Lib) CheckEventStates(token string, ids []string) (result map[string]bool, err error) {
	result = map[string]bool{}
	if this.config.EventManagerUrl == "" || this.config.EventManagerUrl == "-" {
		return result, nil
	}
	req, err := http.NewRequest("GET", this.config.EventManagerUrl+"/event-states?ids="+url.QueryEscape(strings.Join(ids, ",")), nil)
	if err != nil {
		debug.PrintStack()
		return result, err
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("ERROR: GetProcessDeploymentList()::http.DefaultClient.Do(req)", err)
		debug.PrintStack()
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		responseMsg, _ := ioutil.ReadAll(resp.Body)
		log.Println("ERROR: CheckEventState(): unexpected response", resp.StatusCode, string(responseMsg))
		debug.PrintStack()
		return result, errors.New(string(responseMsg))
	}
	if resp.StatusCode != 200 {
		responseMsg, _ := ioutil.ReadAll(resp.Body)
		log.Println("DEBUG: event pipeline not ready:", resp.StatusCode, string(responseMsg))
		debug.PrintStack()
		return result, nil
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
