package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (u *User) PostAsset(a Asset) sc.Response {
	return shim.Success(nil)
}

func (u *User) UnpostAsset(a Asset) sc.Response {
	return shim.Success(nil)
}

func (u *User) AddAsset()
