package entity

import "errors"

var (
	// Erros de usuário
	ErrUserNotFound        = errors.New("usuário não encontrado")
	ErrUserAlreadyExists   = errors.New("usuário já existe")
	ErrInvalidUserType     = errors.New("tipo de usuário inválido")
	ErrMerchantCannotSend  = errors.New("lojistas não podem enviar dinheiro")
	ErrInsufficientBalance = errors.New("saldo insuficiente")
	ErrInvalidCredentials  = errors.New("credenciais inválidas")
	ErrWeakPassword        = errors.New("senha muito fraca")

	// Erros de documento
	ErrInvalidCPF            = errors.New("CPF inválido")
	ErrInvalidCNPJ           = errors.New("CNPJ inválido")
	ErrDocumentAlreadyExists = errors.New("documento já cadastrado")
	ErrEmailAlreadyExists    = errors.New("email já cadastrado")

	// Erros de transação
	ErrTransactionNotFound         = errors.New("transação não encontrada")
	ErrInvalidAmount               = errors.New("valor inválido")
	ErrSelfTransfer                = errors.New("não é possível transferir para si mesmo")
	ErrTransactionNotPending       = errors.New("transação não está pendente")
	ErrTransactionNotAuthorized    = errors.New("transação não está autorizada")
	ErrTransactionAlreadyCompleted = errors.New("transação já foi concluída")
	ErrAmountExceedsLimit          = errors.New("valor excede o limite máximo")

	// Erros de autorização
	ErrAuthorizationFailed  = errors.New("autorização negada")
	ErrAuthorizationTimeout = errors.New("timeout na autorização")
	ErrAuthorizationService = errors.New("serviço de autorização indisponível")

	// Erros de notificação
	ErrNotificationFailed  = errors.New("falha ao enviar notificação")
	ErrNotificationTimeout = errors.New("timeout na notificação")
	ErrNotificationService = errors.New("serviço de notificação indisponível")

	// Erros de validação
	ErrValidationFailed = errors.New("falha na validação")
	ErrRequiredField    = errors.New("campo obrigatório")
	ErrInvalidFormat    = errors.New("formato inválido")

	// Erros de banco de dados
	ErrDatabaseConnection  = errors.New("falha na conexão com o banco de dados")
	ErrDatabaseTransaction = errors.New("falha na transação do banco de dados")
	ErrDatabaseConstraint  = errors.New("violação de restrição do banco de dados")
)

type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e BusinessError) Error() string {
	return e.Message
}

// NewBusinessError cria um novo erro de negócio
func NewBusinessError(code, message, details string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

const (
	ErrorCodeValidation    = "VALIDATION_ERROR"
	ErrorCodeNotFound      = "NOT_FOUND"
	ErrorCodeAlreadyExists = "ALREADY_EXISTS"
	ErrorCodeUnauthorized  = "UNAUTHORIZED"
	ErrorCodeForbidden     = "FORBIDDEN"
	ErrorCodeInternal      = "INTERNAL_ERROR"
	ErrorCodeExternal      = "EXTERNAL_SERVICE_ERROR"
)

var (
	ErrBusinessUserNotFound = NewBusinessError(
		ErrorCodeNotFound,
		"Usuário não encontrado",
		"O usuário solicitado não existe no sistema",
	)

	ErrBusinessDocumentExists = NewBusinessError(
		ErrorCodeAlreadyExists,
		"Documento já cadastrado",
		"Já existe um usuário com este CPF/CNPJ",
	)

	ErrBusinessEmailExists = NewBusinessError(
		ErrorCodeAlreadyExists,
		"Email já cadastrado",
		"Já existe um usuário com este email",
	)

	ErrBusinessMerchantCannotSend = NewBusinessError(
		ErrorCodeForbidden,
		"Lojistas não podem enviar dinheiro",
		"Apenas usuários comuns podem realizar transferências",
	)

	ErrBusinessInsufficientBalance = NewBusinessError(
		ErrorCodeForbidden,
		"Saldo insuficiente",
		"O usuário não possui saldo suficiente para realizar a transação",
	)

	ErrBusinessAuthorizationFailed = NewBusinessError(
		ErrorCodeExternal,
		"Transação não autorizada",
		"A transação foi negada pelo serviço de autorização",
	)
)
