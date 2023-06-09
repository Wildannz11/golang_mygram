package controller

import (
	"net/http"
	database "project4/databases"
	"project4/helpers"
	"project4/models"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func PostPhoto(ctx *gin.Context) {
	db := database.GetDB()

	Photo := models.Photo{}
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["id"].(float64))

	reqHeaders := helpers.GetRequestHeaders(ctx)
	if reqHeaders.ContentType == "application/json" {
		ctx.ShouldBindJSON(&Photo)
	} else {
		ctx.ShouldBind(&Photo)
	}

	Photo.UserID = userID

	err := db.Debug().Create(&Photo).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         Photo.ID,
		"title":      Photo.Title,
		"captions":   Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserID,
		"created_at": Photo.CreatedAt,
	})
}

func GetPhotos(ctx *gin.Context) {
	db := database.GetDB()

	photos := []models.Photo{}
	err := db.Preload(clause.Associations).Find(&photos).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	photosResponse := []models.PhotosResponse{}
	for _, photo := range photos {
		response := models.PhotosResponse{}
		response.ID = photo.ID
		response.CreatedAt = photo.CreatedAt
		response.UpdatedAt = photo.UpdatedAt
		response.Title = photo.Title
		response.Caption = photo.Caption
		response.PhotoUrl = photo.PhotoUrl
		response.UserID = photo.UserID
		response.User.Email = photo.User.Email
		response.User.UserName = photo.User.UserName

		photosResponse = append(photosResponse, response)
	}
	ctx.JSON(http.StatusOK, photosResponse)
}

func GetPhoto(ctx *gin.Context)  {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	photo := models.Photo{}
	
	photoID, err := strconv.Atoi(ctx.Param("photoId"))
	photo.UserID = int(userData["id"].(float64))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   "Bad request",
			"message": "Invalid ID",
		})
		return
	}

	// err = db.Where("id = ?", photoID).First(&photo).Error
	err = db.Preload(clause.Associations).Where("id = ?", photoID).First(&photo).Error
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := models.PhotosResponse{
		GormModel: photo.GormModel,
		Title:     photo.Title,
		Caption:   photo.Caption,
		PhotoUrl:  photo.PhotoUrl,
		UserID:    photo.UserID,
		User: struct {
			Email    string `json:"email"`
			UserName string `json:"user_name"`
		}{
			Email:    photo.User.Email,
			UserName: photo.User.UserName,
		},
	}

	// photosResponse := []models.PhotosResponse{}
	// photosResponse := make([]models.PhotosResponse, 0)
	
	// response := models.PhotosResponse{}
	// response.ID = photo.ID
	// response.CreatedAt = photo.CreatedAt
	// response.UpdatedAt = photo.UpdatedAt
	// response.Title = photo.Title
	// response.Caption = photo.Caption
	// response.PhotoUrl = photo.PhotoUrl
	// response.UserID = photo.UserID
	// response.User.Email = photo.User.Email
	// response.User.UserName = photo.User.UserName

	// photosResponse = append(photosResponse, response)
	

	ctx.JSON(http.StatusOK, response)
}


func UpdatePhoto(ctx *gin.Context) {
	db := database.GetDB()

	Photo := models.Photo{}
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["id"].(float64))
	photoId, err := strconv.Atoi(ctx.Param("photoId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   "Bad request",
			"message": "Invalid ID",
		})
		return
	}

	reqHeaders := helpers.GetRequestHeaders(ctx)
	if reqHeaders.ContentType == "application/json" {
		ctx.ShouldBindJSON(&Photo)
	} else {
		ctx.ShouldBind(&Photo)
	}

	Photo.UserID = userID

	err = db.Model(&Photo).Where("id=?", photoId).Updates(Photo).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         photoId,
		"title":      Photo.Title,
		"captions":   Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserID,
		"updated_at": Photo.UpdatedAt,
	})
}

func DeletePhoto(ctx *gin.Context) {
	db := database.GetDB()

	Photo := models.Photo{}
	photoId, err := strconv.Atoi(ctx.Param("photoId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   "Bad request",
			"message": "Invalid ID",
		})
		return
	}

	err = db.Where("id=?", photoId).Delete(&Photo).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"error":   "Not Found",
			"message": "Photo not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been succsessfuly deleted",
	})

}
