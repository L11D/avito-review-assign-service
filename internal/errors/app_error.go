package errors

type ErrorCode string

const (
	TEAM_EXISTS         ErrorCode = "TEAM_EXISTS"
	USER_EXISTS         ErrorCode = "USER_EXISTS"
	PR_EXISTS           ErrorCode = "PR_EXISTS"
	PR_MERGED           ErrorCode = "PR_MERGED"
	NOT_ASSIGNED        ErrorCode = "NOT_ASSIGNED"
	NO_CANDIDATE        ErrorCode = "NO_CANDIDATE"
	NOT_FOUND           ErrorCode = "NOT_FOUND"
	VALIDATION_FAILED   ErrorCode = "VALIDATION_FAILED"
	QUERY_PARAM_MISSING ErrorCode = "QUERY_PARAM_MISSING"
)

type AppError struct {
	Code       ErrorCode
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return string(e.Code) + " " + e.Message
}

func NewAppError(code ErrorCode, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func NewValidationFailedError(reason string) *AppError {
	return &AppError{
		Code:       VALIDATION_FAILED,
		Message:    reason,
		StatusCode: 400,
	}
}

func NewTeamExistsError(teamName string) *AppError {
	return &AppError{
		Code:       TEAM_EXISTS,
		Message:    "Team '" + teamName + "' already exists",
		StatusCode: 400,
	}
}

func NewUserExistsError(userName string) *AppError {
	return &AppError{
		Code:       USER_EXISTS,
		Message:    "User '" + userName + "' already exists",
		StatusCode: 400,
	}
}

func NewNotFoundError(entity string) *AppError {
	return &AppError{
		Code:       NOT_FOUND,
		Message:    entity + " not found",
		StatusCode: 404,
	}
}

func NewQueryParamMissingError(param string) *AppError {
	return &AppError{
		Code:       QUERY_PARAM_MISSING,
		Message:    "Query parameter '" + param + "' is missing",
		StatusCode: 400,
	}
}
