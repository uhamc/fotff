package fotff

import (
	"fotff/pkg"
	"log"
)

func FindOutTheFirstFail(testCase string, m pkg.Manager, latestSuccessPkg string, failPkg string) string {
	if latestSuccessPkg == "" {
		log.Printf("can not get the latest success package for %s, stop analysing", testCase)
		return ""
	}
	panic("TODO")
}
