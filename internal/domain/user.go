package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint64         `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"column:nama;not null" validate:"required,min=2,max=100"`
	Email           string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Phone           string         `json:"phone" gorm:"column:notelp;uniqueIndex" validate:"required,min=10,max=15"`
	Password        string         `json:"-" gorm:"column:kata_sandi;not null" validate:"required,min=6"`
	DateOfBirth     *time.Time     `json:"date_of_birth" gorm:"column:tanggal_lahir"`
	Gender          string         `json:"gender" gorm:"column:jenis_kelamin"`
	About           string         `json:"about" gorm:"column:tentang"`
	Job             string         `json:"job" gorm:"column:pekerjaan"`
	ProvinceID      *uint64        `json:"province_id" gorm:"column:id_provinsi"`
	CityID          *uint64        `json:"city_id" gorm:"column:id_kota"`
	IsAdmin         bool           `json:"is_admin" gorm:"default:false"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	Status          string         `json:"status" gorm:"default:active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id uint64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByPhone(phone string) (*User, error)
	Update(user *User) error
	Delete(id uint64) error
}

type RegisterRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required,min=10,max=15"`
	Password    string `json:"password" validate:"required,min=6"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdateProfileRequest struct {
	Name        string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone       string  `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
	DateOfBirth string  `json:"date_of_birth,omitempty"`
	Gender      string  `json:"gender,omitempty" validate:"omitempty,oneof=L P"`
	About       string  `json:"about,omitempty"`
	Job         string  `json:"job,omitempty"`
	ProvinceID  *uint64 `json:"province_id,omitempty"`
	CityID      *uint64 `json:"city_id,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}