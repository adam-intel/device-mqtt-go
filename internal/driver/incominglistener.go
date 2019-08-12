// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
// Copyright (c) 2019 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	sdk "github.com/edgexfoundry/device-sdk-go"
	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
)

func startIncomingListening() error {
	var scheme = driver.Config.IncomingSchema
	var brokerUrl = driver.Config.IncomingHost
	var brokerPort = driver.Config.IncomingPort
	var username = driver.Config.IncomingUser
	var password = driver.Config.IncomingPassword
	var mqttClientId = driver.Config.IncomingClientId
	var qos = byte(driver.Config.IncomingQos)
	var keepAlive = driver.Config.IncomingKeepAlive
	var topics = driver.Config.IncomingTopics

	uri := &url.URL{
		Scheme: strings.ToLower(scheme),
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(username, password),
	}

	client, err := createClient(mqttClientId, uri, keepAlive)
	if err != nil {
		return err
	}

	defer func() {
		if client.IsConnected() {
			client.Disconnect(5000)
		}
	}()

	for topic := range topics {
		token := client.Subscribe(topic, qos, onIncomingDataReceived)
		if token.Wait() && token.Error() != nil {
			driver.Logger.Info(fmt.Sprintf("[Incoming listener] Stop incoming data listening. Cause:%v", token.Error()))
			return token.Error()
		}
		driver.Logger.Info(fmt.Sprintf("[Incoming listener] Subscribed to topic: %v", topic))
	}

	driver.Logger.Info("[Incoming listener] Start incoming data listening. ")
	select {}
}

func onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {

	topicInfo, ok := driver.Config.IncomingTopics[message.Topic()]
	if !ok {
		driver.Logger.Error(fmt.Sprintf("[Incoming listener] Topic %s not in IncomingTopics configuration", message.Topic()))
		return
	}

	service := sdk.RunningService()

	deviceObject, ok := service.DeviceResource(topicInfo.DeviceName, topicInfo.Resource, "get")
	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No DeviceObject found for topic=%v, deviceName=%s and resource=%s", message.Topic(), topicInfo.DeviceName, topicInfo.Resource))
		return
	}

	req := sdkModel.CommandRequest{
		DeviceResourceName: topicInfo.Resource,
		Type:               sdkModel.ParseValueType(deviceObject.Properties.Value.Type),
	}

	newPayload := []byte{}
	if sdk.DriverConfigs()[AdditionalPayloadJsonValues] != "" {
		extendedPayload, err := addJSONValuesToPayload(message.Payload())
		if err != nil {
			driver.Logger.Error(fmt.Sprintf("[Incoming listener] Error adding JSON values to payload: %v\n", err))
			return
		}
		newPayload = extendedPayload
	} else {
		newPayload = message.Payload()
	}

	result, err := newResult(req, string(newPayload))
	if err != nil {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. Failed to parse reading to a CommandValue for topic=%v msg=%v error=%v", message.Topic(), string(message.Payload()), err))
		return
	}

	asyncValues := &sdkModel.AsyncValues{
		DeviceName:    topicInfo.DeviceName,
		CommandValues: []*sdkModel.CommandValue{result},
	}

	driver.Logger.Info(fmt.Sprintf("[Incoming listener] Incoming reading received: topic=%v msg=%v", message.Topic(), string(message.Payload())))

	driver.AsyncCh <- asyncValues
}

func addJSONValuesToPayload(incomingPayload []byte) ([]byte, error) {
	var payloadData map[string]interface{}
	err := json.Unmarshal(incomingPayload, &payloadData)
	if err != nil {
		return nil, fmt.Errorf("Couldn't unmarshal incoming message payload: %v", err)
	}

	for jsonKey, jsonValue := range driver.Config.AdditionalPayloadJsonValues {
		payloadData[jsonKey] = jsonValue
	}
	payload, err := json.Marshal(payloadData)
	if err != nil {
		return nil, fmt.Errorf("Couldn't marshal payload data: %v", err)
	}

	return payload, nil
}
