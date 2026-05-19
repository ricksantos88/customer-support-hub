package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/ricksantos88/customer-support-hub/internal/models"
)

type ContactRepositoryMock struct {
	mock.Mock
}

func (m *ContactRepositoryMock) Create(ctx context.Context, contact *models.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *ContactRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*models.Contact, error) {
	args := m.Called(ctx, id)
	if contact, ok := args.Get(0).(*models.Contact); ok {
		return contact, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ContactRepositoryMock) GetByPhone(ctx context.Context, phone string) (*models.Contact, error) {
	args := m.Called(ctx, phone)
	if contact, ok := args.Get(0).(*models.Contact); ok {
		return contact, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ContactRepositoryMock) Update(ctx context.Context, contact *models.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *ContactRepositoryMock) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ContactRepositoryMock) List(ctx context.Context, limit, offset int) ([]models.Contact, error) {
	args := m.Called(ctx, limit, offset)
	if contacts, ok := args.Get(0).([]models.Contact); ok {
		return contacts, args.Error(1)
	}
	return nil, args.Error(1)
}
