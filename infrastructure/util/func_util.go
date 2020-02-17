package util

import (
	"runtime"
)

/**
 * Returns a Name of interface's Func
 */
//func GetFunctionName(i interface{}) string {
//	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
//}

/**
 * Returns Name of a Func
 */
func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}