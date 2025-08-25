package handlers

import (
	"encoding/json"
	"net/http"
	"rwa/internal/models"
	"rwa/internal/services"
)

type UserHandler struct {
	us *services.UserService
	sm *services.SessionManager
}

func NewUserHandler(us *services.UserService, sesManager *services.SessionManager) *UserHandler {
	return &UserHandler{
		us: us,
		sm: sesManager,
	}
}

type RegisterResponse struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UserCreateRequest struct {
	User models.UserCreateInfo `json:"user"`
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	registerData := UserCreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&registerData)
	if err != nil {
		badJsonError(w)
		return
	}

	user, err := uh.us.CreateUser(registerData.User)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res := map[string]interface{}{
		"user": user,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequestData struct {
	User LoginData `json:"user"`
}

type LoginResponseData struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Token     string `json:"token"`
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginData := LoginRequestData{}
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		badJsonError(w)
		return
	}

	if loginData.User.Email == "" || loginData.User.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email or Password is empty"))
		return
	}

	user, err := uh.us.GetByEmail(loginData.User.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// TODO: ошибка может сообщить лишнее о системе
		// следует явно управлять сообщениями ошибок во внешний мир
		w.Write([]byte(err.Error()))
		return
	}

	session, err := uh.sm.Create(*user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	logDataRes := LoginResponseData{
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(timeFormat),
		UpdatedAt: user.UpdatedAt.Format(timeFormat),
		Token:     session.GetSessionId(),
	}

	res := map[string]interface{}{"user": logDataRes}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sesId := GetTokenFromRequest(r)
	err := uh.sm.Delete(sesId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (uh *UserHandler) Info(w http.ResponseWriter, r *http.Request) {
	uId, err := GetUserIdFromRequestCtx(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := uh.us.GetUserById(uId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res := map[string]interface{}{
		"user": user,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

type UpdateRequest struct {
	User models.UserUpdateInfo `json:"user"`
}

func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	uId, err := GetUserIdFromRequestCtx(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	reqData := UpdateRequest{}
	err = json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		badJsonError(w)
		return
	}
	data := reqData.User

	user, err := uh.us.GetUserById(uId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	updatedUser, err := uh.us.UpdateUser(*user, data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res := map[string]interface{}{
		"user": updatedUser,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
