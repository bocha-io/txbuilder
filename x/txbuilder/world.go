package txbuilder

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bocha-io/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type AbiStruct []struct {
	Inputs []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"inputs"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Anonymous bool   `json:"anonymous,omitempty"`
	Outputs   []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"outputs,omitempty"`
	StateMutability string `json:"stateMutability,omitempty"`
}

func NewWorldABI(contract []byte) abi.ABI {
	// We need to remove everything that is type error because it breaks the abi decoder
	var raw AbiStruct
	err := json.Unmarshal(contract, &raw)
	if err != nil {
		logger.LogError(fmt.Sprintf("failed to unmarshal abi json: %s", err))
		panic("could not unmarshal json")
	}

	withoutErrors := make(AbiStruct, 0)
	for _, v := range raw {
		if v.Type != "error" {
			withoutErrors = append(withoutErrors, v)
		}
	}

	fixedToAbi, err := json.Marshal(withoutErrors)
	if err != nil {
		logger.LogError(fmt.Sprintf("failed to marshal the fixed data: %s", err))
		panic("failed to marshal the fixed data")
	}

	abiDecoded, err := abi.JSON(strings.NewReader(string(fixedToAbi)))
	if err != nil {
		logger.LogError(fmt.Sprintf("error decoding IWorld abi: %s", err))
		panic("error decoding IWorld abi")
	}
	return abiDecoded
}

// This is needed when encoding the params to interact with a contract
func StringToSlice(stringID string) ([32]byte, error) {
	id, err := hexutil.Decode(stringID)
	if err != nil {
		return [32]byte{}, fmt.Errorf("error decoding the string %s", err.Error())
	}

	if len(id) != 32 {
		return [32]byte{}, fmt.Errorf("invalid length")
	}

	var idArray [32]byte
	copy(idArray[:], id)
	return idArray, nil
}
