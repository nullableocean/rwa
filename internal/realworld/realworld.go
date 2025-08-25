package realworld

import (
	"net/http"
	"rwa/http/handlers"
	"rwa/http/middleware"
	"rwa/internal/repository/ram"
	"rwa/internal/services"
	"rwa/pkg/passwordcryptor"

	"github.com/gorilla/mux"
)

func GetApp() http.Handler {

	router := mux.NewRouter()
	createApi(router)

	return router
}

func createApi(router *mux.Router) {
	userRepo := ram.NewUserRepository()
	sessionRepo := ram.NewSessionRepository()
	articleRepo := ram.NewArticleRepository()

	userService := services.NewUserService(userRepo, passwordcryptor.PasswordCryptor{})
	sessionService := services.NewSessionManager(sessionRepo, userService)
	articleService := services.NewArticleService(articleRepo)

	userHandler := handlers.NewUserHandler(userService, sessionService)
	articleHandler := handlers.NewArticleHandler(articleService, userService)

	sessionGuard := middleware.NewSessionGuard(sessionService)
	authMiddleware := sessionGuard.GetAuthMiddleware()

	router.HandleFunc("/api/users", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/users/login", userHandler.Register).Methods("POST")

	ur := router.PathPrefix("/api/user").Subrouter()
	ur.Use(authMiddleware)
	ur.HandleFunc("/", userHandler.Info).Methods("GET")
	ur.HandleFunc("/", userHandler.Update).Methods("PUT")
	ur.HandleFunc("/logout", userHandler.Logout).Methods("POST")

	ar := router.PathPrefix("/api/articles").Subrouter()
	ar.Use(authMiddleware)
	ar.HandleFunc("/", articleHandler.Create).Methods("POST")
	ar.HandleFunc("/", articleHandler.Get).Methods("GET")
}
