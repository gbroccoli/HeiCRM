package response

// Response codes for HeiCRM API

const (
	// OK 1xxx - Success operations
	OK      = 1000 // Successful operation
	Created = 1001 // Resource created
	Updated = 1002 // Resource updated
	Deleted = 1003 // Resource deleted

	// InvalidData 2xxx - Data and validation errors
	InvalidData     = 2000 // Invalid data
	InvalidFormat   = 2001 // Invalid format
	MissingRequired = 2002 // Missing required field
	InvalidRange    = 2003 // Value out of range
	InvalidType     = 2004 // Invalid data type

	// AuthRequired 3xxx - Authentication errors
	AuthRequired          = 3000 // Authentication required
	InvalidCredentials    = 3001 // Invalid credentials
	InvalidToken          = 3002 // Invalid token
	ExpiredToken          = 3003 // Token expired
	TokenGenerationFailed = 3004 // Failed to generate token

	// AccessDenied 4xxx - Access errors
	AccessDenied       = 4000 // Access denied
	InsufficientRights = 4001 // Insufficient rights
	OperationForbidden = 4002 // Operation forbidden

	// NotFound 5xxx - Resource errors
	NotFound      = 5000 // Resource not found
	AlreadyExists = 5001 // Resource already exists
	Conflict      = 5002 // Data conflict
	Gone          = 5003 // Resource deleted/unavailable

	// OperationFailed 6xxx - Business logic errors
	OperationFailed     = 6000 // Operation failed
	InvalidState        = 6001 // Invalid state
	ConstraintViolation = 6002 // Constraint violation
	LimitExceeded       = 6003 // Limit exceeded

	// InternalError 9xxx - System errors
	InternalError        = 9000 // Internal error
	DatabaseError        = 9001 // Database error
	ServiceUnavailable   = 9002 // Service unavailable
	ConfigurationError   = 9003 // Configuration error
	ExternalServiceError = 9004 // External service error
)

// Messages contains human-readable messages for response codes
var Messages = map[int]string{
	// Success
	OK:      "Успешно",
	Created: "Ресурс успешно создан",
	Updated: "Ресурс успешно обновлён",
	Deleted: "Ресурс успешно удалён",

	// Data errors
	InvalidData:     "Некорректные данные",
	InvalidFormat:   "Некорректный формат",
	MissingRequired: "Отсутствует обязательное поле",
	InvalidRange:    "Значение вне допустимого диапазона",
	InvalidType:     "Некорректный тип данных",

	// Authentication
	AuthRequired:          "Требуется авторизация",
	InvalidCredentials:    "Неверные учётные данные",
	InvalidToken:          "Недействительный токен",
	ExpiredToken:          "Срок действия токена истёк",
	TokenGenerationFailed: "Не удалось сгенерировать токен",

	// Access
	AccessDenied:       "Доступ запрещён",
	InsufficientRights: "Недостаточно прав",
	OperationForbidden: "Операция запрещена",

	// Resources
	NotFound:      "Не найдено",
	AlreadyExists: "Уже существует",
	Conflict:      "Конфликт данных",
	Gone:          "Ресурс недоступен",

	// Business logic
	OperationFailed:     "Операция не выполнена",
	InvalidState:        "Некорректное состояние",
	ConstraintViolation: "Нарушение ограничения",
	LimitExceeded:       "Превышен лимит",

	// System
	InternalError:        "Внутренняя ошибка",
	DatabaseError:        "Ошибка базы данных",
	ServiceUnavailable:   "Сервис недоступен",
	ConfigurationError:   "Ошибка конфигурации",
	ExternalServiceError: "Ошибка внешнего сервиса",
}

// GetMessage returns message for given code
func GetMessage(code int) string {
	if msg, ok := Messages[code]; ok {
		return msg
	}
	return "Неизвестная ошибка"
}
