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
	"encoding/json"
	"errors"
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/auth"
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/model"
	"net/http"
	"net/url"
)

type Function struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ConceptId   string `json:"concept_id"`
	RdfType     string `json:"rdf_type"`
}

func (this *Lib) GetMeasuringFunctionsForAspect(token auth.Token, aspectId string) (functions []Function, err error, code int) {
	resp, err := get(token.Token, this.config.SemanticRepoUrl+"/aspects/"+url.PathEscape(aspectId)+"/measuring-functions")
	if err != nil {
		return nil, err, http.StatusBadGateway
	}
	if resp.StatusCode > 299 {
		return nil, errors.New("unexpected status code from semantic-repo"), resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(&functions)
	return functions, err, resp.StatusCode
}

func (this *Lib) GetMeasuringFunctions(token auth.Token, functionIds []string) (functions []Function, err error, code int) {
	m, err := this.PermSelectIds(token, "functions", "r", functionIds)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	err = json.Unmarshal(b, &functions)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return functions, nil, http.StatusOK
}

type CharacteristicsWrapper struct {
	Raw model.Characteristic `json:"raw"`
}

func (this *Lib) GetNestedFunctionInfos(token auth.Token) (result []model.FunctionInfo, err error) {
	concepts := []model.Concept{}
	err, _ = this.QueryPermissionsSearch(token.Jwt(), QueryMessage{
		Resource: "concepts",
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:  9999,
				Offset: 0,
				Rights: "r",
			},
		},
	}, &concepts)
	if err != nil {
		return result, err
	}
	characteristics := []CharacteristicsWrapper{}
	err, _ = this.QueryPermissionsSearch(token.Jwt(), QueryMessage{
		Resource: "characteristics",
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:  9999,
				Offset: 0,
				Rights: "r",
			},
		},
	}, &characteristics)
	if err != nil {
		return result, err
	}
	err, _ = this.QueryPermissionsSearch(token.Jwt(), QueryMessage{
		Resource: "functions",
		Find: &QueryFind{
			QueryListCommons: QueryListCommons{
				Limit:  9999,
				Offset: 0,
				Rights: "r",
			},
		},
	}, &result)
	if err != nil {
		return result, err
	}

	conceptIndex := map[string]model.Concept{}
	for _, c := range concepts {
		conceptIndex[c.Id] = c
	}

	characteristicsIndex := map[string]model.Characteristic{}
	for _, c := range characteristics {
		characteristicsIndex[c.Raw.Id] = c.Raw
	}

	for i, f := range result {
		concept, ok := conceptIndex[f.ConceptId]
		if ok {
			f.Concept = model.ConceptInfo{Concept: concept, BaseCharacteristic: characteristicsIndex[concept.BaseCharacteristicId]}
		}
		result[i] = f
	}
	return result, nil
}
