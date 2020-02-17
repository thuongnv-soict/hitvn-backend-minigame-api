package util

import "g-tech.com/constant"

func GetExchangeMobileCardProgramName(value int) string{
	switch value {
	case 10000:
		return constant.ProgramExchangeMobileCard10
	case 20000:
		return constant.ProgramExchangeMobileCard20
	case 50000:
		return constant.ProgramExchangeMobileCard50
	case 100000:
		return constant.ProgramExchangeMobileCard100
	case 200000:
		return constant.ProgramExchangeMobileCard200
	case 500000:
		return constant.ProgramExchangeMobileCard500
	default:
		return "unknown"
	}
}
