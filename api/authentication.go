package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/martikan/carrental_auth-api/db/sqlc"
	"github.com/martikan/carrental_auth-api/util"
	e "github.com/martikan/carrental_common/exception"
	"github.com/martikan/carrental_common/middleware"
	common_utils "github.com/martikan/carrental_common/util"
)

type signInRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signUpRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type signInResponse struct {
	AccessToken string `json:"access_token" binding:"required"`
	User        db.User
}

type userResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}

func (a *Api) currentUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*common_utils.Payload)

	currentUser, err := a.db.GetUserByEmail(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ApiMessage("User not found", 404))
			return
		}

		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	dto := newUserResponse(currentUser)

	ctx.JSON(http.StatusOK, dto)
}

func (a *Api) signIn(ctx *gin.Context) {
	var req signInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	user, err := a.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ApiMessage("User not found", 404))
			return
		}

		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	err = util.PasswordUtils.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, e.ApiMessage("Invalid password", 401))
		return
	}

	accessToken, err := a.tokenMaker.CreateToken(user.Email, a.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	}

	dto := signInResponse{
		AccessToken: accessToken,
		User:        user,
	}

	ctx.JSON(http.StatusOK, dto)
}

func (a *Api) signUp(ctx *gin.Context) {

	var req signUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	passHash, err := util.PasswordUtils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	}

	args := db.CreateUserParams{
		Email:     req.Email,
		Password:  passHash,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := a.db.CreateUser(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, e.ApiMessage("Email address is already exist", 403))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	}

	dto := newUserResponse(user)

	ctx.JSON(http.StatusCreated, dto)

}
