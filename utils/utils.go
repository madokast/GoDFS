package utils

import "github.com/madokast/GoDFS/utils/logger"

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIf(flag bool, infos ...interface{}) {
	if flag {
		logger.Error(infos...)
		panic(infos)
	}
}
