package dto

type Prize struct {
	Id				string	`json:"Id"`
	Name			string	`json:"Name"`
	Value 			int 	`json:"Value"`
	Description 	string 	`json:"Description"`
	CreatedAt 		string 	`json:"CreatedAt"`
	LastUpdatedAt	string 	`json:"LastUpdatedAt"`
}