/*
 * Copyright 2018 InfAI (CC SES)
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

package deprecated

import (
	"github.com/SmartEnergyPlatform/api-aggregator/lib"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"

	"encoding/json"

	"github.com/SmartEnergyPlatform/util/http/cors"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"github.com/SmartEnergyPlatform/util/http/response"
)

func Start(lib lib.Interface) {
	log.Println("start server on port: ", lib.Config().ServerPort)
	httpHandler := GetRoutes(lib)
	corseHandler := cors.New(httpHandler)
	logger := logger.New(corseHandler, lib.Config().LogLevel)
	log.Println(http.ListenAndServe(":"+lib.Config().ServerPort, logger))
}

func GetRoutes(lib lib.Interface) (router *jwt_http_router.Router) {
	router = jwt_http_router.New(jwt_http_router.JwtConfig{
		ForceUser: lib.Config().ForceUser == "true",
		ForceAuth: lib.Config().ForceAuth == "true",
	})

	router.GET("/filter/devices/state/:value/usertag/:tag/orderby/name/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		tag := ps.ByName("tag")
		state := ps.ByName("value")
		result, err := lib.ListDevicesByUserTag(jwt, tag)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result, err = lib.FilterDevicesByState(jwt, result, state)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result = lib.SortByName(result, true)
		response.To(res).Json(result)
	})

	router.GET("/filter/devices/state/:value/usertag/:tag/orderby/name/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		tag := ps.ByName("tag")
		state := ps.ByName("value")
		result, err := lib.ListDevicesByUserTag(jwt, tag)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result, err = lib.FilterDevicesByState(jwt, result, state)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result = lib.SortByName(result, false)
		response.To(res).Json(result)
	})

	router.GET("/filter/devices/state/:value/tag/:tag/orderby/name/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		tag := ps.ByName("tag")
		state := ps.ByName("value")
		result, err := lib.ListDevicesByTag(jwt, tag)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result, err = lib.FilterDevicesByState(jwt, result, state)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result = lib.SortByName(result, true)
		response.To(res).Json(result)
	})

	router.GET("/filter/devices/state/:value/tag/:tag/orderby/name/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		tag := ps.ByName("tag")
		state := ps.ByName("value")
		result, err := lib.ListDevicesByTag(jwt, tag)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result, err = lib.FilterDevicesByState(jwt, result, state)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result = lib.SortByName(result, false)
		response.To(res).Json(result)
	})

	router.GET("/filter/devices/state/:value/name/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		result, err := lib.GetConnectionFilteredDevicesOrder(jwt, value, true)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/filter/devices/state/:value/name/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		result, err := lib.GetConnectionFilteredDevicesOrder(jwt, value, false)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	//old maintenance helper
	/*
		router.GET("/filter/devices/state/:value/search/:searchtext/name/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
			value := ps.ByName("value")
			searchText := ps.ByName("searchtext")
			result, err := lib.GetConnectionFilteredDevicesSearchOrder(jwt, value, searchText, true)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			response.To(res).Json(result)
		})

		router.GET("/filter/devices/state/:value/search/:searchtext/name/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
			value := ps.ByName("value")
			searchText := ps.ByName("searchtext")
			result, err := lib.GetConnectionFilteredDevicesSearchOrder(jwt, value, searchText, false)
			if err != nil {
				log.Println("ERROR: ", err)
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			response.To(res).Json(result)
		})
	*/

	router.GET("/filter/devices/state/:value", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		result, err := lib.GetConnectionFilteredDevices(jwt, value)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/list/devices/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		result, err := lib.ListDevices(jwt, limit, offset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/search/devices/:query/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		query := ps.ByName("query")
		result, err := lib.SearchDevices(jwt, query, limit, offset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/list/devices/:limit/:offset/:orderfeature/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderfeature := ps.ByName("orderfeature")
		result, err := lib.ListDevicesOrdered(jwt, limit, offset, orderfeature, "asc")
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/search/devices/:query/:limit/:offset/:orderfeature/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		query := ps.ByName("query")
		orderfeature := ps.ByName("orderfeature")
		result, err := lib.SearchDevicesOrdered(jwt, query, limit, offset, orderfeature, "asc")
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/list/devices/:limit/:offset/:orderfeature/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderfeature := ps.ByName("orderfeature")
		result, err := lib.ListDevicesOrdered(jwt, limit, offset, orderfeature, "desc")
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/search/devices/:query/:limit/:offset/:orderfeature/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		query := ps.ByName("query")
		orderfeature := ps.ByName("orderfeature")
		result, err := lib.SearchDevicesOrdered(jwt, query, limit, offset, orderfeature, "desc")
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/select/devices/tag/:value", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		result, err := lib.ListDevicesByTag(jwt, value)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/select/devices/tag/:value/:limit/:offset/:orderfeature/:direction", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		orderfeature := ps.ByName("orderfeature")
		direction := ps.ByName("direction")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		result, err := lib.ListOrderdDevicesByTag(jwt, value, limit, offset, orderfeature, direction)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/select/devices/usertag/:value", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		result, err := lib.ListDevicesByUserTag(jwt, value)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/select/devices/usertag/:value/:limit/:offset/:orderfeature/:direction", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		orderfeature := ps.ByName("orderfeature")
		direction := ps.ByName("direction")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		result, err := lib.ListOrderedDevicesByUserTag(jwt, value, limit, offset, orderfeature, direction)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.POST("/select/devices/ids", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		ids := []string{}
		err := json.NewDecoder(r.Body).Decode(&ids)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := lib.CompleteDevices(jwt, ids)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.POST("/select/devices/ids/:limit/:offset/:orderfeature/:direction", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		orderfeature := ps.ByName("orderfeature")
		direction := ps.ByName("direction")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		ids := []string{}
		err := json.NewDecoder(r.Body).Decode(&ids)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := lib.CompleteDevicesOrdered(jwt, ids, limit, offset, orderfeature, direction)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/history/devices/:duration", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		duration := ps.ByName("duration")
		result, err := lib.GetDevicesHistory(jwt, duration)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/history/gateways/:duration", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		duration := ps.ByName("duration")
		result, err := lib.GetGatewaysHistory(jwt, duration)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/list/gateways/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		result, err := lib.ListGateways(jwt, limit, offset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/search/gateways/:query/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		query := ps.ByName("query")
		result, err := lib.SearchGateways(jwt, query, limit, offset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/list/gateways/:limit/:offset/:orderfeature/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderfeature := ps.ByName("orderfeature")
		result, err := lib.ListGatewaysOrdered(jwt, limit, offset, orderfeature, "asc")
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/search/gateways/:query/:limit/:offset/:orderfeature/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		query := ps.ByName("query")
		orderfeature := ps.ByName("orderfeature")
		result, err := lib.SearchGatewaysOrdered(jwt, query, limit, offset, orderfeature, "asc")
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/list/gateways/:limit/:offset/:orderfeature/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderfeature := ps.ByName("orderfeature")
		result, err := lib.ListGatewaysOrdered(jwt, limit, offset, orderfeature, "desc")
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/search/gateways/:query/:limit/:offset/:orderfeature/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		query := ps.ByName("query")
		orderfeature := ps.ByName("orderfeature")
		result, err := lib.SearchGatewaysOrdered(jwt, query, limit, offset, orderfeature, "desc")
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
