package sms

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Provider interface {
	SendCode(ctx context.Context, phone, code string) error
}

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) SendCode(ctx context.Context, phone, code string) error {
	fmt.Printf("[MOCK SMS] send code %s to phone %s\n", code, phone)
	return nil
}

func GenerateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}
