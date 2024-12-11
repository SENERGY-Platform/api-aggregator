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
	"encoding/json"
	"fmt"
	"github.com/SENERGY-Platform/service-commons/pkg/accesslog"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/SENERGY-Platform/api-aggregator/pkg"
	"github.com/SENERGY-Platform/api-aggregator/pkg/api/util"
	"github.com/SENERGY-Platform/api-aggregator/pkg/auth"
	"github.com/SENERGY-Platform/api-aggregator/pkg/model"
	"github.com/julienschmidt/httprouter"
)

func Start(lib pkg.Interface) {
	log.Println("start server on port: ", lib.Config().ServerPort)
	httpHandler := getRoutes(lib)
	corseHandler := util.NewCors(httpHandler)
	logger := accesslog.New(corseHandler)
	log.Println(http.ListenAndServe(":"+lib.Config().ServerPort, logger))
}

func getRoutes(lib pkg.Interface) (router *httprouter.Router) {
	router = httprouter.New()

	//returns device-classes used by user devices
	router.GET("/device-class-uses", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := lib.GetDeviceClassUses(token)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.GET("/devices", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		logDuration := r.URL.Query().Get("log")

		limit, offset = limitOffsetDefault(limit, offset)

		intLimit, err := strconv.Atoi(limit)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, "limit is not a number: "+err.Error(), http.StatusBadRequest)
			return
		}

		intOffset, err := strconv.Atoi(offset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, "offset is not a number: "+err.Error(), http.StatusBadRequest)
			return
		}

		token, err := auth.GetParsedToken(r)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := lib.FindDevices(token, intLimit, intOffset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if logDuration != "" {
			result, err = lib.CompleteDeviceHistory(token, logDuration, result)
		}

		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(res).Encode(result)
	})

	/*
		query-parameter:
			optional:
				limit 	{int} 		may default to 100
				offset 	{int}		may default to 0
				log		{string}	influxdb duration (for example 4h) https://docs.influxdata.com/influxdb/v1.7/query_language/spec/#durations
	*/
	router.GET("/hubs", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		logDuration := r.URL.Query().Get("log")

		token, err := auth.GetParsedToken(r)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		result := []map[string]interface{}{}

		if limit == "" && offset == "" {
			result, err = lib.ListAllGateways(token)
		} else {
			intLimit, err := strconv.ParseInt(limit, 10, 64)
			if err != nil {
				err = fmt.Errorf("limit is not a number: %w", err)
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}

			intOffset, err := strconv.ParseInt(offset, 10, 64)
			if err != nil {
				err = fmt.Errorf("offset is not a number: %w", err)
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			result, err = lib.ListGateways(token, intLimit, intOffset)
		}
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if logDuration != "" {
			result, err = lib.CompleteGatewayHistory(token, logDuration, result)
		}

		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(res).Encode(result)
	})

	//reads query parameter like https://docs.camunda.org/manual/7.5/reference/rest/deployment/get-query/
	router.GET("/processes", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Println("DEBUG: ", r.URL.Query())
		token, err := auth.GetParsedToken(r)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := lib.GetExtendedProcessList(token, r.URL.Query())
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(res).Encode(result)
	})

	router.GET("/aspects/:id/measuring-functions", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		// Get from semantic
		id := params.ByName("id")
		functions, err, code := lib.GetMeasuringFunctionsForAspect(token, id)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(writer, err.Error(), code)
			return
		}

		// Get from Permsearch (import-types)
		node, err := lib.GetAspectNodes([]string{id}, token)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(writer, err.Error(), http.StatusBadGateway)
			return
		}
		if len(node) != 1 {
			log.Println("ERROR: ", err)
			http.Error(writer, "unexpected length of reponse", http.StatusBadGateway)
			return
		}
		ids := append(node[0].DescendentIds, node[0].Id)

		importTypes, err, code := lib.GetImportTypesWithAspect(token, ids)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(writer, err.Error(), code)
			return
		}

		additionalFunctionIds := []string{}
		for _, importType := range importTypes {
			for _, c := range importType.Criteria {
				if isInSlice(ids, c.AspectId) && !isInSlice(additionalFunctionIds, c.FunctionId) && !isFunctionLoaded(functions, c.FunctionId) {
					additionalFunctionIds = append(additionalFunctionIds, c.FunctionId)
				}
			}
		}
		additionalFunctions, err, code := lib.GetMeasuringFunctions(token, additionalFunctionIds)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(writer, err.Error(), code)
			return
		}
		functions = append(functions, additionalFunctions...)

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(functions)
	})

	router.GET("/aspect-nodes", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		function := request.URL.Query().Get("function")
		if function != "measuring-function" {
			http.Error(writer, "May only use function=measuring-function", http.StatusBadRequest)
			return
		}

		var result []model.AspectNode
		var err error

		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		// Get for devices, ancestors already included
		result, err = lib.GetAspectNodesWithMeasuringFunction(token)

		aspectIds := []string{}
		for _, r := range result {
			aspectIds = append(aspectIds, r.Id)
		}

		// Get import types and prepare loading additional nodes
		importTypes, err, code := lib.GetImportTypes(token)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		additionalAspectIds := []string{}
		for _, t := range importTypes {
			for _, c := range t.Criteria {
				if !isInSlice(aspectIds, c.AspectId) {
					additionalAspectIds = append(additionalAspectIds, c.AspectId)
					aspectIds = append(aspectIds, c.AspectId)
				}
			}
		}

		// Get additional nodes if needed
		if len(additionalAspectIds) > 0 {
			importTypeNodes, err := lib.GetAspectNodes(additionalAspectIds, token)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(writer, err.Error(), http.StatusBadGateway)
				return
			}
			result = append(result, importTypeNodes...)

			// Check for ancestors of additional nodes and prepare loading those
			additionalAspectIds = []string{}
			for _, node := range importTypeNodes {
				for _, ancestorId := range node.AncestorIds {
					if !isInSlice(aspectIds, ancestorId) {
						additionalAspectIds = append(additionalAspectIds, ancestorId)
						aspectIds = append(aspectIds, ancestorId)
					}
				}
			}

			// Load ancestors if needed
			if len(additionalAspectIds) > 0 {
				additionalNodes, err := lib.GetAspectNodes(additionalAspectIds, token)
				if err != nil {
					log.Println("ERROR: ", err)
					http.Error(writer, err.Error(), http.StatusBadGateway)
					return
				}
				result = append(result, additionalNodes...)
			}
		}

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
		return
	})

	return

}

func limitOffsetDefault(limit, offset string) (string, string) {
	if limit == "" {
		limit = "100"
	}
	if offset == "" {
		offset = "0"
	}
	return limit, offset
}

func getSortParts(sort string) (orderfeature string, direction string) {
	orderfeature = strings.TrimSuffix(strings.TrimSuffix(sort, ".desc"), ".asc")
	direction = "asc"
	if strings.HasSuffix(sort, ".desc") {
		direction = "desc"
	}
	return
}

func isFunctionLoaded(functions []pkg.Function, id string) bool {
	for i := range functions {
		if functions[i].Id == id {
			return true
		}
	}
	return false
}

func isInSlice(list []string, element string) bool {
	for i := range list {
		if list[i] == element {
			return true
		}
	}
	return false
}
