package store

import (
	"context"

	"github.com/A-pen-app/ai-client/models"
)

type OCR interface {
	ScanName(ctx context.Context, link string) (string, error)
	// ScanCompany reads the company name and tax id off one registration
	// document. Matching them against anything is the caller's business.
	ScanCompany(ctx context.Context, file models.InputFile) (*models.CompanyInfo, error)
	ScanRawInfo(ctx context.Context, userID string, link string, professionType models.PlatformType) (*models.OCRRawInfo, error)
}

type Article interface {
	ExtractTags(ctx context.Context, content string, professionType models.PlatformType) (*models.ExtractTagsResult, error)
	Polish(ctx context.Context, content string, professionType models.PlatformType) (string, error)
}

type AIClient interface {
	Generate(ctx context.Context, message models.AIChatMessage, opts models.AIClientOptions) (string, error)
}
