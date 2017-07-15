package main

import (
	"encoding/json"
	"strings"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type Mortgageplatform struct{
}

func (s *Mortgageplatform) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	APIstub.PutState("postedasset", []byte(""))  // Initialize "posted Asset" value.
	APIstub.PutState("userslist", []byte(""))  // Initialize "userlist" value.
	return shim.Success(nil)
}

func (s *Mortgageplatform) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	fmt.Println("ARGS length is: ", len(args))
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createnewuser" {
		return s.CreateNewUser(APIstub, args)
	} else if function == "updateuser" {
		return s.UpdateUser(APIstub, args)
	} else if function == "login" {
		return s.Login(APIstub, args)
	} else if function == "addasset" {
		return s.AddAsset(APIstub, args)
	} else if function == "addpayment" {
		return s.AddPayment(APIstub, args)
	} else if function == "getpostedassets" {
		return s.GetPostedAssets(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * Call AccountManage.go ---> login
 * Args: 1. username
 * 		 2. password
 */
func (s *Mortgageplatform) Login(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	res, err := login(APIstub, args)
	if err != nil { return shim.Error(err.Error()) }
	return shim.Success(res)
}


/*
 * Call AccountManage.go ---> createNewUser
 * Args: 1. username
 * 		 2. password
 * 		 3. userinfo string: {"fname":"yifan", "lname":"yang", "gender":"male"}
 */
func (s *Mortgageplatform) CreateNewUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	err := createNewUser(APIstub, args)
	if err != nil { return shim.Error(err.Error()) }
	return shim.Success(nil)
}

/*
 * Call AccountManage.go ---> updateUser
 */
func (s *Mortgageplatform) UpdateUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	err := updateUser(APIstub,args)
	if err != nil { return shim.Error(err.Error()) }
	return shim.Success(nil)
}

/*
 * Call AssetManage.go 
 */
func (s *Mortgageplatform) AddAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 { return shim.Error("AddAsset - incorrect number of arguments. Expecting 3")}

	userid := args[0]
	assetkey := args[1]
	assetval := args[2]
	assetcompositeKey, _ := GetAssetKey(APIstub, assetkey)
	assetexist, _ := checkAsset(APIstub, assetcompositeKey)
	if assetexist { return shim.Error("AddAsset, asset exists already --- ") }

	usercompositekey, createUserKeyErr := GetUserKey(APIstub, userid)
	if createUserKeyErr != nil { return shim.Error("AddAsset, --- " + createUserKeyErr.Error()) }
	userinfoAsbytes, _ := APIstub.GetState(usercompositekey)
	userinfo := User{}
	json.Unmarshal(userinfoAsbytes, &userinfo)

	newasset := Asset{}
	json.Unmarshal([]byte(assetval), &newasset)
	newasset.Key = assetcompositeKey
	newasset.Status = 0
	newasset.Owned = usercompositekey
	newasset.Buyerspercent = make(map[string]float64)   // initialize `Buyerspercent`
	createAsset(APIstub, newasset)
	userinfo.Info.Assetlist = userinfo.Info.Assetlist + "|" + assetcompositeKey
	// Update user into ledger
	updatedUserinfoAsbytes, _ := json.Marshal(userinfo)
	if APIstub.PutState(usercompositekey, updatedUserinfoAsbytes) != nil {
		return shim.Error("AddAsset, Update user info fail!")
	}

	return shim.Success(nil)
}

/*
 * Call AssetManage.go 
 */
 func (s *Mortgageplatform) RemoveAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 { return shim.Error("AddAsset - incorrect number of arguments. Expecting 2")}

	userid := args[0]
	assetid := args[1]
	assetcompositeKey, _ := GetAssetKey(APIstub, assetid)
	assetexist, _ := checkAsset(APIstub, assetcompositeKey)
	if !assetexist { return shim.Error("RemoveAsset, asset NOT exists yet --- ") }

	usercompositekey, createUserKeyErr := GetUserKey(APIstub, userid)
	if createUserKeyErr != nil { return shim.Error("RemoveAsset, --- " + createUserKeyErr.Error()) }
	userinfoAsbytes, _ := APIstub.GetState(usercompositekey)
	userinfo := User{}
	json.Unmarshal(userinfoAsbytes, &userinfo)

	removeAsset(APIstub, assetcompositeKey) // remove from ledger
	splited := strings.Split(userinfo.Info.Assetlist, "|" + assetcompositeKey)
	userinfo.Info.Assetlist = splited[0] + splited[1]
	// Update user into ledger
	updatedUserinfoAsbytes, _ := json.Marshal(userinfo)
	if APIstub.PutState(usercompositekey, updatedUserinfoAsbytes) != nil {
		return shim.Error("AddAsset, Update user info fail!")
	}

	return shim.Success(nil)
}

/*
 * Call AssetManage.go 
 * Args: user := args[0]
 *		 assetid := args[1]
 *		 interestrate := args[2]
 *		 worth := args[3]
 *		 period := args[4]
 */
 func (s *Mortgageplatform) PostAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	res, err := postAsset(APIstub, args)
	if err != nil { return shim.Error("PostAsset, " + err.Error())}

	return shim.Success(res)
 }

 /*
  * Call AssetManage.go 
  */
 func (s *Mortgageplatform) UnpostAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	res, err := unpostAsset(APIstub, args)
	if err != nil { return shim.Error("UnpostAsset, " + err.Error())}

	return shim.Success(res)
 }

/*
 * Add payment to user
 * Args: 1. userid
 * 		 2. payment id
 *		 3. payment info   {"Paymenttype":"Credit card", "addr":{"addr1":"6650 Corporate","addr2":"apt.1815", "city":"Jax"}}
 */
func (s *Mortgageplatform) AddPayment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 { return shim.Error("AddPayment - incorrect number of arguments. Expecting 3") }

	userid := args[0]
	paymentid := args[1]
	paymentinfo := args[2]
	// Get user info
	userexist, userinfoAsbytes := CheckUser(APIstub, userid)
	if !userexist { return shim.Error("AddPayment - User not exists")}
	user := User{}
	json.Unmarshal(userinfoAsbytes, &user)
	// Get payment info
	payment := Paymentmethod{}
	if err := json.Unmarshal([]byte(paymentinfo), &payment); err != nil {
		return shim.Error("AddPayment - can NOT unmarshal payment info")
	}
	payment.Accountnumber = paymentid
	// Add payment method to user
	if err := user.AddPaymentmethod(APIstub, payment); err != nil {
		return shim.Error("AddPayment - " + err.Error())
	}

	res, _ := json.Marshal(user)
	return shim.Success(res)   // return the updated user info
}

func (s *Mortgageplatform)GetPostedAssets(APIstub shim.ChaincodeStubInterface) sc.Response {
	res, err := APIstub.GetState("postedasset")
	if err != nil {
		return shim.Error("GetPostedAssets, --- " + err.Error())
	}

	return shim.Success(res)
}

func (s *Mortgageplatform)BuyAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 { return shim.Error("AddPayment - incorrect number of arguments. Expecting 4") }
	userid := args[0]
	assetid := args[1]
	// get user
	userCompositekey, _ := GetUserKey(APIstub, userid)
	userAsbytes, _ := APIstub.GetState(userCompositekey)
	user := User{}
	if err := json.Unmarshal(userAsbytes, &user); err != nil {
		return shim.Error("BuyAsset - " + err.Error())
	}
	// get asset
	assetCompositekey, _ := GetAssetKey(APIstub, assetid)
	assetAsbytes, _ := APIstub.GetState(assetCompositekey)
	asset := Asset{}
	if err := json.Unmarshal(assetAsbytes, &asset); err != nil {
		return shim.Error("BuyAsset - " + err.Error())
	}
	if user.BuyAsset(APIstub, asset, args[2:]) != nil {
		return shim.Error("BuyAsset, error on BuyAsset")
	}

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(Mortgageplatform))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
