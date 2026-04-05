package controller

import (
	"errors"
	"net/http"

	"github.com/OtavMacedo/url-shortener-golang/internal/model"
	"github.com/OtavMacedo/url-shortener-golang/internal/repository"
	"github.com/OtavMacedo/url-shortener-golang/internal/service"
	"github.com/gin-gonic/gin"
)

type UrlController struct {
	service *service.UrlService
}

func NewUrlService(service *service.UrlService) *UrlController {
	return &UrlController{service: service}
}

type createUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required,url"`
	Slug        string `json:"slug"`
}

func (uc *UrlController) Create(c *gin.Context) {
	var request createUrlRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	url := model.UrlModel{
		OriginalUrl: request.OriginalUrl,
		Slug:        request.Slug,
		UserID:      "019d5557-5e7f-7363-ac5a-52f4696a2afe",
	}

	if err := uc.service.Create(c.Request.Context(), url); err != nil {
		if errors.Is(err, repository.ErrUrlSlugConflict) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "slug already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create slug",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "slug created successfully",
	})
}
