package main

import (
)

type Account struct {
	UserID    string `json:"id"`
	Firstname string `json:"fname"`
	Lastname  string `json:"lname"`
	Addr      string `json:"addr"` // Address
}

type Address struct {
	Address1 string `json:"addr1"`
	Address2 string `json:"addr2"`
	Apt      string `json:"apt"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
}

type Asset struct {
	Key          string   	`json:"key"`   // restore composite asset key
	Status       float64   	`json:"status"` // bought percentage = sum of Buyerspercent
	Owned        string   	`json:"owned"`  // owner name
	Desc		 string  	`json:"desc"`	// Asset description
	Interestrate float64   	`json:"interestrate"`
	Worth		 float64  	`json:"worth"`
	Period		 int64 		`json:"period"`
	Buyerspercent map[string]float64 `json:"buyers"`
}

type Paymentmethod struct {
	Accountnumber string  				`json:"accnumber"`
	Paymenttype   string  				`json:"paymenttype"`
	Addr          Address 				`json:"addr"`

	Payassets     map[string]Payassetrecord 	`json:"payassets"`   // [assetid]
}

type Payassetrecord struct {
	Period int			`json:"round"`
	Amountper float64 	`json:"everypay"`
}

type User struct {
	Username string   // restore hashed username
	Password []byte  
	Info     UserInfo 
}

type UserInfo struct {
	Fname             string          	`json:"fname"`
	Lname             string          	`json:"lname"`
	Gender            string          	`json:"gender"`
	Addr              Address         	`json:"addr"`
	Money			  float64			`json:"money"`
	Assetlist         string         	`json:"assetlist"`  // own assets list: restore hashed asset keys
	Paymentmethodlist map[string]Paymentmethod 	`json:"paymentmethodlist"` // paymentmethodlist
}
