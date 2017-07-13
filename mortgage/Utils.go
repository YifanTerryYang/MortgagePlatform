package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"io"
)

// generate 32-bit MD5 string
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// generate GUID string
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func GetAssetKey(APIstub shim.ChaincodeStubInterface, assetKey string) (string, error) {
	result, err := APIstub.CreateCompositeKey("asset", []string{assetKey})
	if err != nil {
		return "", errors.New("Utils.GetAssetKey---" + err.Error())
	}
	return result, nil
}

func GetUserKey(APIstub shim.ChaincodeStubInterface, uname string) (string, error) {
	result, err := APIstub.CreateCompositeKey("user", []string{uname})
	if err != nil {
		return "", errors.New("Utils.GetUserKey---" + err.Error())
	}
	return result, nil
}
