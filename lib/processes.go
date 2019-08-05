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

package lib

import (
	"encoding/json"
	"errors"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
)

func (this *Lib) GetExtendedProcessList(jwt jwt_http_router.Jwt, query url.Values) (result []map[string]interface{}, err error) {
	processes, err := this.GetProcessDeploymentList(jwt, query)
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
	metadata, err := this.GetProcessDependencyList(jwt, ids)
	metadataIndex := map[string]Metadata{}
	for _, m := range metadata {
		metadataIndex[m.Process] = m
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

func (this *Lib) GetProcessDeploymentList(jwt jwt_http_router.Jwt, query url.Values) (result []map[string]interface{}, err error) {
	req, err := http.NewRequest("GET", this.config.CamundaWrapperUrl+"/deployment?"+query.Encode(), nil)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	req.Header.Set("Authorization", string(jwt.Impersonate))
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

func (this *Lib) GetProcessDependencyList(jwt jwt_http_router.Jwt, processIds []string) (result []Metadata, err error) {
	err = jwt.Impersonate.GetJSON(this.config.ProcessDeploymentUrl+"/dependencies?deployments="+strings.Join(processIds, ","), &result)
	return
}

func getOfflineReasons(metadata Metadata) (result []OfflineReason, err error) {
	for _, param := range metadata.Abstract.AbstractTasks {
		if param.State != "unknown" && param.State != "connected" {
			result = append(result, OfflineReason{
				Type:           "device-offline",
				Id:             param.Selected.Id,
				AdditionalInfo: map[string]interface{}{"name": param.Selected.Name, "tasks": param.Tasks},
				Description:    "device " + param.Selected.Name + " is " + param.State,
			})
		}
	}
	for _, event := range metadata.Abstract.MsgEvents {
		if event.State != "running" {
			result = append(result, OfflineReason{
				Type:           "event-filter-offline",
				Id:             event.FilterId,
				AdditionalInfo: map[string]string{"shape_id": event.ShapeId},
				Description:    "event-filter " + event.FilterId + " is " + event.State,
			})
		}
	}
	for _, event := range metadata.Abstract.ReceiveTasks {
		if event.State != "running" {
			result = append(result, OfflineReason{
				Type:           "event-filter-offline",
				Id:             event.FilterId,
				AdditionalInfo: map[string]string{"shape_id": event.ShapeId},
				Description:    "event-filter " + event.FilterId + " for shape " + event.ShapeId + " is " + event.State,
			})
		}
	}
	return
}

type OfflineReason struct {
	Type           string      `json:"type"`
	Id             string      `json:"id"`
	AdditionalInfo interface{} `json:"additional_info,omitempty"`
	Description    string      `json:"description"`
}

type Metadata struct {
	Process  string          `json:"process"`
	Abstract AbstractProcess `json:"abstract"`
	Online   bool            `json:"online"`
	Owner    string          `json:"owner"`
}

type AbstractProcess struct {
	AbstractTasks []AbstractTask `json:"abstract_tasks"`
	ReceiveTasks  []MsgEvent     `json:"receive_tasks"`
	MsgEvents     []MsgEvent     `json:"msg_events"`
}

type AbstractTask struct {
	Selected DeviceInstance `json:"selected"`
	State    string         `json:"state" bson:"-"`
	Tasks    []Task         `json:"tasks"`
}

type Task struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type MsgEvent struct {
	FilterId string `json:"filter_id,omitempty"`
	ShapeId  string `json:"shape_id"`
	State    string `json:"state,omitempty" bson:"-"`
}

type DeviceInstance struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
