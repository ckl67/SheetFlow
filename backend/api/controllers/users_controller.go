package controllers

import (
	"backend/api/auth"
	"backend/api/config"
	"backend/api/forms"
	"backend/api/models"
	"backend/api/utils"
	"backend/api/utils/formaterror"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (server *Server) CreateUser(c *gin.Context) {
	// Check for authentication
	token := utils.ExtractToken(c)
	uid, err := auth.ExtractTokenID(token, config.Config().ApiSecret)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	if uid != config.ADMIN_UID {
		c.String(http.StatusUnauthorized, "only Admins are able to persue this command")
		return
	}

	var user models.User
	err = c.BindJSON(&user)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	user.PrepareForCreate()
	if err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}
	// Dans SaveUser() un Hoock va appeler BeforeSave()
	log.Printf("User struct: %+v\n", user)
	userCreated, err := user.Save(server.DB)
	if err != nil {
		// log.Println("RAW ERROR:", err.Error())
		formattedError := formaterror.FormatError(err.Error())
		c.String(http.StatusUnprocessableEntity, formattedError.Error())
		return
	}
	c.Header("Location", fmt.Sprintf("%s%s/%d", c.Request.Host, c.Request.RequestURI, userCreated.ID))
	c.JSON(http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(c *gin.Context) {
	// Check for authentication
	token := utils.ExtractToken(c)
	uid, err := auth.ExtractTokenID(token, config.Config().ApiSecret)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	if uid != config.ADMIN_UID {
		c.String(http.StatusUnauthorized, "Only admins are able to persue this command")
		return
	}

	users, err := models.GetAllUsers(server.DB)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

func (server *Server) GetUser(c *gin.Context) {
	uidString := c.Param("id")
	uid, err := strconv.ParseUint(uidString, 10, 32)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	token := utils.ExtractToken(c)
	userId, err := auth.ExtractTokenID(token, config.Config().ApiSecret)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	var newUid uint32 = uint32(uid)
	if uid == 0 {
		/*
			make it the details about own user
		*/
		newUid = userId
	}

	// Check for admin
	if uid != 0 && userId != config.ADMIN_UID {
		c.String(http.StatusUnauthorized, "Only admins are able to look at user that aren't themselves. Try the endpoint /users/0 to look at your own user details")
		return
	}

	uid = uint64(newUid)
	user := models.User{}
	userGotten, err := user.FindByID(server.DB, uint32(uid))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, userGotten)
}

// UpdateUser updates a user. Only admins and the user itself are able to update a user.
// Admins are able to update the role of a user, while users themselves are not able to update their own role.
// To update a user go to endpoint:
// PUT: /api/users/:id
//
//	Body: {
//		"email": "your-updated-email"
//		"password": "your-updated-password"
//	}
//
// Admins can also update the role of a user by adding the "role" field to the body with either "admin" or "user" as value.
// Example body for admin updating a user:
//
//	{
//		"email": "your-updated-email"
//		"password": "your-updated-password"
//		"role": "admin"
//	}
func (server *Server) UpdateUser(c *gin.Context) {
	uidString := c.Param("id")
	uid, err := strconv.ParseUint(uidString, 10, 32)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	// Set updated At param
	user.UpdatedAt = time.Now()

	// Gin remplit le struct user avec les données JSON envoyées.
	// Si { "role": 0 }, user.Role devient 0.
	err = c.BindJSON(&user)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	token := utils.ExtractToken(c)
	tokenID, err := auth.ExtractTokenID(token, config.Config().ApiSecret)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	if tokenID != uint32(uid) && tokenID != config.ADMIN_UID {
		c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	user.PrepareForUpdate()

	if err := utils.ValidateStruct(user); err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	var updatedUser *models.User
	log.Println("Updating user with ID:", uid, "by user with ID:", tokenID)

	if tokenID == config.ADMIN_UID {
		updatedUser, err = user.Update(server.DB, uint32(uid), tokenID)
	} else {
		updatedUser, err = user.Update(server.DB, uint32(uid), tokenID)
	}

	if err != nil {
		fmt.Println(err)
		formattedError := formaterror.FormatError(err.Error())
		c.String(http.StatusUnprocessableEntity, formattedError.Error())
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (server *Server) DeleteUser(c *gin.Context) {
	uidString := c.Param("id")
	uid, err := strconv.ParseUint(uidString, 10, 32)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	token := utils.ExtractToken(c)
	tokenID, err := auth.ExtractTokenID(token, config.Config().ApiSecret)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	if tokenID != 1 && tokenID != uint32(uid) {
		c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	var user models.User
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Entity", fmt.Sprint(uid))
	c.JSON(http.StatusNoContent, gin.H{})
}

func (server *Server) ResetPassword(c *gin.Context) {
	var form forms.ResetPasswordRequest
	if err := c.ShouldBind(&form); err != nil {
		utils.DoError(c, http.StatusBadRequest, err)
		return
	}
	if err := utils.ValidateStruct(form); err != nil {
		utils.DoError(c, http.StatusBadRequest, err)
		return
	}

	user, err, statusCode := models.ResetPassword(server.DB, form.PasswordResetId, form.Password)
	if err != nil {
		c.JSON(statusCode, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (server *Server) RequestPasswordReset(c *gin.Context) {
	var form forms.RequestResetPasswordRequest

	if err := c.ShouldBind(&form); err != nil {
		utils.DoError(c, http.StatusBadRequest, err)
		return
	}

	if err := utils.ValidateStruct(form); err != nil {
		utils.DoError(c, http.StatusBadRequest, err)
		return
	}

	resetToken, err := models.RequestPasswordReset(server.DB, form.Email)
	if err != nil {
		// ⚠️ On ne révèle pas si l'email existe
		c.JSON(http.StatusOK, "If the email exists, a reset link has been sent.")
		return
	}

	cfg := config.Config().Smtp
	if cfg.Enabled != "true" {
		c.JSON(http.StatusBadGateway, "SMTP backend not configured.")
		return
	}

	err = utils.SendMail(
		form.Email,
		"Password Reset",
		buildResetBody(resetToken),
	)
	if err != nil {
		utils.DoError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "If the email exists, a reset link has been sent.")
}

func buildResetBody(token string) string {
	return fmt.Sprintf(
		"Click the link below to reset your password:\n\nhttps://yourdomain.com/reset?token=%s",
		token,
	)
}
