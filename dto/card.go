package dto

type MobileCard struct {
	Id				string	`json:"Id"`
	Name			string  `json:"Name"`
	VendorCode		string	`json:"VendorCode"`
	Serial 			string 	`json:"Serial"`
	Code 			string 	`json:"Code"`
	Value 			int 	`json:"Value"`
	Status			int 	`json:"Status"`
	CreatedAt		string	`json:"CreatedAt"`
	LastUpdatedAt	string	`json:"LastUpdatedAt"`
}

type MobileCardFilter struct {
	Name			string	`json:"Name"`
	Value 			int 	`json:"Value"`
	Status			int 	`json:"Status"`
}


type UserExchange struct {
	UserId 			string 		`json:"UserId"`
	VendorName		string 		`json:"Vendor"`
	Value 			int 		`json:"Value"`
	Quantity		int			`json:"Quantity"`
}

type MobileCardVendor struct {
	Id 				string		`json:"Id"`
	Name 			string 		`json:"Name"`
	VendorCode		string		`json:"VendorCode"`
	Status 			int 		`json:"Status"`
	Value  			int			`json:"Value"`
	Quantity 		int 		`json:"Quantity"`
	CreatedAt		string		`json:"CreatedAt"`
	LastUpdatedAt	string		`json:"LastUpdatedAt"`
}