package util

import (
	"github.com/golang/glog"
)

// PanicErr used to panic when err != nil
func PanicErr(err error) {

	if err != nil {
		panic(err)
	}
}

// RecoverErr used to recover err and display msg
func RecoverErr(frontMsg string) {

	if r := recover(); r != nil {
		glog.Warning(frontMsg, r)
	}
}
