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
	"github.com/SmartEnergyPlatform/api-aggregator/lib"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/SmartEnergyPlatform/util/http/cors"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"github.com/SmartEnergyPlatform/util/http/response"
	"log"
	"net/http"
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
		ForceUser: lib.Config().ForceUser == "true",
		ForceAuth: lib.Config().ForceAuth == "true",
	})

	/*
		query-parameter:
			mutual exclusive:
				usertag {string} 	filters result by user-tag
				tag		{string} 	filters result by tag
				search  {string}	filters by partial text match
				ids		{string,string...} returns by ids
			optional:
				log		{string}	influxdb duration (for example 4h) https://docs.influxdata.com/influxdb/v1.7/query_language/spec/#durations
				state 	{string} 	filters result by device state
				sort 	{string} 	sorts result by filed; if data-source does not support sorting, it will be performed locally
										name | name.asc | name.desc
				limit 	{int} 		may default to 100; no effect when used with usertag and tag
				offset 	{int}		may default to 0; no effect when used with usertag and tag
	*/
	router.GET("/devices", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		usertag := r.URL.Query().Get("usertag")
		tag := r.URL.Query().Get("tag")
		state := r.URL.Query().Get("state")
		sort := r.URL.Query().Get("sort")
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		search := r.URL.Query().Get("search")
		ids := r.URL.Query().Get("ids")
		logDuration := r.URL.Query().Get("log")

		sorted := false
		result := []map[string]interface{}{}
		var err error
		switch {
		case ids != "":
			idList := strings.Split(strings.Replace(ids, " ", "", -1), ",")
			if limit != "" || offset != "" {
				limit, offset = limitOffsetDefault(limit, offset)
				if sort == "" {
					sort = "name"
				}
				orderfeature, direction := getSortParts(sort)
				sorted = true
				result, err = lib.CompleteDevicesOrdered(jwt, idList, limit, offset, orderfeature, direction)
			} else {
				result, err = lib.CompleteDevices(jwt, idList)
			}
		case usertag != "":
			if limit != "" || offset != "" {
				limit, offset = limitOffsetDefault(limit, offset)
				if sort == "" {
					sort = "name"
				}
				orderfeature, direction := getSortParts(sort)
				sorted = true
				result, err = lib.ListOrderedDevicesByUserTag(jwt, usertag, limit, offset, orderfeature, direction)
			} else {
				result, err = lib.ListDevicesByUserTag(jwt, usertag)
			}
		case tag != "":
			if limit != "" || offset != "" {
				limit, offset = limitOffsetDefault(limit, offset)
				if sort == "" {
					sort = "name"
				}
				orderfeature, direction := getSortParts(sort)
				sorted = true
				result, err = lib.ListOrderdDevicesByTag(jwt, tag, limit, offset, orderfeature, direction)
			} else {
				result, err = lib.ListDevicesByTag(jwt, tag)
			}
		case search != "":
			limit, offset = limitOffsetDefault(limit, offset)
			if sort != "" {
				orderfeature, direction := getSortParts(sort)
				sorted = true
				result, err = lib.SearchDevicesOrdered(jwt, search, limit, offset, orderfeature, direction)
			} else {
				limit, offset = limitOffsetDefault(limit, offset)
				result, err = lib.SearchDevices(jwt, search, limit, offset)
			}
		default:
			if limit != "" || offset != "" {
				limit, offset = limitOffsetDefault(limit, offset)
				if sort != "" {
					orderfeature, direction := getSortParts(sort)
					sorted = true
					result, err = lib.ListDevicesOrdered(jwt, limit, offset, orderfeature, direction)
				} else {
					limit, offset = limitOffsetDefault(limit, offset)
					result, err = lib.ListDevices(jwt, limit, offset)
				}
			} else {
				result, err = lib.ListAllDevices(jwt)
			}
		}
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

		if !sorted && sort != "" {
			switch sort {
			case "name.asc":
				result = lib.SortByName(result, true)
			case "name.desc":
				result = lib.SortByName(result, false)
			case "name":
				result = lib.SortByName(result, true)
			default:
				http.Error(res, "unable to interpret sort "+sort, http.StatusBadRequest)
				return
			}
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
	router.GET("/gateways", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
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
