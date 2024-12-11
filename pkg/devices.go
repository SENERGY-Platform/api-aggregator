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
	"github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/models/go/models"
	"log"
)

func (this *Lib) FindDevices(token auth.Token, limit int, offset int) (devices []map[string]interface{}, err error) {
	devicesFromRepo, _, err, _ := this.deviceRepo.ListExtendedDevices(token.Jwt(), client.ExtendedDeviceListOptions{
		Limit:      int64(limit),
		Offset:     int64(offset),
		SortBy:     "name.asc",
		Permission: client.READ,
		FullDt:     true,
	})
	if err != nil {
		return nil, err
	}
	return this.legacyDeviceTransformations(token, devicesFromRepo)
}

func (this *Lib) legacyDeviceTransformations(token auth.Token, devices []models.ExtendedDevice) (result []map[string]interface{}, err error) {
	for _, device := range devices {
		element := map[string]interface{}{}
		temp, err := json.Marshal(device)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(temp, &element)
		if err != nil {
			return nil, err
		}
		switch device.ConnectionState {
		case *client.ConnectionStateOnline:
			element["log_state"] = "connected"
		case *client.ConnectionStateOffline:
			element["log_state"] = "disconnected"
		case *client.ConnectionStateUnknown:
			element["log_state"] = "unknown"
		default:
			element["log_state"] = "unknown"
		}
		element["creator"] = device.OwnerId
		result = append(result, element)
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
	logHistory, err := this.GetDeviceLogHistory(token, ids, duration)
	if err != nil {
		log.Println("ERROR legacyDeviceTransformations.GetDeviceLogHistory()", err)
		return result, err
	}
	logEdges, err := this.GetLogedges(token, "device", ids, duration)
	if err != nil {
		log.Println("ERROR legacyDeviceTransformations.GetLogedges()", err)
		return result, err
	}
	for _, id := range ids {
		device := deviceMap[id]
		device["log_history"] = logHistory[id]
		device["log_edge"] = logEdges[id]
		result = append(result, device)
	}
	return
}
