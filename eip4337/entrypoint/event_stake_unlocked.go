package entrypoint_interface

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EntryPointStakeUnlockedIterator struct {
	Event *EntryPointStakeUnlocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EntryPointStakeUnlockedIterator) Next() bool {
	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EntryPointStakeUnlocked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EntryPointStakeUnlocked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EntryPointStakeUnlockedIterator) Error() error {
	return it.fail
}

func (it *EntryPointStakeUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EntryPointStakeUnlocked struct {
	Account      common.Address
	WithdrawTime *big.Int
	Raw          types.Log
}
