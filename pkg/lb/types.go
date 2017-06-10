// Copyright 2016-2017 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lb

import (
	"fmt"
	"net"
	"strings"
)

// DatapathID is the cluster wide unique identifier used in the datapath
// FIXME: This needs to be uint32 in the future
type DatapathID uint16

// Frontend is a single frontend consisting of an IP and ports
type Frontend struct {
	IP       net.IP
	Port     uint16
	Protocol Protocol
	ID       DatapathID
}

// Backend is an L3/L4 endpoint which a frontend can loadbalance to
type Backend struct {
	// Name of the backport port
	Name string

	// IP to map to when loadbalancing to this backend
	IP net.IP

	// Port to map to when loadbalancing to this backend. Only used when
	Port uint16

	// Protocol is the L4 protocol of the port or NONE. It must match
	// the frontend protocol.
	Protocol Protocol

	// Weight is the proportion this backend should be used in relation
	// to other backends
	Weight uint16
}

// Service consists of a frontend which is loadbalanced to a list of backends.
type Service struct {
	Id       DatapathID
	Sha256   string
	Frontend Frontend
	Backends []*Backend
}

// NewService returns a new empty service
func NewService() *Service {
	return &Service{
		Backends: []*Backend{},
	}
}

// ParseIP returns the parsed net.IP or an error if the provided IP is not
// a valid service IP
func ParseServiceIP(serviceIP string) (net.IP, error) {
	ip := net.ParseIP(serviceIP)
	if ip == nil {
		return nil, fmt.Errorf("invalid service ip '%s'", serviceIP)
	}

	return ip, nil
}

// Protocol is the protocol name used for a service
type Protocol string

// Valid loadbalancer protocol types
const (
	// NONE type
	NONE = Protocol("NONE")
	// TCP type.
	TCP = Protocol("TCP")
	// UDP type.
	UDP = Protocol("UDP")
)

// ValidateProtocol returns an error if protocol is not a valid loadbalancer protocol
func ValidateProtocol(proto string) error {
	switch NormalizeProtocol(proto) {
	case TCP, UDP, NONE:
	default:
		return fmt.Errorf("unknown protocol type '%s'", proto)
	}

	return nil
}

// NormalizeProtocol takes a case insensitive string and returns the Protocol
func NormalizeProtocol(protocol string) Protocol {
	return Protocol(strings.ToUpper(protocol))
}
