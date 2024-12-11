/*
 * Copyright 2022 InfAI (CC SES)
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
	"github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/models/go/models"
	"maps"
	"slices"
)

func (this *Lib) GetDeviceClassUses(token auth.Token) (result interface{}, err error) {
	allDevices := []models.ExtendedDevice{}
	deviceClassToDevices := map[string][]string{}
	deviceTypeToDevice := map[string][]string{}
	for {
		var limit int64 = 9999
		var offset int64 = 9999
		devices, _, err, _ := this.deviceRepo.ListExtendedDevices(token.Jwt(), client.ExtendedDeviceListOptions{
			Limit:      limit,
			Offset:     offset,
			SortBy:     "name.asc",
			Permission: client.READ,
			FullDt:     true,
		})
		if err != nil {
			return result, err
		}
		allDevices = append(allDevices, devices...)
		for _, device := range devices {
			deviceClassToDevices[device.DeviceType.DeviceClassId] = append(deviceClassToDevices[device.DeviceType.DeviceClassId], device.DeviceType.DeviceClassId)
			deviceTypeToDevice[device.DeviceTypeId] = append(deviceTypeToDevice[device.DeviceTypeId], device.Id)
		}
		if int64(len(devices)) < limit {
			break
		}
		offset = offset + limit
	}
	deviceClassIds := slices.Collect(maps.Keys(deviceClassToDevices))
	deviceClasses, _, err, _ := this.deviceRepo.ListDeviceClasses(client.DeviceClassListOptions{
		Ids:    deviceClassIds,
		Limit:  int64(len(deviceClassIds)),
		Offset: 0,
		SortBy: "name.asc",
	})
	if err != nil {
		return result, err
	}

	return map[string]interface{}{"device-classes": deviceClasses, "used-devices": deviceClassToDevices}, nil
}
