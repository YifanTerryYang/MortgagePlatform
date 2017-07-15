package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type Mortgageplatform struct{
	//var assetsListkey, usersListkey, postedAssetsListkey string
	assetsList, usersList, postedAssetsList []string
}

func (s *Mortgageplatform) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	// init Assets list
	//s.assetsListkey = "assetslist"
	//s.assetsListval = make([]string,0,2000)
	// init Users list
	//s.usersListkey = "userslist"
	//s.usersListval = make([]string,0,2000)
	// init Posted assets list
	//s.postedAssetsListkey = "postedassetslist"
	//s.postedAssetsListval = make([]string,0,2000)

	// Initialize "posted Asset" value.
	APIstub.PutState("postedasset", []byte(""))
	return shim.Success(nil)
}

func (s *Mortgageplatform) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	fmt.Println("ARGS length is: ", len(args))
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createNewUser" {
		return s.createNewUser(APIstub, args)
	} else if function == "getAccountsList" {
		return s.getAccountsList(APIstub)
	} else if function == "findAccount" {
		return s.findAccount(APIstub, args)
	} else if function == "updateAccount" {
		return s.updateAccount(APIstub, args)
	} else if function == "random" {
		return s.random(APIstub)
	} else if function == "helloWorld" {
		return s.helloWorld(APIstub)
	} else if function == "test" {
		return test(APIstub, "args")
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func test(APIstub shim.ChaincodeStubInterface, s string) sc.Response {
	fmt.Print("second para is: " + s)
	fmt.Print("test function")

	return shim.Success(nil)
}

func (s *Mortgageplatform) createNewUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var newAccount = Account{Firstname: args[0], Lastname: args[1], Addr: args[2]}

	newAccount.UserID = UniqueId()
	accountAsBytes, _ := json.Marshal(newAccount)
	APIstub.PutState(newAccount.UserID, accountAsBytes)
	fmt.Println("Added", newAccount)

	return shim.Success(nil)
}

func (s *Mortgageplatform) getAccountsList(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "00000000000000000000000000000000"
	endKey := "ffffffffffffffffffffffffffffffff"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllAccounts:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *Mortgageplatform) helloWorld(APIstub shim.ChaincodeStubInterface) sc.Response {
	var key, err = APIstub.CreateCompositeKey("user", []string{"yifan", "Yang", "male"})
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("key is " + key)
	//return shim.Success([]byte("{\"Hello world!\"}"))
	return shim.Success([]byte(key))
}

func (s *Mortgageplatform) findAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	accountAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(accountAsBytes)
}

func (s *Mortgageplatform) updateAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	accountAsBytes, _ := APIstub.GetState(args[0])
	account := Account{}

	json.Unmarshal(accountAsBytes, &account)
	account.Firstname = args[1]
	account.Lastname = args[2]
	account.Addr = args[3]

	accountAsBytes, _ = json.Marshal(account)
	APIstub.PutState(args[0], accountAsBytes)

	return shim.Success(nil)
}

func (s *Mortgageplatform) random(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success([]byte("{\"Random function not implemented yet!\"}"))
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(Mortgageplatform))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
