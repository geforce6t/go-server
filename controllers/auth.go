package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
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

	td, err := utils.CreateToken(uint64(dbResponse.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	if saveErr := utils.CreateAuth(uint64(dbResponse.ID), td); saveErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		log.Fatalf("Redis error dude: %v", saveErr.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Success!",
		"data":    tokens,
	})
}

func Refresh(c *gin.Context, db *gorm.DB) {

	res := map[string]string{}
	var userId uint64

	if err := c.ShouldBindJSON(&res); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid response!",
		})
		return
	}

	refreshToken := res["refresh_token"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("secret")), nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		log.Fatalf("Invalid token: %v", err.Error())
	}

	refreshData := &utils.RefreshDetails{}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		}
		id, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			log.Fatalf("Invalid token: %v", err.Error())
		}

		refreshData.RefreshUuid = refreshUuid
		userId = id
	}

	if _, err = utils.Client.Get(refreshData.RefreshUuid).Result(); err == nil {
		_, err = utils.Client.Del(refreshData.RefreshUuid).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server Error",
			})
			log.Fatalf("error %v", err.Error())
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server Error",
		})
		log.Fatalf("error %v", err.Error())
	}

	td, err := utils.CreateToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		log.Fatalf("error %v", err.Error())
		return
	}

	if saveErr := utils.CreateAuth(userId, td); saveErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		log.Fatalf("Redis error dude: %v", saveErr.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Success!",
		"data":    tokens,
	})
}
