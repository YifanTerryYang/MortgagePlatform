package main

import ()

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
	Zip      string `json:'zip"`
}

type Asset struct {
	Key          string   `json:"key"`
	Status       string   `json:"status"` // "own", "pending"
	Owned        string   `json:"owned"`
	Interestrate string   `json:"interestrate"`
	Pending      []string `json:"pending"`
}

type Paymentmethod struct {
	Accountnumber string  `json:"accnumber"`
	Paymenttype   string  `json:"paymenttype"`
	Addr          Address `json:"addr"`
	Payassets     []Asset `json:"payassets"`
}

type User struct {
	Password []byte   `json:"password"`
	Info     UserInfo `json:"userinfo"`
}

type UserInfo struct {
	Fname             string          `json:"fname"`
	Lname             string          `json:"lname"`
	Gender            string          `json:"gender"`
	Addr              Address         `json:"addr"`
	Assetlist         []Asset         `json:"assetlist"`
	Paymentmethodlist []Paymentmethod `json:"paymentmethodlist"`
}
