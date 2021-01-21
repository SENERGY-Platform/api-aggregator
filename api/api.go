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
	"github.com/SmartEnergyPlatform/api-aggregator/lib"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/SmartEnergyPlatform/util/http/cors"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"github.com/SmartEnergyPlatform/util/http/response"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Start(lib lib.Interface) {
	log.Println("start server on port: ", lib.Config().ServerPort)
	httpHandler := getRoutes(lib)
	corseHandler := cors.New(httpHandler)
	logger := logger.New(corseHandler, lib.Config().LogLevel)
	log.Println(http.ListenAndServe(":"+lib.Config().ServerPort, logger))
}

func getRoutes(lib lib.Interface) (router *jwt_http_router.Router) {
	router = jwt_http_router.New(jwt_http_router.JwtConfig{
		ForceUser: lib.Config().ForceUser,
		ForceAuth: lib.Config().ForceAuth,
	})

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
	router.GET("/devices", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
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

		orderfeature, direction := getSortParts(sort)
		result, err := lib.FindDevices(jwt, search, idList, intLimit, intOffset, orderfeature, direction, location)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if state != "" {
			result, err = lib.FilterDevicesByState(jwt, result, state)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if logDuration != "" {
			result, err = lib.CompleteDeviceHistory(jwt, logDuration, result)
		}

		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
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
	router.GET("/hubs", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		logDuration := r.URL.Query().Get("log")
		search := r.URL.Query().Get("search")
		sort := r.URL.Query().Get("sort")

		result := []map[string]interface{}{}
		var err error

		switch {
		case limit == "" && offset == "" && sort == "" && search == "":
			result, err = lib.ListAllGateways(jwt)
		case search != "" && sort == "":
			limit, offset = limitOffsetDefault(limit, offset)
			result, err = lib.SearchGateways(jwt, search, limit, offset)
		case search != "" && sort != "":
			orderfeature, direction := getSortParts(sort)
			limit, offset = limitOffsetDefault(limit, offset)
			result, err = lib.SearchGatewaysOrdered(jwt, search, limit, offset, orderfeature, direction)
		case search == "" && sort == "":
			limit, offset = limitOffsetDefault(limit, offset)
			result, err = lib.ListGateways(jwt, limit, offset)
		case search == "" && sort != "":
			limit, offset = limitOffsetDefault(limit, offset)
			orderfeature, direction := getSortParts(sort)
			result, err = lib.ListGatewaysOrdered(jwt, limit, offset, orderfeature, direction)
		}

		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if logDuration != "" {
			result, err = lib.CompleteGatewayHistory(jwt, logDuration, result)
		}

		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		response.To(res).Json(result)
	})

	//reads query parameter like https://docs.camunda.org/manual/7.5/reference/rest/deployment/get-query/
	router.GET("/processes", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		log.Println("DEBUG: ", r.URL.Query())
		result, err := lib.GetExtendedProcessList(jwt, r.URL.Query())
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/hubs/:id/devices", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		sort := request.URL.Query().Get("sort")
		limit := request.URL.Query().Get("limit")
		offset := request.URL.Query().Get("offset")
		state := request.URL.Query().Get("state")

		id := params.ByName("id")
		idList, err := lib.GetGatewayDevices(jwt, id)
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
			result, err = lib.CompleteDevicesOrdered(jwt, idList, limit, offset, orderfeature, direction)
		} else {
			result, err = lib.CompleteDevices(jwt, idList)
		}

		if state != "" {
			result, err = lib.FilterDevicesByState(jwt, result, state)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		response.To(writer).Json(result)
	})

	router.GET("/device-types/:id/devices", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
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

		idList, err := lib.GetDeviceTypeDevices(jwt, id, limit, offset, orderfeature, direction)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := lib.CompleteDevices(jwt, idList)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if state != "" {
			result, err = lib.FilterDevicesByState(jwt, result, state)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		response.To(writer).Json(result)
	})

	router.POST("/device-types-devices", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		state := request.URL.Query().Get("state")
		direction := request.URL.Query().Get("direction")
		orderfeature, direction := getSortParts("name.asc")

		type deviceTypesDevicesBody struct {
			Ids []string `json:"ids"`
		}
		var body deviceTypesDevicesBody
		err := json.NewDecoder(request.Body).Decode(&body)
		if err != nil {
			http.Error(writer, "unable to parse request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		idList := []string{}
		for _, id := range body.Ids {
			deviceIds, err := lib.GetDeviceTypeDevices(jwt, id, "-1", "-1", orderfeature, direction)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			idList = append(idList, deviceIds...)
		}

		result, err := lib.CompleteDevices(jwt, idList)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if state != "" {
			result, err = lib.FilterDevicesByState(jwt, result, state)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		result = lib.SortByName(result, true)

		response.To(writer).Json(result)
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
