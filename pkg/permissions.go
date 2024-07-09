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
	"github.com/SENERGY-Platform/api-aggregator/pkg/auth"
	"github.com/SENERGY-Platform/permission-search/lib/client"
	"github.com/SENERGY-Platform/permission-search/lib/model"
	"net/http"
	"strconv"
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
	return this.PermSelectIds(token, "devices", right, ids)
}

func (this *Lib) PermDeviceIdListOrdered(token auth.Token, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	return this.PermIdListOrdered(token, "devices", ids, right, limit, offset, orderfeature, direction)
}

func (this *Lib) PermIdListOrdered(token auth.Token, kind string, ids []string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		ListIds: &QueryListIds{
			QueryListCommons: QueryListCommons{
				Limit:    len(ids),
				Rights:   right,
				SortBy:   orderfeature,
				SortDesc: direction == "desc",
			},
			Ids: ids,
		},
	}, &result)
	return
}

func (this *Lib) PermListAllGateways(token auth.Token, right string) (result []map[string]interface{}, err error) {
	return this.PermListAll(token, "hubs", right)
}

func (this *Lib) PermListAll(token auth.Token, kind string, right string) (result []map[string]interface{}, err error) {
	return this.permissionsearch.List(token.Jwt(), kind, model.ListOptions{
		QueryListCommons: model.QueryListCommons{
			Limit:  9999,
			Rights: right,
		},
	})
}

func (this *Lib) PermList(token auth.Token, kind string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:  l,
				Offset: o,
				Rights: right,
			},
		},
	}, &result)
	return
}

func (this *Lib) PermListOrdered(token auth.Token, kind string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:    l,
				Offset:   o,
				Rights:   right,
				SortBy:   orderfeature,
				SortDesc: direction == "desc",
			},
		},
	}, &result)
	return
}

func (this *Lib) PermSearch(token auth.Token, kind string, query string, right string, limit string, offset string) (result []map[string]interface{}, err error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:  l,
				Offset: o,
				Rights: right,
			},
			Search: query,
		},
	}, &result)
	return
}

func (this *Lib) PermSearchOrdered(token auth.Token, kind string, query string, right string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:    l,
				Offset:   o,
				Rights:   right,
				SortBy:   orderfeature,
				SortDesc: direction == "desc",
			},
			Search: query,
		},
	}, &result)
	return
}

func (this *Lib) PermSelectIds(token auth.Token, kind string, right string, ids []string) (result []map[string]interface{}, err error) {
	err, _ = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: kind,
		ListIds: &QueryListIds{
			QueryListCommons: QueryListCommons{
				Limit:  len(ids),
				Rights: right,
			},
			Ids: ids,
		},
	}, &result)
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
	result, _, err = client.Query[map[string]bool](this.permissionsearch, token.String(), model.QueryMessage{
		Resource: kind,
		CheckIds: &model.QueryCheckIds{
			Ids:    ids,
			Rights: right,
		},
	})
	return
}

func QueryPermissionSearchFindAll[T any](lib *Lib, token string, query QueryMessage, idGetter func(e T) string) (result []T, err error, code int) {
	var after *ListAfter
	temp := []T{}
	limit := 9999
	for {
		if query.Find != nil {
			query.Find.QueryListCommons.Limit = limit
			query.Find.QueryListCommons.Offset = 0
			query.Find.QueryListCommons.After = after
		}
		err, code = lib.QueryPermissionsSearch(token, query, &temp)
		if err != nil {
			return result, err, code
		}
		result = append(result, temp...)
		if len(temp) < limit {
			return result, nil, http.StatusOK
		}
		if len(temp) > 0 {
			after = &ListAfter{
				Id: idGetter(temp[len(temp)-1]),
			}
		}
		temp = []T{}
	}
}

func (this *Lib) QueryPermissionsSearch(token string, query QueryMessage, result interface{}) (err error, code int) {
	temp, code, err := this.permissionsearch.Query(token, query)
	if err != nil {
		return err, code
	}
	b, err := json.Marshal(temp)
	if err != nil {
		return err, 500
	}
	err = json.Unmarshal(b, result)
	if err != nil {
		return err, 500
	}
	return nil, 200
}

type QueryMessage = client.QueryMessage
type QueryFind = client.QueryFind

type QueryListIds = client.QueryListIds

type QueryCheckIds = client.QueryCheckIds

type QueryListCommons = client.QueryListCommons

type ListAfter = client.ListAfter

type QueryOperationType = client.QueryOperationType

const (
	QueryEqualOperation             = client.QueryEqualOperation
	QueryUnequalOperation           = client.QueryUnequalOperation
	QueryAnyValueInFeatureOperation = client.QueryAnyValueInFeatureOperation
)

type ConditionConfig = client.ConditionConfig

type Selection = client.Selection
