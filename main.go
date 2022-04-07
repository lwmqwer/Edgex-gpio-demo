//
// Copyright (C) 2022 Wuming Liu
//
// SPDX-License-Identifier: Apache-2.0

// This package provides device service of a SwissArmyKnife.
package main

import (
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/startup"
	"github.com/lwmqwer/Edgex-gpio-demo/driver"
)

const (
	serviceName string = "device-SwissArmyKnife"
)

func main() {
	d := driver.Driver{}
	startup.Bootstrap(serviceName, Version, &d)
}
