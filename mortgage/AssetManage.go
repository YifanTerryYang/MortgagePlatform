package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func CreateAsset(APIstub shim.ChaincodeStubInterface, a Asset) (bool, Asset, error) {
	exist, asset, err := CheckAsset(APIstub, a.Key)
	if exist {
		return false, nil, errors.New("Asset exist")
	}

	if createKeyErr != nil {
		return false, nil, createKeyErr
	}
	val, _ := json.Marshal(a)
	//APIstub.PutState(a.)

	return new(Asset)
}

func RemoveAsset(APIstub shim.ChaincodeStubInterface, assetKey string) bool {
	err := APIstub.DelState(assetKey)
	if err != nil {
		return false
	}
	return true
}

func CheckAsset(APIstub shim.ChaincodeStubInterface, assetKey string) (bool, Asset) {
	key, err := APIstub.CreateCompositeKey("asset", assetKey)
	if err != nil {
		return false, nil
	}

	val, err1 := APIstub.GetState(key)
	if err1 != nil {
		return false, nil
	}

	resultAsset := new(Asset)

	err2 := json.Unmarshal(val, &resultAsset)
	if err2 != nil {
		return false, nil
	}

	return true, resultAsset

}
