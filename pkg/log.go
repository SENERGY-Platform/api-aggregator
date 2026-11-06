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
	"github.com/SENERGY-Platform/api-aggregator/pkg/auth"
)

func (this *Lib) SetOnlineState(token auth.Token, dependencies []Dependencies) (result []Dependencies, err error) {

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
	devicestates, err = this.GetDeviceLogStates(token, deviceids)
	if err != nil {
		return result, err
	}

	//get event states
	eventstates := map[string]bool{}
	eventstates, err = this.CheckEventStates(token.Token, eventids)
	if err != nil {
		return result, err
	}

	//translate device states and event states to dependencies state
	for _, dependency := range dependencies {
		dependency.Online = true
		for index, device := range dependency.Devices {
			device.Online = true
			temp, ok := devicestates[device.DeviceId]
			if ok && !temp {
				device.Online = false
				dependency.Online = false
			}
			dependency.Devices[index] = device
		}
		for index, event := range dependency.Events {
			event.Online = true
			temp, ok := eventstates[event.EventId]
			if ok && !temp {
				event.Online = false
				dependency.Online = false
			}
			dependency.Events[index] = event
		}
		result = append(result, dependency)
	}
	return result, nil
}

func (this *Lib) GetDeviceLogStates(token auth.Token, deviceIds []string) (result map[string]bool, err error) {
	result = map[string]bool{}
	if this.Config().ConnectionLogUrl == "" || this.Config().ConnectionLogUrl == "-" {
		this.Config().GetLogger().Warn("no connectionlog url configured")
		return
	}
	err = postJson(token.Token, this.config.ConnectionLogUrl+"/intern/state/device/check", deviceIds, &result)
	return
}

func (this *Lib) GetGatewayLogStates(token auth.Token, ids []string) (result map[string]bool, err error) {
	result = map[string]bool{}
	if this.Config().ConnectionLogUrl == "" || this.Config().ConnectionLogUrl == "-" {
		this.Config().GetLogger().Warn("no connectionlog url configured")
		for _, id := range ids {
			result[id] = true
		}
		return
	}
	err = postJson(token.Token, this.config.ConnectionLogUrl+"/intern/state/gateway/check", ids, &result)
	return
}
