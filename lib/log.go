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
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
)

func (this *Lib) SetOnlineState(jwt jwt_http_router.Jwt, dependencies []Dependencies) (result []Dependencies, err error) {

	//create device id list
	//use map to prevent duplicate ids
	deviceidset := map[string]bool{}
	deviceids := []string{}
	for _, dependency := range dependencies {
		for _, device := range dependency.Devices {
			deviceidset[device.DeviceId] = true
		}
	}
	for id, _ := range deviceidset {
		deviceids = append(deviceids, id)
	}

	//create event id list
	//use map to prevent duplicate ids
	eventidset := map[string]bool{}
	eventids := []string{}
	for _, dependency := range dependencies {
		for _, event := range dependency.Events {
			eventidset[event.EventId] = true
		}
	}
	for id, _ := range eventidset {
		eventids = append(eventids, id)
	}

	//get device states
	devicestates := map[string]bool{}
	devicestates, err = this.GetDeviceLogStates(jwt, deviceids)
	if err != nil {
		return result, err
	}

	//get event states
	eventstates := map[string]bool{}
	eventstates, err = this.CheckEventStates(string(jwt.Impersonate), eventids)
	if err != nil {
		return result, err
	}

	//translate device states and event states to dependencies state
	for _, dependency := range dependencies {
		dependency.Online = true
		for index, device := range dependency.Devices {
			device.Online = devicestates[device.DeviceId]
			if !device.Online {
				dependency.Online = false
			}
			dependency.Devices[index] = device
		}
		for index, event := range dependency.Events {
			event.Online = eventstates[event.EventId]
			if event.Online {
				dependency.Online = false
			}
			dependency.Events[index] = event
		}
		result = append(result, dependency)
	}
	return result, nil
}

func (this *Lib) GetDeviceLogStates(jwt jwt_http_router.Jwt, deviceIds []string) (result map[string]bool, err error) {
	result = map[string]bool{}
	if this.Config().ConnectionLogUrl == "" || this.Config().ConnectionLogUrl == "-" {
		log.Println("WARNING: no connectionlog url configured")
		return
	}
	err = jwt.Impersonate.PostJSON(this.config.ConnectionLogUrl+"/intern/state/device/check", deviceIds, &result)
	return
}

func (this *Lib) GetGatewayLogStates(jwt jwt_http_router.Jwt, deviceIds []string) (result map[string]bool, err error) {
	result = map[string]bool{}
	if this.Config().ConnectionLogUrl == "" || this.Config().ConnectionLogUrl == "-" {
		log.Println("WARNING: no connectionlog url configured")
		for _, id := range deviceIds {
			result[id] = true
		}
		return
	}
	err = jwt.Impersonate.PostJSON(this.config.ConnectionLogUrl+"/intern/state/gateway/check", deviceIds, &result)
	return
}

func (this *Lib) GetDeviceLogHistory(jwt jwt_http_router.Jwt, deviceIds []string, duration string) (result map[string]HistorySeries, err error) {
	if this.Config().ConnectionLogUrl == "" || this.Config().ConnectionLogUrl == "-" {
		log.Println("WARNING: no connectionlog url configured")
		result = map[string]HistorySeries{}
		return
	}
	return this.GetLogHistory(jwt, "device", deviceIds, duration)
}

func (this *Lib) GetGatewayLogHistory(jwt jwt_http_router.Jwt, ids []string, duration string) (result map[string]HistorySeries, err error) {
	return this.GetLogHistory(jwt, "gateway", ids, duration)
}

type HistoryResult struct {
	Series []HistorySeries `json:"Series"`
}

type HistorySeries struct {
	Name    string            `json:"name"`
	Tags    map[string]string `json:"tags"`
	Columns []string          `json:"columns"`
	Values  [][]interface{}   `json:"values"`
}

func (this *Lib) GetLogHistory(jwt jwt_http_router.Jwt, kind string, ids []string, duration string) (result map[string]HistorySeries, err error) {
	result = map[string]HistorySeries{}
	if this.Config().ConnectionLogUrl == "" || this.Config().ConnectionLogUrl == "-" {
		log.Println("WARNING: no connectionlog url configured")
		return
	}
	temp := []HistoryResult{}
	err = jwt.Impersonate.PostJSON(this.config.ConnectionLogUrl+"/intern/history/"+kind+"/"+duration, ids, &temp)
	if err != nil {
		return result, err
	}
	for _, series := range temp[0].Series {
		result[series.Tags[kind]] = series
	}
	return
}

func (this *Lib) GetLogstarts(jwt jwt_http_router.Jwt, kind string, ids []string) (result map[string]interface{}, err error) {
	result = map[string]interface{}{}
	if this.Config().ConnectionLogUrl == "" || this.Config().ConnectionLogUrl == "-" {
		log.Println("WARNING: no connectionlog url configured")
		return
	}
	err = jwt.Impersonate.PostJSON(this.config.ConnectionLogUrl+"/intern/logstarts/"+kind, ids, &result)
	return
}

func (this *Lib) GetLogedges(jwt jwt_http_router.Jwt, kind string, ids []string, duration string) (result map[string]interface{}, err error) {
	result = map[string]interface{}{}
	if this.Config().ConnectionLogUrl == "" || this.Config().ConnectionLogUrl == "-" {
		log.Println("WARNING: no connectionlog url configured")
		return
	}
	err = jwt.Impersonate.PostJSON(this.config.ConnectionLogUrl+"/intern/logedge/"+kind+"/"+duration, ids, &result)
	return
}
