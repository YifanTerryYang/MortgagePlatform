package main

import(
	"encoding/json"
	"golang.org/x/crypto/bcrypt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func CheckUser(APIstub shim.ChaincodeStubInterface, username string) (bool, []byte) {
	userkey, err := GetUserKey(APIstub,username)
	if err != nil {
		return false, nil
	}

	val,_ := APIstub.GetState(userkey)
	if val != nil {
		return true, val
	}

	return false, nil
}

func CreateNewUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	username := args[0]
	password := args[1]
	if password == "" {return shim.Error("Password is empty")}
	hashedPassword,_ := bcrypt.GenerateFromPassword([]byte(password), 10)
	userinfo := args[2]

	// Check if user exists
	check, _ := CheckUser(APIstub, username)
	if check {   // if user exists
		return shim.Error("User exist")
	}

	compositeKey, _ := GetUserKey(APIstub, username)
	userinfoEntity := UserInfo{}
	json.Unmarshal([]byte(userinfo), &userinfoEntity)
	newUser := User{Password: hashedPassword, Info: userinfoEntity}
	newUserByte, _ := json.Marshal(newUser)
	if APIstub.PutState(compositeKey, newUserByte) != nil {
		return shim.Error("Error on shim.PutState")
	}
	return shim.Success(nil)
}

func Login(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {  // Should contain "username" and "password"
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	
	username := args[0]
	password := args[1]
	if password == "" {return shim.Error("Password is empty")}
	hashedPassword,_ := bcrypt.GenerateFromPassword([]byte(password), 10)

	// Check if user exists
	check, userval := CheckUser(APIstub, username)
	if !check {   // if user exists
		return shim.Error("User not exists or password incorrect")
	}

	user := User{}
	json.Unmarshal(userval, &user)
	userinfo, _ := json.Marshal(user.Info)

	return shim.Success(userinfo)
}

func UpdateUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	username := args[0]
	usernewinfo := args[1]

	check, useroldval := CheckUser(APIstub, username)
	if !check {
		return shim.Error("User not exist")
	}
	
	olduserinfo := User{}
	newuserinfo := UserInfo{}
	json.Unmarshal(useroldval, &olduserinfo)
	json.Unmarshal([]byte(usernewinfo),&newuserinfo)
	if newuserinfo.Fname != "" {olduserinfo.Info.Fname = newuserinfo.Fname}
	if newuserinfo.Lname != "" {olduserinfo.Info.Lname = newuserinfo.Lname}
	if newuserinfo.Gender != "" {olduserinfo.Info.Gender = newuserinfo.Gender}
	// update address
	if newuserinfo.Addr.Address1 != "" {olduserinfo.Info.Addr.Address1 = newuserinfo.Addr.Address1}
	if newuserinfo.Addr.Address2 != "" {olduserinfo.Info.Addr.Address2 = newuserinfo.Addr.Address2}
	if newuserinfo.Addr.Apt != "" {olduserinfo.Info.Addr.Apt = newuserinfo.Addr.Apt}
	if newuserinfo.Addr.City != "" {olduserinfo.Info.Addr.City = newuserinfo.Addr.City}
	if newuserinfo.Addr.State != "" {olduserinfo.Info.Addr.State = newuserinfo.Addr.State}
	if newuserinfo.Addr.Zip != "" {olduserinfo.Info.Addr.Zip = newuserinfo.Addr.Zip}

	olduserinfoByte,_ := json.Marshal(olduserinfo)
	userkey, err := GetUserKey(APIstub,username)
	if APIstub.PutState(userkey, olduserinfoByte) != nil {
		return shim.Error("Fail to update user info")
	}

	return shim.Success(nil)
}
