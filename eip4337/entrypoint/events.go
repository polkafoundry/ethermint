package entrypoint_interface

import (
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

type EntryPointEventIterator struct {
	Event interface{}

	abi abi.ABI

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EntryPointEventIterator) unpackEventLog(log types.Log) (interface{}, error) {
	topic := log.Topics[0].String()
	var abiEvent *abi.Event
	for _, e := range it.abi.Events {
		if e.ID.String() == topic {
			abiEvent = &e
			break
		}
	}

	if abiEvent == nil {
		return nil, fmt.Errorf("event signature not found in abi")
	}

	// FIXME: find a better solution
	var out interface{}
	switch abiEvent.RawName {
	case "AccountDeployed":
		out = EntryPointAccountDeployed{Raw: log}
	case "Deposited":
		out = EntryPointDeposited{Raw: log}
	case "SignatureAggregatorChanged":
		out = EntryPointSignatureAggregatorChanged{Raw: log}
	case "StakeLocked":
		out = EntryPointStakeLocked{Raw: log}
	case "StakeUnlocked":
		out = EntryPointStakeUnlocked{Raw: log}
	case "StakeWithdrawn":
		out = EntryPointStakeWithdrawn{Raw: log}
	case "UserOperationEvent":
		out = EntryPointUserOperationEvent{Raw: log}
	case "UserOperationRevertReason":
		out = EntryPointUserOperationRevertReason{Raw: log}
	case "Withdrawn":
		out = EntryPointWithdrawn{Raw: log}
	default:
		return nil, fmt.Errorf("unknown abi")
	}

	if len(log.Data) > 0 {
		unpacked, err := abiEvent.Inputs.Unpack(log.Data)
		if err != nil {
			return nil, err
		}
		err = abiEvent.Inputs.Copy(out, unpacked)
		if err != nil {
			return nil, err
		}
	}

	var indexed abi.Arguments
	for _, arg := range abiEvent.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	err := abi.ParseTopics(out, indexed, log.Topics[1:])
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (it *EntryPointEventIterator) Next() bool {
	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			out, err := it.unpackEventLog(log)
			if err != nil {
				it.fail = err
				return false
			}
			it.Event = out
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		out, err := it.unpackEventLog(log)
		if err != nil {
			it.fail = err
			return false
		}
		it.Event = out
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EntryPointEventIterator) Error() error {
	return it.fail
}

func (it *EntryPointEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}
