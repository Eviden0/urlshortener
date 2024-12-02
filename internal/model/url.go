package model

import "time"

type CreateURLRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	CustomCode  string `json:"custom_code,omitempty" validate:"omitempty,min=4,max=10,alphanum"`
	Duration    *int   `json:"duration,omitempty"`
}

type CreateURLResponse struct {
	ShortURL  string    `json:"short_url"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
