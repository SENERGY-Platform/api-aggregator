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
	"github.com/SENERGY-Platform/api-aggregator/pkg/auth"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
)

func (this *Lib) GetExtendedProcessList(token auth.Token, query url.Values) (result []map[string]interface{}, err error) {
	processes, err := this.GetProcessDeploymentList(token, query)
	if err != nil {
		return result, err
	}
	ids := []string{}
	for _, process := range processes {
		id, ok := process["id"].(string)
		if !ok {
			log.Println("ERROR: unable to read process id", process)
			return result, errors.New("unable to read process id")
		}
		ids = append(ids, id)
	}
	metadata, err := this.GetProcessDependencyList(token, ids)
	if err != nil {
		return result, err
	}
	metadata, err = this.SetOnlineState(token, metadata)
	if err != nil {
		return result, err
	}
	metadataIndex := map[string]Dependencies{}
	for _, m := range metadata {
		metadataIndex[m.DeploymentId] = m
	}
	for _, process := range processes {
		id, ok := process["id"].(string)
		if !ok {
			log.Println("ERROR: unable to read process id", process)
			return result, errors.New("unable to read process id")
		}
		process["online"] = true
		process["offline_reasons"] = []OfflineReason{}
		if !metadataIndex[id].Online {
			process["online"] = false
			process["offline_reasons"], err = getOfflineReasons(metadataIndex[id])
		}
		result = append(result, process)
	}
	return
}

func (this *Lib) GetProcessDeploymentList(token auth.Token, query url.Values) (result []map[string]interface{}, err error) {
	if this.Config().CamundaWrapperUrl == "" || this.Config().CamundaWrapperUrl == "-" {
		log.Println("WARNING: no CamundaWrapperUrl url configured")
		return
	}
	req, err := http.NewRequest("GET", this.config.CamundaWrapperUrl+"/deployment?"+query.Encode(), nil)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	req.Header.Set("Authorization", token.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("ERROR: GetProcessDeploymentList()::http.DefaultClient.Do(req)", err)
		debug.PrintStack()
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		responseMsg, _ := ioutil.ReadAll(resp.Body)
		log.Println("ERROR: GetProcessDeploymentList(): unexpected response", resp.StatusCode, string(responseMsg))
		debug.PrintStack()
		return result, errors.New("unexpected response")
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		debug.PrintStack()
	}
	return result, err
}

func (this *Lib) GetProcessDependencyList(token auth.Token, processIds []string) (result []Dependencies, err error) {
	if this.Config().ProcessDeploymentUrl == "" || this.Config().ProcessDeploymentUrl == "-" {
		log.Println("WARNING: no ProcessDeploymentUrl url configured")
		return
	}
	err = GetJson(token.Token, this.config.ProcessDeploymentUrl+"/dependencies?ids="+strings.Join(processIds, ","), &result)
	return
}

func getOfflineReasons(metadata Dependencies) (result []OfflineReason, err error) {
	for _, device := range metadata.Devices {
		if !device.Online {
			result = append(result, OfflineReason{
				Type:           "device-offline",
				Id:             device.DeviceId,
				AdditionalInfo: map[string]interface{}{"name": device.Name, "tasks": device.BpmnResources},
				Description:    "device " + device.Name + " is offline",
			})
		}
	}
	for _, event := range metadata.Events {
		if !event.Online {
			result = append(result, OfflineReason{
				Type:           "event-filter-offline",
				Id:             event.EventId,
				AdditionalInfo: map[string]interface{}{"tasks": event.BpmnResources},
				Description:    "event-filter " + event.EventId + " is offline",
			})
		}
	}
	return
}
