package contract

type Contract struct {
	Account_Name  string
	Activated_By  string
	Activate_Date string

	Billing_City    string
	Billing_Country string
	Billing_State   string
	Billing_Street  string
	Billing_Zipcode uint32

	Company_Signed_By   string
	Company_Signed_Date string
	Contract_Division   string
	Contract_End_Date   string
	Contract_Name       string
	Contract_Number     string
	Contract_Owner      string
	Contract_Start_Date string
	Contract_Term       int

	Created_By            string
	Customer_Signed_By    string
	Customer_Signed_Date  string
	Customer_Signed_Title string

	Description string

	Last_Modified_By string

	Owner_Expiration_Notice string

	Shipping_City    string
	Shipping_Country string
	Shipping_State   string
	Shipping_Street  string
	Shipping_Zipcode string

	Special_Terms string

	Status string //Draft, In Approval Process, and Activated
}
