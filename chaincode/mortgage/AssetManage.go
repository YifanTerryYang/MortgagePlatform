package main

import (
	"encoding/json"
	"strings"
	"errors"
	"strconv"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func createAsset(APIstub shim.ChaincodeStubInterface, a Asset) error {
	val, _ := json.Marshal(a)
	return APIstub.PutState(a.Key, []byte(val))
	
}

func removeAsset(APIstub shim.ChaincodeStubInterface, compositeAssetKey string) error {
	return APIstub.DelState(compositeAssetKey)
}

func checkAsset(APIstub shim.ChaincodeStubInterface, compositeAssetKey string) (bool, Asset) {
	val, err1 := APIstub.GetState(compositeAssetKey)
	if err1 != nil {
		return false, Asset{}
	}

	resultAsset := Asset{}
	err2 := json.Unmarshal(val, &resultAsset)
	if err2 != nil {
		return false, Asset{}
	}

	return true, resultAsset
}

func getAssetInfo(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	assetid := args[0]
	compositeAssetKey, _ := GetAssetKey(APIstub, assetid)
	check, asset := checkAsset(APIstub, compositeAssetKey)
	if !check {   // if user exists
		return nil, errors.New("User not exists or password incorrect")
	}
	assetvalAsbytes, _ := json.Marshal(asset)
	return assetvalAsbytes, nil
}

// args: username and assetid
func postAsset(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 { return []byte{}, errors.New("PostAsset, args should be 5 --- ")}

	user := args[0]
	assetid := args[1]
	interestrate := args[2]
	worth := args[3]
	period := args[4]
	assetKey, _ := GetAssetKey(APIstub, assetid)
	assetval, _ := APIstub.GetState(assetKey)
	asset := Asset{}
	if err := json.Unmarshal(assetval, &asset); err != nil {
		return []byte{}, errors.New("postAsset --- " + err.Error())
	}
	fmt.Println(asset.Owned)
	fmt.Println(user)
	compositeUserKey, _ := GetUserKey(APIstub, user)
	if asset.Owned != compositeUserKey {
		return []byte{}, errors.New("PostAsset, asset not owned by this user --- ")
	}
	asset.Interestrate, _ = strconv.ParseFloat(interestrate, 32)
	asset.Worth, _ = strconv.ParseFloat(worth, 32)
	asset.Period, _ = strconv.ParseInt(period, 10, 32)
	updatedassetAsbytes, _ := json.Marshal(asset)
	if puterr := APIstub.PutState(assetKey, updatedassetAsbytes); puterr != nil { return []byte{}, puterr }
	// add this asset to "postedasset"
	postedassetAsbytes, _ := APIstub.GetState("postedasset")
	if strings.Contains(string(postedassetAsbytes), assetKey) {
		return []byte{}, errors.New("asset been posted already")
	}
	newval := string(postedassetAsbytes)
	newval = newval + "|" + assetKey
	if puterr := APIstub.PutState("postedasset", []byte(newval)); puterr != nil {return []byte{}, puterr}
	return []byte(newval), nil  // insert successfully
}

// args: username and assetid
func unpostAsset(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 { return []byte{}, errors.New("UnpostAsset, args should be 2 --- ")}

	user := args[0]
	assetid := args[1]
	assetKey, _ := GetAssetKey(APIstub, assetid)
	assetval, _ := APIstub.GetState(assetKey)
	asset := Asset{}
	if err := json.Unmarshal(assetval, &asset); err != nil {
		return []byte{}, err
	}
	compositeUserKey, _ := GetUserKey(APIstub, user)
	if asset.Owned != compositeUserKey {
		return []byte{}, errors.New("UnpostAsset, asset not owned by this user --- ")
	}
	if len(asset.Buyerspercent) > 0 {
		return []byte{}, errors.New("UnpostAsset, asset can not be unposted. Loan not clear --- ")
	}
	// remove this asset from "postedasset"
	getpostedassetAsbytes, _ := APIstub.GetState("postedasset")
	if !strings.Contains(string(getpostedassetAsbytes), assetKey) {
		return []byte{}, errors.New("asset NOT been posted yet!")
	}
	splited := strings.Split(string(getpostedassetAsbytes), "|" + assetKey)
	var newvalstr string
	if len(splited) == 1 { 
		newvalstr = splited[0] 
	} else { 
		newvalstr = splited[0] + splited[1] 
	}
	if puterr := APIstub.PutState("postedasset", []byte(newvalstr)); puterr != nil { return []byte{}, puterr}
	return []byte(newvalstr), nil // remove successfully, return updated list
}