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
	"log"
	"net/http"

	"github.com/SmartEnergyPlatform/util/http/cors"
	"github.com/SmartEnergyPlatform/util/http/logger"
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

	return

}
