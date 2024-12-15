package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/chilipizdrick/muzek-server/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const AUTH_API_ROUTE = "/auth"

func assignAuthRouteHandlers(parentGroup *gin.RouterGroup, db *gorm.DB) {
	group := parentGroup.Group(AUTH_API_ROUTE)

	group.POST("/register", registerUserWrapper(db))
	group.POST("/login", loginUserWrapper(db))
	group.POST("/logout", logoutUserWrapper(db))
}

func hashPassword(username string, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s%s", username, password)), 14)
	return string(bytes), err
}

func checkPasswordHash(username string, password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(fmt.Sprintf("%s%s", username, password)))
	return err == nil
}

type RegisterUserRequest struct {
	Username  string `validate:"required,username"`
	Password1 string `validate:"required,password"`
	Password2 string `validate:"required,eqfield=Password1"`
}

func registerUserWrapper(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		request := RegisterUserRequest{
			Username:  c.Query("username"),
			Password1: c.Query("password1"),
			Password2: c.Query("password2"),
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
			value := fl.Field().Interface().(string)
			if len(value) < 3 || len(value) > 32 {
				return false
			}
			if strings.Contains("_.", string(value[0])) || strings.Contains("_.", string(value[len(value)-1])) {
				return false
			}
			for _, char := range value {
				if char > unicode.MaxASCII {
					return false
				}
				if !(unicode.IsLetter(char) || unicode.IsDigit(char) || strings.Contains("_.", string(char))) {
					return false
				}
			}
			return true
		})
		validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
			value := fl.Field().Interface().(string)
			var (
				hasMinLength = false
				hasLetter    = false
				hasDigit     = false
			)
			if len(value) >= 8 && len(value) <= 32 {
				hasMinLength = true
			}
			for _, char := range value {
				switch {
				case unicode.IsLetter(char):
					hasLetter = true
				case unicode.IsNumber(char):
					hasDigit = true
				}
			}
			return hasMinLength && hasLetter && hasDigit
		})

		if err := validate.Struct(request); err != nil {
			log.Printf("[TRACE] Failed to validate user registration request: %s", err)
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				switch validationErrors[0].Field() {
				case "Username":
					switch validationErrors[0].Tag() {
					case "required":
						badRequestResponse(c, "\"username\" is required.")
					case "username":
						badRequestResponse(c, "\"username\" is invalid (a-z, A-Z, 0-9, `.`, `_` (not first or last character) are allowed characters).")
					}
				case "Password1":
					switch validationErrors[0].Tag() {
					case "required":
						badRequestResponse(c, "\"password1\" is required.")
					case "password":
						badRequestResponse(c, "\"password1\" is invalid (8-32 characters long, must contain at least one character and digit).")
					}
				case "Password2":
					switch validationErrors[0].Tag() {
					case "required":
						badRequestResponse(c, "\"password2\" is required.")
					case "eqfield":
						badRequestResponse(c, "\"password2\" does not match \"password1\".")
					}
				}
				return
			} else {
				internalServerErrorResponse(c, "Internal server error.")
				panic("unreachable")
			}
		}

		var userModel database.UserModel
		err := db.First(&userModel, "username = ?", request.Username).Error
		if err == nil {
			badRequestResponse(c, fmt.Sprintf("User with username \"%s\" already exists.", request.Username))
			return
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			internalServerErrorResponse(c, "Internal server error.")
			return

		}
		passwordHash, err := hashPassword(request.Username, request.Password1)
		if err != nil {
			internalServerErrorResponse(c, "Internal server error.")
			return
		}

		userModel.Username = request.Username
		userModel.PasswordHash = passwordHash
		userModel.PlaylistIDs = make(pq.Int64Array, 0)
		userModel.PlayerState = nil
		userModel.IsPlayerStatePublic = false

		err = db.Create(&userModel).Error
		if err != nil {
			log.Printf("[ERROR] Error creating user database entry: %s", err)
			internalServerErrorResponse(c, "Internal server error.")
			return
		}

		var playlistModel database.PlaylistModel
		playlistModel.Title = fmt.Sprintf("%s's Liked Songs", request.Username)
		playlistModel.OwnerID = &userModel.ID
		playlistModel.IsPublic = false
		playlistModel.TrackIDs = make(pq.Int64Array, 0)
		playlistModel.Deletable = false

		err = db.Create(&playlistModel).Error
		if err != nil {
			log.Printf("[ERROR] Error creating playlist database entry: %s", err)
			internalServerErrorResponse(c, "Internal server error.")
			return
		}

		userModel.PlaylistIDs = pq.Int64Array{int64(playlistModel.ID)}
		err = db.Save(&userModel).Error
		if err != nil {
			log.Printf("[ERROR] Error updating user database entry: %s", err)
			internalServerErrorResponse(c, "Internal server error.")
			return
		}
	}
}

type LoginUserRequest struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

func loginUserWrapper(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		request := LoginUserRequest{
			Username: c.Query("username"),
			Password: c.Query("password"),
		}
		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(request); err != nil {
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				badRequestResponse(c, "\"username\" and \"password\" are required.")
				return
			} else {
				internalServerErrorResponse(c, "Internal server error.")
				return
			}
		}

		var userModel database.UserModel
		err := db.First(&userModel, "username = ?", request.Username).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				badRequestResponse(c, "User does not exist.")
				return
			}
			log.Printf("[ERROR] Error fetching user from db: %s", err)
			internalServerErrorResponse(c, "Internal server error.")
			return
		}

		if !checkPasswordHash(request.Username, request.Password, userModel.PasswordHash) {
			badRequestResponse(c, "User does not exist.")
			return
		}

		session := sessions.Default(c)
		session.Set("userId", userModel.ID)
		if err = session.Save(); err != nil {
			internalServerErrorResponse(c, "Internal server error.")
		}
		c.JSON(http.StatusOK, OKResponse{OK{
			Status:  http.StatusOK,
			Message: "Successfully logged in.",
		}})
	}
}

func logoutUserWrapper(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionId := session.Get("userId")
		if sessionId == nil {
			badRequestResponse(c, "Invalid session token.")
			return
		}
		session.Delete("userId")
		if err := session.Save(); err != nil {
			internalServerErrorResponse(c, "Internal server error.")
			return
		}
		c.JSON(http.StatusOK, OKResponse{OK{
			Status:  http.StatusOK,
			Message: "Successfully logged out.",
		}})
	}
}
