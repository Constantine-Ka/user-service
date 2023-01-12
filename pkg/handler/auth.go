package handler

import (
	"errors"
	"fmt"
	"github.com/Constantine-Ka/user-service/model"
	"github.com/Constantine-Ka/user-service/tools/logging"
	"github.com/Constantine-Ka/user-service/tools/mailer"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type singInInput struct {
	id        int    `json:"id,omitempty" :"id"`
	Confirm   string `json:"confirm"`
	Login     string `json:"login"  json:"login,omitempty":"login"`
	Password  string `json:"password" binding:"required" json:"password,omitempty":"password"`
	Password2 string `json:"password2" json:"password2,omitempty":"password2"`
}
type updateToken struct {
	Token  string `json:"token" :"token"`
	UserId int    `json:"user_id"`
}

func (h *Handler) singUp(c *gin.Context) {
	var input model.UserCreator
	err := c.BindJSON(&input)
	if input.Password == "" {
		newErrorResponse(c, http.StatusBadRequest, "Password is empty")
		return
	}
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, confirm, err := h.services.Authorization.CreateUser(input)
	logrus.Info(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = mailer.AuthMail(input.EMail, confirm, false)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":               id,
		"ConfirmationCode": confirm,
	})

}
func (h *Handler) singIn(c *gin.Context) {
	logger := logging.GetLogger()
	var (
		input      singInInput
		userGetter model.UserGet
	)
	err := c.BindJSON(&input)

	if input.Login == "" {
		err := errors.New("login is empty")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	logger.Info(input)
	userGetter, token, err := h.services.Authorization.GenerateToken(input.Login, input.Password)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token":     token,
		"user_info": userGetter,
	})
}
func (h *Handler) resetLink(c *gin.Context) {
	var userdata model.UserGet
	userEmail := c.Query("email")
	userdata, err := h.services.GetConfirm(userEmail)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = mailer.AuthMail(userdata.EMail, userdata.Confirmation.String, true)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"ConfirmationCode": userdata.Confirmation,
		"Message":          fmt.Sprintf("Link to send of email %s", userdata.EMail),
	})
}
func (h *Handler) confirm(c *gin.Context) {
	logger := logging.GetLogger()
	logger.Info(c)
	code := c.Query("code")
	logger.Warn(code)
	id, err := h.services.ConfirmUser(code)
	fmt.Println(id)
	logger.Error(err)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":        id,
		"code":      code,
		"isConfirm": true,
	})
}
func (h *Handler) resetPassword(c *gin.Context) {
	logger := logging.GetLogger()
	var (
		input singInInput
		//userdata model.UserCreator
	)
	err := c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if input.Password != input.Password2 {
		err := errors.New("passwords is different")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.ResetPassword(input.Confirm, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Reset password of user",
	})
	logger.Info(fmt.Sprintf("User reset password"))
}
func (h *Handler) updateToken(c *gin.Context) {
	var updateInput updateToken
	errors := c.BindJSON(&updateInput)
	if errors != nil {
		fmt.Println(errors)
		return
	}
	id, token, err := h.services.RefreshToken(updateInput.Token, updateInput.UserId)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"new_token": token,
		"message":   "token is update",
		"user_id":   id,
	})
}
