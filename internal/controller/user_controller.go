package controller

import (
	"errors"
	"net/http"

	"github.com/OtavMacedo/url-shortener-golang/internal/model"
	"github.com/OtavMacedo/url-shortener-golang/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.UserService
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) Create(c *gin.Context) {
	var request createUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user := model.UserModel{
		Email:        request.Email,
		PasswordHash: request.Password,
	}

	if err := uc.service.Create(c.Request.Context(), user); err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "user already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}
