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
	"net/url"
)

type Interface interface {
	Config() Config
	ListDevicesByUserTag(jwt jwt_http_router.Jwt, value string) (result []map[string]interface{}, err error)
	FilterDevicesByState(jwt jwt_http_router.Jwt, devices []map[string]interface{}, state string) (result []map[string]interface{}, err error)
	SortByName(input []map[string]interface{}, sortAsc bool) (output []map[string]interface{})
	ListDevicesByTag(jwt jwt_http_router.Jwt, value string) (result []map[string]interface{}, err error)
	GetConnectionFilteredDevicesOrder(jwt jwt_http_router.Jwt, value string, sortAsc bool) (result []map[string]interface{}, err error)
	GetConnectionFilteredDevices(jwt jwt_http_router.Jwt, value string) (result []map[string]interface{}, err error)
	ListAllDevices(jwt jwt_http_router.Jwt) (result []map[string]interface{}, err error)
	ListDevices(jwt jwt_http_router.Jwt, limit string, offset string) (result []map[string]interface{}, err error)
	SearchDevices(jwt jwt_http_router.Jwt, query string, limit string, offset string) (result []map[string]interface{}, err error)
	ListDevicesOrdered(jwt jwt_http_router.Jwt, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error)
	SearchDevicesOrdered(jwt jwt_http_router.Jwt, query string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error)
	ListOrderdDevicesByTag(jwt jwt_http_router.Jwt, value string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error)
	ListOrderedDevicesByUserTag(jwt jwt_http_router.Jwt, value string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error)
	CompleteDevices(jwt jwt_http_router.Jwt, ids []string) (result []map[string]interface{}, err error)
	CompleteDevicesOrdered(jwt jwt_http_router.Jwt, ids []string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error)
	GetDevicesHistory(jwt jwt_http_router.Jwt, duration string) (result []map[string]interface{}, err error)
	GetGatewaysHistory(jwt jwt_http_router.Jwt, duration string) (result []map[string]interface{}, err error)
	ListGateways(jwt jwt_http_router.Jwt, limit string, offset string) (result []map[string]interface{}, err error)
	SearchGateways(jwt jwt_http_router.Jwt, query string, limit string, offset string) (result []map[string]interface{}, err error)
	ListGatewaysOrdered(jwt jwt_http_router.Jwt, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error)
	SearchGatewaysOrdered(jwt jwt_http_router.Jwt, query string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error)
	GetExtendedProcessList(jwt jwt_http_router.Jwt, query url.Values) (result []map[string]interface{}, err error)
	CompleteDeviceHistory(jwt jwt_http_router.Jwt, duration string, devices []map[string]interface{}) (result []map[string]interface{}, err error)
}

type Lib struct {
	config Config
}

func (this *Lib) Config() Config {
	return this.config
}

func New(config Config) *Lib {
	return &Lib{config: config}
}
