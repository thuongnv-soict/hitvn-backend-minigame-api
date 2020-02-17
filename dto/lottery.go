package dto

type LotteryResult struct {
	Special		string 	`json:"special"`
	First   	string 	`json:"first"`
	Second1 	string 	`json:"second1"`
	Second2 	string 	`json:"second2"`
	Third1  	string 	`json:"third1"`
	Third2  	string 	`json:"third2"`
	Third3  	string 	`json:"third3"`
	Third4  	string 	`json:"third4"`
	Third5  	string 	`json:"third5"`
	Third6  	string 	`json:"third6"`
	Fourth1 	string 	`json:"fourth1"`
	Fourth2 	string 	`json:"fourth2"`
	Fourth3 	string 	`json:"fourth3"`
	Fourth4 	string 	`json:"fourth4"`
	Fifth1 		string 	`json:"fifth1"`
	Fifth2  	string 	`json:"fifth2"`
	Fifth3  	string 	`json:"fifth3"`
	Fifth4		string 	`json:"fifth4"`
	Fifth5		string 	`json:"fifth5"`
	Fifth6		string 	`json:"fifth6"`
	Sixth1		string 	`json:"sixth1"`
	Sixth2		string 	`json:"sixth2"`
	Sixth3		string 	`json:"sixth3"`
	Seventh1	string 	`json:"seventh1"`
	Seventh2	string 	`json:"seventh2"`
	Seventh3	string 	`json:"seventh3"`
}

type LotteryPlayer struct {
	Id 				string 	`json:"Id"`
	UserId 			string	`json:"UserId"`
	NumberSelected	string	`json:"NumberSelected"`
	Date 			string  `json:"Date"`
	WalletId		string 	`json:"WalletId"`
}
