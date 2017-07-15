package main

import (
	"errors"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)


/************* payment section *************/
func (u *User) AddPaymentmethod(APIstub shim.ChaincodeStubInterface, payment Paymentmethod) error {
	u.Info.Paymentmethodlist[payment.Accountnumber] = payment
	updatedUserAsbytes, _ := json.Marshal(u)
	if APIstub.PutState(u.Username, updatedUserAsbytes) != nil {
		return errors.New("AddPaymentmethod, ledger update fail")
	}
	return nil
}