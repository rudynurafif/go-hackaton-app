// Package dto berisi bentuk request body beserta aturan validasinya.
// Padanan dari folder dto/ + class-validator di NestJS. Tag `binding`
// dievaluasi oleh validator bawaan Gin saat ShouldBindJSON dipanggil.
package dto

// RegisterDTO — body untuk POST /auth/register.
// Role sengaja tidak ada di sini: sama seperti `input: false` pada
// auth.config.ts, role selalu jatuh ke default PARTICIPANT dan hanya
// bisa diubah lewat jalur server-side (mis. SQL oleh admin).
type RegisterDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginDTO — body untuk POST /auth/login.
type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
