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
	"github.com/SmartEnergyPlatform/api-aggregator/pkg"
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/api/util"
	"github.com/SmartEnergyPlatform/api-aggregator/pkg/auth"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Start(lib pkg.Interface) {
	log.Println("start server on port: ", lib.Config().ServerPort)
	httpHandler := getRoutes(lib)
	corseHandler := util.NewCors(httpHandler)
	logger := util.NewLogger(corseHandler)
	log.Println(http.ListenAndServe(":"+lib.Config().ServerPort, logger))
}

func getRoutes(lib pkg.Interface) (router *httprouter.Router) {
	router = httprouter.New()

	/*
		query-parameter:
				search  {string}	filters by partial text match
				ids		{string,string...} returns by ids
				location {string}	id of location in which the found devices must be located
				log		{string}	influxdb duration (for example 4h) https://docs.influxdata.com/influxdb/v1.7/query_language/spec/#durations
				state 	{string} 	filters result by device state
				sort 	{string} 	sorts result by filed; if data-source does not support sorting, it will be performed locally
										name | name.asc | name.desc
				limit 	{int} 		may default to 100
				offset 	{int}		may default to 0
	*/
	router.GET("/devices", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		state := r.URL.Query().Get("state")
		sort := r.URL.Query().Get("sort")
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		search := r.URL.Query().Get("search")
		ids := r.URL.Query().Get("ids")
		location := r.URL.Query().Get("location")
		logDuration := r.URL.Query().Get("log")

		idList := []string{}
		if ids != "" {
			idList = strings.Split(strings.Replace(ids, " ", "", -1), ",")
		}
		limit, offset = limitOffsetDefault(limit, offset)
		if sort == "" {
			sort = "name"
		}

		intLimit, err := strconv.Atoi(limit)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, "limit is not a number: "+err.Error(), http.StatusBadRequest)
			return
			return
		}

		intOffset, err := strconv.Atoi(offset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, "offset is not a number: "+err.Error(), http.StatusBadRequest)
			return
			return
		}

		token, err := auth.GetParsedToken(r)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		orderfeature, direction := getSortParts(sort)
		result, err := lib.FindDevices(token, search, idList, intLimit, intOffset, orderfeature, direction, location)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if state != "" {
			result, err = lib.FilterDevicesByState(token, result, state)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
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
				search  {string}	filters by partial text match
				limit 	{int} 		may default to 100
				offset 	{int}		may default to 0
				log		{string}	influxdb duration (for example 4h) https://docs.influxdata.com/influxdb/v1.7/query_language/spec/#durations
				sort 	{string} 	sorts result by filed
										name | name.asc | name.desc
	*/
	router.GET("/hubs", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		logDuration := r.URL.Query().Get("log")
		search := r.URL.Query().Get("search")
		sort := r.URL.Query().Get("sort")

		token, err := auth.GetParsedToken(r)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		result := []map[string]interface{}{}

		switch {
		case limit == "" && offset == "" && sort == "" && search == "":
			result, err = lib.ListAllGateways(token)
		case search != "" && sort == "":
			limit, offset = limitOffsetDefault(limit, offset)
			result, err = lib.SearchGateways(token, search, limit, offset)
		case search != "" && sort != "":
			orderfeature, direction := getSortParts(sort)
			limit, offset = limitOffsetDefault(limit, offset)
			result, err = lib.SearchGatewaysOrdered(token, search, limit, offset, orderfeature, direction)
		case search == "" && sort == "":
			limit, offset = limitOffsetDefault(limit, offset)
			result, err = lib.ListGateways(token, limit, offset)
		case search == "" && sort != "":
			limit, offset = limitOffsetDefault(limit, offset)
			orderfeature, direction := getSortParts(sort)
			result, err = lib.ListGatewaysOrdered(token, limit, offset, orderfeature, direction)
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

	router.GET("/hubs/:id/devices", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		sort := request.URL.Query().Get("sort")
		limit := request.URL.Query().Get("limit")
		offset := request.URL.Query().Get("offset")
		state := request.URL.Query().Get("state")

		id := params.ByName("id")

		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		idList, err := lib.GetGatewayDevices(token, id)
		var result []map[string]interface{}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if limit != "" || offset != "" || sort != "" {
			limit, offset = limitOffsetDefault(limit, offset)
			if sort == "" {
				sort = "name"
			}
			orderfeature, direction := getSortParts(sort)
			result, err = lib.CompleteDevicesOrdered(token, idList, limit, offset, orderfeature, direction)
		} else {
			result, err = lib.CompleteDevices(token, idList)
		}

		if state != "" {
			result, err = lib.FilterDevicesByState(token, result, state)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.GET("/device-types/:id/devices", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		sort := request.URL.Query().Get("sort")
		limit := request.URL.Query().Get("limit")
		offset := request.URL.Query().Get("offset")
		state := request.URL.Query().Get("state")

		id := params.ByName("id")
		limit, offset = limitOffsetDefault(limit, offset)

		if sort == "" {
			sort = "name"
		}
		orderfeature, direction := getSortParts(sort)

		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		idList, err := lib.GetDeviceTypeDevices(token, id, limit, offset, orderfeature, direction)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := lib.CompleteDevices(token, idList)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if state != "" {
			result, err = lib.FilterDevicesByState(token, result, state)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.POST("/device-types-devices", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		state := request.URL.Query().Get("state")
		direction := request.URL.Query().Get("direction")
		orderfeature, direction := getSortParts("name.asc")

		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		type deviceTypesDevicesBody struct {
			Ids []string `json:"ids"`
		}
		var body deviceTypesDevicesBody
		err = json.NewDecoder(request.Body).Decode(&body)
		if err != nil {
			http.Error(writer, "unable to parse request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		idList := []string{}
		for _, id := range body.Ids {
			deviceIds, err := lib.GetDeviceTypeDevices(token, id, "-1", "-1", orderfeature, direction)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			idList = append(idList, deviceIds...)
		}

		result, err := lib.CompleteDevices(token, idList)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if state != "" {
			result, err = lib.FilterDevicesByState(token, result, state)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		result = lib.SortByName(result, true)

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
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
	parts := strings.Split(sort, ".")
	orderfeature = parts[0]
	direction = "asc"
	if len(parts) > 1 {
		direction = parts[1]
	}
	return
}
