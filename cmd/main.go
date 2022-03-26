package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-snndmr/internal/domain/author"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-snndmr/internal/domain/book"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-snndmr/internal/infrastructure"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	addr             = "127.0.0.1:3306"
	userID           = "root"
	password         = "password123!"
	database         = "library"
	connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True&loc=Local", userID, password, addr, database)
	pathBooksCSV     = "resources/books.csv"
	pathAuthorsCSV   = "resources/authors.csv"

	db               = infrastructure.NewMySQL(connectionString)
	bookRepository   = book.NewRepository(db)
	authorRepository = author.NewRepository(db)
)

func main() {
	authorRepository.Migration()
	authorRepository.InitializeWithSampleData(infrastructure.GetAuthorsFromCSV(pathAuthorsCSV))

	bookRepository.Migration()
	bookRepository.InitializeWithSampleData(infrastructure.GetBooksFromCSV(pathBooksCSV))

	InitializeServer()
}

func InitializeServer() {
	router := mux.NewRouter()
	setCorsOptions()
	router.Use(loggingMiddleware)
	router.Use(authenticationMiddleware)

	sub := router.PathPrefix("/books").Subrouter()
	sub.HandleFunc("/list", ProductListHandler)
	sub.HandleFunc("/buy", ProductBuyHandler)
	sub.HandleFunc("/delete", ProductDeleteHandler)
	sub.HandleFunc("/search", ProductSearchHandler)
	sub.HandleFunc("/create", ProductCreateHandler)

	server := &http.Server{
		Addr:         "0.0.0.0:8090",
		Handler:      router,
		IdleTimeout:  time.Second * 60,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ShutdownServer(server, time.Second*10)
}

func setCorsOptions() {
	handlers.AllowedHeaders([]string{"Content Type", "Authorization"})
	handlers.AllowedMethods([]string{"POST", "GET", "PUT", "PATCH"})
}

func ShutdownServer(server *http.Server, duration time.Duration) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		return
	}
	log.Println("Shutting Down..!")
	os.Exit(0)
}

func ProductListHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json")

	resp, _ := json.Marshal(bookRepository.GetBooks())
	_, err := w.Write(resp)
	if err != nil {
		return
	}
}

func ProductBuyHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("id") != "" && r.FormValue("count") != "" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			return
		}

		count, err := strconv.Atoi(r.FormValue("count"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			return
		}

		err, result := bookRepository.GetById(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			return
		}

		err = result.DecreaseAmount(count)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			return
		}

		err = bookRepository.Update(result)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write([]byte("Book purchased successfully!"))
		if err != nil {
			return
		}
	}
}

func ProductDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("id") != "" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			return
		}

		err, result := bookRepository.GetById(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			return
		}

		err = result.Delete()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			return
		}

		err = bookRepository.Update(result)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write([]byte("Book deleted successfully!"))
		if err != nil {
			return
		}
	}
}

func ProductSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("substr") != "" {
		resp, _ := json.Marshal(bookRepository.Search(r.FormValue("substr")))
		_, err := w.Write(resp)
		if err != nil {
			return
		}
	}
}

func ProductCreateHandler(w http.ResponseWriter, r *http.Request) {
	pageCount, _ := strconv.Atoi(r.FormValue("pageCount"))
	stockCount, _ := strconv.Atoi(r.FormValue("stockCount"))
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	authorID, _ := strconv.Atoi(r.FormValue("authorID"))

	err := bookRepository.Create(
		book.NewBook(
			r.FormValue("title"), r.FormValue("stockId"), r.FormValue("isbn"), pageCount, stockCount, price, false,
			uint32(authorID),
		),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write([]byte("Book created successfully!"))
	if err != nil {
		return
	}
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if strings.HasPrefix(r.URL.Path, "/books/create") {
				if token != "" {
					next.ServeHTTP(w, r)
				} else {
					http.Error(w, "Token not found", http.StatusUnauthorized)
				}
			} else {
				next.ServeHTTP(w, r)
			}
		},
	)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("Requested: ", r.RequestURI)
			next.ServeHTTP(w, r)
		},
	)
}
