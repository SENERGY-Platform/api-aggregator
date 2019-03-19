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

func (this *Lib) PermListGateways(jwt jwt_http_router.Jwt, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return this.PermList(jwt, "gateway", right, limit, offset)
}

func (this *Lib) PermListGatewaysOrdered(jwt jwt_http_router.Jwt, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermListOrdered(jwt, "gateway", right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermSearchGateways(jwt jwt_http_router.Jwt, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return this.PermSearch(jwt, "gateway", query, right, limit, offset)
}

func (this *Lib) PermSearchGatewaysOrdered(jwt jwt_http_router.Jwt, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermSearchOrdered(jwt, "gateway", query, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermDeviceIdList(jwt jwt_http_router.Jwt, ids []string, right string) (result []map[string]interface{}, err error) {
	return this.PermIdList(jwt, "deviceinstance", ids, right)
}

func (this *Lib) PermDeviceIdListOrdered(jwt jwt_http_router.Jwt, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermIdListOrdered(jwt, "deviceinstance", ids, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermIdList(jwt jwt_http_router.Jwt, kind string, ids []string, right string) (result []map[string]interface{}, err error) {
	err = jwt.Impersonate.PostJSON(this.config.PermissionsUrl+"/ids/select/"+url.PathEscape(kind)+"/"+right, ids, &result)
	return
}

func (this *Lib) PermIdListOrdered(jwt jwt_http_router.Jwt, kind string, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	err = jwt.Impersonate.PostJSON(this.config.PermissionsUrl+"/ids/select/"+url.PathEscape(kind)+"/"+right+"/"+limit+"/"+offset+"/"+url.PathEscape(orderfeature)+"/"+direction, ids, &result)
	return
}

func (this *Lib) PermListDevices(jwt jwt_http_router.Jwt, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return this.PermList(jwt, "deviceinstance", right, limit, offset)
}

func (this *Lib) PermListDevicesOrdered(jwt jwt_http_router.Jwt, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermListOrdered(jwt, "deviceinstance", right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermListAllDevices(jwt jwt_http_router.Jwt, right string) (result []map[string]interface{}, err error) {
	return this.PermListAll(jwt, "deviceinstance", right)
}

func (this *Lib) PermListAllGateways(jwt jwt_http_router.Jwt, right string) (result []map[string]interface{}, err error) {
	return this.PermListAll(jwt, "gateway", right)
}

func (this *Lib) PermListAll(jwt jwt_http_router.Jwt, kind string, right string) (result []map[string]interface{}, err error) {
	resp, err := jwt.Impersonate.Get(this.config.PermissionsUrl + "/jwt/list/" + url.PathEscape(kind) + "/" + right)
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

func (this *Lib) PermList(jwt jwt_http_router.Jwt, kind string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	//"/jwt/list/:resource_kind/:right"
	resp, err := jwt.Impersonate.Get(this.config.PermissionsUrl + "/jwt/list/" + url.PathEscape(kind) + "/" + right + "/" + limit + "/" + offset)
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

func (this *Lib) PermListOrdered(jwt jwt_http_router.Jwt, kind string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	resp, err := jwt.Impersonate.Get(this.config.PermissionsUrl + "/jwt/list/" + url.PathEscape(kind) + "/" + right + "/" + limit + "/" + offset + "/" + url.PathEscape(orderfeature) + "/" + direction)
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

func (this *Lib) PermSearchDevices(jwt jwt_http_router.Jwt, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return this.PermSearch(jwt, "deviceinstance", query, right, limit, offset)
}

func (this *Lib) PermSearchDevicesOrdered(jwt jwt_http_router.Jwt, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermSearchOrdered(jwt, "deviceinstance", query, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermSearch(jwt jwt_http_router.Jwt, kind string, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	//"/jwt/search/:resource_kind/:query/:right/:limit/:offset"
	resp, err := jwt.Impersonate.Get(this.config.PermissionsUrl + "/jwt/search/" + url.PathEscape(kind) + "/" + url.PathEscape(query) + "/" + right + "/" + limit + "/" + offset)
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

func (this *Lib) PermSearchOrdered(jwt jwt_http_router.Jwt, kind string, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	resp, err := jwt.Impersonate.Get(this.config.PermissionsUrl + "/jwt/search/" + url.PathEscape(kind) + "/" + url.PathEscape(query) + "/" + right + "/" + limit + "/" + offset + "/" + url.PathEscape(orderfeature) + "/" + direction)
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

func (this *Lib) PermSelectUserTagDevices(jwt jwt_http_router.Jwt, value string, right string) (result []map[string]interface{}, err error) {
	return this.PermSelect(jwt, "deviceinstance", "usertag", value, right)
}

func (this *Lib) PermSelectUserTagDevicesOrdered(jwt jwt_http_router.Jwt, value string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermSelectOrdered(jwt, "deviceinstance", "usertag", value, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermSelectTagDevices(jwt jwt_http_router.Jwt, value string, right string) (result []map[string]interface{}, err error) {
	return this.PermSelect(jwt, "deviceinstance", "tag", value, right)
}

func (this *Lib) PermSelectTagDevicesOrdered(jwt jwt_http_router.Jwt, value string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermSelectOrdered(jwt, "deviceinstance", "tag", value, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermSelect(jwt jwt_http_router.Jwt, kind string, field string, value string, right string) (result []map[string]interface{}, err error) {
	resp, err := jwt.Impersonate.Get(this.config.PermissionsUrl + "/jwt/select/" + url.PathEscape(kind) + "/" + url.PathEscape(field) + "/" + url.PathEscape(value) + "/" + right)
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

func (this *Lib) PermSelectOrdered(jwt jwt_http_router.Jwt, kind string, field string, value string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	resp, err := jwt.Impersonate.Get(this.config.PermissionsUrl + "/jwt/select/" + url.PathEscape(kind) + "/" + url.PathEscape(field) + "/" + url.PathEscape(value) + "/" + right + "/" + limit + "/" + offset + "/" + url.PathEscape(orderfeature) + "/" + direction)
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

func (this *Lib) PermCheckDeviceAdmin(jwt jwt_http_router.Jwt, ids []string) (result map[string]bool, err error) {
	return this.PermCheck(jwt, "deviceinstance", ids, "a")
}

func (this *Lib) PermCheckDeviceRead(jwt jwt_http_router.Jwt, ids []string) (result map[string]bool, err error) {
	return this.PermCheck(jwt, "deviceinstance", ids, "r")
}

func (this *Lib) PermCheck(jwt jwt_http_router.Jwt, kind string, ids []string, right string) (result map[string]bool, err error) {
	result = map[string]bool{}
	err = jwt.Impersonate.PostJSON(this.config.PermissionsUrl+"/ids/check/"+url.PathEscape(kind)+"/"+right, ids, &result)
	return
}
