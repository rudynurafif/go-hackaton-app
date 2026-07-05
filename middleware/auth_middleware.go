// Package middleware — padanan guards di NestJS (AuthGuard global milik
// Better Auth + decorator @Roles).
package middleware

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"hackaton-management-app/config"
	"hackaton-management-app/models"
	"hackaton-management-app/utils"
)

// contextUserKey — key untuk menyimpan user login di gin.Context,
// padanan dari @Session() session: UserSession.
const contextUserKey = "currentUser"

// RequireAuth memverifikasi header "Authorization: Bearer <jwt>" lalu memuat
// user dari database. User diambil dari DB (bukan hanya dari claims) supaya
// perubahan role langsung berlaku dan user yang sudah dihapus ditolak —
// perilakunya setara dengan session lookup milik Better Auth.
func RequireAuth(db *sql.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			utils.RespondError(c, utils.NewUnauthorized("Unauthorized"))
			return
		}

		claims, err := utils.ParseToken(token, cfg.JWTSecret)
		if err != nil {
			utils.RespondError(c, utils.NewUnauthorized("Unauthorized"))
			return
		}

		var user models.User
		err = db.QueryRow(
			`SELECT id, name, email, role, created_at, updated_at
			 FROM users WHERE id = $1`,
			claims.Subject,
		).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondError(c, utils.NewUnauthorized("Unauthorized"))
			return
		}
		if err != nil {
			utils.RespondError(c, err)
			return
		}

		c.Set(contextUserKey, &user)
		c.Next()
	}
}

// RequireRoles membatasi akses ke role tertentu — padanan @Roles(['ADMIN']).
// Harus dipasang setelah RequireAuth.
func RequireRoles(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil {
			utils.RespondError(c, utils.NewUnauthorized("Unauthorized"))
			return
		}

		for _, role := range roles {
			if user.Role == role {
				c.Next()
				return
			}
		}
		utils.RespondError(c, utils.NewForbidden("Forbidden resource"))
	}
}

// CurrentUser mengambil user login yang disimpan RequireAuth — padanan
// decorator @Session() di controller NestJS.
func CurrentUser(c *gin.Context) *models.User {
	value, exists := c.Get(contextUserKey)
	if !exists {
		return nil
	}
	user, _ := value.(*models.User)
	return user
}
