package services

import (
	"context"

	"github.com/septalfauzan/saku-api/app/server/domain"
)

type ReceiptService interface {
	ExtractReceipt(ctx context.Context, request domain.OCRRequest) (domain.Receipt, error)
}

type DefaultReceiptService struct {
	extractor domain.ReceiptExtractor
}

func NewReceiptService(extractor domain.ReceiptExtractor) *DefaultReceiptService {
	return &DefaultReceiptService{
		extractor: extractor,
	}
}

func (s *DefaultReceiptService) ExtractReceipt(ctx context.Context, request domain.OCRRequest) (domain.Receipt, error) {
	return s.extractor.ExtractImageOCR(ctx, request)
}
