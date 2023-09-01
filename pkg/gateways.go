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

package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/api-aggregator/pkg/auth"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
)

func (this *Lib) GetGatewaysHistory(token auth.Token, duration string) (result []map[string]interface{}, err error) {
	result, err = this.PermListAllGateways(token, "r")
	if err != nil {
		log.Println("ERROR PermListAllGateways()", err)
		return result, err
	}
	result, err = this.CompleteGatewayHistory(token, duration, result)
	return
}

func (this *Lib) CompleteGatewayHistory(token auth.Token, duration string, gateways []map[string]interface{}) (result []map[string]interface{}, err error) {
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
	logStates, err := this.GetGatewayLogStates(token, ids)
	if err != nil {
		log.Println("ERROR completeGatewayList.GetGatewayLogStates()", err)
		return result, err
	}
	logHistory, err := this.GetGatewayLogHistory(token, ids, duration)
	if err != nil {
		log.Println("ERROR completeGatewayList.GetGatewayLogHistory()", err)
		return result, err
	}
	logEdges, err := this.GetLogedges(token, "gateway", ids, duration)
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

func (this *Lib) ListGateways(token auth.Token, limit string, offset string) (result []map[string]interface{}, err error) {
	gateways, err := this.PermListGateways(token, "r", limit, offset)
	if err != nil {
		log.Println("ERROR ListGateways.PermListGateways()", err)
		return result, err
	}
	return this.completeGatewayList(token, gateways)
}

func (this *Lib) ListAllGateways(token auth.Token) (result []map[string]interface{}, err error) {
	return this.PermListAllGateways(token, "r")
}

func (this *Lib) ListGatewaysOrdered(token auth.Token, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	gateways, err := this.PermListGatewaysOrdered(token, "r", limit, offset, orderfeature, direction)
	if err != nil {
		log.Println("ERROR ListGateways.PermListGateways()", err)
		return result, err
	}
	return this.completeGatewayList(token, gateways)
}

func (this *Lib) SearchGateways(token auth.Token, query string, limit string, offset string) (result []map[string]interface{}, err error) {
	gateways, err := this.PermSearchGateways(token, query, "r", limit, offset)
	if err != nil {
		return result, err
	}
	return this.completeGatewayList(token, gateways)
}

func (this *Lib) SearchGatewaysOrdered(token auth.Token, query string, limit string, offset string, orderfeature string, direction string) (result []map[string]interface{}, err error) {
	gateways, err := this.PermSearchGatewaysOrdered(token, query, "r", limit, offset, orderfeature, direction)
	if err != nil {
		return result, err
	}
	return this.completeGatewayList(token, gateways)
}

func (this *Lib) completeGatewayList(token auth.Token, gateways []map[string]interface{}) (result []map[string]interface{}, err error) {
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
	logStates, err := this.GetGatewayLogStates(token, ids)
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

func (this *Lib) GetGatewayDevices(token auth.Token, id string) (ids []string, err error) {
	req, err := http.NewRequest("GET", this.config.IotUrl+"/hubs/"+url.PathEscape(id)+"/devices?as=id", nil)
	if err != nil {
		debug.PrintStack()
		return ids, err
	}
	req.Header.Set("Authorization", token.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return ids, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		err = errors.New(buf.String())
		debug.PrintStack()
		return ids, err
	}
	err = json.NewDecoder(resp.Body).Decode(&ids)
	if err != nil {
		debug.PrintStack()
		return ids, err
	}
	return ids, nil
}
