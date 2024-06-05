package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-go/internal/config/environment"
	"todo-go/internal/handlers/home"
	"todo-go/internal/handlers/login"
	"todo-go/internal/handlers/logout"
	"todo-go/internal/handlers/notfound"
	"todo-go/internal/handlers/register"
	"todo-go/internal/handlers/todos"
	"todo-go/internal/handlers/users"
	database "todo-go/internal/store/db"
	"todo-go/internal/store/dbstore"
	"todo-go/internal/util/hash/passwordhash"

	m "todo-go/internal/middleware"

	_ "todo-go/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Todo Web API
//	@version		1.0
//	@description	This is a sample todo app.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
/*
* Set to prod at build time
* used to determine what assets to load
 */
const Environment = "dev"

func init() {
	os.Setenv("env", Environment)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	envConfig := environment.NewEnvConfig(&environment.DefaultEnvConfigProcessor{})
	cfg := envConfig.MustLoadConfig()

	db, err := database.MustOpenMysql(cfg.DatabaseName, cfg.DatabasePassword)
	if err != nil {
		logger.Error("failed mysql", err)
		db = database.MustOpenSqlite(cfg.DatabaseName)
	}

	passwordhash := passwordhash.NewHPasswordHash()
	/* stores */

	userStore := dbstore.NewUserStore(db, passwordhash)

	sessionStore := dbstore.NewSessionStore(db)

	todoStore := dbstore.NewTodoStore(db)

	/* middlewares */
	authMiddleware := m.NewAuthMiddleware(sessionStore, cfg.SessionCookieName)

	/* handlers */
	registerHandler := register.NewRegisterHandler(userStore)
	loginHandler := login.NewLoginHandler(userStore, sessionStore, passwordhash, cfg.SessionCookieName)
	logoutHandler := logout.NewLogoutHandler(cfg.SessionCookieName)

	todosHandler := todos.NewTodosHandler(todoStore)
	usersHandler := users.NewUsersHandler(userStore)
	homeHandler := home.NewHomeHandler(todoStore)

	notfoundHandler := notfound.NewNotFoundHandler()
	/* router */
	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	fileServerSwagger := http.FileServer(http.Dir("./docs"))
	r.Handle("/docs/*", http.StripPrefix("/docs/", fileServerSwagger))

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			authMiddleware.AddUserToContext,
			m.TextHTMLMiddleware,
		)

		r.NotFound(notfoundHandler.NotFound)

		r.Get("/", homeHandler.Get)
		r.Group(func(r chi.Router) {
			r.Use(m.UserLoggedInMiddlewareTemplate)
			r.Post("/add-todo", todosHandler.AddTodoTemplate)
			r.Patch("/finish-todo", todosHandler.FinishTodoTemplate)
			r.Patch("/unfinish-todo", todosHandler.UnFinishTodoTemplate)
			r.Post("/delete-todo", todosHandler.DeleteTodoTemplate)
		})

		r.Get("/register", registerHandler.Get)
		r.Post("/register", registerHandler.Post)

		r.Get("/login", loginHandler.Get)
		r.Post("/login", loginHandler.Post)

		r.Post("/logout", logoutHandler.Post)

	})

	/* api */
	r.Route("/api/users", func(r chi.Router) {
		r.Use(
			m.JSONMiddleware,
		)
		// r.Get("/", getUsers)
		// r.Post("/", createUser)
		r.Route("/{userId}", func(r chi.Router) {
			r.Get("/", usersHandler.GetUserById)
			// r.Put("/", updateUser)
			// r.Delete("/", deleteUser)
			r.Route("/todos", func(r chi.Router) {
				r.Get("/", todosHandler.GetAllTodos)
				r.Post("/", todosHandler.AddTodo)
				r.Get("/search", todosHandler.GetTodosBySearch)
				r.Route("/{todoId}", func(r chi.Router) {
					r.Get("/", todosHandler.GetTodo)
					r.Delete("/", todosHandler.DeleteTodo)
					r.Put("/", todosHandler.UpdateTodo)

				})

			})

		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/swagger.json")))

	/* exit if channel */
	killSigChan := make(chan os.Signal, 1)

	signal.Notify(killSigChan, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Info("Server started", slog.String("port", cfg.Port), slog.String("env", Environment))
	<-killSigChan

	logger.Info("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
