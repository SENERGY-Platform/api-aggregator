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
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/api-aggregator/pkg/auth"
	"github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/models/go/models"
	"log"
)

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
	logHistory, err := this.GetGatewayLogHistory(token, ids, duration)
	if err != nil {
		log.Println("ERROR legacyHubTransformations.GetGatewayLogHistory()", err)
		return result, err
	}
	logEdges, err := this.GetLogedges(token, "gateway", ids, duration)
	if err != nil {
		log.Println("ERROR legacyDeviceTransformations.GetLogedges()", err)
		return result, err
	}
	for _, id := range ids {
		gateway := gatewayMap[id]
		gateway["log_history"] = logHistory[id]
		gateway["log_edge"] = logEdges[id]
		result = append(result, gateway)
	}
	return
}

func (this *Lib) ListGateways(token auth.Token, limit int64, offset int64) (result []map[string]interface{}, err error) {
	hubs, _, err, _ := this.deviceRepo.ListExtendedHubs(token.Jwt(), client.HubListOptions{
		Limit:      limit,
		Offset:     offset,
		SortBy:     "name.asc",
		Permission: client.READ,
	})
	if err != nil {
		return nil, err
	}
	return this.legacyHubTransformations(token, hubs)
}

func (this *Lib) ListAllGateways(token auth.Token) (result []map[string]interface{}, err error) {
	var limit int64 = 0
	var offset int64 = 0
	for {
		temp, err := this.ListGateways(token, limit, offset)
		if err != nil {
			return nil, err
		}
		result = append(result, temp...)
		if int64(len(temp)) < limit {
			return result, nil
		}
	}
}

func (this *Lib) legacyHubTransformations(token auth.Token, hubs []models.ExtendedHub) (result []map[string]interface{}, err error) {
	for _, hub := range hubs {
		element := map[string]interface{}{}
		temp, err := json.Marshal(hub)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(temp, &element)
		if err != nil {
			return nil, err
		}
		switch hub.ConnectionState {
		case *client.ConnectionStateOnline:
			element["log_state"] = "connected"
		case *client.ConnectionStateOffline:
			element["log_state"] = "disconnected"
		case *client.ConnectionStateUnknown:
			element["log_state"] = "unknown"
		default:
			element["log_state"] = "unknown"
		}
		//TODO: perm-search transformations for creator, permissions etc
		result = append(result, element)
	}
	return
}
