package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func AddBuyer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	return shim.Success(nil)
}

func RemoveBuyer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	return shim.Success(nil)
}
