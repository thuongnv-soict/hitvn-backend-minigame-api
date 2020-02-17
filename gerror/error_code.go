package gerror

/********************************************************************/
/* Client-side Error Code											*/
/********************************************************************/
const (
	ErrorBindData			int = 40000
	ErrorValidData			int = 40001
)

/********************************************************************/
/* Server-side Error Code											*/
/********************************************************************/
const (
	ErrorConnect			int = 50000
	ErrorSaveData			int = 50001
	ErrorRetrieveData		int = 50002
	ErrorLogin				int = 50003
	ErrorNotFound			int = 50004
	ErrorDeleteDataInCache	int = 50005
)

/********************************************************************/
/* Logic Error Code											*/
/********************************************************************/
const (
	ErrorUserHasInsertedCode				int = 40010		// Invited User has inputted some code before
	ErrorInvitingUserDoesNotExisted			int = 40011
	ErrorInvitingUserIsYou					int = 40012
	ErrorInvitingProgramNotFound			int = 40013
	ErrorInvitedProgramNotFound				int = 40014


	ErrorNotEnoughCoin						int = 40020
	ErrorMobileCardNotExisted				int = 40021
	ErrorNotEnoughAvailableMobileCard		int = 40022
	ErrorMobileCardProgramNotFound			int = 40023

	ErrorReadDailyProgramNotFound			int = 40030
	ErrorReadDailyUserHasReceivedCoinToday	int = 40031

	ErrorLotteryProgramNotFound				int = 40040
	ErrorLotteryExceedNumberOfSelected		int = 40041
	ErrorLotteryDuplicatedSelectedNumber	int = 40042
	ErrorLotteryTimeUp						int = 40043
)