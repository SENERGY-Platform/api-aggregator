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

package main

import (
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"

	"encoding/json"

	"github.com/SmartEnergyPlatform/util/http/cors"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"github.com/SmartEnergyPlatform/util/http/response"
)

func StartApi() {
	log.Println("start server on port: ", Config.ServerPort)
	httpHandler := getRoutes()
	corseHandler := cors.New(httpHandler)
	logger := logger.New(corseHandler, Config.LogLevel)
	log.Println(http.ListenAndServe(":"+Config.ServerPort, logger))
}

func getRoutes() (router *jwt_http_router.Router) {
	router = jwt_http_router.New(jwt_http_router.JwtConfig{
		ForceUser: Config.ForceUser == "true",
		ForceAuth: Config.ForceAuth == "true",
	})

	router.GET("/filter/devices/state/:value/usertag/:tag", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		tag := ps.ByName("tag")
		state := ps.ByName("value")
		result, err := ListDevicesByUserTag(jwt, tag)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result, err = FilterDevicesByState(jwt, result, state)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/filter/devices/state/:value/tag/:tag", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		tag := ps.ByName("tag")
		state := ps.ByName("value")
		result, err := ListDevicesByTag(jwt, tag)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result, err = FilterDevicesByState(jwt, result, state)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})


	router.GET("/filter/devices/state/:value/name/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		result, err := GetConnectionFilteredDevicesOrder(jwt, value, true)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/filter/devices/state/:value/name/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		result, err := GetConnectionFilteredDevicesOrder(jwt, value, false)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/filter/devices/state/:value/search/:searchtext/name/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		searchText := ps.ByName("searchtext")
		result, err := GetConnectionFilteredDevicesSearchOrder(jwt, value, searchText, true)
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
		result, err := GetConnectionFilteredDevicesSearchOrder(jwt, value, searchText, false)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})


	router.GET("/filter/devices/state/:value", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		result, err := GetConnectionFilteredDevices(jwt, value)
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
		result, err := ListDevices(jwt, limit, offset)
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
		result, err := SearchDevices(jwt, query, limit, offset)
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
		result, err := ListDevicesOrdered(jwt, limit, offset, orderfeature, "asc")
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
		result, err := SearchDevicesOrdered(jwt, query, limit, offset, orderfeature, "asc")
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
		result, err := ListDevicesOrdered(jwt, limit, offset, orderfeature, "desc")
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
		result, err := SearchDevicesOrdered(jwt, query, limit, offset, orderfeature, "desc")
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/select/devices/tag/:value", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		result, err := ListDevicesByTag(jwt, value)
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
		result, err := ListOrderdDevicesByTag(jwt, value, limit, offset, orderfeature, direction)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})


	router.GET("/select/devices/usertag/:value", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		value := ps.ByName("value")
		result, err := ListDevicesByUserTag(jwt, value)
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
		result, err := ListOrderedDevicesByUserTag(jwt, value, limit, offset, orderfeature, direction)
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
		result, err := CompleteDevices(jwt, ids)
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
		result, err := CompleteDevicesOrdered(jwt, ids, limit, offset, orderfeature, direction)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/history/devices/:duration", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		duration := ps.ByName("duration")
		result, err := GetDevicesHistory(jwt, duration)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/history/gateways/:duration", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		duration := ps.ByName("duration")
		result, err := GetGatewaysHistory(jwt, duration)
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
		result, err := ListGateways(jwt, limit, offset)
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
		result, err := SearchGateways(jwt, query, limit, offset)
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
		result, err := ListGatewaysOrdered(jwt, limit, offset, orderfeature, "asc")
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
		result, err := SearchGatewaysOrdered(jwt, query, limit, offset, orderfeature, "asc")
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
		result, err := ListGatewaysOrdered(jwt, limit, offset, orderfeature, "desc")
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
		result, err := SearchGatewaysOrdered(jwt, query, limit, offset, orderfeature, "desc")
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	return
}
