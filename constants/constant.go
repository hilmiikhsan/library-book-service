package constants

const (
	SuccessMessage                = "success"
	ErrFailedBadRequest           = "failed to parse request"
	ErrAuthorizationIsEmpty       = "authorization is empty"
	ErrInvalidAuthorizationFormat = "invalid authorization format"
	ErrInvalidAuthorization       = "invalid authorization"
	ErrBookNotFound               = "book not found"
	ErrParamIdIsRequired          = "param id is required"
	ErrIdIsNotValidUUID           = "id is not valid uuid"
	ErrInvalidFormatDate          = "invalid format date"
	ErrAuthRolePermission         = "you do not have permission to access this endpoint"
	ErrIsbnAlreadyExist           = "isbn already exist"
	ErrCategoryNotFound           = "category not found"
	ErrAuthorNotFound             = "author not found"
	ErrBookStockNotFound          = "book stock not found"
	ErrBookStockAlreadyExist      = "book stock already exist"
	ErrBookAlreadyBorrowed        = "book already borrowed"
	ErrInsufficientStock          = "insufficient stock"
	ErrBookAlreadyReturned        = "book already returned"
)

const (
	HeaderAuthorization = "Authorization"
	TokenTypeAccess     = "token"
	DateTimeFormat      = "2006-01-02"
	AuthRoleUser        = "User"
	AuthRoleAdmin       = "Admin"
)
