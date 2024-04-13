package main

import (
	authHandler "avitoTech/internal/pkg/auth/delivery/http"
	authRepo "avitoTech/internal/pkg/auth/repo"
	authUsecase "avitoTech/internal/pkg/auth/usecase"
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	_ = godotenv.Load()
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Println("fail ping postgres")
		err = fmt.Errorf("error happened in db.Ping: %w", err)
		log.Println(err)
	}

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	r.HandleFunc("/ping", pingPongHandler).Methods(http.MethodGet)

	aRepo := authRepo.NewRepository(db)
	aUsecase := authUsecase.NewUsecase(aRepo)
	aHandler := authHandler.NewHandler(aUsecase)

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign_in", aHandler.SignIn).Methods(http.MethodPost, http.MethodOptions)
	auth.HandleFunc("/sign_up", aHandler.SignUp).Methods(http.MethodPost, http.MethodOptions)
	auth.HandleFunc("/sign_out", aHandler.SignOut).Methods(http.MethodGet, http.MethodOptions)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Start server on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	sig := <-signalCh
	log.Printf("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}
}

func pingPongHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "pong")
}
