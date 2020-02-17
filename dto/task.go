package dto

type Task struct {
	MessageType string      `json:"MessageType"`
	Data		interface{} `json:"Data"`
}