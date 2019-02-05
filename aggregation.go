/*
 * Copyright 2018 InfAI (CC SES)
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

package main

import (
	"errors"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"sort"
	"strings"

	"log"
)

func GetConnectionFilteredDevicesOrder(jwt jwt_http_router.Jwt, value string, sortAsc bool) (result []map[string]interface{}, err error) {
	result, err = GetConnectionFilteredDevices(jwt,value)
	if err != nil {
		log.Println("ERROR GetConnectionFilteredDevices", err)
		return result, err
	}

	result = sortArray(result, sortAsc)

	return
}

func GetConnectionFilteredDevicesSearchOrder(jwt jwt_http_router.Jwt, value string, searchText string, sortAsc bool) (result []map[string]interface{}, err error) {
	result, err = GetConnectionFilteredDevices(jwt,value)

	if err != nil {
		log.Println("ERROR GetConnectionFilteredDevices", err)
		return result, err

	}

	result = filter(result, "name", searchText)
	result = sortArray(result, sortAsc)

	return
}

func sortArray(input []map[string]interface{}, sortAsc bool) (output []map[string]interface{})  {
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

func filter(list []map[string]interface{}, key string, value string)(result []map[string]interface{}) {
	for _, element := range list{
		str, ok := element[key].(string)
		if ok && strings.Contains(str, value) {
			result = append(result, element)
		}
	}
	return
}

func GetConnectionFilteredDevices(jwt jwt_http_router.Jwt, value string) (result []map[string]interface{}, err error) {
	devices, err := PermListAllDevices(jwt, "r")
	if err != nil {
		log.Println("ERROR GetConnectionFilteredDevices.PermListAllDevices()", err)
		return result, err
	}
	devicesWithOnlineState, err := completeDeviceList(jwt, devices)
	if err != nil {
		log.Println("ERROR GetConnectionFilteredDevices.completeDeviceList()", err)
		return result, err
	}

	for _, device := range devicesWithOnlineState {
		state, ok := device["log_state"]
		if value == "connected" && ok && state.(bool) {
			result = append(result, device)
		}
		if value == "disconnected" && ok && !state.(bool) {
			result = append(result, device)
		}
		if value == "unknown" && !ok {
			result = append(result, device)
		}
	}
	return
}

func ListDevices(jwt jwt_http_router.Jwt, limit string, offset string) (result []map[string]interface{}, err error) {
	devices, err := PermListDevices(jwt, "r", limit, offset)
	if err != nil {
		log.Println("ERROR ListDevices.PermListDevices()", err)
		return result, err
	}
	return completeDeviceList(jwt, devices)
}

func ListDevicesOrdered(jwt jwt_http_router.Jwt, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	devices, err := PermListDevicesOrdered(jwt, "r", limit, offset, orderfeature, direction)
	if err != nil {
		log.Println("ERROR ListDevices.PermListDevices()", err)
		return result, err
	}
	return completeDeviceList(jwt, devices)
}

func SearchDevices(jwt jwt_http_router.Jwt, query string, limit string, offset string) (result []map[string]interface{}, err error) {
	devices, err := PermSearchDevices(jwt, query, "r", limit, offset)
	if err != nil {
		return result, err
	}
	return completeDeviceList(jwt, devices)
}

func SearchDevicesOrdered(jwt jwt_http_router.Jwt, query string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	devices, err := PermSearchDevicesOrdered(jwt, query, "r", limit, offset, orderfeature, direction)
	if err != nil {
		return result, err
	}
	return completeDeviceList(jwt, devices)
}

func ListDevicesByTag(jwt jwt_http_router.Jwt, value string) (result []map[string]interface{}, err error) {
	devices, err := PermSelectTagDevices(jwt, value, "r")
	if err != nil {
		return result, err
	}
	return completeDeviceList(jwt, devices)
}

func ListOrderdDevicesByTag(jwt jwt_http_router.Jwt, value string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	devices, err := PermSelectTagDevicesOrdered(jwt, value, "r", limit, offset, orderfeature, direction)
	if err != nil {
		return result, err
	}
	return completeDeviceList(jwt, devices)
}

func ListDevicesByUserTag(jwt jwt_http_router.Jwt, value string) (result []map[string]interface{}, err error) {
	devices, err := PermSelectUserTagDevices(jwt, value, "r")
	if err != nil {
		return result, err
	}
	return completeDeviceList(jwt, devices)
}

func ListOrderedDevicesByUserTag(jwt jwt_http_router.Jwt, value string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	devices, err := PermSelectUserTagDevicesOrdered(jwt, value, "r", limit, offset, orderfeature, direction)
	if err != nil {
		return result, err
	}
	return completeDeviceList(jwt, devices)
}

func CompleteDevices(jwt jwt_http_router.Jwt, ids []string) (result []map[string]interface{}, err error) {
	devices, err := PermDeviceIdList(jwt, ids, "r")
	if err != nil {
		return result, err
	}
	return completeDeviceList(jwt, devices)
}

func CompleteDevicesOrdered(jwt jwt_http_router.Jwt, ids []string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	devices, err := PermDeviceIdListOrdered(jwt, ids, "r", limit, offset, orderfeature, direction)
	if err != nil {
		return result, err
	}
	return completeDeviceList(jwt, devices)
}

func completeDeviceList(jwt jwt_http_router.Jwt, devices []map[string]interface{}) (result []map[string]interface{}, err error) {
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
	logStates, err := GetDeviceLogStates(jwt, ids)
	if err != nil {
		log.Println("ERROR completeDeviceList.GetDeviceLogStates()", err)
		return result, err
	}
	/*
		gateways, err := GatewayNames(jwt, ids)
		if err != nil {
			log.Println("ERROR completeDeviceList.GatewayNames()", err)
			return result, err
		}
	*/
	for _, id := range ids {
		device := deviceMap[id]
		logState, logExists := logStates[id]
		if logExists {
			device["log_state"] = logState
		}
		//device["gateway_name"] = gateways[id]
		result = append(result, device)
	}
	return
}

func GetDevicesHistory(jwt jwt_http_router.Jwt, duration string) (result []map[string]interface{}, err error) {
	result, err = PermListAllDevices(jwt, "r")
	if err != nil {
		log.Println("ERROR PermListAllDevices()", err)
		return result, err
	}
	result, err = completeDeviceHistory(jwt, duration, result)
	return
}

func GetGatewaysHistory(jwt jwt_http_router.Jwt, duration string) (result []map[string]interface{}, err error) {
	result, err = PermListAllGateways(jwt, "r")
	if err != nil {
		log.Println("ERROR PermListAllGateways()", err)
		return result, err
	}
	result, err = completeGatewayHistory(jwt, duration, result)
	return
}

func completeDeviceHistory(jwt jwt_http_router.Jwt, duration string, devices []map[string]interface{}) (result []map[string]interface{}, err error) {
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
	logStates, err := GetDeviceLogStates(jwt, ids)
	if err != nil {
		log.Println("ERROR completeDeviceList.GetDeviceLogStates()", err)
		return result, err
	}
	logHistory, err := GetDeviceLogHistory(jwt, ids, duration)
	if err != nil {
		log.Println("ERROR completeDeviceList.GetDeviceLogHistory()", err)
		return result, err
	}
	logEdges, err := GetLogedges(jwt, "device", ids, duration)
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

func completeGatewayHistory(jwt jwt_http_router.Jwt, duration string, gateways []map[string]interface{}) (result []map[string]interface{}, err error) {
	ids := []string{}
	gatewayMap := map[string]map[string]interface{}{}
	for _, gateway := range gateways {
		id, ok := gateway["id"]
		if !ok {
			err = errors.New("unable to get gateway id")
			return
		}
		idStr, ok := id.(string)
		if !ok {
			err = errors.New("unable to cast gateway id to string")
			return
		}
		ids = append(ids, idStr)
		gatewayMap[idStr] = gateway
	}
	logStates, err := GetGatewayLogStates(jwt, ids)
	if err != nil {
		log.Println("ERROR completeGatewayList.GetGatewayLogStates()", err)
		return result, err
	}
	logHistory, err := GetGatewayLogHistory(jwt, ids, duration)
	if err != nil {
		log.Println("ERROR completeGatewayList.GetGatewayLogHistory()", err)
		return result, err
	}
	logEdges, err := GetLogedges(jwt, "gateway", ids, duration)
	if err != nil {
		log.Println("ERROR completeDeviceList.GetLogedges()", err)
		return result, err
	}
	for _, id := range ids {
		gateway := gatewayMap[id]
		logState, logExists := logStates[id]
		if !logExists {
			gateway["log_state"] = "unknown"
		} else {
			if logState {
				gateway["log_state"] = "connected"
			} else {
				gateway["log_state"] = "disconnected"
			}
		}
		gateway["log_history"] = logHistory[id]
		gateway["log_edge"] = logEdges[id]
		result = append(result, gateway)
	}
	return
}

func ListGateways(jwt jwt_http_router.Jwt, limit string, offset string) (result []map[string]interface{}, err error) {
	gateways, err := PermListGateways(jwt, "r", limit, offset)
	if err != nil {
		log.Println("ERROR ListGateways.PermListGateways()", err)
		return result, err
	}
	return completeGatewayList(jwt, gateways)
}

func ListGatewaysOrdered(jwt jwt_http_router.Jwt, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	gateways, err := PermListGatewaysOrdered(jwt, "r", limit, offset, orderfeature, direction)
	if err != nil {
		log.Println("ERROR ListGateways.PermListGateways()", err)
		return result, err
	}
	return completeGatewayList(jwt, gateways)
}

func SearchGateways(jwt jwt_http_router.Jwt, query string, limit string, offset string) (result []map[string]interface{}, err error) {
	gateways, err := PermSearchGateways(jwt, query, "r", limit, offset)
	if err != nil {
		return result, err
	}
	return completeGatewayList(jwt, gateways)
}

func SearchGatewaysOrdered(jwt jwt_http_router.Jwt, query string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	gateways, err := PermSearchGatewaysOrdered(jwt, query, "r", limit, offset, orderfeature, direction)
	if err != nil {
		return result, err
	}
	return completeGatewayList(jwt, gateways)
}

func completeGatewayList(jwt jwt_http_router.Jwt, gateways []map[string]interface{}) (result []map[string]interface{}, err error) {
	ids := []string{}
	gatewayMap := map[string]map[string]interface{}{}
	for _, gateway := range gateways {
		id, ok := gateway["id"]
		if !ok {
			err = errors.New("unable to get gateway id")
			return
		}
		idStr, ok := id.(string)
		if !ok {
			err = errors.New("unable to cast gateway id to string")
			return
		}
		ids = append(ids, idStr)
		gatewayMap[idStr] = gateway
	}
	logStates, err := GetGatewayLogStates(jwt, ids)
	if err != nil {
		log.Println("ERROR completeGatewayList.GetGatewayLogStates()", err)
		return result, err
	}
	for _, id := range ids {
		gateway := gatewayMap[id]
		logState, logExists := logStates[id]
		if logExists {
			gateway["log_state"] = logState
		}
		//gateway["gateway_name"] = gateways[id]
		result = append(result, gateway)
	}
	return
}
