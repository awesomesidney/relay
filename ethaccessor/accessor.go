/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package ethaccessor

import (
	"github.com/Loopring/relay/config"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthNodeAccessor struct {
	Erc20Abi      *abi.ABI
	ProtocolImpls map[common.Address]*ProtocolImpl
	ks            *keystore.KeyStore
	*rpc.Client
}

func NewAccessor(accessorOptions config.AccessorOptions, commonOptions config.CommonOptions, ks *keystore.KeyStore) (*EthNodeAccessor, error) {
	var err error
	accessor := &EthNodeAccessor{}
	accessor.Client, err = rpc.Dial(accessorOptions.RawUrl)
	if nil != err {
		return nil, err
	}

	accessor.ks = ks
	if accessor.Erc20Abi, err = NewAbi(commonOptions.Erc20Abi); nil != err {
		return nil, err
	}
	accessor.ProtocolImpls = make(map[common.Address]*ProtocolImpl)

	for version, opts := range commonOptions.ProtocolImpls {
		impl := &ProtocolImpl{Version: version, ContractAddress: common.HexToAddress(opts.Address)}
		if protocolImplAbi, err := NewAbi(opts.ImplAbi); nil != err {
			return nil, err
		} else {
			impl.ProtocolImplAbi = protocolImplAbi
		}
		if registryAbi, err := NewAbi(opts.RegistryAbi); nil != err {
			return nil, err
		} else {
			impl.RinghashRegistryAbi = registryAbi
		}
		if transferDelegateAbi, err := NewAbi(opts.DelegateAbi); nil != err {
			return nil, err
		} else {
			impl.DelegateAbi = transferDelegateAbi
		}
		if tokenRegistryAbi, err := NewAbi(opts.TokenRegistryAbi); nil != err {
			return nil, err
		} else {
			impl.TokenRegistryAbi = tokenRegistryAbi
		}

		callMethod := accessor.ContractCallMethod(impl.ProtocolImplAbi, impl.ContractAddress)
		var addr string
		if err := callMethod(&addr, "lrcTokenAddress", "latest"); nil != err {
			return nil, err
		} else {
			impl.LrcTokenAddress = common.HexToAddress(addr)
		}
		if err := callMethod(&addr, "ringhashRegistryAddress", "latest"); nil != err {
			return nil, err
		} else {
			impl.RinghashRegistryAddress = common.HexToAddress(addr)
		}
		if err := callMethod(&addr, "tokenRegistryAddress", "latest"); nil != err {
			return nil, err
		} else {
			impl.TokenRegistryAddress = common.HexToAddress(addr)
		}
		if err := callMethod(&addr, "delegateAddress", "latest"); nil != err {
			return nil, err
		} else {
			impl.DelegateAddress = common.HexToAddress(addr)
		}
		accessor.ProtocolImpls[impl.ContractAddress] = impl
	}

	return accessor, nil
}