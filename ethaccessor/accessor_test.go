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

package ethaccessor_test

import (
	"github.com/Loopring/relay/config"
	"github.com/Loopring/relay/ethaccessor"
	"github.com/Loopring/relay/test"
	"github.com/Loopring/relay/types"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestNewAccessor(t *testing.T) {
	cfg := config.LoadConfig("/Users/yuhongyu/Desktop/service/go/src/github.com/Loopring/relay/config/relay.toml")
	accessor, err := ethaccessor.NewAccessor(cfg.Accessor, cfg.Common, nil)
	if nil != err {
		println(err.Error())
	}
	var b types.Big
	if err := accessor.Call(&b, "eth_getBalance", common.HexToAddress("0x750ad4351bb728cec7d639a9511f9d6488f1e259"), "pending"); nil != err {
		t.Error(err.Error())
	}

	t.Log(b.BigInt().String())

	balance, _ := accessor.Erc20Balance(common.HexToAddress("0x937ff659c8a9d85aac39dfa84c4b49bb7c9b226e"), common.HexToAddress("0xb5fab0b11776aad5ce60588c16bd59dcfd61a1c2"), "pending")
	t.Log(balance.String())

	reqs := []*ethaccessor.BatchErc20Req{&ethaccessor.BatchErc20Req{
		Owner:          common.HexToAddress("0xb5fab0b11776aad5ce60588c16bd59dcfd61a1c2"),
		Token:          common.HexToAddress("0x937ff659c8a9d85aac39dfa84c4b49bb7c9b226e"),
		BlockParameter: "pending",
	},
		&ethaccessor.BatchErc20Req{
			Owner:          common.HexToAddress("0x48ff2269e58a373120FFdBBdEE3FBceA854AC30A"),
			Token:          common.HexToAddress("0x937ff659c8a9d85aac39dfa84c4b49bb7c9b226e"),
			BlockParameter: "pending",
		}}
	accessor.BatchErc20BalanceAndAllowance(reqs)

	t.Log("balance", reqs[0].Balance.BigInt().String())
	t.Log("balance", reqs[1].Balance.BigInt().String())
}

func TestEthNodeAccessor_Erc20Balance(t *testing.T) {
	c := test.LoadConfig()
	accessor, err := test.GenerateAccessor(c)
	if err != nil {
		t.Fatalf("generate accessor error:%s", err.Error())
	}

	tokenAddress := common.HexToAddress("0x478d07f3cBE07f01B5c7D66b4Ba57e5a3c520564")
	owner := common.HexToAddress("0x1b978a1d302335a6f2ebe4b8823b5e17c3c84135")
	balance, err := accessor.Erc20Balance(tokenAddress, owner, "latest")
	if err != nil {
		t.Fatalf("accessor get erc20 balance error:%s", err.Error())
	}

	t.Log(balance.String())
}