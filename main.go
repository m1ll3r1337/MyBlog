package main

import (
	"blog/controllers"
	"blog/models"
	"blog/models/templates"
	"blog/views"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no PSQL config provided")
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	//Setup db
	db, err := models.Open(cfg.PSQL)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	//Setup services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)

	postService := &models.PostService{
		DB:               db,
		Mr:               models.FileReader{},
		Mw:               models.FileWriter{},
	}
	//Setup middleware

	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Secure(cfg.CSRF.Secure), csrf.Path("/"))
	//Setup controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	postsC := controllers.Posts{
		PostService: postService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS,
		"check-your-email.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS,
		"reset-pw.gohtml", "tailwind.gohtml",
	))
	postsC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"posts/new.gohtml", "tailwind.gohtml"))
	postsC.Templates.Edit = views.Must(views.ParseFS(
		templates.FS,
		"posts/edit.gohtml", "tailwind.gohtml"))
	postsC.Templates.Index = views.Must(views.ParseFS(
		templates.FS,
		"posts/index.gohtml", "tailwind.gohtml"))
	postsC.Templates.Show = views.Must(views.ParseFS(
		templates.FS,
		"posts/show.gohtml", "tailwind.gohtml"))

	//Setup router
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Route("/posts", func(r chi.Router) {
		r.Get("/{id}", postsC.Show)
		r.Get("/{id}/images/{filename}", postsC.Image)
		r.Get("/", postsC.Index)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/new", postsC.New)
			r.Get("/{id}/edit", postsC.Edit)
			r.Post("/{id}", postsC.Update)
			r.Post("/{id}/delete", postsC.Delete)
			r.Post("/", postsC.Create)
			r.Post("/{id}/images/{filename}/delete", postsC.DeleteImage)
			r.Post("/{id}/images", postsC.UploadImage)
		})
	})

	//Start the server
	fmt.Printf("Listening on %s...", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
