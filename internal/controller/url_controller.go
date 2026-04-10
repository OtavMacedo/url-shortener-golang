package controller

import (
	"errors"
	"net/http"

	"github.com/OtavMacedo/url-shortener-golang/internal/apperr"
	"github.com/OtavMacedo/url-shortener-golang/internal/model"
	"github.com/OtavMacedo/url-shortener-golang/internal/repository"
	"github.com/OtavMacedo/url-shortener-golang/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	userID := c.MustGet("userID").(uuid.UUID)
	url := model.UrlModel{
		OriginalUrl: request.OriginalUrl,
		Slug:        request.Slug,
		UserID:      userID.String(),
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

func (uc *UrlController) Redirect(c *gin.Context) {
	slug := c.Param("slug")
	originalUrl, err := uc.service.FindBySlug(c, slug)
	if err != nil {
		if errors.Is(err, apperr.ErrSlugNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "url not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.Redirect(http.StatusFound, originalUrl)
}
