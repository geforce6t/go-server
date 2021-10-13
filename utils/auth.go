package utils

import (
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/geforce6t/go-server/models"
	"github.com/twinj/uuid"
)

func CreateToken(id uint) (*TokenDetails, error) {
	td := &TokenDetails{}
	var mySigningKey = []byte(models.GetEnvValue("secret"))

	// access token details
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	// refresh token details
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = id
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = at.SignedString(mySigningKey)
	if err != nil {
		log.Fatalf("Something Went Wrong: %s", err.Error())
		return td, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = id
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(mySigningKey)
	if err != nil {
		log.Fatalf("Something Went Wrong: %s", err.Error())
		return td, err
	}

	return td, nil
}

func CreateAuth(id uint, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(id)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(id)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
