package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents standard API response structure
type Response struct {
	Code    int         `json:"code"`              // Application code
	Message string      `json:"message"`           // Human-readable message
	Data    interface{} `json:"data,omitempty"`    // Response data (optional)
	Error   string      `json:"error,omitempty"`   // Error details (optional)
}

// JSON sends response with custom code and message
func JSON(c *gin.Context, httpStatus int, code int, message string, data interface{}) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// JSONError sends error response
func JSONError(c *gin.Context, httpStatus int, code int, message string, err error) {
	resp := Response{
		Code:    code,
		Message: message,
	}

	if err != nil {
		resp.Error = err.Error()
	}

	c.AbortWithStatusJSON(httpStatus, resp)
}

// JSONErrorString sends error response with string error
func JSONErrorString(c *gin.Context, httpStatus int, code int, message string, errStr string) {
	c.AbortWithStatusJSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Error:   errStr,
	})
}

// Success sends 200 OK with custom message
func Success(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusOK, OK, message, data)
}

// SuccessOK sends 200 OK with default message
func SuccessOK(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, OK, GetMessage(OK), data)
}

// SuccessCreated sends 201 Created
func SuccessCreated(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, Created, GetMessage(Created), data)
}

// SuccessUpdated sends 200 OK with updated message
func SuccessUpdated(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, Updated, GetMessage(Updated), data)
}

// SuccessDeleted sends 200 OK with deleted message
func SuccessDeleted(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, Deleted, GetMessage(Deleted), data)
}

// BadRequest sends 400 Bad Request with custom message
func BadRequest(c *gin.Context, message string) {
	JSONErrorString(c, http.StatusBadRequest, InvalidData, message, "")
}

// BadRequestError sends 400 Bad Request with error
func BadRequestError(c *gin.Context, message string, err error) {
	JSONError(c, http.StatusBadRequest, InvalidData, message, err)
}

// ValidationError sends 400 with validation error
func ValidationError(c *gin.Context, message string) {
	JSONErrorString(c, http.StatusBadRequest, InvalidFormat, message, "")
}

// Unauthorized sends 401 Unauthorized
func Unauthorized(c *gin.Context, message string) {
	JSONErrorString(c, http.StatusUnauthorized, AuthRequired, message, "")
}

// InvalidCredentialsError sends 401 with invalid credentials
func InvalidCredentialsError(c *gin.Context) {
	JSONErrorString(c, http.StatusUnauthorized, InvalidCredentials, GetMessage(InvalidCredentials), "")
}

// InvalidTokenError sends 401 with invalid token
func InvalidTokenError(c *gin.Context) {
	JSONErrorString(c, http.StatusUnauthorized, InvalidToken, GetMessage(InvalidToken), "")
}

// ExpiredTokenError sends 401 with expired token
func ExpiredTokenError(c *gin.Context) {
	JSONErrorString(c, http.StatusUnauthorized, ExpiredToken, GetMessage(ExpiredToken), "")
}

// Forbidden sends 403 Forbidden
func Forbidden(c *gin.Context, message string) {
	JSONErrorString(c, http.StatusForbidden, AccessDenied, message, "")
}

// InsufficientRightsError sends 403 with insufficient rights
func InsufficientRightsError(c *gin.Context) {
	JSONErrorString(c, http.StatusForbidden, InsufficientRights, GetMessage(InsufficientRights), "")
}

// NotFoundError sends 404 Not Found
func NotFoundError(c *gin.Context, message string) {
	JSONErrorString(c, http.StatusNotFound, NotFound, message, "")
}

// AlreadyExistsError sends 409 Conflict
func AlreadyExistsError(c *gin.Context, message string) {
	JSONErrorString(c, http.StatusConflict, AlreadyExists, message, "")
}

// ConflictError sends 409 Conflict with error
func ConflictError(c *gin.Context, message string, err error) {
	JSONError(c, http.StatusConflict, Conflict, message, err)
}

// InternalErrorResponse sends 500 Internal Server Error
func InternalErrorResponse(c *gin.Context, message string, err error) {
	JSONError(c, http.StatusInternalServerError, InternalError, message, err)
}

// DatabaseErrorResponse sends 500 with database error
func DatabaseErrorResponse(c *gin.Context, err error) {
	JSONError(c, http.StatusInternalServerError, DatabaseError, GetMessage(DatabaseError), err)
}

// TokenGenerationError sends 500 with token generation error
func TokenGenerationError(c *gin.Context, err error) {
	JSONError(c, http.StatusInternalServerError, TokenGenerationFailed, GetMessage(TokenGenerationFailed), err)
}
