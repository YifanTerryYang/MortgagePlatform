package main

import (
	"errors"
	"strconv"
	"math"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*
 * Args: 1. paymentid
 *    	 2. pay percentage
 */
func (u *User)BuyAsset(APIstub shim.ChaincodeStubInterface, a Asset, args []string) error {
	if len(args) != 2 { return errors.New("BuyAsset - incorrect number of arguments. Expecting 2")}

	paymentid := args[0]
	paypercent := args[1]
	percentage, parseErr := strconv.ParseFloat(paypercent, 32)
	if parseErr != nil { return errors.New("BuyAsset - " + parseErr.Error()) }
	
	if a.Owned == u.Username {
		return errors.New("BuyAsset, user can NOT buy its own asset --- ")
	}
	// add buyer info into Asset
	remaining := 1 - a.Status
	if percentage > remaining { 
		return errors.New("BuyAsset - No enough share for this asset '" + a.Key + "'")
	}
	a.Status = remaining - percentage   // update asset remaining share
	a.Buyerspercent[u.Username] = a.Buyerspercent[u.Username] + percentage
	// calculate payment plan
	totalPayWithInterest := percentage * math.Pow((1+a.Interestrate), float64(a.Period)) * a.Worth
	payPer := totalPayWithInterest / float64(a.Period)

	// add asset to user's paymentmethod's Payssets list
	paymentmethod := u.Info.Paymentmethodlist[paymentid]
	if paymentmethod.Accountnumber == "" {
		return errors.New("BuyAsset - payment not exists")
	}

	var record Payassetrecord
	record = paymentmethod.Payassets[a.Key];
	if (Payassetrecord{}) != record {
		return errors.New("BuyAsset - you have bought this asset already.")
	} 
	//p, _ := strconv.ParseInt(a.Period, 10, 32)
	record = Payassetrecord{Period: int(a.Period), Amountper:payPer}
	paymentmethod.Payassets[a.Key] = record

	// Put asset and user back to ledger
	assetAsbytes, _ := json.Marshal(a)
	userAsbytes, _ := json.Marshal(u)

	if err := APIstub.PutState(a.Key, assetAsbytes); err != nil {
		return errors.New("BuyAsset - " + err.Error())
	}
	if err := APIstub.PutState(u.Username, userAsbytes); err != nil {
		return errors.New("BuyAsset - " + err.Error())
	}

	return nil
}

