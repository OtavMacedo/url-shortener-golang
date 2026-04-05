package controller

import (
	"errors"
	"net/http"

	"github.com/OtavMacedo/url-shortener-golang/internal/apperr"
	"github.com/OtavMacedo/url-shortener-golang/internal/dto"
	"github.com/OtavMacedo/url-shortener-golang/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) Create(c *gin.Context) {
	var request dto.CreateUserInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := uc.service.Create(c.Request.Context(), request); err != nil {
		if errors.Is(err, apperr.ErrUserAlreadyExists) {
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

func (uc *UserController) Login(c *gin.Context) {
	var request dto.LoginInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	token, err := uc.service.Login(c, request)
	if err != nil {
		if errors.Is(err, apperr.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if errors.Is(err, apperr.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
