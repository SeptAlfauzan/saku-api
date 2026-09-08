package services

import (
	"context"

	"github.com/septalfauzan/saku-api/app/server/datasources"
	"github.com/septalfauzan/saku-api/app/server/datasources/remote"
	"github.com/septalfauzan/saku-api/app/server/domain"
)

type ReceiptService interface {
	ExtractReceipt(ctx context.Context, request remote.OCRRequest) (domain.Receipt, error)
}

type DefaultReceiptService struct {
	geminiAPI remote.GeminiAPI
}

func NewReceiptService(ds *datasources.Datasources) *DefaultReceiptService {
	return &DefaultReceiptService{
		geminiAPI: ds.Remote.GeminiAPI,
	}
}

func (s *DefaultReceiptService) ExtractReceipt(ctx context.Context, request remote.OCRRequest) (domain.Receipt, error) {
	return s.geminiAPI.ExtractImageOCR(ctx, request)
}
