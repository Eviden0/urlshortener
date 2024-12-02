package shortcode

import "math/rand"

type ShortCodeGenerator interface {
	NextID() string
}

type shortCodeGenerator struct {
	minLength int
}

func NewShortCodeGenerator(minLength int) ShortCodeGenerator {
	return &shortCodeGenerator{
		minLength: minLength,
	}
}

const chars = "abcdefghijklmnopqrstuvwsyzABCDEFJHIJKLMNOKPRSTUVWSVZ0123456789"

func (s *shortCodeGenerator) NextID() string {
	length := len(chars)
	id := make([]byte, s.minLength)

	for i := 0; i < s.minLength; i++ {
		id[i] = chars[rand.Intn(length)]
	}

	return string(id)
}
