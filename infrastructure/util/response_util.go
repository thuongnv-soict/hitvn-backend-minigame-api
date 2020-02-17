package util

import (
	"g-tech.com/constant"
)

func GetMessage(statusCode int) string{
	message := ""

	switch statusCode {
	// SQL Error
	case constant.ExecuteQueryFailed:
		message = "ServerError"
	case constant.CreatePrepareStatementFailed:
		message = "ServerError"
	case constant.ParseResultSetFailed:
		message = "ServerError"
	case constant.GetRowsAffectFailed:
		message = "ServerError"

	/*
		ERROR CODE
	 */
	// Inviting
	case constant.UserHasInsertedCode:
		message = "UserHasInsertedCode"

	case constant.InvitingUserIsYou:
		message = "InvitingUserIsYou"

	case constant.InvitingUserDoesNotExisted:
		message = "InvitingUserDoesNotExisted"

	case constant.InvitingProgramNotFound:
		message = "InvitingProgramNotFound"


	/*
		RESULT
	 */
	// Inviting Result
	case constant.Successfully:
		message = "Successfully"
	case constant.Failed:
		message = "Failed"







	default:
		message = "Unknown"
	}

	return message
}
