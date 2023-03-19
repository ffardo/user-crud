package controllers

import (
	"net/http"

	"github.com/ffardo/user-crud/interfaces"
	"github.com/ffardo/user-crud/services"
	"github.com/gin-gonic/gin"
)

const DATE_FORMAT = "2006-01-02"

type UserController struct {
	interfaces.UserService
}

type UserRequest struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Password  string `json:"password"`
}

type UserResponse struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Password  string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleError(ctx *gin.Context, err error) {
	r := ErrorResponse{
		Error: err.Error(),
	}
	switch err {
	case services.ErrEmailRegistered:
		ctx.IndentedJSON(http.StatusBadRequest, r)
	case services.ErrInvalidEmailFormat:
		ctx.IndentedJSON(http.StatusBadRequest, r)
	case services.ErrInvalidDateFormat:
		ctx.IndentedJSON(http.StatusBadRequest, r)
	case services.ErrInvalidUuidFormat:
		ctx.IndentedJSON(http.StatusBadRequest, r)
	case services.ErrUserNotFound:
		ctx.IndentedJSON(http.StatusNotFound, r)
	}

}

func (c *UserController) Post(ctx *gin.Context) {
	var UserParam UserRequest

	err := ctx.BindJSON(&UserParam)

	if err != nil {
		handleError(ctx, err)
		return
	}

	user, err := c.CreateUser(
		UserParam.Name,
		UserParam.BirthDate,
		UserParam.Email,
		UserParam.Address,
		UserParam.Password,
	)

	if err != nil {
		handleError(ctx, err)
		return
	}

	r := UserResponse{
		UUID:      user.UUID.String(),
		Name:      user.Name,
		BirthDate: user.BirthDate.Format(DATE_FORMAT),
		Email:     user.Email,
		Password:  user.Password,
		Address:   user.Address,
	}

	ctx.IndentedJSON(http.StatusCreated, r)

}

func (c *UserController) Get(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	user, err := c.GetUser(uuid)

	if err != nil {
		handleError(ctx, err)
		return
	}

	r := UserResponse{
		UUID:      user.UUID.String(),
		Name:      user.Name,
		BirthDate: user.BirthDate.Format(DATE_FORMAT),
		Email:     user.Email,
		Password:  user.Password,
		Address:   user.Address,
	}

	ctx.IndentedJSON(200, r)

}

func (c *UserController) Patch(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	userParams := make(map[string]string)
	ctx.BindJSON(&userParams)

	user, err := c.UpdateUser(uuid, userParams)

	if err != nil {
		handleError(ctx, err)
		return
	}

	r := UserResponse{
		UUID:      user.UUID.String(),
		Name:      user.Name,
		BirthDate: user.BirthDate.Format(DATE_FORMAT),
		Email:     user.Email,
		Password:  user.Password,
		Address:   user.Address,
	}

	ctx.IndentedJSON(200, r)

}

func (c *UserController) Delete(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := c.DeleteUser(uuid)

	if err != nil {
		handleError(ctx, err)
		return
	}

}
