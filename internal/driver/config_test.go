// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"strings"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

func TestCreateConnectionInfo(t *testing.T) {
	schema := "tcp"
	host := "0.0.0.0"
	port := "1883"
	user := "admin"
	password := "password"
	clientId := "CommandPublisher"
	topic := "CommandTopic"
	protocols := map[string]models.ProtocolProperties{
		Protocol: {
			Schema:   schema,
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			ClientId: clientId,
			Topic:    topic,
		},
	}

	connectionInfo, err := CreateConnectionInfo(protocols)

	if err != nil {
		t.Fatalf("Fail to create connectionIfo. Error: %v", err)
	}
	if connectionInfo.Schema != schema || connectionInfo.Host != host || connectionInfo.Port != port ||
		connectionInfo.User != user || connectionInfo.Password != password || connectionInfo.ClientId != clientId ||
		connectionInfo.Topic != topic {
		t.Fatalf("Unexpect test result. %v should match to %v ", connectionInfo, protocols)
	}
}

func TestCreateConnectionInfo_fail(t *testing.T) {
	protocols := map[string]models.ProtocolProperties{
		Protocol: {},
	}

	_, err := CreateConnectionInfo(protocols)
	if err == nil || !strings.Contains(err.Error(), "unable to load config") {
		t.Fatalf("Unexpect test result, config should be fail to load")
	}
}

func TestCreateDriverConfig_SingleTopic(t *testing.T) {
	configs := map[string]string{
		IncomingSchema: "tcp", IncomingHost: "0.0.0.0", IncomingPort: "1883",
		IncomingUser: "admin", IncomingPassword: "public", IncomingQos: "0",
		IncomingKeepAlive: "3600", IncomingClientId: "IncomingDataSubscriber", IncomingTopics: "testTopic:TestResourceJSON:TestDevice",

		ResponseSchema: "tcp", ResponseHost: "0.0.0.0", ResponsePort: "1883",
		ResponseUser: "admin", ResponsePassword: "public", ResponseQos: "0",
		ResponseKeepAlive: "3600", ResponseClientId: "CommandResponseSubscriber", ResponseTopic: "ResponseTopic",

		AdditionalPayloadJsonValues: "",
	}
	driverConfig, err := CreateDriverConfig(configs)
	if err != nil {
		t.Fatalf("Fail to load config, %v", err)
	}
	topicInfo, ok := driverConfig.IncomingTopics["testTopic"]
	if !ok {
		t.Fatalf("Unexpect test result, driver config doesn't correctly load")
	}
	if driverConfig.IncomingSchema != configs[IncomingSchema] || driverConfig.IncomingHost != configs[IncomingHost] ||
		driverConfig.IncomingPort != 1883 || driverConfig.IncomingUser != configs[IncomingUser] ||
		driverConfig.IncomingPassword != configs[IncomingPassword] || driverConfig.IncomingQos != 0 ||
		driverConfig.IncomingKeepAlive != 3600 || driverConfig.IncomingClientId != configs[IncomingClientId] ||
		topicInfo.Resource != "TestResourceJSON" ||
		driverConfig.ResponseSchema != configs[ResponseSchema] || driverConfig.ResponseHost != configs[ResponseHost] ||
		driverConfig.ResponsePort != 1883 || driverConfig.ResponseUser != configs[ResponseUser] ||
		driverConfig.ResponsePassword != configs[ResponsePassword] || driverConfig.ResponseQos != 0 ||
		driverConfig.ResponseKeepAlive != 3600 || driverConfig.ResponseClientId != configs[ResponseClientId] ||
		driverConfig.ResponseTopic != configs[ResponseTopic] {

		t.Fatalf("Unexpect test result, driver config doesn't correct load")
	}
}

func TestCreateDriverConfig_IncomingTopics_MultiTopics(t *testing.T) {
	configs := map[string]string{
		IncomingSchema: "tcp", IncomingHost: "0.0.0.0", IncomingPort: "1883",
		IncomingUser: "admin", IncomingPassword: "public", IncomingQos: "0",
		IncomingKeepAlive: "3600", IncomingClientId: "IncomingDataSubscriber", IncomingTopics: "testTopic:TestResourceJSON:TestDevice, testTopic1: TestResourceJSON1 :TestDevice1",

		ResponseSchema: "tcp", ResponseHost: "0.0.0.0", ResponsePort: "1883",
		ResponseUser: "admin", ResponsePassword: "public", ResponseQos: "0",
		ResponseKeepAlive: "3600", ResponseClientId: "CommandResponseSubscriber", ResponseTopic: "ResponseTopic",

		AdditionalPayloadJsonValues: "",
	}
	driverConfig, err := CreateDriverConfig(configs)
	if err != nil {
		t.Fatalf("Fail to load config, %v", err)
	}
	topicInfo, ok := driverConfig.IncomingTopics["testTopic"]
	if !ok {
		t.Fatalf("Unexpect test result, driver config doesn't correctly load IncomingTopics for testTopic")
	}

	topicInfo1, ok := driverConfig.IncomingTopics["testTopic1"]
	if !ok {
		t.Fatalf("Unexpect test result, driver config doesn't correctly load IncomingTopics for testTopic1")
	}
	if driverConfig.IncomingSchema != configs[IncomingSchema] || driverConfig.IncomingHost != configs[IncomingHost] ||
		driverConfig.IncomingPort != 1883 || driverConfig.IncomingUser != configs[IncomingUser] ||
		driverConfig.IncomingPassword != configs[IncomingPassword] || driverConfig.IncomingQos != 0 ||
		driverConfig.IncomingKeepAlive != 3600 || driverConfig.IncomingClientId != configs[IncomingClientId] ||
		topicInfo.Resource != "TestResourceJSON" || topicInfo1.Resource != "TestResourceJSON1" ||
		driverConfig.ResponseSchema != configs[ResponseSchema] || driverConfig.ResponseHost != configs[ResponseHost] ||
		driverConfig.ResponsePort != 1883 || driverConfig.ResponseUser != configs[ResponseUser] ||
		driverConfig.ResponsePassword != configs[ResponsePassword] || driverConfig.ResponseQos != 0 ||
		driverConfig.ResponseKeepAlive != 3600 || driverConfig.ResponseClientId != configs[ResponseClientId] ||
		driverConfig.ResponseTopic != configs[ResponseTopic] {

		t.Fatalf("Unexpect test result, driver config doesn't correct load")
	}
}

func TestCreateDriverConfig_IncomingTopics_KeyValueNotColonSeparated(t *testing.T) {
	configs := map[string]string{
		IncomingSchema: "tcp", IncomingHost: "0.0.0.0", IncomingPort: "1883",
		IncomingUser: "admin", IncomingPassword: "public", IncomingQos: "0",
		IncomingKeepAlive: "3600", IncomingClientId: "IncomingDataSubscriber", IncomingTopics: "testTopic@TestResourceJSON @TestDevice",

		ResponseSchema: "tcp", ResponseHost: "0.0.0.0", ResponsePort: "1883",
		ResponseUser: "admin", ResponsePassword: "public", ResponseQos: "0",
		ResponseKeepAlive: "3600", ResponseClientId: "CommandResponseSubscriber", ResponseTopic: "ResponseTopic",

		AdditionalPayloadJsonValues: "",
	}
	_, err := CreateDriverConfig(configs)
	if err == nil {
		t.Fatalf("Expected to get error during load config for IncomingTopics")
	}
}

func TestCreateDriverConfig_AdditionalPayloadJsonValues(t *testing.T) {
	configs := map[string]string{
		IncomingSchema: "tcp", IncomingHost: "0.0.0.0", IncomingPort: "1883",
		IncomingUser: "admin", IncomingPassword: "public", IncomingQos: "0",
		IncomingKeepAlive: "3600", IncomingClientId: "IncomingDataSubscriber", IncomingTopics: "testTopic:TestResourceJSON:TestDevice",

		ResponseSchema: "tcp", ResponseHost: "0.0.0.0", ResponsePort: "1883",
		ResponseUser: "admin", ResponsePassword: "public", ResponseQos: "0",
		ResponseKeepAlive: "3600", ResponseClientId: "CommandResponseSubscriber", ResponseTopic: "ResponseTopic",

		AdditionalPayloadJsonValues: "key1:value1, key2: value2",
	}
	driverConfig, err := CreateDriverConfig(configs)
	if err != nil {
		t.Fatalf("Fail to load config, %v", err)
	}
	topicInfo, ok := driverConfig.IncomingTopics["testTopic"]
	if !ok {
		t.Fatalf("Unexpected test result, driver config doesn't correctly load IncominTopics")
	}

	additionalJsonValue1, ok := driverConfig.AdditionalPayloadJsonValues["key1"]
	if !ok {
		t.Fatalf("Unexpected test result, driver config doesn't correctly load AdditionalPayloadJsonValues")
	}

	additionalJsonValue2, ok := driverConfig.AdditionalPayloadJsonValues["key2"]
	if !ok {
		t.Fatalf("Unexpected test result, driver config doesn't correctly load AdditionalPayloadJsonValues")
	}

	if driverConfig.IncomingSchema != configs[IncomingSchema] || driverConfig.IncomingHost != configs[IncomingHost] ||
		driverConfig.IncomingPort != 1883 || driverConfig.IncomingUser != configs[IncomingUser] ||
		driverConfig.IncomingPassword != configs[IncomingPassword] || driverConfig.IncomingQos != 0 ||
		driverConfig.IncomingKeepAlive != 3600 || driverConfig.IncomingClientId != configs[IncomingClientId] ||
		topicInfo.Resource != "TestResourceJSON" ||
		driverConfig.ResponseSchema != configs[ResponseSchema] || driverConfig.ResponseHost != configs[ResponseHost] ||
		driverConfig.ResponsePort != 1883 || driverConfig.ResponseUser != configs[ResponseUser] ||
		driverConfig.ResponsePassword != configs[ResponsePassword] || driverConfig.ResponseQos != 0 ||
		driverConfig.ResponseKeepAlive != 3600 || driverConfig.ResponseClientId != configs[ResponseClientId] ||
		driverConfig.ResponseTopic != configs[ResponseTopic] ||
		additionalJsonValue1 != "value1" || additionalJsonValue2 != "value2" {

		t.Fatalf("Unexpect test result, driver config doesn't correct load")
	}
}

func TestCreateDriverConfig_AdditionalPayloadJsonValues_KeyValueNotColonSeparated(t *testing.T) {
	configs := map[string]string{
		IncomingSchema: "tcp", IncomingHost: "0.0.0.0", IncomingPort: "1883",
		IncomingUser: "admin", IncomingPassword: "public", IncomingQos: "0",
		IncomingKeepAlive: "3600", IncomingClientId: "IncomingDataSubscriber", IncomingTopics: "testTopic:TestResourceJSON:TestDevice",

		ResponseSchema: "tcp", ResponseHost: "0.0.0.0", ResponsePort: "1883",
		ResponseUser: "admin", ResponsePassword: "public", ResponseQos: "0",
		ResponseKeepAlive: "3600", ResponseClientId: "CommandResponseSubscriber", ResponseTopic: "ResponseTopic",

		AdditionalPayloadJsonValues: "key1@value1, key2: value2",
	}
	_, err := CreateDriverConfig(configs)
	if err == nil {
		t.Fatalf("Expected to get error during load config for AdditionalPayloadJsonValues")
	}
}

func TestCreateDriverConfig_fail(t *testing.T) {
	configs := map[string]string{}
	_, err := CreateDriverConfig(configs)
	if err == nil || !strings.Contains(err.Error(), "unable to load config") {
		t.Fatalf("Unexpect test result, config should be fail to load")
	}
}
