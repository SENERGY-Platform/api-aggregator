/*
 * Copyright 2021 InfAI (CC SES)
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

package environment

import (
	"context"
	"github.com/SENERGY-Platform/api-aggregator/pkg/tests/environment/docker"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

func New(ctx context.Context, wg *sync.WaitGroup) (repoUrl string, publisher *Publisher, err error) {
	_, zk, err := docker.Zookeeper(ctx, wg)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", nil, err
	}
	zkUrl := zk + ":2181"

	kafkaUrl, err := docker.Kafka(ctx, wg, zkUrl)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", nil, err
	}

	_, mongoIp, err := docker.MongoDB(ctx, wg)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", nil, err
	}
	mongoUrl := "mongodb://" + mongoIp + ":27017"

	_, permV2Ip, err := docker.PermissionsV2(ctx, wg, mongoUrl, kafkaUrl)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", nil, err
	}
	permv2Url := "http://" + permV2Ip + ":8080"

	_, repoIp, err := docker.DeviceRepo(ctx, wg, kafkaUrl, mongoUrl, permv2Url)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", nil, err
	}
	repoUrl = "http://" + repoIp + ":8080"

	time.Sleep(2 * time.Second)

	publisher, err = NewPublisher(zkUrl, "devices", "locations")

	return
}
