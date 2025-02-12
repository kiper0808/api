package v1

// Errors
const (
	UnknownErrorCode    = 0
	UnknownErrorMessage = "unknown error"

	StorageAlreadyExistsCode    = 1001
	StorageAlreadyExistsMessage = "storage already exists"
)

type ErrorCode int
type ErrorMessage string

type ErrorStruct struct {
	ErrorCode    `json:"error_code"`
	ErrorMessage `json:"error_message"`
} // @name ErrorStruct

func getErrorStruct(code ErrorCode) *ErrorStruct {
	errorStruct := &ErrorStruct{
		ErrorCode:    UnknownErrorCode,
		ErrorMessage: UnknownErrorMessage,
	}

	switch code {
	case StorageAlreadyExistsCode:
		errorStruct.ErrorCode = StorageAlreadyExistsCode
		errorStruct.ErrorMessage = StorageAlreadyExistsMessage
	}

	return errorStruct
}
