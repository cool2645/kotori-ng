package auth

import (
	"github.com/cool2645/kotori-ng/model"
	"github.com/cool2645/kotori-ng/database"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"time"
	. "github.com/cool2645/kotori-ng/config"
	"fmt"
	"github.com/satori/go.uuid"
)

func checkCredential(username string, password string) (isValid bool, user model.User, err error) {
	user, err = model.GetUserByUsername(database.DB, username)
	if err != nil {
		if err.Error() == "GetUserByUsername: record not found" {
			err = nil
		} else {
			err = errors.Wrap(err, "checkCredential")
			return
		}
	}
	e := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if e != nil {
		isValid = false
	} else {
		isValid = true
	}
	return
}

func GenerateToken(username string, password string) (ok bool, tokenString string, msg string) {
	ok, user, err := checkCredential(username, password)
	if err != nil {
		msg = errors.Wrap(err, "Login").Error()
		ok = false
		return
	}

	if !ok {
		msg = "Invalid credential"
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":      time.Now().Add(time.Hour * time.Duration(GlobCfg.JWT_EXPIRETIME)).Unix(),
		"nbf":      time.Now().Unix(),
		"uuid":     user.UUID,
		"username": user.Username,
		"name":     user.Name,
		"is_admin": user.IsAdmin,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString([]byte(GlobCfg.JWT_KEY))
	if err != nil {
		msg = errors.Wrap(err, "Login").Error()
		ok = false
		return
	}
	return
}

func CheckToken(tokenString string) (ok bool, user model.User, msg string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(GlobCfg.JWT_KEY), nil
	})
	if err != nil {
		msg = errors.Wrap(err, "CheckToken").Error()
		ok = false
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		user, _ = model.GetUserByUUID(database.DB, claims["uuid"].(string))
	} else {
		msg = errors.Wrap(err, "CheckToken").Error()
	}
	return
}

func MakeUser(user model.User) (newUser model.User, err error) {
	u2, err := uuid.NewV4()
	if err != nil {
		err = errors.Wrap(err, "MakeUser")
		return
	}
	user.UUID = u2.String()
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Wrap(err, "MakeUser")
		return
	}
	user.Password = string(hash)
	user.IsAdmin = false
	newUser, err = model.StoreUser(database.DB, user)
	if err != nil {
		err = errors.Wrap(err, "MakeUser")
		return
	}
	return
}
