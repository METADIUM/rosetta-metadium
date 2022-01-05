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
	"os"
	"testing"

	"github.com/metadium/rosetta-metadium/metadium"
	// "github.com/metadium/rosetta-metadium/params"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	tests := map[string]struct {
		Mode          string
		Network       string
		Port          string
		Gmet          string
		SkipGmetAdmin string

		cfg *Configuration
		err error
	}{
		"no envs set": {
			err: errors.New("MODE must be populated"),
		},
		"only mode set": {
			Mode: string(Online),
			err:  errors.New("NETWORK must be populated"),
		},
		"only mode and network set": {
			Mode:    string(Online),
			Network: Mainnet,
			err:     errors.New("PORT must be populated"),
		},
		"all set (mainnet)": {
			Mode:          string(Online),
			Network:       Mainnet,
			Port:          "1000",
			SkipGmetAdmin: "FALSE",
			cfg: &Configuration{
				Mode: Online,
				Network: &types.NetworkIdentifier{
					Network:    metadium.MainnetNetwork,
					Blockchain: metadium.Blockchain,
				},
				Params:                 params.MetadiumMainnetChainConfig,
				GenesisBlockIdentifier: metadium.MainnetGenesisBlockIdentifier,
				Port:          1000,
				GmetURL:       DefaultGmetURL,
				GmetArguments: metadium.MainnetGmetArguments,
				SkipGmetAdmin: false,
			},
		},
		"all set (mainnet) + gmet": {
			Mode:          string(Online),
			Network:       Mainnet,
			Port:          "1000",
			Gmet:          "http://blah",
			SkipGmetAdmin: "TRUE",
			cfg: &Configuration{
				Mode: Online,
				Network: &types.NetworkIdentifier{
					Network:    metadium.MainnetNetwork,
					Blockchain: metadium.Blockchain,
				},
				Params:                 params.MetadiumMainnetChainConfig,
				GenesisBlockIdentifier: metadium.MainnetGenesisBlockIdentifier,
				Port:          1000,
				GmetURL:       "http://blah",
				RemoteGmet:    true,
				GmetArguments: metadium.MainnetGmetArguments,
				SkipGmetAdmin: true,
			},
		},
		"all set (testnet)": {
			Mode:          string(Online),
			Network:       Testnet,
			Port:          "1000",
			SkipGmetAdmin: "TRUE",
			cfg: &Configuration{
				Mode: Online,
				Network: &types.NetworkIdentifier{
					Network:    metadium.TestnetNetwork,
					Blockchain: metadium.Blockchain,
				},
				Params:                 params.MetadiumTestnetChainConfig,
				GenesisBlockIdentifier: metadium.TestnetGenesisBlockIdentifier,
				Port:          1000,
				GmetURL:       DefaultGmetURL,
				GmetArguments: metadium.TestnetGmetArguments,
				SkipGmetAdmin: true,
			},
		},
		"invalid mode": {
			Mode:    "bad mode",
			Network: Testnet,
			Port:    "1000",
			err:     errors.New("bad mode is not a valid mode"),
		},
		"invalid network": {
			Mode:    string(Offline),
			Network: "bad network",
			Port:    "1000",
			err:     errors.New("bad network is not a valid network"),
		},
		"invalid port": {
			Mode:    string(Offline),
			Network: Testnet,
			Port:    "bad port",
			err:     errors.New("unable to parse port bad port"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			os.Setenv(ModeEnv, test.Mode)
			os.Setenv(NetworkEnv, test.Network)
			os.Setenv(PortEnv, test.Port)
			os.Setenv(GmetEnv, test.Gmet)
			os.Setenv(SkipGmetAdminEnv, test.SkipGmetAdmin)

			cfg, err := LoadConfiguration()
			if test.err != nil {
				assert.Nil(t, cfg)
				assert.Contains(t, err.Error(), test.err.Error())
			} else {
				assert.Equal(t, test.cfg, cfg)
				assert.NoError(t, err)
			}
		})
	}
}
