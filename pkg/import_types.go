/*
 * Copyright 2021 InfAI (CC SES)
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
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/auth"
	"net/http"
)

type ImportTypePermissionSearch struct {
	AspectFunctions    []string `json:"aspect_functions"`
	ContentAspectIds   []string `json:"content_aspect_ids"`
	Creator            string   `json:"creator"`
	DefaultRestart     bool     `json:"default_restart"`
	Description        string   `json:"description"`
	ContentFunctionIds []string `json:"content_function_ids"`
	Id                 string   `json:"id"`
	Image              string   `json:"image"`
	Name               string   `json:"name"`
	PermissionHolders  struct {
		AdminUsers   []string `json:"admin_users"`
		ExecuteUsers []string `json:"execute_users"`
		ReadUsers    []string `json:"read_users"`
		WriteUsers   []string `json:"write_users"`
	} `json:"permission_holders"`
	Permissions struct {
		A bool `json:"a"`
		R bool `json:"r"`
		W bool `json:"w"`
		X bool `json:"x"`
	} `json:"permissions"`
	Shared bool `json:"shared"`
}

func (this *Lib) GetImportTypesWithAspect(token auth.Token, aspectId string) (importTypes []ImportTypePermissionSearch, err error, code int) {
	node, err := this.GetAspectNodes([]string{aspectId}, token)
	if err != nil {
		return nil, err, http.StatusBadGateway
	}
	ids := append(node[0].AncestorIds, node[0].Id)

	err, code = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: "import-types",
		Find: &QueryFind{
			Filter: &Selection{
				Condition: ConditionConfig{
					Feature:   "features.content_aspect_ids",
					Operation: QueryAnyValueInFeatureOperation,
					Value:     ids,
				}}},
	}, &importTypes)
	return importTypes, err, code
}

func (this *Lib) GetImportTypes(token auth.Token) (importTypes []ImportTypePermissionSearch, err error, code int) {
	err, code = this.QueryPermissionsSearch(token.Token, QueryMessage{
		Resource: "import-types",
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit: 9999,
			},
		},
	}, &importTypes)
	return importTypes, err, code
}
