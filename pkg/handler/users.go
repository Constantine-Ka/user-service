package handler

import (
	"github.com/Constantine-Ka/user-service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) meUser(c *gin.Context) {
	value, _ := c.Get("userId")
	i := value.(int)

	user, err := h.services.Users.MeUser(i)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Me Handler",
		"user":    user,
	})
}

func (h *Handler) updateUser(c *gin.Context) {
	var input model.UserCreator
	value, _ := c.Get("userId")
	id := value.(int)
	errBinding := c.BindJSON(&input)
	if errBinding != nil {
		newErrorResponse(c, http.StatusBadRequest, errBinding.Error())
		return
	}
	input.Password = "___"
	user, errService := h.services.UpdateUser(id, input)
	if errService != nil {
		newErrorResponse(c, http.StatusBadRequest, errService.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Update Handler",
		"user":    user,
	})
}
func (h *Handler) getUserAll(c *gin.Context) {
	userAll, err := h.services.GetUserAll()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, userAll)
}
func (h *Handler) updateLink(c *gin.Context) {
	var input []model.Link
	value, _ := c.Get("userId")
	id := value.(int)
	errBinding := c.BindJSON(&input)
	if errBinding != nil {
		newErrorResponse(c, http.StatusBadRequest, "incorrect request. You may have specified an object without an array")
		return
	}
	errService := h.services.UpdateLink(input, id)
	if errService != nil {
		newErrorResponse(c, http.StatusBadRequest, errService.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "UpdateLink Handler",
		"Input":   input,
	})
}
