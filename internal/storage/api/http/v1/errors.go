package v1

// Errors
const (
	UnknownErrorCode    = 0
	UnknownErrorMessage = "unknown error"

	UserAlreadyExistsCode    = 1001
	UserAlreadyExistsMessage = "user already exists"

	UserNotFoundCode    = 1002
	UserNotFoundMessage = "user not found"

	UserRefreshTokenCookieNotFoundCode    = 1003
	UserRefreshTokenCookieNotFoundMessage = "user refresh token cookie not found"
	UserRefreshTokenExpiredCode           = 1004
	UserRefreshTokenExpiredMessage        = "user refresh token expired"

	UserEmailCodeDoesntExistsCode    = 1005
	UserEmailCodeDoesntExistsMessage = "email code doesnt exists"

	UserAlreadyHasCompanyErrorCode    = 1006
	UserAlreadyHasCompanyErrorMessage = "user already has a company"
	UserCompanyRelationErrorCode      = 1007
	UserCompanyRelationErrorMessage   = "user has no relation with the company"

	UserVerificationCodeDoesntExistsCode    = 1008
	UserVerificationCodeDoesntExistsMessage = "user verification code doesnt exists"

	CityAlreadyExistsCode    = 2001
	CityAlreadyExistsMessage = "city already exists"

	CompanyAlreadyExistsCode              = 3001
	CompanyAlreadyExistsMessage           = "company already exists"
	CompanyNotFoundErrorCode              = 3002
	CompanyNotFoundErrorMessage           = "company not found"
	CompanyAlreadySetCategoryErrorCode    = 3003
	CompanyAlreadySetCategoryErrorMessage = "company already has a category"
	CompanyUpdateProhibitedErrorCode      = 3004
	CompanyUpdateProhibitedErrorMessage   = "company update is prohibited"
	CompanyDeletedErrorCode               = 3005
	CompanyDeletedErrorMessage            = "company deleted"
	CompanyNotVerifiedCode                = 3006
	CompanyNotVerifiedMessage             = "company not verified"

	CategoryAlreadyExistsCode    = 4001
	CategoryAlreadyExistsMessage = "category already exists"
	CategoryNotFoundErrorCode    = 4002
	CategoryNotFoundErrorMessage = "category not found"

	AddressNotFoundErrorCode    = 5001
	AddressNotFoundErrorMessage = "address not found"

	OrderAlreadyExistsCode                       = 6001
	OrderAlreadyExistsMessage                    = "order already exists"
	UserOrderRelationErrorCode                   = 6002
	UserOrderRelationErrorMessage                = "user has no relation with the order"
	OrderNotFoundErrorCode                       = 6003
	OrderNotFoundErrorMessage                    = "order not found"
	OrderCreationErrorCode                       = 6004
	OrderCreationErrorMessage                    = "order creation error"
	OrderInappropriateEventErrorCode             = 6005
	OrderInappropriateEventErrorMessage          = "event inappropriate in current stage"
	InsufficientFundsToCreateAnOrderErrorCode    = 6006
	InsufficientFundsToCreateAnOrderErrorMessage = "insufficient funds to create an order"
	ActiveOrdersLimitExceededErrorCode           = 6007
	ActiveOrdersLimitExceededErrorMessage        = "active orders limit exceeded"
	OrderStatusAlreadyChangedErrorCode           = 6012
	OrderStatusAlreadyChangedErrorMessage        = "order status already changed"
	OrderEditCooldownErrorCode                   = 6013
	OrderEditCooldownErrorMessage                = "order edit cooldown"
	OrderAmountIsLessThanCityMinCode             = 6014
	OrderAmountIsLessThanCityMinMessage          = "order amount is less than city min"
	OrderNotActiveErrorCode                      = 6015
	OrderNotActiveErrorMessage                   = "order not active"
	OrderNotInProcessErrorCode                   = 6016
	OrderNotInProcessErrorMessage                = "order not in process"
	OrdersNotFoundErrorCode                      = 6017
	OrdersNotFoundErrorMessage                   = "orders not found"

	InvoiceCreationErrorCode    = 7001
	InvoiceCreationErrorMessage = "invoice creation error"

	EmployeeAlreadyExistsErrorCode                      = 8000
	EmployeeAlreadyExistsErrorMessage                   = "employee already exists"
	EmployeeDriveeUserNotFoundErrorCode                 = 8001
	EmployeeDriveeUserNotFoundErrorMessage              = "drivee user not found"
	EmployeeDriveeUserAlreadyBoundToCompanyErrorCode    = 8002
	EmployeeDriveeUserAlreadyBoundToCompanyErrorMessage = "drivee user already bound to company"

	DriverLocationNotFoundErrorCode    = 9000
	DriverLocationNotFoundErrorMessage = "driver location not found"

	CorpCompanyCardNotFoundErrorCode                    = 10000
	CorpCompanyCardNotFoundErrorMessage                 = "corp company card not found"
	CorpCompanyCardUnbindActiveOrdersExistsErrorCode    = 10001
	CorpCompanyCardUnbindActiveOrdersExistsErrorMessage = "there is one or more active orders at the moment"
	CorpCompanyCardAlreadyExistsErrorCode               = 10002
	CorpCompanyCardAlreadyExistsErrorMessage            = "company already have bound card"
)

// Highrate errors
const (
	OrderNoOneIsAvailableErrorCode        = 6008
	OrderNoOneIsAvailableErrorMessage     = "unfortunately, no one is available nearby at the moment"
	OrderUnattractivePriceErrorCode       = 6009
	OrderUnattractivePriceErrorMessage    = "the drivers aren't interested in the offered fare"
	OrderOfferAHigherFareErrorCode        = 6010
	OrderOfferAHigherFareErrorMessage     = "offer a higher fare"
	OrderRequestIsNotAcceptedErrorCode    = 6011
	OrderRequestIsNotAcceptedErrorMessage = "unfortunately, request is not accepted"
)

type ErrorCode int
type ErrorMessage string

type ErrorStruct struct {
	ErrorCode    `json:"error_code"`
	ErrorMessage `json:"error_message"`
} // @name ErrorStruct

type ValidationErrorStruct struct {
	ErrorCode    int               `json:"error_code"`
	ErrorMessage string            `json:"error_message"`
	Errors       []ValidationError `json:"validation_errors"`
}

type ValidationError struct {
	FieldKey     string `json:"field_key"`
	ErrorMessage string `json:"error_message"`
}

func getErrorStruct(code ErrorCode) *ErrorStruct {
	errorStruct := &ErrorStruct{
		ErrorCode:    UnknownErrorCode,
		ErrorMessage: UnknownErrorMessage,
	}

	switch code {
	case CityAlreadyExistsCode:
		errorStruct.ErrorCode = CityAlreadyExistsCode
		errorStruct.ErrorMessage = CityAlreadyExistsMessage

	case CompanyAlreadyExistsCode:
		errorStruct.ErrorCode = CompanyAlreadyExistsCode
		errorStruct.ErrorMessage = CompanyAlreadyExistsMessage
	case CompanyNotFoundErrorCode:
		errorStruct.ErrorCode = CompanyNotFoundErrorCode
		errorStruct.ErrorMessage = CompanyNotFoundErrorMessage
	case CompanyAlreadySetCategoryErrorCode:
		errorStruct.ErrorCode = CompanyAlreadySetCategoryErrorCode
		errorStruct.ErrorMessage = CompanyAlreadySetCategoryErrorMessage
	case CompanyUpdateProhibitedErrorCode:
		errorStruct.ErrorCode = CompanyUpdateProhibitedErrorCode
		errorStruct.ErrorMessage = CompanyUpdateProhibitedErrorMessage
	case CompanyDeletedErrorCode:
		errorStruct.ErrorCode = CompanyDeletedErrorCode
		errorStruct.ErrorMessage = CompanyDeletedErrorMessage
	case CompanyNotVerifiedCode:
		errorStruct.ErrorCode = CompanyNotVerifiedCode
		errorStruct.ErrorMessage = CompanyNotVerifiedMessage

	case CategoryAlreadyExistsCode:
		errorStruct.ErrorCode = CategoryAlreadyExistsCode
		errorStruct.ErrorMessage = CategoryAlreadyExistsMessage
	case CategoryNotFoundErrorCode:
		errorStruct.ErrorCode = CategoryNotFoundErrorCode
		errorStruct.ErrorMessage = CategoryNotFoundErrorMessage
	case UserAlreadyExistsCode:
		errorStruct.ErrorCode = UserAlreadyExistsCode
		errorStruct.ErrorMessage = UserAlreadyExistsMessage
	case UserEmailCodeDoesntExistsCode:
		errorStruct.ErrorCode = UserEmailCodeDoesntExistsCode
		errorStruct.ErrorMessage = UserEmailCodeDoesntExistsMessage
	case UserNotFoundCode:
		errorStruct.ErrorCode = UserNotFoundCode
		errorStruct.ErrorMessage = UserNotFoundMessage
	case UserRefreshTokenCookieNotFoundCode:
		errorStruct.ErrorCode = UserRefreshTokenCookieNotFoundCode
		errorStruct.ErrorMessage = UserRefreshTokenCookieNotFoundMessage
	case UserRefreshTokenExpiredCode:
		errorStruct.ErrorCode = UserRefreshTokenExpiredCode
		errorStruct.ErrorMessage = UserRefreshTokenExpiredMessage
	case UserAlreadyHasCompanyErrorCode:
		errorStruct.ErrorCode = UserAlreadyHasCompanyErrorCode
		errorStruct.ErrorMessage = UserAlreadyHasCompanyErrorMessage
	case UserCompanyRelationErrorCode:
		errorStruct.ErrorCode = UserCompanyRelationErrorCode
		errorStruct.ErrorMessage = UserCompanyRelationErrorMessage
	case UserVerificationCodeDoesntExistsCode:
		errorStruct.ErrorCode = UserVerificationCodeDoesntExistsCode
		errorStruct.ErrorMessage = UserVerificationCodeDoesntExistsMessage
	case AddressNotFoundErrorCode:
		errorStruct.ErrorCode = AddressNotFoundErrorCode
		errorStruct.ErrorMessage = AddressNotFoundErrorMessage
	case OrderAlreadyExistsCode:
		errorStruct.ErrorCode = OrderAlreadyExistsCode
		errorStruct.ErrorMessage = OrderAlreadyExistsMessage
	case UserOrderRelationErrorCode:
		errorStruct.ErrorCode = UserOrderRelationErrorCode
		errorStruct.ErrorMessage = UserOrderRelationErrorMessage
	case OrderNotFoundErrorCode:
		errorStruct.ErrorCode = OrderNotFoundErrorCode
		errorStruct.ErrorMessage = OrderNotFoundErrorMessage
	case OrderCreationErrorCode:
		errorStruct.ErrorCode = OrderCreationErrorCode
		errorStruct.ErrorMessage = OrderCreationErrorMessage
	case OrderInappropriateEventErrorCode:
		errorStruct.ErrorCode = OrderInappropriateEventErrorCode
		errorStruct.ErrorMessage = OrderInappropriateEventErrorMessage
	case InsufficientFundsToCreateAnOrderErrorCode:
		errorStruct.ErrorCode = InsufficientFundsToCreateAnOrderErrorCode
		errorStruct.ErrorMessage = InsufficientFundsToCreateAnOrderErrorMessage
	case ActiveOrdersLimitExceededErrorCode:
		errorStruct.ErrorCode = ActiveOrdersLimitExceededErrorCode
		errorStruct.ErrorMessage = ActiveOrdersLimitExceededErrorMessage
	case OrderNoOneIsAvailableErrorCode:
		errorStruct.ErrorCode = OrderNoOneIsAvailableErrorCode
		errorStruct.ErrorMessage = OrderNoOneIsAvailableErrorMessage
	case OrderUnattractivePriceErrorCode:
		errorStruct.ErrorCode = OrderUnattractivePriceErrorCode
		errorStruct.ErrorMessage = OrderUnattractivePriceErrorMessage
	case OrderOfferAHigherFareErrorCode:
		errorStruct.ErrorCode = OrderOfferAHigherFareErrorCode
		errorStruct.ErrorMessage = OrderOfferAHigherFareErrorMessage
	case OrderRequestIsNotAcceptedErrorCode:
		errorStruct.ErrorCode = OrderRequestIsNotAcceptedErrorCode
		errorStruct.ErrorMessage = OrderRequestIsNotAcceptedErrorMessage
	case OrderStatusAlreadyChangedErrorCode:
		errorStruct.ErrorCode = OrderStatusAlreadyChangedErrorCode
		errorStruct.ErrorMessage = OrderStatusAlreadyChangedErrorMessage
	case OrderEditCooldownErrorCode:
		errorStruct.ErrorCode = OrderEditCooldownErrorCode
		errorStruct.ErrorMessage = OrderEditCooldownErrorMessage
	case OrderAmountIsLessThanCityMinCode:
		errorStruct.ErrorCode = OrderAmountIsLessThanCityMinCode
		errorStruct.ErrorMessage = OrderAmountIsLessThanCityMinMessage
	case OrderNotActiveErrorCode:
		errorStruct.ErrorCode = OrderNotActiveErrorCode
		errorStruct.ErrorMessage = OrderNotActiveErrorMessage
	case OrderNotInProcessErrorCode:
		errorStruct.ErrorCode = OrderNotInProcessErrorCode
		errorStruct.ErrorMessage = OrderNotInProcessErrorMessage
	case OrdersNotFoundErrorCode:
		errorStruct.ErrorCode = OrdersNotFoundErrorCode
		errorStruct.ErrorMessage = OrdersNotFoundErrorMessage

	case InvoiceCreationErrorCode:
		errorStruct.ErrorCode = InvoiceCreationErrorCode
		errorStruct.ErrorMessage = InvoiceCreationErrorMessage

	case EmployeeAlreadyExistsErrorCode:
		errorStruct.ErrorCode = EmployeeAlreadyExistsErrorCode
		errorStruct.ErrorMessage = EmployeeAlreadyExistsErrorMessage
	case EmployeeDriveeUserNotFoundErrorCode:
		errorStruct.ErrorCode = EmployeeDriveeUserNotFoundErrorCode
		errorStruct.ErrorMessage = EmployeeDriveeUserNotFoundErrorMessage
	case EmployeeDriveeUserAlreadyBoundToCompanyErrorCode:
		errorStruct.ErrorCode = EmployeeDriveeUserAlreadyBoundToCompanyErrorCode
		errorStruct.ErrorMessage = EmployeeDriveeUserAlreadyBoundToCompanyErrorMessage

	case DriverLocationNotFoundErrorCode:
		errorStruct.ErrorCode = DriverLocationNotFoundErrorCode
		errorStruct.ErrorMessage = DriverLocationNotFoundErrorMessage

	case CorpCompanyCardNotFoundErrorCode:
		errorStruct.ErrorCode = CorpCompanyCardNotFoundErrorCode
		errorStruct.ErrorMessage = CorpCompanyCardNotFoundErrorMessage
	case CorpCompanyCardUnbindActiveOrdersExistsErrorCode:
		errorStruct.ErrorCode = CorpCompanyCardUnbindActiveOrdersExistsErrorCode
		errorStruct.ErrorMessage = CorpCompanyCardUnbindActiveOrdersExistsErrorMessage
	case CorpCompanyCardAlreadyExistsErrorCode:
		errorStruct.ErrorCode = CorpCompanyCardAlreadyExistsErrorCode
		errorStruct.ErrorMessage = CorpCompanyCardAlreadyExistsErrorMessage
	}

	return errorStruct
}
