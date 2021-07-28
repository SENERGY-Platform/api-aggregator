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
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/auth"
	"net/http"
	"net/url"
)

func (this *Lib) PermListGateways(token auth.Token, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return this.PermList(token, "hubs", right, limit, offset)
}

func (this *Lib) PermListGatewaysOrdered(token auth.Token, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermListOrdered(token, "hubs", right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermSearchGateways(token auth.Token, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	return this.PermSearch(token, "hubs", query, right, limit, offset)
}

func (this *Lib) PermSearchGatewaysOrdered(token auth.Token, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermSearchOrdered(token, "hubs", query, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermDeviceIdList(token auth.Token, ids []string, right string) (result []map[string]interface{}, err error) {
	return this.PermIdList(token, "devices", ids, right)
}

func (this *Lib) PermDeviceIdListOrdered(token auth.Token, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermIdListOrdered(token, "devices", ids, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermIdList(token auth.Token, kind string, ids []string, right string) (result []map[string]interface{}, err error) {
	err = postJson(token.Token, this.config.PermissionsUrl+"/ids/select/"+url.PathEscape(kind)+"/"+right, ids, &result)
	return
}

func (this *Lib) PermIdListOrdered(token auth.Token, kind string, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	err = postJson(token.Token, this.config.PermissionsUrl+"/ids/select/"+url.PathEscape(kind)+"/"+right+"/"+limit+"/"+offset+"/"+url.PathEscape(orderfeature)+"/"+direction, ids, &result)
	return
}

func (this *Lib) PermListAllGateways(token auth.Token, right string) (result []map[string]interface{}, err error) {
	return this.PermListAll(token, "hubs", right)
}

func (this *Lib) PermListAll(token auth.Token, kind string, right string) (result []map[string]interface{}, err error) {
	resp, err := get(token.Token, this.config.PermissionsUrl+"/jwt/list/"+url.PathEscape(kind)+"/"+right)
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

func (this *Lib) PermList(token auth.Token, kind string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	//"/jwt/list/:resource_kind/:right"
	resp, err := get(token.Token, this.config.PermissionsUrl+"/jwt/list/"+url.PathEscape(kind)+"/"+right+"/"+limit+"/"+offset)
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

func (this *Lib) PermListOrdered(token auth.Token, kind string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	resp, err := get(token.Token, this.config.PermissionsUrl+"/jwt/list/"+url.PathEscape(kind)+"/"+right+"/"+limit+"/"+offset+"/"+url.PathEscape(orderfeature)+"/"+direction)
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

func (this *Lib) PermSearch(token auth.Token, kind string, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	//"/jwt/search/:resource_kind/:query/:right/:limit/:offset"
	resp, err := get(token.Token, this.config.PermissionsUrl+"/jwt/search/"+url.PathEscape(kind)+"/"+url.PathEscape(query)+"/"+right+"/"+limit+"/"+offset)
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

func (this *Lib) PermSearchOrdered(token auth.Token, kind string, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	resp, err := get(token.Token, this.config.PermissionsUrl+"/jwt/search/"+url.PathEscape(kind)+"/"+url.PathEscape(query)+"/"+right+"/"+limit+"/"+offset+"/"+url.PathEscape(orderfeature)+"/"+direction)
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

func (this *Lib) PermSelectIds(token auth.Token, kind string, right string, ids []string) (result []map[string]interface{}, err error) {
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(ids)
	if err != nil {
		return
	}
	resp, err := post(token.Token, this.config.PermissionsUrl+"/ids/select/"+url.PathEscape(kind)+"/"+right, "application/json", b)
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

func (this *Lib) PermSelectDeviceTypesByIdRead(token auth.Token, ids []string) (result []map[string]interface{}, err error) {
	return this.PermSelectIds(token, "device-types", "r", ids)
}

func (this *Lib) PermCheckDeviceAdmin(token auth.Token, ids []string) (result map[string]bool, err error) {
	return this.PermCheck(token, "devices", ids, "a")
}

func (this *Lib) PermCheckDeviceRead(token auth.Token, ids []string) (result map[string]bool, err error) {
	return this.PermCheck(token, "devices", ids, "r")
}

func (this *Lib) PermCheck(token auth.Token, kind string, ids []string, right string) (result map[string]bool, err error) {
	result = map[string]bool{}
	err = postJson(token.Token, this.config.PermissionsUrl+"/ids/check/"+url.PathEscape(kind)+"/"+right, ids, &result)
	return
}
