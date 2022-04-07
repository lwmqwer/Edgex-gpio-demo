//
// Copyright (C) 2022 Wuming Liu
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"time"

	dsModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	sakshat "github.com/lwmqwer/SAKS-SDK-GO"
)

type Driver struct {
	lc      logger.LoggingClient
	asyncCh chan<- *dsModels.AsyncValues
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues, deviceCh chan<- []dsModels.DiscoveredDevice) error {
	s.lc = lc
	s.asyncCh = asyncCh

	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *Driver) HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	s.lc.Debug(fmt.Sprintf("protocols: %v", protocols))

	now := time.Now().UnixNano()

	for i, req := range reqs {
		s.lc.Debug(fmt.Sprintf("request: %d resource: %v attributes: %v", i, req.DeviceResourceName, req.Attributes))
		switch req.DeviceResourceName {
		case "TemperatureSensor":
			{
				value := sakshat.Ds18b20.Temperature(0)
				cv, _ := dsModels.NewFloat64Value(reqs[0].DeviceResourceName, now, value)
				res = append(res, cv)
			}
		}
	}

	return res, nil
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *Driver) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {

	s.lc.Debug(fmt.Sprintf("protocols: %v", protocols))

	for _, param := range params {
		s.lc.Debug(fmt.Sprintf("param: %v", param))
		switch param.DeviceResourceName {
		case "DigitalDisplay":
			{
				str, err := param.StringValue()
				if err != nil {
					return err
				}
				sakshat.DigitalDisplay.Show(str)
			}
		case "LED":
			{
				mask, err := param.Uint8Value()
				if err != nil {
					return err
				}
				for i := 0; i < 8; i++ {
					if mask&(1<<i) != 0 {
						sakshat.LEDRow.OnForIndex(uint(i))
					} else {
						sakshat.LEDRow.OffForIndex(uint(i))
					}
				}
			}
		}
	}
	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *Driver) Stop(force bool) error {
	s.lc.Debug(fmt.Sprintf("Driver.Stop called: force=%v", force))
	sakshat.Clean()
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (s *Driver) AddDevice(deviceName string, protocols map[string]contract.ProtocolProperties, adminState contract.AdminState) error {
	s.lc.Debug(fmt.Sprintf("a new Device is added: %s", deviceName))
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (s *Driver) UpdateDevice(deviceName string, protocols map[string]contract.ProtocolProperties, adminState contract.AdminState) error {
	s.lc.Debug(fmt.Sprintf("Device %s is updated", deviceName))
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (s *Driver) RemoveDevice(deviceName string, protocols map[string]contract.ProtocolProperties) error {
	s.lc.Debug(fmt.Sprintf("Device %s is removed", deviceName))
	return nil
}
