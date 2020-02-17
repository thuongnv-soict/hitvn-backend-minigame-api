package dto

type User struct{
	Id				string	`json:"Id"`
	UserId			string	`json:"UserId"`
	Code 			string  `json:"Code"`
	IsInvited		bool 	`json:"IsInvited"`
	Wallet 			int 	`json:"Wallet"`
	ReadDaysOfWeek	[]int	`json:"ReadDaysOfWeek"`
	DayOfWeek		int 	`json:"DayOfWeek"`
	Description 	string 	`json:"Description"`
	Value 			int 	`json:"Value"`
	LastUpdatedAt	string 	`json:"LastUpdatedAt"`
}

type Invitation struct {
	UserId 		string 		`json:"UserId"`
	Code 		string 		`json:"Code"`
}

type Generation struct {
	PhoneNumber		string 		`json:"PhoneNumber"`
}

type Response struct {
	Message		string
	Data 		interface{}
}


