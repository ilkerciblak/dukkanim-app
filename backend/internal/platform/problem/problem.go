package problem

import "net/http"

type Problem struct {
	StatusCode int
	Code       string
	Errors     map[string]string
	Detail     string
}

func (p *Problem) Error() string {

	return p.Detail
}

var (
	BadRequest                   = &Problem{StatusCode: http.StatusBadRequest, Detail: http.StatusText(http.StatusBadRequest)}
	Unauthorized                 = &Problem{StatusCode: http.StatusUnauthorized, Detail: http.StatusText(http.StatusUnauthorized)}
	PaymentRequired              = &Problem{StatusCode: http.StatusPaymentRequired, Detail: http.StatusText(http.StatusPaymentRequired)}
	Forbidden                    = &Problem{StatusCode: http.StatusForbidden, Detail: http.StatusText(http.StatusForbidden)}
	NotFound                     = &Problem{StatusCode: http.StatusNotFound, Detail: http.StatusText(http.StatusNotFound)}
	MethodNotAllowed             = &Problem{StatusCode: http.StatusMethodNotAllowed, Detail: http.StatusText(http.StatusMethodNotAllowed)}
	NotAcceptable                = &Problem{StatusCode: http.StatusNotAcceptable, Detail: http.StatusText(http.StatusNotAcceptable)}
	ProxyAuthRequired            = &Problem{StatusCode: http.StatusProxyAuthRequired, Detail: http.StatusText(http.StatusProxyAuthRequired)}
	RequestTimeout               = &Problem{StatusCode: http.StatusRequestTimeout, Detail: http.StatusText(http.StatusRequestTimeout)}
	Conflict                     = &Problem{StatusCode: http.StatusConflict, Detail: http.StatusText(http.StatusConflict)}
	Gone                         = &Problem{StatusCode: http.StatusGone, Detail: http.StatusText(http.StatusGone)}
	LengthRequired               = &Problem{StatusCode: http.StatusLengthRequired, Detail: http.StatusText(http.StatusLengthRequired)}
	PreconditionFailed           = &Problem{StatusCode: http.StatusPreconditionFailed, Detail: http.StatusText(http.StatusPreconditionFailed)}
	RequestEntityTooLarge        = &Problem{StatusCode: http.StatusRequestEntityTooLarge, Detail: http.StatusText(http.StatusRequestEntityTooLarge)}
	RequestURITooLong            = &Problem{StatusCode: http.StatusRequestURITooLong, Detail: http.StatusText(http.StatusRequestURITooLong)}
	UnsupportedMediaType         = &Problem{StatusCode: http.StatusUnsupportedMediaType, Detail: http.StatusText(http.StatusUnsupportedMediaType)}
	RequestedRangeNotSatisfiable = &Problem{StatusCode: http.StatusRequestedRangeNotSatisfiable, Detail: http.StatusText(http.StatusRequestedRangeNotSatisfiable)}
	ExpectationFailed            = &Problem{StatusCode: http.StatusExpectationFailed, Detail: http.StatusText(http.StatusExpectationFailed)}
	Teapot                       = &Problem{StatusCode: http.StatusTeapot, Detail: http.StatusText(http.StatusTeapot)}
	MisdirectedRequest           = &Problem{StatusCode: http.StatusMisdirectedRequest, Detail: http.StatusText(http.StatusMisdirectedRequest)}
	UnprocessableEntity          = &Problem{StatusCode: http.StatusUnprocessableEntity, Detail: http.StatusText(http.StatusUnprocessableEntity)}
	Locked                       = &Problem{StatusCode: http.StatusLocked, Detail: http.StatusText(http.StatusLocked)}
	FailedDependency             = &Problem{StatusCode: http.StatusFailedDependency, Detail: http.StatusText(http.StatusFailedDependency)}
	TooEarly                     = &Problem{StatusCode: http.StatusTooEarly, Detail: http.StatusText(http.StatusTooEarly)}
	UpgradeRequired              = &Problem{StatusCode: http.StatusUpgradeRequired, Detail: http.StatusText(http.StatusUpgradeRequired)}
	PreconditionRequired         = &Problem{StatusCode: http.StatusPreconditionRequired, Detail: http.StatusText(http.StatusPreconditionRequired)}
	TooManyRequests              = &Problem{StatusCode: http.StatusTooManyRequests, Detail: http.StatusText(http.StatusTooManyRequests)}
	RequestHeaderFieldsTooLarge  = &Problem{StatusCode: http.StatusRequestHeaderFieldsTooLarge, Detail: http.StatusText(http.StatusRequestHeaderFieldsTooLarge)}
	UnavailableForLegalReasons   = &Problem{StatusCode: http.StatusUnavailableForLegalReasons, Detail: http.StatusText(http.StatusUnavailableForLegalReasons)}

	InternalServerError           = &Problem{StatusCode: http.StatusInternalServerError, Detail: http.StatusText(http.StatusInternalServerError)}
	NotImplemented                = &Problem{StatusCode: http.StatusNotImplemented, Detail: http.StatusText(http.StatusNotImplemented)}
	BadGateway                    = &Problem{StatusCode: http.StatusBadGateway, Detail: http.StatusText(http.StatusBadGateway)}
	ServiceUnavailable            = &Problem{StatusCode: http.StatusServiceUnavailable, Detail: http.StatusText(http.StatusServiceUnavailable)}
	GatewayTimeout                = &Problem{StatusCode: http.StatusGatewayTimeout, Detail: http.StatusText(http.StatusGatewayTimeout)}
	HTTPVersionNotSupported       = &Problem{StatusCode: http.StatusHTTPVersionNotSupported, Detail: http.StatusText(http.StatusHTTPVersionNotSupported)}
	VariantAlsoNegotiates         = &Problem{StatusCode: http.StatusVariantAlsoNegotiates, Detail: http.StatusText(http.StatusVariantAlsoNegotiates)}
	InsufficientStorage           = &Problem{StatusCode: http.StatusInsufficientStorage, Detail: http.StatusText(http.StatusInsufficientStorage)}
	LoopDetected                  = &Problem{StatusCode: http.StatusLoopDetected, Detail: http.StatusText(http.StatusLoopDetected)}
	NotExtended                   = &Problem{StatusCode: http.StatusNotExtended, Detail: http.StatusText(http.StatusNotExtended)}
	NetworkAuthenticationRequired = &Problem{StatusCode: http.StatusNetworkAuthenticationRequired, Detail: http.StatusText(http.StatusNetworkAuthenticationRequired)}
)

func (p *Problem) WithMessage(msg string) *Problem {
	p.Detail = msg
	return p
}

func (p *Problem) WithValidation(errorMap map[string]string) *Problem {

	p.StatusCode = http.StatusUnprocessableEntity
	p.Errors = errorMap

	return p
}

func (p *Problem) WithCode(code string) *Problem {
	p.Code = code

	return p
}
