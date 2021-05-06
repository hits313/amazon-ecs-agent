// +build windows

// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package ecscni

import (
	"time"

	"github.com/containernetworking/cni/pkg/types"
)

const (
	// ECSVPCENIPluginName is the name of the vpc-eni plugin.
	ECSVPCENIPluginName = "vpc-eni"
	// ECSVPCENIPluginExecutable is the name of vpc-eni executable.
	ECSVPCENIPluginExecutable = "vpc-eni.exe"
	setupNSBackoffMin         = time.Second
	setupNSBackoffMax         = time.Second * 3
	setupNSBackoffJitter      = 0.2
	setupNSBackoffMultiple    = 1.3
	setupNSMaxRetryCount      = 3
)

// VPCENIPluginConfig contains all the information required to invoke the vpc-eni plugin.
type VPCENIPluginConfig struct {
	// Type is the cni plugin name.
	Type string `json:"type,omitempty"`
	// CNIVersion is the cni spec version to use.
	CNIVersion string `json:"cniVersion,omitempty"`
	// DNS is used to pass DNS information to the plugin.
	DNS types.DNS `json:"dns"`

	// ENIName is the name of the eni on the instance.
	ENIName string `json:"eniName"`
	// ENIMACAddress is the MAC address of the eni.
	ENIMACAddress string `json:"eniMACAddress"`
	// ENIIPAddress is the is the ipv4 of eni.
	ENIIPAddress string `json:"eniIPAddress"`
	// GatewayIPAddress specifies the IPv4 address of the subnet gateway for the eni.
	GatewayIPAddress string `json:"gatewayIPAddress"`
	// NoInfraContainer specifies if HCN Namespace is being used for networking setup.
	NoInfraContainer bool `json:"noInfraContainer"`
	// UseExistingNetwork specifies if existing network should be used instead of creating a new one.
	UseExistingNetwork bool `json:"useExistingNetwork"`
}
