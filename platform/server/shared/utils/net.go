/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package utils

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	defaultNetwork = "tcp"
	allEths        = "0.0.0.0"
	envPodIP       = "POD_IP"
	consulTags     = "consul_tags"
)

var _ net.Addr = (*NetAddr)(nil)

// NetAddr implements the net.Addr interface.
type NetAddr struct {
	network string
	address string
}

// NewNetAddr creates a new NetAddr object with the network and address provided.
func NewNetAddr(network, address string) net.Addr {
	return &NetAddr{network, address}
}

// Network implements the net.Addr interface.
func (na *NetAddr) Network() string {
	return na.network
}

func FigureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIP)
	if len(ip) == 0 {
		ip = InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}

// String implements the net.Addr interface.
func (na *NetAddr) String() string {
	return na.address
}

func ParseAddr(addr net.Addr) (host string, port int, err error) {
	host, portStr, err := net.SplitHostPort(addr.String())
	if err != nil {
		return "", 0, err
	}

	if host == "" || host == "::" {
		host, err = getLocalIPv4Address()
		if err != nil {
			return "", 0, fmt.Errorf("get local ipv4 error, cause %w", err)
		}
	}
	port, err = net.LookupPort(defaultNetwork, portStr)
	if err != nil {
		return "", 0, err
	}
	if port == 0 {
		return "", 0, fmt.Errorf("invalid port %s", portStr)
	}

	return host, port, nil
}

func getLocalIPv4Address() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addr {
		ipNet, isIpNet := addr.(*net.IPNet)
		if isIpNet && !ipNet.IP.IsLoopback() {
			ipv4 := ipNet.IP.To4()
			if ipv4 != nil {
				return ipv4.String(), nil
			}
		}
	}
	return "", errors.New("not found ipv4 address")
}

func InternalIp() string {
	infs, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, inf := range infs {
		if isEthDown(inf.Flags) || isLoopback(inf.Flags) {
			continue
		}

		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return ""
}

func isEthDown(f net.Flags) bool {
	return f&net.FlagUp != net.FlagUp
}

func isLoopback(f net.Flags) bool {
	return f&net.FlagLoopback == net.FlagLoopback
}
