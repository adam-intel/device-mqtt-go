// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
// Copyright (c) 2019 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type ConnectionInfo struct {
	Schema   string
	Host     string
	Port     string
	User     string
	Password string
	ClientId string
	Topic    string
}

type topicInfo struct {
	Resource   string
	DeviceName string
}

type configuration struct {
	IncomingSchema    string
	IncomingHost      string
	IncomingPort      int
	IncomingUser      string
	IncomingPassword  string
	IncomingQos       int
	IncomingKeepAlive int
	IncomingClientId  string
	IncomingTopics    map[string]topicInfo // topic is the key

	ResponseSchema    string
	ResponseHost      string
	ResponsePort      int
	ResponseUser      string
	ResponsePassword  string
	ResponseQos       int
	ResponseKeepAlive int
	ResponseClientId  string
	ResponseTopic     string

	AdditionalJsonValues map[string]string
}

// CreateDriverConfig use to load driver config for incoming listener and response listener
func CreateDriverConfig(configMap map[string]string) (*configuration, error) {
	config := new(configuration)
	err := load(configMap, config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// CreateConnectionInfo use to load MQTT connectionInfo for read and write command
func CreateConnectionInfo(protocols map[string]models.ProtocolProperties) (*ConnectionInfo, error) {
	info := new(ConnectionInfo)
	protocol, ok := protocols[Protocol]
	if !ok {
		return info, fmt.Errorf("unable to load config, '%s' not exist", Protocol)
	}

	err := load(protocol, info)
	if err != nil {
		return info, err
	}
	return info, nil
}

// load by reflect to check map key and then fetch the value
func load(config map[string]string, des interface{}) error {
	errorMessage := "unable to load config, '%s' not exist"
	val := reflect.ValueOf(des).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)

		configVal, ok := config[typeField.Name]
		if !ok {
			return fmt.Errorf(errorMessage, typeField.Name)
		}

		switch valueField.Interface().(type) {
		case int:
			intVal, err := strconv.Atoi(configVal)
			if err != nil {
				return err
			}
			valueField.SetInt(int64(intVal))
		case string:
			valueField.SetString(configVal)
		case map[string]topicInfo:
			if strings.TrimSpace(configVal) == "" {
				continue
			}
			incomingTopicsPairs := strings.Split(configVal, ",")

			topics := make(map[string]topicInfo)

			for _, incomingTopicPair := range incomingTopicsPairs {
				values := strings.Split(incomingTopicPair, ":")

				if len(values) != 3 {
					return fmt.Errorf("wrong number of elements in %v expecting 3 received %v", incomingTopicPair, len(values))
				}

				topic := strings.TrimSpace(values[0])

				topics[topic] = topicInfo{
					Resource:   strings.TrimSpace(values[1]),
					DeviceName: strings.TrimSpace(values[2]),
				}
			}
			valueField.Set(reflect.ValueOf(topics))
		case map[string]string:
			if strings.TrimSpace(configVal) == "" {
				continue
			}

			additionalValuePairs := strings.Split(configVal, ",")
			additionalValues := make(map[string]string)

			for _, additionalValuePair := range additionalValuePairs {
				values := strings.Split(additionalValuePair, ":")

				if len(values) != 2 {
					return fmt.Errorf("wrong number of elements in %v: expecting 2, received %v", additionalValuePair, len(values))
				}

				jsonKey := strings.TrimSpace(values[0])
				jsonValue := strings.TrimSpace(values[1])
				additionalValues[jsonKey] = jsonValue

				valueField.Set(reflect.ValueOf(additionalValues))
			}

		default:
			return fmt.Errorf("non supported value type %v ,%v", valueField.Kind(), typeField.Name)
		}
	}
	return nil
}
