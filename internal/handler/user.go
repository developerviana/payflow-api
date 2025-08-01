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

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

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

func (h *UserHandler) ListUsers(c *gin.Context) {
	filters := &entity.UserFilters{}

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

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || str[0:len(substr)] == substr || str[len(str)-len(substr):] == substr)
}
