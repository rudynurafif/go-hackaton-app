// Package routes merakit controller, service, dan middleware menjadi satu
// router — padanan dari app.module.ts + module-module fitur di NestJS.
package routes

import (
	"database/sql"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"hackaton-management-app/config"
	"hackaton-management-app/controllers"
	"hackaton-management-app/middleware"
	"hackaton-management-app/models"
	"hackaton-management-app/services"
	"hackaton-management-app/utils"
)

// Setup membangun router Gin lengkap dengan semua route dan middleware.
func Setup(db *sql.DB, cfg *config.Config) *gin.Engine {
	registerValidators()

	router := gin.New()
	router.Use(gin.Logger())
	// Panic menjadi 500 dengan format error standar, bukan koneksi putus.
	router.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		utils.RespondError(c, &utils.AppError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		})
	}))

	// Route tak dikenal → 404 dengan format Nest: "Cannot GET /xyz".
	router.NoRoute(func(c *gin.Context) {
		utils.RespondError(c, utils.NewNotFound(
			"Cannot "+c.Request.Method+" "+c.Request.URL.Path,
		))
	})

	// --- Wiring dependensi (pengganti dependency injection NestJS) ---
	authService := &services.AuthService{DB: db, Cfg: cfg}
	userService := &services.UserService{DB: db}
	hackathonService := &services.HackathonService{DB: db}

	appCtrl := &controllers.AppController{}
	authCtrl := &controllers.AuthController{Service: authService}
	userCtrl := &controllers.UserController{Service: userService}
	hackathonCtrl := &controllers.HackathonController{Service: hackathonService}

	requireAuth := middleware.RequireAuth(db, cfg)
	adminOnly := middleware.RequireRoles(models.RoleAdmin)
	participantOnly := middleware.RequireRoles(models.RoleParticipant)

	// --- App (padanan app.controller.ts) ---
	router.GET("/", appCtrl.Hello)                 // @AllowAnonymous
	router.GET("/me", requireAuth, appCtrl.Me)

	// --- Auth (pengganti /api/auth/* milik Better Auth) ---
	auth := router.Group("/auth")
	{
		auth.POST("/register", authCtrl.Register)
		auth.POST("/login", authCtrl.Login)
	}

	// --- User (padanan user.controller.ts) — semua butuh login ---
	user := router.Group("/user", requireAuth)
	{
		// Didaftarkan sebelum ":id" agar "/user/all" tidak tertangkap
		// oleh route parameter (komentar yang sama ada di NestJS-nya).
		user.GET("/all", adminOnly, userCtrl.FindAll)
		user.GET("/:id", userCtrl.FindOne)
	}

	// --- Hackathon (padanan hackaton.controller.ts) ---
	hackathon := router.Group("/hackaton")
	{
		hackathon.GET("", hackathonCtrl.FindAll)     // @AllowAnonymous
		hackathon.GET("/:id", hackathonCtrl.FindOne) // @AllowAnonymous
		hackathon.POST("", requireAuth, adminOnly, hackathonCtrl.Create)
		hackathon.PATCH("/:id", requireAuth, adminOnly, hackathonCtrl.Update)
		hackathon.DELETE("/:id", requireAuth, adminOnly, hackathonCtrl.Remove)
		hackathon.POST("/:id/join", requireAuth, participantOnly, hackathonCtrl.Join)
	}

	return router
}

// registerValidators menyetel validator bawaan Gin:
//   - nama field pada pesan error memakai tag json (name, startsAt, ...)
//     alih-alih nama field Go (Name, StartsAt, ...);
//   - validator custom "future" — padanan @MinDate(() => new Date()),
//     dievaluasi per request sehingga "sekarang" selalu waktu request.
func registerValidators() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	_ = v.RegisterValidation("future", func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		return ok && date.After(time.Now())
	})
}
