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

package lib

import (
	"errors"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
)

func (this *Lib) GetGatewaysHistory(jwt jwt_http_router.Jwt, duration string) (result []map[string]interface{}, err error) {
	result, err = this.PermListAllGateways(jwt, "r")
	if err != nil {
		log.Println("ERROR PermListAllGateways()", err)
		return result, err
	}
	result, err = this.CompleteGatewayHistory(jwt, duration, result)
	return
}

func (this *Lib) CompleteGatewayHistory(jwt jwt_http_router.Jwt, duration string, gateways []map[string]interface{}) (result []map[string]interface{}, err error) {
	ids := []string{}
	gatewayMap := map[string]map[string]interface{}{}
	for _, gateway := range gateways {
		id, ok := gateway["id"]
		if !ok {
			err = errors.New("unable to get gateway id")
			return
		}
		idStr, ok := id.(string)
		if !ok {
			err = errors.New("unable to cast gateway id to string")
			return
		}
		ids = append(ids, idStr)
		gatewayMap[idStr] = gateway
	}
	logStates, err := this.GetGatewayLogStates(jwt, ids)
	if err != nil {
		log.Println("ERROR completeGatewayList.GetGatewayLogStates()", err)
		return result, err
	}
	logHistory, err := this.GetGatewayLogHistory(jwt, ids, duration)
	if err != nil {
		log.Println("ERROR completeGatewayList.GetGatewayLogHistory()", err)
		return result, err
	}
	logEdges, err := this.GetLogedges(jwt, "gateway", ids, duration)
	if err != nil {
		log.Println("ERROR completeDeviceList.GetLogedges()", err)
		return result, err
	}
	for _, id := range ids {
		gateway := gatewayMap[id]
		logState, logExists := logStates[id]
		if !logExists {
			gateway["log_state"] = "unknown"
		} else {
			if logState {
				gateway["log_state"] = "connected"
			} else {
				gateway["log_state"] = "disconnected"
			}
		}
		gateway["log_history"] = logHistory[id]
		gateway["log_edge"] = logEdges[id]
		result = append(result, gateway)
	}
	return
}

func (this *Lib) ListGateways(jwt jwt_http_router.Jwt, limit string, offset string) (result []map[string]interface{}, err error) {
	gateways, err := this.PermListGateways(jwt, "r", limit, offset)
	if err != nil {
		log.Println("ERROR ListGateways.PermListGateways()", err)
		return result, err
	}
	return this.completeGatewayList(jwt, gateways)
}

func (this *Lib) ListAllGateways(jwt jwt_http_router.Jwt) (result []map[string]interface{}, err error) {
	return this.PermListAllGateways(jwt, "r")
}

func (this *Lib) ListGatewaysOrdered(jwt jwt_http_router.Jwt, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	gateways, err := this.PermListGatewaysOrdered(jwt, "r", limit, offset, orderfeature, direction)
	if err != nil {
		log.Println("ERROR ListGateways.PermListGateways()", err)
		return result, err
	}
	return this.completeGatewayList(jwt, gateways)
}

func (this *Lib) SearchGateways(jwt jwt_http_router.Jwt, query string, limit string, offset string) (result []map[string]interface{}, err error) {
	gateways, err := this.PermSearchGateways(jwt, query, "r", limit, offset)
	if err != nil {
		return result, err
	}
	return this.completeGatewayList(jwt, gateways)
}

func (this *Lib) SearchGatewaysOrdered(jwt jwt_http_router.Jwt, query string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	gateways, err := this.PermSearchGatewaysOrdered(jwt, query, "r", limit, offset, orderfeature, direction)
	if err != nil {
		return result, err
	}
	return this.completeGatewayList(jwt, gateways)
}

func (this *Lib) completeGatewayList(jwt jwt_http_router.Jwt, gateways []map[string]interface{}) (result []map[string]interface{}, err error) {
	ids := []string{}
	gatewayMap := map[string]map[string]interface{}{}
	for _, gateway := range gateways {
		id, ok := gateway["id"]
		if !ok {
			err = errors.New("unable to get gateway id")
			return
		}
		idStr, ok := id.(string)
		if !ok {
			err = errors.New("unable to cast gateway id to string")
			return
		}
		ids = append(ids, idStr)
		gatewayMap[idStr] = gateway
	}
	logStates, err := this.GetGatewayLogStates(jwt, ids)
	if err != nil {
		log.Println("ERROR completeGatewayList.GetGatewayLogStates()", err)
		return result, err
	}
	for _, id := range ids {
		gateway := gatewayMap[id]
		logState, logExists := logStates[id]
		if logExists {
			gateway["log_state"] = logState
		}
		//gateway["gateway_name"] = gateways[id]
		result = append(result, gateway)
	}
	return
}
