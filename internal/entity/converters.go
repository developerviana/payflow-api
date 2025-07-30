package entity

// ToCreateUserResponse converte User para CreateUserResponse
func (u *User) ToCreateUserResponse() *CreateUserResponse {
	return &CreateUserResponse{
		ID:       u.ID,
		FullName: u.FullName,
		Document: u.maskDocument(),
		Email:    u.Email,
		UserType: u.UserType,
		Balance:  u.Balance.StringFixed(2),
	}
}

// ToGetUserResponse converte User para GetUserResponse
func (u *User) ToGetUserResponse() *GetUserResponse {
	return &GetUserResponse{
		ID:        u.ID,
		FullName:  u.FullName,
		Document:  u.maskDocument(),
		Email:     u.Email,
		UserType:  u.UserType,
		Balance:   u.Balance.StringFixed(2),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUserSummary converte User para UserSummary
func (u *User) ToUserSummary() *UserSummary {
	return &UserSummary{
		ID:       u.ID,
		FullName: u.FullName,
		Email:    u.Email,
		UserType: u.UserType,
	}
}

// ToBalanceResponse converte User para BalanceResponse
func (u *User) ToBalanceResponse() *BalanceResponse {
	return &BalanceResponse{
		UserID:    u.ID,
		Balance:   u.Balance.StringFixed(2),
		UpdatedAt: u.UpdatedAt,
	}
}

// maskDocument mascara o documento para exibição
func (u *User) maskDocument() string {
	if len(u.Document) == 11 { // CPF
		return u.Document[:3] + ".***.***-" + u.Document[9:]
	}
	if len(u.Document) == 14 { // CNPJ
		return u.Document[:2] + ".***.***/" + u.Document[8:12] + "-" + u.Document[12:]
	}
	return u.Document
}

// ToCreateTransactionResponse converte Transaction para CreateTransactionResponse
func (t *Transaction) ToCreateTransactionResponse() *CreateTransactionResponse {
	return &CreateTransactionResponse{
		ID:         t.ID,
		PayerID:    t.PayerID,
		PayeeID:    t.PayeeID,
		Amount:     t.Amount.StringFixed(2),
		Status:     t.Status,
		StatusDesc: t.GetStatusDescription(),
		CreatedAt:  t.CreatedAt,
	}
}

// ToGetTransactionResponse converte Transaction para GetTransactionResponse
func (t *Transaction) ToGetTransactionResponse() *GetTransactionResponse {
	response := &GetTransactionResponse{
		ID:               t.ID,
		PayerID:          t.PayerID,
		PayeeID:          t.PayeeID,
		Amount:           t.Amount.StringFixed(2),
		Status:           t.Status,
		StatusDesc:       t.GetStatusDescription(),
		AuthorizationID:  t.AuthorizationID,
		NotificationSent: t.NotificationSent,
		FailureReason:    t.FailureReason,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
		CompletedAt:      t.CompletedAt,
	}

	// Adiciona dados dos usuários se disponíveis
	if t.Payer != nil {
		response.Payer = t.Payer.ToUserSummary()
	}
	if t.Payee != nil {
		response.Payee = t.Payee.ToUserSummary()
	}

	return response
}

// FromCreateUserRequest converte CreateUserRequest para User
func FromCreateUserRequest(req *CreateUserRequest) (*User, error) {
	return NewUser(
		req.FullName,
		req.Document,
		req.Email,
		req.Password,
		req.UserType,
	)
}

// FromCreateTransactionRequest converte CreateTransactionRequest para Transaction
func FromCreateTransactionRequest(req *CreateTransactionRequest, payerID string) (*Transaction, error) {
	return NewTransaction(
		payerID,
		req.PayeeID,
		req.Amount,
	)
}

// ApplyUpdateUserRequest aplica as alterações do UpdateUserRequest ao User
func (u *User) ApplyUpdateUserRequest(req *UpdateUserRequest) error {
	if req.FullName != "" {
		u.FullName = req.FullName
	}
	
	if req.Email != "" {
		u.Email = req.Email
	}
	
	u.UpdatedAt = u.CreatedAt // time.Now() seria chamado no usecase
	
	return u.Validate()
}

// NewErrorResponse cria uma nova resposta de erro
func NewErrorResponse(message, code, details, field string, value interface{}) *ErrorResponse {
	return &ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
		Field:   field,
		Value:   value,
	}
}

// NewSuccessResponse cria uma nova resposta de sucesso
func NewSuccessResponse(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Message: message,
		Data:    data,
	}
}
