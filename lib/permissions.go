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
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"net/url"
)

func (this *Lib) PermListGateways(jwt jwt_http_router.Jwt, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return this.PermList(jwt, "hubs", right, limit, offset)
}

func (this *Lib) PermListGatewaysOrdered(jwt jwt_http_router.Jwt, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermListOrdered(jwt, "hubs", right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermSearchGateways(jwt jwt_http_router.Jwt, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return this.PermSearch(jwt, "hubs", query, right, limit, offset)
}

func (this *Lib) PermSearchGatewaysOrdered(jwt jwt_http_router.Jwt, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermSearchOrdered(jwt, "hubs", query, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermDeviceIdList(jwt jwt_http_router.Jwt, ids []string, right string) (result []map[string]interface{}, err error) {
	return this.PermIdList(jwt, "devices", ids, right)
}

func (this *Lib) PermDeviceIdListOrdered(jwt jwt_http_router.Jwt, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermIdListOrdered(jwt, "devices", ids, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermIdList(jwt jwt_http_router.Jwt, kind string, ids []string, right string) (result []map[string]interface{}, err error) {
	err = jwt.Impersonate.PostJSON(this.config.PermissionsUrl+"/ids/select/"+url.PathEscape(kind)+"/"+right, ids, &result)
	return
}

func (this *Lib) PermIdListOrdered(jwt jwt_http_router.Jwt, kind string, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	err = jwt.Impersonate.PostJSON(this.config.PermissionsUrl+"/ids/select/"+url.PathEscape(kind)+"/"+right+"/"+limit+"/"+offset+"/"+url.PathEscape(orderfeature)+"/"+direction, ids, &result)
	return
}

func (this *Lib) PermListAllGateways(jwt jwt_http_router.Jwt, right string) (result []map[string]interface{}, err error) {
	return this.PermListAll(jwt, "hubs", right)
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

func (this *Lib) PermSelectIds(jwt jwt_http_router.Jwt, kind string, right string, ids []string) (result []map[string]interface{}, err error) {
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(ids)
	if err != nil {
		return
	}
	resp, err := jwt.Impersonate.Post(this.config.PermissionsUrl+"/ids/select/"+url.PathEscape(kind)+"/"+right, "application/json", b)
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

func (this *Lib) PermSelectDeviceTypesByIdRead(jwt jwt_http_router.Jwt, ids []string) (result []map[string]interface{}, err error) {
	return this.PermSelectIds(jwt, "device-types", "r", ids)
}

func (this *Lib) PermCheckDeviceAdmin(jwt jwt_http_router.Jwt, ids []string) (result map[string]bool, err error) {
	return this.PermCheck(jwt, "devices", ids, "a")
}

func (this *Lib) PermCheckDeviceRead(jwt jwt_http_router.Jwt, ids []string) (result map[string]bool, err error) {
	return this.PermCheck(jwt, "devices", ids, "r")
}

func (this *Lib) PermCheck(jwt jwt_http_router.Jwt, kind string, ids []string, right string) (result map[string]bool, err error) {
	result = map[string]bool{}
	err = jwt.Impersonate.PostJSON(this.config.PermissionsUrl+"/ids/check/"+url.PathEscape(kind)+"/"+right, ids, &result)
	return
}
