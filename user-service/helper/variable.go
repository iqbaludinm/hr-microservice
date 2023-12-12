package helper

import (
	"strconv"

	"github.com/iqbaludinm/hr-microservice/user-service/utils"
)

var (
	PermissionAllJob, _    = strconv.Atoi(utils.GetEnv("PERMISSION_ALL_JOB"))
	PermissionListJob, _   = strconv.Atoi(utils.GetEnv("PERMISSION_VIEW_JOB"))
	PermissionCreateJob, _ = strconv.Atoi(utils.GetEnv("PERMISSION_CREATE_JOB"))
	PermissionUpdateJob, _ = strconv.Atoi(utils.GetEnv("PERMISSION_UPDATE_JOB"))
	PermissionDeleteJob, _ = strconv.Atoi(utils.GetEnv("PERMISSION_DELETE_JOB"))

	RoleAdmin          = utils.GetEnv("ROLE_ADMIN_ID")
	RoleAdmin2         = utils.GetEnv("ROLE_ADMIN_2_ID")
	RoleSupervisor     = utils.GetEnv("ROLE_SUPERVISOR_ID")
	RoleSuperintendent = utils.GetEnv("ROLE_SUPERINTENDENT_ID")
	RoleSCM001M        = utils.GetEnv("ROLE_SCM_001M_ID")
	RoleMechanic       = utils.GetEnv("ROLE_MECHANIC_ID")
	DefaultLimit, _    = strconv.Atoi(utils.GetEnv("DEFAULT_LIMIT"))
	SafetyCheckID      = utils.GetEnv("SAFETY_CHECK_ID")
	gotenbergEndpoint  = utils.GetEnv("GOTENBERG_ENDPOINT")

	// reset-pass
	UrlReset = utils.GetEnv("URL_RESET_PASSWORD_LOCAL")
)
