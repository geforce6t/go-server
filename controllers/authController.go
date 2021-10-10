package controllers

import (
	"net/http"

	"github.com/geforce6t/go-server/models"
	"github.com/geforce6t/go-server/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// controller to register a user
func RegisterUser(c *gin.Context, db *gorm.DB) {
	var res = RegisterResponse{}

	if err := c.ShouldBindJSON(&res); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid response!",
		})
		return
	}

	result := db.Where("email = ?", res.Email).First(&models.User{})
	if result.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(res.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Hashing issues!!",
		})
		return
	}

	res.Password = string(hashedPassword)

	var data = models.User{}

	data.Name = res.Name
	data.Email = res.Email
	data.Password = res.Password

	dbResponse := db.Create(&data)
	if dbResponse.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in creating entry!",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Created Successfully",
	})
}

//controller to login a user
func LoginUser(c *gin.Context, db *gorm.DB) {
	var res = LoginResponse{}

	if err := c.ShouldBindJSON(&res); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid response!",
		})
		return
	}

	dbResponse := models.User{}
	if result := db.Where("email = ?", res.Email).First(&dbResponse); result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No user with this email exists!",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbResponse.Password), []byte(res.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Password doesn't match",
		})
		return
	}

	validToken, err := utils.GenerateJwt(dbResponse.ID, dbResponse.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server eroor",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Success!",
		"data":    validToken,
	})
}
