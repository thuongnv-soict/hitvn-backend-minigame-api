package gerror

func T(errorCode int) string {
	switch errorCode {
	//////////////////////////
	// Client-side
	//////////////////////////
	case ErrorBindData:
		return "Failed to bind data"
	case ErrorValidData:
		return "Failed to valid data"
	case ErrorNotFound:
		return "Item not found"
	//////////////////////////
	// Server-side
	//////////////////////////
	case ErrorConnect:
		return "Failed to connect database"
	case ErrorSaveData:
		return "Failed to save data"
	case ErrorRetrieveData:
		return "Failed to retrieve data"
	case ErrorLogin:
		return "Failed to login. Please try again!"
	case ErrorDeleteDataInCache:
		return "Failed to delete data in cache"


	//////////////////////////
	// Logic-side
	//////////////////////////
	//ErrorUserHasInsertedCode				int = 40010		// Invited User has inputted some code before
	//ErrorInvitingUserDoesNotExisted			int = 40011
	//ErrorInvitingUserIsYou					int = 40012
	//ErrorInvitingProgramNotFound			int = 40013
	case ErrorUserHasInsertedCode:
		return "Bạn đã nhập mã giới thiệu trước đó"
	case ErrorInvitingUserDoesNotExisted:
		return "Mã giới thiệu không đúng"
	case ErrorInvitingUserIsYou:
		return "Không nhập mã giới thiệu của mình"
	case ErrorInvitingProgramNotFound:
		return "Chương trình giới thiệu bạn bè không tồn tại"
	case ErrorMobileCardProgramNotFound:
		return "Chương trình đổi thẻ nạp không tồn tại"
	case ErrorNotEnoughCoin:
		return "Không đủ xu"
	case ErrorMobileCardNotExisted:
		return "Thẻ nạp này đã hết"
	case ErrorNotEnoughAvailableMobileCard:
		return "Không đủ số lượng thẻ nạp được yêu cầu"
	case ErrorReadDailyUserHasReceivedCoinToday:
		return "Bạn đã nhận xu ngày hôm nay"
	case ErrorReadDailyProgramNotFound:
		return "Chương trình đọc báo hằng ngày không tồn tại"

	case ErrorLotteryProgramNotFound:
		return "Chương trình xổ số không tồn tại"
	case ErrorLotteryExceedNumberOfSelected:
		return "Số lượng số chọn vượt quá cho phép"
	case ErrorLotteryDuplicatedSelectedNumber:
		return "Số này đã được chọn trước đó"
	case ErrorLotteryTimeUp:
		return "Đã hết thời gian chọn số trong ngày"
	}

	return "Unknown error"
}