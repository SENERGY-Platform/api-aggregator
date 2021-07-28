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

package environment

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/wvanbergen/kazoo-go"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"time"
)

type Publisher struct {
	deviceTopic   string
	locationTopic string
	devices       *kafka.Writer
	locations     *kafka.Writer
}

func NewPublisher(zkUrl string, deviceTopic string, locationTopic string) (*Publisher, error) {
	log.Println("ensure kafka topics")
	broker, err := GetBroker(zkUrl)
	if err != nil {
		return nil, err
	}
	if len(broker) == 0 {
		return nil, errors.New("missing kafka broker")
	}
	devices, err := getProducer(broker, deviceTopic, true)
	if err != nil {
		return nil, err
	}
	location, err := getProducer(broker, locationTopic, true)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		deviceTopic:   deviceTopic,
		locationTopic: locationTopic,
		devices:       devices,
		locations:     location,
	}, nil
}

func GetBroker(zk string) (brokers []string, err error) {
	return getBroker(zk)
}

func getBroker(zkUrl string) (brokers []string, err error) {
	zookeeper := kazoo.NewConfig()
	zookeeper.Logger = log.New(ioutil.Discard, "", 0)
	zk, chroot := kazoo.ParseConnectionString(zkUrl)
	zookeeper.Chroot = chroot
	if kz, err := kazoo.NewKazoo(zk, zookeeper); err != nil {
		return brokers, err
	} else {
		defer kz.Close()
		return kz.BrokerList()
	}
}

func getProducer(broker []string, topic string, debug bool) (writer *kafka.Writer, err error) {
	var logger *log.Logger
	if debug {
		logger = log.New(os.Stdout, "[KAFKA-PRODUCER] ", 0)
	} else {
		logger = log.New(ioutil.Discard, "", 0)
	}
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:     broker,
		Topic:       topic,
		MaxAttempts: 10,
		Logger:      logger,
	})
	return writer, err
}

func (this *Publisher) PublishLocation(Location Location, userId string) (err error) {
	cmd := LocationCommand{Command: "PUT", Id: Location.Id, Location: Location, Owner: userId}
	return this.PublishLocationCommand(cmd)
}

func (this *Publisher) PublishLocationDelete(id string, userId string) error {
	cmd := LocationCommand{Command: "DELETE", Id: id, Owner: userId}
	return this.PublishLocationCommand(cmd)
}

func (this *Publisher) PublishLocationCommand(cmd LocationCommand) error {
	log.Println("DEBUG: produce Location", cmd)
	message, err := json.Marshal(cmd)
	if err != nil {
		debug.PrintStack()
		return err
	}
	err = this.locations.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(cmd.Id),
			Value: message,
			Time:  time.Now(),
		},
	)
	if err != nil {
		debug.PrintStack()
	}
	return err
}

func (this *Publisher) PublishDevice(device Device, userId string) (err error) {
	cmd := DeviceCommand{Command: "PUT", Id: device.Id, Device: device, Owner: userId}
	return this.PublishDeviceCommand(cmd)
}

func (this *Publisher) PublishDeviceDelete(id string, userId string) error {
	cmd := DeviceCommand{Command: "DELETE", Id: id, Owner: userId}
	return this.PublishDeviceCommand(cmd)
}

func (this *Publisher) PublishDeviceCommand(cmd DeviceCommand) error {
	log.Println("DEBUG: produce device", cmd)
	message, err := json.Marshal(cmd)
	if err != nil {
		debug.PrintStack()
		return err
	}
	err = this.devices.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(cmd.Id),
			Value: message,
			Time:  time.Now(),
		},
	)
	if err != nil {
		debug.PrintStack()
	}
	return err
}
