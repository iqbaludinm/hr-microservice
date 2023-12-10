package helper

import (
	"github.com/iqbaludinm/hr-microservice/user-service/utils"
)

var logger = utils.NewLogger()

func FatalIfError(err error) {
	if err != nil {
		logger.Error(err.Error())
	}
}
