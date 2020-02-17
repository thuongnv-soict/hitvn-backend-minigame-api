package constant

const(
	ProgramInvitingName					string = "Inviting"
	ProgramInvitedName					string = "Invited"
	ProgramReadHITDaily					string = "ReadHITDaily"
	ProgramLotteryWinFirstPrize			string = "LotteryWinFirstPrize"
	ProgramLotteryWinAnyPrize			string = "LotteryWinAnyPrize"
	ProgramExchangeMobileCard10			string = "ExchangeMobileCard10"
	ProgramExchangeMobileCard20			string = "ExchangeMobileCard20"
	ProgramExchangeMobileCard50			string = "ExchangeMobileCard50"
	ProgramExchangeMobileCard100		string = "ExchangeMobileCard100"
	ProgramExchangeMobileCard200		string = "ExchangeMobileCard200"
	ProgramExchangeMobileCard500		string = "ExchangeMobileCard500"


	RedisPrefixKeyAllTransaction		string = "hitvn_bk_minigame_v1_transaction_all_"
	RedisPrefixKeyBoughtMobileCards		string = "hitvn_bk_minigame_v1_bought_mobile_card_"
	RedisPrefixKeyAllPrize				string = "hitvn_bk_minigame_v1_all_prize"
	RedisPrefixKeyAllVendor				string = "hitvn_bk_minigame_v1_all_vendor"
	RedisPrefixKeyUserWallet			string = "hitvn_bk_minigame_v1_user_wallet_"

	// RabbitMQ
	RbSuperExchange						string = "super_exchange"
	RbRouteResult 						string = "minigame.result.daily.lottery"

	/*
		STATUS FOR API
	 */

	// Status API Add Prize
	//AddPrizeSuccessfully				int = 890
	AddPrizeFailed						int = 891

	Successfully						int = 100
	Failed								int = 101

	// Status API Inviting Friends (Logic)
	UserHasInsertedCode					int = 801		// Invited User has inputted some code before
	InvitingUserDoesNotExisted			int = 802
	InvitingUserIsYou					int = 803
	InvitingProgramNotFound				int = 804

	/*
		ERROR CODE
	 */
	// Error Code SQL
	ServerError							int = 500
	CreatePrepareStatementFailed 		int = 501
	ExecuteQueryFailed					int = 502
	ParseResultSetFailed				int = 503
	GetRowsAffectFailed					int = 504


	DefaultMaximumSelectedLotteryNumbers	int = 3
	/*
		Mobile Card
	 */

	// Inviting Status Flag
	StatusMobileCardNotReady				int = 0
	StatusMobileCardReady					int = 1
	StatusMobileCardIsUsed					int = 2
	StatusMobileCardError					int = 10

	StatusMobileCardVendorNotActive			int = 0
	StatusMobileCardVendorActive			int = 1

	/*
		Date format
	 */
	DateTimeLayout							string = "02/01/2006"


)