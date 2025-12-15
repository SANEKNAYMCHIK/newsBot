package models

import "time"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type AddSubscriptionRequest struct {
	SourceID int64 `json:"source_id" binding:"required"`
}

type SubscriptionResponse struct {
	SourceID   int64  `json:"source_id"`
	SourceName string `json:"source_name"`
	CategoryID *int64 `json:"category_id,omitempty"`
	IsActive   bool   `json:"is_active"`
}

type NewsResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Content     *string   `json:"content,omitempty"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	SourceID    int64     `json:"source_id"`
	SourceName  string    `json:"source_name"`
	CategoryID  *int64    `json:"category_id,omitempty"`
}

type CreateSourceRequest struct {
	Name       string `json:"name" binding:"required"`
	URL        string `json:"url" binding:"required,url"`
	CategoryID *int64 `json:"category_id,omitempty"`
	IsActive   bool   `json:"is_active"`
}

type UpdateSourceRequest struct {
	Name       *string `json:"name,omitempty"`
	URL        *string `json:"url,omitempty" binding:"omitempty,url"`
	CategoryID *int64  `json:"category_id,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}
