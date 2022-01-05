// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configuration

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/metadium/rosetta-metadium/metadium"
)

// Mode is the setting that determines if
// the implementation is "online" or "offline".
type Mode string

const (
	// Online is when the implementation is permitted
	// to make outbound connections.
	Online Mode = "ONLINE"

	// Offline is when the implementation is not permitted
	// to make outbound connections.
	Offline Mode = "OFFLINE"

	// Mainnet is the Metadium Mainnet.
	Mainnet string = "MAINNET"

	// Testnet is the Metadium Mainnet.
	Testnet string = "TESTNET"

	// DataDirectory is the default location for all
	// persistent data.
	DataDirectory = "/data"

	// ModeEnv is the environment variable read
	// to determine mode.
	ModeEnv = "MODE"

	// NetworkEnv is the environment variable
	// read to determine network.
	NetworkEnv = "NETWORK"

	// PortEnv is the environment variable
	// read to determine the port for the Rosetta
	// implementation.
	PortEnv = "PORT"

	// GmetEnv is an optional environment variable
	// used to connect rosetta-metadium to an already
	// running gmet node.
	GmetEnv = "GMET"

	// DefaultGmetURL is the default URL for
	// a running gmet node. This is used
	// when GmetEnv is not populated.
	DefaultGmetURL = "http://localhost:8588"

	// SkipGmetAdminEnv is an optional environment variable
	// to skip gmet `admin` calls which are typically not supported
	// by hosted node services. When not set, defaults to false.
	SkipGmetAdminEnv = "SKIP_GMET_ADMIN"

	// MiddlewareVersion is the version of rosetta-metadium.
	MiddlewareVersion = "0.0.4"
)

// Configuration determines how
type Configuration struct {
	Mode                   Mode
	Network                *types.NetworkIdentifier
	GenesisBlockIdentifier *types.BlockIdentifier
	GmetURL                string
	RemoteGmet             bool
	Port                   int
	GmetArguments          string
	SkipGmetAdmin          bool

	// Block Reward Data
	Params *params.ChainConfig
}

// LoadConfiguration attempts to create a new Configuration
// using the ENVs in the environment.
func LoadConfiguration() (*Configuration, error) {
	config := &Configuration{}

	modeValue := Mode(os.Getenv(ModeEnv))
	switch modeValue {
	case Online:
		config.Mode = Online
	case Offline:
		config.Mode = Offline
	case "":
		return nil, errors.New("MODE must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid mode", modeValue)
	}

	networkValue := os.Getenv(NetworkEnv)
	switch networkValue {
	case Mainnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: metadium.Blockchain,
			Network:    metadium.MainnetNetwork,
		}
		config.GenesisBlockIdentifier = metadium.MainnetGenesisBlockIdentifier
		config.Params = params.MetadiumMainnetChainConfig
		config.GmetArguments = metadium.MainnetGmetArguments
	case Testnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: metadium.Blockchain,
			Network:    metadium.TestnetNetwork,
		}
		config.GenesisBlockIdentifier = metadium.TestnetGenesisBlockIdentifier
		config.Params = params.MetadiumTestnetChainConfig
		config.GmetArguments = metadium.TestnetGmetArguments

	case "":
		return nil, errors.New("NETWORK must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid network", networkValue)
	}

	config.GmetURL = DefaultGmetURL
	envGmetURL := os.Getenv(GmetEnv)
	if len(envGmetURL) > 0 {
		config.RemoteGmet = true
		config.GmetURL = envGmetURL
	}

	config.SkipGmetAdmin = false
	envSkipGmetAdmin := os.Getenv(SkipGmetAdminEnv)
	if len(envSkipGmetAdmin) > 0 {
		val, err := strconv.ParseBool(envSkipGmetAdmin)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to parse SKIP_GMET_ADMIN %s", err, envSkipGmetAdmin)
		}
		config.SkipGmetAdmin = val
	}

	portValue := os.Getenv(PortEnv)
	if len(portValue) == 0 {
		return nil, errors.New("PORT must be populated")
	}

	port, err := strconv.Atoi(portValue)
	if err != nil || len(portValue) == 0 || port <= 0 {
		return nil, fmt.Errorf("%w: unable to parse port %s", err, portValue)
	}
	config.Port = port

	return config, nil
}
