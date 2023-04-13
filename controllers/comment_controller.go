package controller

import (
	"fmt"
	"net/http"
	database "project4/databases"
	"project4/helpers"
	"project4/models"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func PostComment(ctx *gin.Context) {
	db := database.GetDB()

	Comment := models.Comment{}
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["id"].(float64))

	reqHeaders := helpers.GetRequestHeaders(ctx)
	if reqHeaders.ContentType == "application/json" {
		ctx.ShouldBindJSON(&Comment)
	} else {
		ctx.ShouldBind(&Comment)
	}

	Comment.UserID = userID

	err := db.Debug().Create(&Comment).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         Comment.ID,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoID,
		"user_id":    Comment.UserID,
		"created_at": Comment.CreatedAt,
		"updated_at": Comment.UpdatedAt,
	})
}

func GetComments(ctx *gin.Context) {
	fmt.Println("ini diatas")
	db := database.GetDB()
	fmt.Println("ini dibawah")
	comments := []models.Comment{}
	err := db.Preload(clause.Associations).Find(&comments).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	commentsResponse := []models.CommentsResponse{}

	for _, comment := range comments {
		response := models.CommentsResponse{}

		response.GormModel = comment.GormModel
		response.Message = comment.Message
		response.PhotoID = comment.PhotoID
		response.UserID = comment.UserID
		response.Photo.ID = comment.Photo.ID
		response.Photo.Title = comment.Photo.Title
		response.Photo.Caption = comment.Photo.Caption
		response.Photo.PhotoUrl = comment.Photo.PhotoUrl
		response.User.UserName = comment.User.UserName
		response.User.Email = comment.User.Email

		commentsResponse = append(commentsResponse, response)
	}

	ctx.JSON(http.StatusOK, commentsResponse)
}

func GetComment(ctx *gin.Context)  {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	comment := models.Comment{}
	
	commentID, err := strconv.Atoi(ctx.Param("commentID"))
	comment.UserID = int(userData["id"].(float64))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   "Bad request",
			"message": "Invalid ID",
		})
		return
	}

	// err = db.Where("id = ?", commentID).First(&comment).Error
	err = db.Preload(clause.Associations).Where("id = ?", commentID).First(&comment).Error
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := models.CommentsResponse{
		GormModel: comment.GormModel,
		Message:     comment.Message,
		PhotoID:   comment.PhotoID,
		UserID:  comment.UserID,
		Photo: struct{
			ID int `json:"photo_id"`; 
			Title string "json:\"title\""; 
			Caption string "json:\"caption\""; 
			PhotoUrl string " json:\"photo_url\""
		}{
			ID: comment.Photo.ID,
			Title: comment.Photo.Title,
			Caption: comment.Photo.Caption,
			PhotoUrl: comment.Photo.PhotoUrl,
		},
		User: struct {
			Email    string `json:"email"`
			UserName string `json:"user_name"`
		}{
			Email:    comment.User.Email,
			UserName: comment.User.UserName,
		},
	}

	// commentsResponse := []models.CommentsResponse{}

	// response := models.CommentsResponse{}

	// response.GormModel = comment.GormModel
	// response.Message = comment.Message
	// response.PhotoID = comment.PhotoID
	// response.UserID = comment.UserID
	// response.Photo.ID = comment.Photo.ID
	// response.Photo.Title = comment.Photo.Title
	// response.Photo.Caption = comment.Photo.Caption
	// response.Photo.PhotoUrl = comment.Photo.PhotoUrl
	// response.User.UserName = comment.User.UserName
	// response.User.Email = comment.User.Email

	// commentsResponse = append(commentsResponse, response)
	

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func UpdateComment(ctx *gin.Context) {
	db := database.GetDB()

	Comment := models.Comment{}
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["id"].(float64))
	CommentId, err := strconv.Atoi(ctx.Param("commentId"))
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
		ctx.ShouldBindJSON(&Comment)
	} else {
		ctx.ShouldBind(&Comment)
	}

	Comment.UserID = userID

	err = db.Model(&Comment).Where("id=?", CommentId).Updates(Comment).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         CommentId,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoID,
		"user_id":    Comment.UserID,
		"updated_at": Comment.UpdatedAt,
	})
}

func DeleteComment(ctx *gin.Context) {
	db := database.GetDB()

	Comment := models.Comment{}
	CommentId, err := strconv.Atoi(ctx.Param("commentId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   "Bad request",
			"message": "Invalid ID",
		})
		return
	}

	err = db.Where("id=?", CommentId).Delete(&Comment).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"error":   "Not Found",
			"message": "Comment not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your Comment has been succsessfuly deleted",
	})

}
