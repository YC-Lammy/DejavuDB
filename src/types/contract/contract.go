package contract

import (
	"errors"
	"strings"

	json "github.com/goccy/go-json"
)

type Contract struct {
	Account_Name  []byte
	Activated_By  []byte
	Activate_Date []byte

	Billing_City    [3]byte // iso 3166-2
	Billing_Country [2]byte // alpha 2 country code
	Billing_State   []byte
	Billing_Street  []byte
	Billing_Zipcode [20]byte

	Company_Signed_By   []byte
	Company_Signed_Date []byte
	Contract_Division   []byte
	Contract_End_Date   []byte
	Contract_Name       []byte
	Contract_Number     uint64
	Contract_Owner      []byte
	Contract_Start_Date []byte
	Contract_Term       int

	Created_By            []byte
	Customer_Signed_By    []byte
	Customer_Signed_Date  []byte
	Customer_Signed_Title []byte

	Description []byte

	Last_Modified_By []byte

	Owner_Expiration_Notice []byte

	Shipping_City    [3]byte // iso 3166-2
	Shipping_Country [2]byte // alpha 2 country code
	Shipping_State   []byte
	Shipping_Street  []byte
	Shipping_Zipcode [20]byte

	Special_Terms []byte

	Status []byte //Draft,Approval, and Activated
}

func NewContract(jsawn string) (*Contract, error) {
	c := Contract{}
	err := json.Unmarshal(jsawn, &c)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(string(c.Billing_Country[:])) == "cn" || strings.ToLower(string(c.Shipping_Country[:])) == "cn" {
		if strings.ToLower(string(c.Billing_City[:3])) == "tw" || strings.ToLower(string(c.Shipping_City[:3])) == "tw" {
			return nil, errors.New("Taiwan is a country")
		}
	}
	return &c, nil
}
