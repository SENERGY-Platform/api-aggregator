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
	"github.com/SmartEnergyPlatform/api-aggregator/tests/environment/docker"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

func New(ctx context.Context, wg *sync.WaitGroup) (permSearchUrl string, publisher *Publisher, err error) {
	_, zk, err := docker.Zookeeper(ctx, wg)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", nil, err
	}
	zkUrl := zk + ":2181"

	err = docker.Kafka(ctx, wg, zkUrl)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", nil, err
	}

	_, elasticIp, err := docker.ElasticSearch(ctx, wg)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", nil, err
	}

	_, permIp, err := docker.PermSearch(ctx, wg, zkUrl, elasticIp)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", nil, err
	}
	permSearchUrl = "http://" + permIp + ":8080"

	time.Sleep(2 * time.Second)

	publisher, err = NewPublisher(zkUrl, "devices", "locations")

	return
}
