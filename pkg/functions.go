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
