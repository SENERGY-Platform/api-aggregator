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

package api

import (
	"github.com/SmartEnergyPlatform/api-aggregator/api/deprecated"
	"github.com/SmartEnergyPlatform/api-aggregator/lib"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http/httptest"
	"net/url"
	"reflect"
)

type MethodCallLog struct {
	Name      string
	Parameter []interface{}
}

type mock struct {
	CallLog []MethodCallLog
}

func (this *mock) Log(name string, parameter ...interface{}) {
	this.CallLog = append(this.CallLog, MethodCallLog{Name: name, Parameter: parameter})
}

func (this *mock) Config() lib.Config {
	this.Log("Config")
	return lib.Config{}
}

func (this *mock) ListDevicesByUserTag(jwt jwt_http_router.Jwt, value string) (result []map[string]interface{}, err error) {
	this.Log("ListDevicesByUserTag", jwt, value)
	return []map[string]interface{}{}, nil
}

func (this *mock) FilterDevicesByState(jwt jwt_http_router.Jwt, devices []map[string]interface{}, state string) (result []map[string]interface{}, err error) {
	this.Log("FilterDevicesByState", jwt, devices, state)
	return []map[string]interface{}{}, nil
}

func (this *mock) SortByName(input []map[string]interface{}, sortAsc bool) (output []map[string]interface{}) {
	this.Log("SortByName", input, sortAsc)
	return []map[string]interface{}{}
}

func (this *mock) ListDevicesByTag(jwt jwt_http_router.Jwt, value string) (result []map[string]interface{}, err error) {
	this.Log("ListDevicesByTag", jwt, value)
	return []map[string]interface{}{}, nil
}

func (this *mock) GetConnectionFilteredDevicesOrder(jwt jwt_http_router.Jwt, value string, sortAsc bool) (result []map[string]interface{}, err error) {
	result, err = this.GetConnectionFilteredDevices(jwt, value)
	result = this.SortByName(result, sortAsc)
	return
}

func (this *mock) GetConnectionFilteredDevices(jwt jwt_http_router.Jwt, value string) (result []map[string]interface{}, err error) {
	devices, err := this.PermListAllDevices(jwt, "r")
	if err != nil {
		log.Println("ERROR GetConnectionFilteredDevices.PermListAllDevices()", err)
		return result, err
	}
	return this.FilterDevicesByState(jwt, devices, value)
}

func (this *mock) ListDevices(jwt jwt_http_router.Jwt, limit string, offset string) (result []map[string]interface{}, err error) {
	this.Log("ListDevices", jwt, limit, offset)
	return []map[string]interface{}{}, nil
}

func (this *mock) SearchDevices(jwt jwt_http_router.Jwt, query string, limit string, offset string) (result []map[string]interface{}, err error) {
	this.Log("SearchDevices", jwt, query, limit, offset)
	return []map[string]interface{}{}, nil
}

func (this *mock) ListDevicesOrdered(jwt jwt_http_router.Jwt, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	this.Log("ListDevicesOrdered", jwt, limit, offset)
	return []map[string]interface{}{}, nil
}

func (this *mock) SearchDevicesOrdered(jwt jwt_http_router.Jwt, query string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	this.Log("SearchDevicesOrdered", jwt, query, limit, offset, orderfeature, direction)
	return []map[string]interface{}{}, nil
}

func (this *mock) ListOrderdDevicesByTag(jwt jwt_http_router.Jwt, value string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	this.Log("ListOrderdDevicesByTag", jwt, value, limit, offset, orderfeature, direction)
	return []map[string]interface{}{}, nil
}

func (this *mock) ListOrderedDevicesByUserTag(jwt jwt_http_router.Jwt, value string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	this.Log("ListOrderedDevicesByUserTag", jwt, value, limit, offset, orderfeature, direction)
	return []map[string]interface{}{}, nil
}

func (this *mock) CompleteDevices(jwt jwt_http_router.Jwt, ids []string) (result []map[string]interface{}, err error) {
	this.Log("CompleteDevices", jwt, ids)
	return []map[string]interface{}{}, nil
}

func (this *mock) CompleteDevicesOrdered(jwt jwt_http_router.Jwt, ids []string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	this.Log("CompleteDevicesOrdered", jwt, ids, limit, offset, orderfeature, direction)
	return []map[string]interface{}{}, nil
}

func (this *mock) GetDevicesHistory(jwt jwt_http_router.Jwt, duration string) (result []map[string]interface{}, err error) {
	result, err = this.PermListAllDevices(jwt, "r")
	result, err = this.CompleteDeviceHistory(jwt, duration, result)
	return
}

func (this *mock) GetGatewaysHistory(jwt jwt_http_router.Jwt, duration string) (result []map[string]interface{}, err error) {
	this.Log("GetGatewaysHistory", jwt, duration)
	return []map[string]interface{}{}, nil
}

func (this *mock) ListGateways(jwt jwt_http_router.Jwt, limit string, offset string) (result []map[string]interface{}, err error) {
	this.Log("ListGateways", jwt, limit, offset)
	return []map[string]interface{}{}, nil
}

func (this *mock) SearchGateways(jwt jwt_http_router.Jwt, query string, limit string, offset string) (result []map[string]interface{}, err error) {
	this.Log("SearchGateways", jwt, query, limit, offset)
	return []map[string]interface{}{}, nil
}

func (this *mock) ListGatewaysOrdered(jwt jwt_http_router.Jwt, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	this.Log("ListGatewaysOrdered", jwt, limit, offset, orderfeature, direction)
	return []map[string]interface{}{}, nil
}

func (this *mock) SearchGatewaysOrdered(jwt jwt_http_router.Jwt, query string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	this.Log("SearchGatewaysOrdered", jwt, query, limit, offset, orderfeature, direction)
	return []map[string]interface{}{}, nil
}

func (this *mock) GetExtendedProcessList(jwt jwt_http_router.Jwt, query url.Values) (result []map[string]interface{}, err error) {
	this.Log("GetExtendedProcessList", jwt, query)
	return []map[string]interface{}{}, nil
}

func (this *mock) ListAllDevices(jwt jwt_http_router.Jwt) (result []map[string]interface{}, err error) {
	return this.PermListAllDevices(jwt, "r")
}

func (this *mock) Compare(other *mock) bool {
	return reflect.DeepEqual(this.CallLog, other.CallLog)
}

func (this *mock) PermListAllDevices(jwt jwt_http_router.Jwt, s string) (result []map[string]interface{}, err error) {
	this.Log("PermListAllDevices", jwt, s)
	return []map[string]interface{}{}, nil
}

func (this *mock) CompleteDeviceHistory(jwt jwt_http_router.Jwt, duration string, devices []map[string]interface{}) (result []map[string]interface{}, err error) {
	this.Log("CompleteDeviceHistory", jwt, duration, devices)
	return []map[string]interface{}{}, nil
}

func newMock() (oldApi string, newApi string, libForOld *mock, libForNew *mock, stop func()) {
	libForNew, libForOld = &mock{}, &mock{}
	newApiServer, oldApiServer := httptest.NewServer(getRoutes(libForNew)), httptest.NewServer(deprecated.GetRoutes(libForOld))
	newApi, oldApi = newApiServer.URL, oldApiServer.URL
	stop = func() {
		newApiServer.Close()
		oldApiServer.Close()
	}
	return
}
