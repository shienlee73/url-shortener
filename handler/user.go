package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shienlee73/url-shortener/store"
	"github.com/shienlee73/url-shortener/util"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserResponse(user store.User) UserResponse {
	return UserResponse{
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) CreateUser(c *gin.Context) {
	var createdUserRequest UserRequest
	if err := c.ShouldBindJSON(&createdUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if createdUserRequest.Username == "" || createdUserRequest.Password == "" {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("username and password are required")))
		return
	}

	hashedPassword, err := util.HashPassword(createdUserRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user := store.User{
		ID:             uuid.NewString(),
		Username:       createdUserRequest.Username,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}

	err = server.store.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := NewUserResponse(user)

	c.JSON(http.StatusCreated, res)
}

type LoginUserResponse struct {
	SessionID             string       `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

func (server *Server) LoginUser(c *gin.Context) {
	// check if session already exists, then return
	_, err := c.Cookie("session_id")
	if err == nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("session already exists")))
		return
	}

	var loginUserRequest UserRequest
	if err := c.ShouldBindJSON(&loginUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if loginUserRequest.Username == "" || loginUserRequest.Password == "" {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("username and password are required")))
		return
	}	

	user, err := server.store.RetrieveUser(loginUserRequest.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.VerifyPassword(user.HashedPassword, loginUserRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session := store.Session{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
		CreatedAt:    time.Now(),
	}

	if err := server.store.CreateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  NewUserResponse(user),
	}

	// c.SetCookie("session_id", session.ID, int(time.Until(refreshPayload.ExpiresAt.Time).Seconds()), "/", "", false, true)
	// c.SetCookie("access_token", accessToken, int(time.Until(accessPayload.ExpiresAt.Time).Seconds()), "/", "", false, true)
	// c.SetCookie("refresh_token", refreshToken, int(time.Until(refreshPayload.ExpiresAt.Time).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, res)
}

func (server *Server) LogoutUser(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, nil)
}
