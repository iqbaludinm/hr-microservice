package helper

import (
	"strconv"

	"github.com/iqbaludinm/hr-microservice/profile-service/utils"
)

var (
	//host
	ServerHost = utils.GetEnv("SERVER_URI")
	ServerPort = utils.GetEnv("SERVER_PORT")

	//jwt
	SecretKey           = utils.GetEnv("SECRET_KEY")
	SecretKeyRefresh    = utils.GetEnv("SECRET_KEY_REFRESH")
	SessionLogin        = utils.GetEnv("SESSION_LOGIN")
	SessionRefreshToken = utils.GetEnv("SESSION_REFRESH_TOKEN")

	//variable
	PrefixFinger      = utils.GetEnv("PREFIX_FINGER")
	PrefixFace        = utils.GetEnv("PREFIX_FACE")
	VariableRoleAdmin = utils.GetEnv("ROLE_ADMIN")
	VariableRoleUser  = utils.GetEnv("ROLE_USER")
	MinioPrefix       = utils.GetEnv("MINIO_PREFIX")
	CategoryService   = utils.GetEnv("CATEGORY_SERVICE")
	
	DefaultLimit, _    = strconv.Atoi(utils.GetEnv("DEFAULT_LIMIT"))
)
