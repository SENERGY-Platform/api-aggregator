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
	"github.com/SENERGY-Platform/api-aggregator/pkg/model"
)

func (this *Lib) GetDeviceClassUses(token auth.Token) (result interface{}, err error) {
	var devices []model.ShortDevice
	devices, err, _ = QueryPermissionSearchFindAll(this, token.Jwt(), QueryMessage{
		Resource: "devices",
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Rights: "r",
				SortBy: "id",
			},
		},
	}, func(e model.ShortDevice) string {
		return e.Id
	})
	if err != nil {
		return result, err
	}

	deviceTypeToDevice := map[string][]string{}
	for _, device := range devices {
		deviceTypeToDevice[device.DeviceTypeId] = append(deviceTypeToDevice[device.DeviceTypeId], device.Id)
	}

	deviceTypeIds := []string{}
	for id, _ := range deviceTypeToDevice {
		deviceTypeIds = append(deviceTypeIds, id)
	}

	deviceTypes := []model.ShortDeviceType{}
	err, _ = this.QueryPermissionsSearch(token.Jwt(), QueryMessage{
		Resource: "device-types",
		ListIds: &QueryListIds{
			QueryListCommons: QueryListCommons{
				Limit:  len(deviceTypeIds),
				Offset: 0,
				Rights: "r",
				SortBy: "name",
			},
			Ids: deviceTypeIds,
		},
	}, &deviceTypes)
	if err != nil {
		return result, err
	}

	deviceClassToDevices := map[string][]string{}
	for _, dt := range deviceTypes {
		deviceClassToDevices[dt.DeviceClassId] = append(deviceClassToDevices[dt.DeviceClassId], deviceTypeToDevice[dt.Id]...)
	}

	deviceClassIds := []string{}
	for id, _ := range deviceClassToDevices {
		deviceClassIds = append(deviceClassIds, id)
	}

	deviceClasses := []model.DeviceClass{}
	err, _ = this.QueryPermissionsSearch(token.Jwt(), QueryMessage{
		Resource: "device-classes",
		ListIds: &QueryListIds{
			QueryListCommons: QueryListCommons{},
			Ids:              deviceClassIds,
		},
	}, &deviceClasses)
	if err != nil {
		return result, err
	}

	return map[string]interface{}{"device-classes": deviceClasses, "used-devices": deviceClassToDevices}, nil
}
