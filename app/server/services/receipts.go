package services

import (
	"context"

	"github.com/septalfauzan/saku-api/app/server/domain"
)

type ReceiptService interface {
	ExtractReceipt(ctx context.Context) (domain.Receipt, error)
}

type DefaultReceiptService struct {
	receipt domain.Receipt
}

func NewReceiptService() *DefaultReceiptService {
	return &DefaultReceiptService{
		receipt: domain.Receipt{},
	}
}

func (s *DefaultReceiptService) ExtractReceipt(ctx context.Context) (domain.Receipt, error) {
	return s.receipt, nil
}
