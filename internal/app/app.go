package app

import (
	"avitoTech/internal/middleware"
	authHandler "avitoTech/internal/pkg/auth/delivery/http"
	authRepo "avitoTech/internal/pkg/auth/repo"
	authUsecase "avitoTech/internal/pkg/auth/usecase"
	bannerHandler "avitoTech/internal/pkg/banners/delivery/http"
	bannerRepo "avitoTech/internal/pkg/banners/repo"
	bannerUsecase "avitoTech/internal/pkg/banners/usecase"
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	logger *logrus.Logger
}

func NewApp(logger *logrus.Logger) *App {
	return &App{logger: logger}
}

func (a *App) Run() error {
	_ = godotenv.Load()
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))
	if err != nil {
		a.logger.Error("failed to connect database" + err.Error())
	}

	if err = db.Ping(); err != nil {
		a.logger.Error("failed to ping database" + err.Error())
	}
	defer db.Close()

	//TODO: redis connect

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
		//TODO: cfg + time
	}

	aRepo := authRepo.NewRepository(db)
	aUsecase := authUsecase.NewUsecase(aRepo)
	aHandler := authHandler.NewHandler(aUsecase)

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign_in", aHandler.SignIn).Methods(http.MethodPost, http.MethodOptions)
	auth.HandleFunc("/sign_up", aHandler.SignUp).Methods(http.MethodPost, http.MethodOptions)
	auth.HandleFunc("/sign_out", aHandler.SignOut).Methods(http.MethodGet, http.MethodOptions)

	bRepo := bannerRepo.NewPostgresRepo(db)
	bUsecase := bannerUsecase.NewUseCase(bRepo)
	bHandler := bannerHandler.NewHandler(bUsecase, a.logger)

	r.Handle("/user_banner", middleware.Auth(http.HandlerFunc(bHandler.GetUserBanners), a.logger, false)).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/banner", middleware.Auth(http.HandlerFunc(bHandler.GetBanners), a.logger, true)).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/banner", middleware.Auth(http.HandlerFunc(bHandler.CreateBanner), a.logger, true)).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/banner/{id}", middleware.Auth(http.HandlerFunc(bHandler.UpdateBanner), a.logger, true)).Methods(http.MethodPatch, http.MethodOptions)
	r.Handle("/banner/{id}", middleware.Auth(http.HandlerFunc(bHandler.DeleteBanner), a.logger, true)).Methods(http.MethodDelete, http.MethodOptions)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Info("Start server on ", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("Error in listen: ", err.Error())
		}
	}()

	sig := <-signalCh
	a.logger.Info("Received signal: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.logger.Fatal("Server shutdown failed: ", err.Error())
	}
	return nil
}
