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
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/auth"
	"log"
	"runtime/debug"
	"sort"
)

func (this *Lib) FindDevicesCommon(token auth.Token, search string, deviceIds []string, queryCommons QueryListCommons, location string, state string) (devices []map[string]interface{}, err error) {
	var listIds *QueryListIds
	var find *QueryFind

	filterById := len(deviceIds) > 0 || location != ""

	filteredIds := []string{}
	if location != "" {
		locationDevices, err := this.GetDevicesInLocation(token, location)
		if err != nil {
			return nil, err
		}
		if len(deviceIds) == 0 {
			filteredIds = locationDevices
		} else {
			filteredIds = intersection(deviceIds, locationDevices)
		}
	}

	if search == "" && filterById {
		listIds = &QueryListIds{
			QueryListCommons: queryCommons,
			Ids:              filteredIds,
		}
	} else {
		var filter *Selection
		if filterById {
			filter = &Selection{
				Condition: ConditionConfig{
					Feature:   "id",
					Operation: QueryAnyValueInFeatureOperation,
					Value:     filteredIds,
				},
			}
		}
		if state != "" {
			stateFilter := Selection{
				Condition: ConditionConfig{
					Feature:   "annotations.connected",
					Operation: QueryEqualOperation,
				},
			}
			switch state {
			case "connected":
				stateFilter.Condition.Value = true
			case "disconnected":
				stateFilter.Condition.Value = false
			case "unknown":
				stateFilter.Condition.Value = nil
			default:
				return devices, errors.New("unknown state in query: " + state)
			}
			if filter != nil {
				temp := Selection{
					And: []Selection{
						*filter,
						stateFilter,
					},
				}
				filter = &temp
			} else {
				filter = &stateFilter
			}
		}
		find = &QueryFind{
			QueryListCommons: queryCommons,
			Search:           search,
			Filter:           filter,
		}
	}

	query := QueryMessage{
		Resource: "devices",
		Find:     find,
		ListIds:  listIds,
	}

	err, _ = this.QueryPermissionsSearch(token.Token, query, &devices)
	if err != nil {
		return nil, err
	}
	return this.completeDeviceList(token, devices)
}

func (this *Lib) FindDevices(token auth.Token, search string, deviceIds []string, limit int, offset int, orderfeature string, direction string, location string, state string) (devices []map[string]interface{}, err error) {
	queryCommons := QueryListCommons{
		Limit:    limit,
		Offset:   offset,
		Rights:   "r",
		SortBy:   orderfeature,
		SortDesc: direction == "desc",
	}
	return this.FindDevicesCommon(token, search, deviceIds, queryCommons, location, state)
}

func (this *Lib) FindDevicesAfter(token auth.Token, search string, deviceIds []string, limit int, afterId string, afterSortValue string, orderfeature string, direction string, location string, state string) ([]map[string]interface{}, error) {
	after := ListAfter{
		Id: afterId,
	}
	err := json.Unmarshal([]byte(afterSortValue), &after.SortFieldValue)
	if err != nil {
		return nil, err
	}
	queryCommons := QueryListCommons{
		Limit:    limit,
		Rights:   "r",
		SortBy:   orderfeature,
		SortDesc: direction == "desc",
		After:    &after,
	}
	return this.FindDevicesCommon(token, search, deviceIds, queryCommons, location, state)
}

func (this *Lib) SortByName(input []map[string]interface{}, sortAsc bool) (output []map[string]interface{}) {
	output = input
	if sortAsc == true {
		sort.Slice(output, func(i, j int) bool {
			return output[i]["name"].(string) < output[j]["name"].(string)
		})
	} else {
		sort.Slice(output, func(i, j int) bool {
			return output[i]["name"].(string) > output[j]["name"].(string)
		})
	}
	return
}

func intersection(a []string, b []string) (result []string) {
	result = []string{}
	aIndex := map[string]bool{}
	for _, element := range a {
		aIndex[element] = true
	}
	for _, element := range b {
		if aIndex[element] {
			result = append(result, element)
		}
	}
	return result
}

func (this *Lib) FilterDevicesByState(token auth.Token, devices []map[string]interface{}, state string) (result []map[string]interface{}, err error) {
	devicesWithOnlineState, err := this.completeDeviceList(token, devices)
	if err != nil {
		log.Println("ERROR GetConnectionFilteredDevices.completeDeviceList()", err)
		return devices, err
	}

	for _, device := range devicesWithOnlineState {
		devicestate, ok := device["log_state"]
		if state == "connected" && ok && devicestate.(bool) {
			result = append(result, device)
		}
		if state == "disconnected" && ok && !devicestate.(bool) {
			result = append(result, device)
		}
		if state == "unknown" && !ok {
			result = append(result, device)
		}
	}
	return
}

func (this *Lib) CompleteDevices(token auth.Token, ids []string) (result []map[string]interface{}, err error) {
	devices, err := this.PermDeviceIdList(token, ids, "r")
	if err != nil {
		return result, err
	}
	return this.completeDeviceList(token, devices)
}

func (this *Lib) CompleteDevicesOrdered(token auth.Token, ids []string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	devices, err := this.PermDeviceIdListOrdered(token, ids, "r", limit, offset, orderfeature, direction)
	if err != nil {
		return result, err
	}
	return this.completeDeviceList(token, devices)
}

func (this *Lib) completeDeviceList(token auth.Token, devices []map[string]interface{}) (result []map[string]interface{}, err error) {
	ids := []string{}
	deviceMap := map[string]map[string]interface{}{}
	for _, device := range devices {
		id, ok := device["id"]
		if !ok {
			err = errors.New("unable to get device id")
			return
		}
		idStr, ok := id.(string)
		if !ok {
			err = errors.New("unable to cast device id to string")
			return
		}
		ids = append(ids, idStr)
		deviceMap[idStr] = device
	}

	logStates := map[string]bool{}
	if this.config.UseAnnotationsForConnectionState {
		for _, device := range devices {
			id, idok := device["id"].(string)
			if idok {
				annotations, aok := device["annotations"].(map[string]interface{})
				if aok {
					state, sok := annotations["connected"].(bool)
					if sok {
						logStates[id] = state
					}
				}
			}
		}
	} else {
		logStates, err = this.GetDeviceLogStates(token, ids)
		if err != nil {
			log.Println("ERROR completeDeviceList.GetDeviceLogStates()", err)
			return result, err
		}
	}

	deviceTypes, err := this.getDeviceDeviceTypeInfos(token, devices)
	if err != nil {
		return result, err
	}

	for _, id := range ids {
		device := deviceMap[id]
		logState, logExists := logStates[id]
		if logExists {
			device["log_state"] = logState
		}

		device["device_type"] = deviceTypes[id]
		// delete(device, "device_type_id") //device_type_id is not needed

		//device["gateway_name"] = gateways[id]
		result = append(result, device)
	}
	return
}

func (this *Lib) CompleteDeviceHistory(token auth.Token, duration string, devices []map[string]interface{}) (result []map[string]interface{}, err error) {
	ids := []string{}
	deviceMap := map[string]map[string]interface{}{}
	for _, device := range devices {
		id, ok := device["id"]
		if !ok {
			err = errors.New("unable to get device id")
			return
		}
		idStr, ok := id.(string)
		if !ok {
			err = errors.New("unable to cast device id to string")
			return
		}
		ids = append(ids, idStr)
		deviceMap[idStr] = device
	}
	logStates, err := this.GetDeviceLogStates(token, ids)
	if err != nil {
		log.Println("ERROR completeDeviceList.GetDeviceLogStates()", err)
		return result, err
	}
	logHistory, err := this.GetDeviceLogHistory(token, ids, duration)
	if err != nil {
		log.Println("ERROR completeDeviceList.GetDeviceLogHistory()", err)
		return result, err
	}
	logEdges, err := this.GetLogedges(token, "device", ids, duration)
	if err != nil {
		log.Println("ERROR completeDeviceList.GetLogedges()", err)
		return result, err
	}
	for _, id := range ids {
		device := deviceMap[id]
		logState, logExists := logStates[id]
		if !logExists {
			device["log_state"] = "unknown"
		} else {
			if logState {
				device["log_state"] = "connected"
			} else {
				device["log_state"] = "disconnected"
			}
		}
		device["log_history"] = logHistory[id]
		device["log_edge"] = logEdges[id]
		result = append(result, device)
	}
	return
}

func (this *Lib) getDeviceDeviceTypeInfos(token auth.Token, devices []map[string]interface{}) (deviceToDeviceType map[string]map[string]interface{}, err error) {
	deviceToDeviceType = map[string]map[string]interface{}{}

	//ensure no device type id duplicates
	dtIdSet := map[string]bool{}
	for _, device := range devices {
		dtIdInterface, dtIdExists := device["device_type_id"]
		if !dtIdExists {
			log.Println("WARNING: unable to find device type id field in device")
			continue
		}

		dtId, dtIdIsString := dtIdInterface.(string)
		if !dtIdIsString {
			log.Println("WARNING: device type id field is not string in device")
			continue
		}
		dtIdSet[dtId] = true
	}
	ids := []string{}
	for id, _ := range dtIdSet {
		ids = append(ids, id)
	}

	//get device types
	deviceTypes, err := this.PermSelectDeviceTypesByIdRead(token, ids)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return deviceToDeviceType, err
	}

	//index device types by its own id
	dtIndex := map[string]map[string]interface{}{}
	for _, deviceType := range deviceTypes {
		idInterface, idExists := deviceType["id"]
		if !idExists {
			log.Println("WARNING: unable to find device type id field")
			continue
		}

		id, idIsString := idInterface.(string)
		if !idIsString {
			log.Println("WARNING: device type id field is not string")
		}
		dtIndex[id] = deviceType
	}

	//result: index device types by device id
	for _, device := range devices {
		idInterface, idExists := device["id"]
		if !idExists {
			log.Println("WARNING: unable to find device id field")
			continue
		}

		id, idIsString := idInterface.(string)
		if !idIsString {
			log.Println("WARNING: device id field is not string")
		}

		dtIdInterface, dtIdExists := device["device_type_id"]
		if !dtIdExists {
			log.Println("WARNING: unable to find device type id field in device")
			continue
		}

		dtId, dtIdIsString := dtIdInterface.(string)
		if !dtIdIsString {
			log.Println("WARNING: device type id field is not string in device")
			continue
		}
		deviceToDeviceType[id] = dtIndex[dtId]
	}
	return
}

func (this *Lib) GetDeviceTypeDevices(token auth.Token, id string, limit string, offset string, orderFeature string, direction string) (ids []string, err error) {
	type device struct {
		Id string `json:"id"`
	}
	var devices []device
	err = GetJson(token.Token, this.config.PermissionsUrl+"/jwt/select/devices/device_type_id/"+id+"/r/"+limit+"/"+offset+"/"+orderFeature+"/"+direction, &devices)

	for _, d := range devices {
		ids = append(ids, d.Id)
	}
	return
}
