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
	"github.com/SENERGY-Platform/api-aggregator/pkg/auth"
	"github.com/SENERGY-Platform/api-aggregator/pkg/model"
	"github.com/SENERGY-Platform/device-repository/lib/client"
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
	resp, err := get(token.Token, this.config.IotUrl+"/aspects/"+url.PathEscape(aspectId)+"/measuring-functions")
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
	temp, _, err, _ := this.deviceRepo.ListFunctions(client.FunctionListOptions{
		Ids:    functionIds,
		Limit:  int64(len(functionIds)),
		Offset: 0,
		SortBy: "name.asc",
	})
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	for _, function := range temp {
		functions = append(functions, Function{
			Id:          function.Id,
			Name:        function.Name,
			Description: function.Description,
			ConceptId:   function.ConceptId,
			RdfType:     function.RdfType,
		})
	}
	return functions, nil, http.StatusOK
}

type CharacteristicsWrapper struct {
	Raw model.Characteristic `json:"raw"`
}
