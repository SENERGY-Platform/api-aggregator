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

package lib

import (
	"encoding/json"
	"errors"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"net/url"
)

func PermListGateways(jwt jwt_http_router.Jwt, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return PermList(jwt, "gateway", right, limit, offset)
}

func PermListGatewaysOrdered(jwt jwt_http_router.Jwt, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return PermListOrdered(jwt, "gateway", right, limit, offset, orderfeature, direction)
}

func PermSearchGateways(jwt jwt_http_router.Jwt, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return PermSearch(jwt, "gateway", query, right, limit, offset)
}

func PermSearchGatewaysOrdered(jwt jwt_http_router.Jwt, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return PermSearchOrdered(jwt, "gateway", query, right, limit, offset, orderfeature, direction)
}

func PermDeviceIdList(jwt jwt_http_router.Jwt, ids []string, right string) (result []map[string]interface{}, err error) {
	return PermIdList(jwt, "deviceinstance", ids, right)
}

func PermDeviceIdListOrdered(jwt jwt_http_router.Jwt, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return PermIdListOrdered(jwt, "deviceinstance", ids, right, limit, offset, orderfeature, direction)
}

func PermIdList(jwt jwt_http_router.Jwt, kind string, ids []string, right string) (result []map[string]interface{}, err error) {
	err = jwt.Impersonate.PostJSON(Config.PermissionsUrl+"/ids/select/"+url.PathEscape( kind)+"/"+right, ids, &result)
	return
}

func PermIdListOrdered(jwt jwt_http_router.Jwt, kind string, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	err = jwt.Impersonate.PostJSON(Config.PermissionsUrl+"/ids/select/"+url.PathEscape( kind)+"/"+right+"/"+limit+"/"+offset+"/"+url.PathEscape( orderfeature)+"/"+direction, ids, &result)
	return
}

func PermListDevices(jwt jwt_http_router.Jwt, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return PermList(jwt, "deviceinstance", right, limit, offset)
}

func PermListDevicesOrdered(jwt jwt_http_router.Jwt, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return PermListOrdered(jwt, "deviceinstance", right, limit, offset, orderfeature, direction)
}

func PermListAllDevices(jwt jwt_http_router.Jwt, right string) (result []map[string]interface{}, err error) {
	return PermListAll(jwt, "deviceinstance", right)
}

func PermListAllGateways(jwt jwt_http_router.Jwt, right string) (result []map[string]interface{}, err error) {
	return PermListAll(jwt, "gateway", right)
}

func PermListAll(jwt jwt_http_router.Jwt, kind string, right string) (result []map[string]interface{}, err error) {
	resp, err := jwt.Impersonate.Get(Config.PermissionsUrl + "/jwt/list/" + url.PathEscape( kind) + "/" + right)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("access denied")
		return result, err
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func PermList(jwt jwt_http_router.Jwt, kind string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	//"/jwt/list/:resource_kind/:right"
	resp, err := jwt.Impersonate.Get(Config.PermissionsUrl + "/jwt/list/" + url.PathEscape( kind) + "/" + right + "/" + limit + "/" + offset)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("access denied")
		return result, err
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func PermListOrdered(jwt jwt_http_router.Jwt, kind string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	resp, err := jwt.Impersonate.Get(Config.PermissionsUrl + "/jwt/list/" + url.PathEscape( kind) + "/" + right + "/" + limit + "/" + offset + "/" + url.PathEscape( orderfeature) + "/" + direction)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("access denied")
		return result, err
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func PermSearchDevices(jwt jwt_http_router.Jwt, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return PermSearch(jwt, "deviceinstance", query, right, limit, offset)
}

func PermSearchDevicesOrdered(jwt jwt_http_router.Jwt, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return PermSearchOrdered(jwt, "deviceinstance", query, right, limit, offset, orderfeature, direction)
}

func PermSearch(jwt jwt_http_router.Jwt, kind string, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	//"/jwt/search/:resource_kind/:query/:right/:limit/:offset"
	resp, err := jwt.Impersonate.Get(Config.PermissionsUrl + "/jwt/search/" + url.PathEscape( kind) + "/" + url.PathEscape( query) + "/" + right + "/" + limit + "/" + offset)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("access denied")
		return result, err
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func PermSearchOrdered(jwt jwt_http_router.Jwt, kind string, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	resp, err := jwt.Impersonate.Get(Config.PermissionsUrl + "/jwt/search/" + url.PathEscape( kind) + "/" + url.PathEscape( query) + "/" + right + "/" + limit + "/" + offset + "/" + url.PathEscape( orderfeature) + "/" + direction)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("access denied")
		return result, err
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func PermSelectUserTagDevices(jwt jwt_http_router.Jwt, value string, right string) (result []map[string]interface{}, err error) {
	return PermSelect(jwt, "deviceinstance", "usertag", value, right)
}

func PermSelectUserTagDevicesOrdered(jwt jwt_http_router.Jwt, value string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return PermSelectOrdered(jwt, "deviceinstance", "usertag", value, right, limit, offset, orderfeature, direction)
}

func PermSelectTagDevices(jwt jwt_http_router.Jwt, value string, right string) (result []map[string]interface{}, err error) {
	return PermSelect(jwt, "deviceinstance", "tag", value, right)
}

func PermSelectTagDevicesOrdered(jwt jwt_http_router.Jwt, value string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return PermSelectOrdered(jwt, "deviceinstance", "tag", value, right, limit, offset, orderfeature, direction)
}

func PermSelect(jwt jwt_http_router.Jwt, kind string, field string, value string, right string) (result []map[string]interface{}, err error) {
	resp, err := jwt.Impersonate.Get(Config.PermissionsUrl + "/jwt/select/" + url.PathEscape( kind) + "/" + url.PathEscape( field) + "/" + url.PathEscape(value)   + "/" + right)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("access denied")
		return result, err
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func PermSelectOrdered(jwt jwt_http_router.Jwt, kind string, field string, value string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	resp, err := jwt.Impersonate.Get(Config.PermissionsUrl + "/jwt/select/" + url.PathEscape( kind) + "/" + url.PathEscape( field) + "/" + url.PathEscape(value) + "/" + right + "/" + limit + "/" + offset + "/" + url.PathEscape( orderfeature) + "/" + direction)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("access denied")
		return result, err
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func PermCheckDeviceAdmin(jwt jwt_http_router.Jwt, ids []string) (result map[string]bool, err error) {
	return PermCheck(jwt, "deviceinstance", ids, "a")
}

func PermCheckDeviceRead(jwt jwt_http_router.Jwt, ids []string) (result map[string]bool, err error) {
	return PermCheck(jwt, "deviceinstance", ids, "r")
}

func PermCheck(jwt jwt_http_router.Jwt, kind string, ids []string, right string) (result map[string]bool, err error) {
	result = map[string]bool{}
	err = jwt.Impersonate.PostJSON(Config.PermissionsUrl+"/ids/check/"+url.PathEscape( kind)+"/"+right, ids, &result)
	return
}
