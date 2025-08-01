package handler

import (
	"net/http"
	"strconv"

	"payflow-api/internal/entity"
	"payflow-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler cria uma nova instância do handler
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUser godoc
// @Summary Criar um novo usuário
// @Description Cria um novo usuário no sistema com validações de CPF/CNPJ e email únicos
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.CreateUserRequest true "Dados do usuário"
// @Success 201 {object} entity.CreateUserResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 409 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req entity.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.NewErrorResponse(
			"Dados inválidos",
			"INVALID_REQUEST",
			err.Error(),
			"",
			nil,
		))
		return
	}

	response, err := h.userUseCase.CreateUser(c.Request.Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		// Determinar código de status baseado no erro
		if contains(err.Error(), "já existe") {
			statusCode = http.StatusConflict
			code = "USER_ALREADY_EXISTS"
		} else if contains(err.Error(), "validar dados") {
			statusCode = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		}

		c.JSON(statusCode, entity.NewErrorResponse(
			err.Error(),
			code,
			"",
			"",
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser godoc
// @Summary Buscar usuário por ID
// @Description Retorna os dados de um usuário específico
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} entity.GetUserResponse
// @Failure 404 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, entity.NewErrorResponse(
			"ID do usuário é obrigatório",
			"MISSING_ID",
			"",
			"id",
			nil,
		))
		return
	}

	response, err := h.userUseCase.GetUser(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		if contains(err.Error(), "não encontrado") {
			statusCode = http.StatusNotFound
			code = "USER_NOT_FOUND"
		}

		c.JSON(statusCode, entity.NewErrorResponse(
			err.Error(),
			code,
			"",
			"",
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUser godoc
// @Summary Atualizar usuário
// @Description Atualiza os dados de um usuário
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param user body entity.UpdateUserRequest true "Dados para atualização"
// @Success 200 {object} entity.GetUserResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 404 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, entity.NewErrorResponse(
			"ID do usuário é obrigatório",
			"MISSING_ID",
			"",
			"id",
			nil,
		))
		return
	}

	var req entity.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.NewErrorResponse(
			"Dados inválidos",
			"INVALID_REQUEST",
			err.Error(),
			"",
			nil,
		))
		return
	}

	response, err := h.userUseCase.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		if contains(err.Error(), "não encontrado") {
			statusCode = http.StatusNotFound
			code = "USER_NOT_FOUND"
		} else if contains(err.Error(), "validar dados") {
			statusCode = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		}

		c.JSON(statusCode, entity.NewErrorResponse(
			err.Error(),
			code,
			"",
			"",
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListUsers godoc
// @Summary Listar usuários
// @Description Lista usuários com paginação e filtros
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Limite por página" default(20)
// @Param user_type query string false "Tipo de usuário" Enums(common, merchant)
// @Param email query string false "Filtro por email (busca parcial)"
// @Success 200 {object} entity.ListUsersResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	filters := &entity.UserFilters{}

	// Parse query parameters
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filters.Page = p
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filters.Limit = l
		}
	}

	if userType := c.Query("user_type"); userType != "" {
		filters.UserType = entity.UserType(userType)
	}

	if email := c.Query("email"); email != "" {
		filters.Email = email
	}

	response, err := h.userUseCase.ListUsers(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.NewErrorResponse(
			err.Error(),
			"INTERNAL_ERROR",
			"",
			"",
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Deletar usuário
// @Description Remove um usuário do sistema
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 204 "No Content"
// @Failure 404 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, entity.NewErrorResponse(
			"ID do usuário é obrigatório",
			"MISSING_ID",
			"",
			"id",
			nil,
		))
		return
	}

	err := h.userUseCase.DeleteUser(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		if contains(err.Error(), "não encontrado") {
			statusCode = http.StatusNotFound
			code = "USER_NOT_FOUND"
		}

		c.JSON(statusCode, entity.NewErrorResponse(
			err.Error(),
			code,
			"",
			"",
			nil,
		))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetBalance godoc
// @Summary Consultar saldo
// @Description Retorna o saldo atual de um usuário
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} entity.BalanceResponse
// @Failure 404 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /users/{id}/balance [get]
func (h *UserHandler) GetBalance(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, entity.NewErrorResponse(
			"ID do usuário é obrigatório",
			"MISSING_ID",
			"",
			"id",
			nil,
		))
		return
	}

	response, err := h.userUseCase.GetBalance(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		if contains(err.Error(), "não encontrado") {
			statusCode = http.StatusNotFound
			code = "USER_NOT_FOUND"
		}

		c.JSON(statusCode, entity.NewErrorResponse(
			err.Error(),
			code,
			"",
			"",
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// contains verifica se uma string contém uma substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || str[0:len(substr)] == substr || str[len(str)-len(substr):] == substr)
}
