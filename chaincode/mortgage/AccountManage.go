package main

import(
	"bytes"
	"encoding/json"
	//"golang.org/x/crypto/bcrypt"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func CheckUser(APIstub shim.ChaincodeStubInterface, username string) (bool, []byte) {
	userkey, err := GetUserKey(APIstub,username)
	if err != nil {
		return false, nil
	}

	val,_ := APIstub.GetState(userkey)
	fmt.Println("CheckUser, val is " + string(val))
	if val != nil {
		return true, val
	}

	return false, nil
}

func getUserInfo(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	username := args[0]
	check, userval := CheckUser(APIstub, username)
	if !check {   // if user exists
		return nil, errors.New("User not exists or password incorrect")
	}

	user := User{}
	json.Unmarshal(userval, &user)
	
	userinfoAsbytes, _ := json.Marshal(user.Info)

	return userinfoAsbytes, nil
}

func createNewUser(APIstub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 3 {
		return errors.New("Incorrect number of arguments. Expecting 3")
	}

	username := args[0]
	password := args[1]
	if password == "" {return errors.New("Password is empty")}
	//hashedPasswordAsbytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	hashedPasswordAsbytes := []byte(password)
	userinfo := args[2]

	// Check if user exists
	check, _ := CheckUser(APIstub, username)
	if check {   // if user exists
		return errors.New("User exist")
	}

	compositeKey, _ := GetUserKey(APIstub, username)
	userinfoEntity := UserInfo{}
	json.Unmarshal([]byte(userinfo), &userinfoEntity)
	newUser := User{Username: compositeKey, Password: hashedPasswordAsbytes, Info: userinfoEntity}
	newUserByte, _ := json.Marshal(newUser)
	if APIstub.PutState(compositeKey, newUserByte) != nil {
		return errors.New("Error on shim.PutState")
	}
	return nil
}

// Will return user info after login successfully
func login(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {  // Should contain "username" and "password"
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	username := args[0]
	password := args[1]
	if password == "" {return nil, errors.New("Password is empty")}
	//hashedPasswordAsbytes,_ := bcrypt.GenerateFromPassword([]byte(password), 10)
	hashedPasswordAsbytes := []byte(password)

	// Check if user exists
	check, userval := CheckUser(APIstub, username)
	if !check {   // if user exists
		return nil, errors.New("User not exists or password incorrect")
	}

	user := User{}
	json.Unmarshal(userval, &user)
	fmt.Println(hashedPasswordAsbytes)
	fmt.Println(user.Password)
	fmt.Println(string(hashedPasswordAsbytes))
	fmt.Println(string(user.Password))
	if bytes.Compare(hashedPasswordAsbytes, user.Password) != 0 {
		return nil, errors.New("User not exists or password incorrect")
	}
	userinfoAsbytes, _ := json.Marshal(user.Info)

	return userinfoAsbytes, nil
}

func updateUser(APIstub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 2 {
		return errors.New("Incorrect number of arguments. Expecting 2")
	}
	username := args[0]
	usernewinfo := args[1]

	check, useroldval := CheckUser(APIstub, username)
	if !check {
		return errors.New("User not exist")
	}
	
	olduserinfo := User{}
	newuserinfo := UserInfo{}
	json.Unmarshal(useroldval, &olduserinfo)
	json.Unmarshal([]byte(usernewinfo),&newuserinfo)
	olduserinfo.Info.Fname = newuserinfo.Fname
	olduserinfo.Info.Lname = newuserinfo.Lname
	olduserinfo.Info.Gender = newuserinfo.Gender
	// update address
	olduserinfo.Info.Addr.Address1 = newuserinfo.Addr.Address1
	olduserinfo.Info.Addr.Address2 = newuserinfo.Addr.Address2
	olduserinfo.Info.Addr.Apt = newuserinfo.Addr.Apt
	olduserinfo.Info.Addr.City = newuserinfo.Addr.City
	olduserinfo.Info.Addr.State = newuserinfo.Addr.State
	olduserinfo.Info.Addr.Zip = newuserinfo.Addr.Zip

	olduserinfoByte,_ := json.Marshal(olduserinfo)  // convert to byte stream after update
	userkey, err := GetUserKey(APIstub,username)
	if APIstub.PutState(userkey, olduserinfoByte) != nil {
		return errors.New("Fail to update user info" + err.Error())
	}

	return nil
}
