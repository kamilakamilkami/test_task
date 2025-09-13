package service

import (
	"context"
	"project/internal/domain"
	"project/repository"

	"github.com/google/uuid"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"strings"
	"log"
)

type TemplateService interface {
	Create(ctx context.Context, template *domain.Template) (*domain.Template, error)
	GetAll(ctx context.Context, templateType *string, active *bool ) ([]domain.Template, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error)
	Update(ctx context.Context, id uuid.UUID, template *domain.Template) (*domain.Template, error)
	FromHtmlToPdf(ctx context.Context, html string) ([]byte, error)
	// Delete(ctx context.Context, id uuid.UUID) error
}

type templateService struct {
	repo repository.TemplateRepository
}

func NewTemplateService(repo repository.TemplateRepository) TemplateService {
	return &templateService{repo: repo}
}

func (s *templateService) Create(ctx context.Context, template *domain.Template) (*domain.Template, error) {
	// template.ID = uuid.New()
	if err := s.repo.Create(ctx, template); err != nil {
		return nil, err
	}
	return template, nil
}

func (s *templateService) GetAll(ctx context.Context, templateType *string, active *bool) ([]domain.Template, error) {
	return s.repo.GetAll(ctx, templateType, active)
}

func (s *templateService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *templateService) Update(ctx context.Context, id uuid.UUID, template *domain.Template) (*domain.Template, error) {
	return s.repo.Update(ctx, id, template)
}

func (s *templateService) FromHtmlToPdf(ctx context.Context, html string) ([]byte, error) {
    pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)
	pdfg.NoCollate.Set(false)
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)
	pdfg.MarginTop.Set(10)
	pdfg.MarginBottom.Set(10)

	page := wkhtmltopdf.NewPageReader(strings.NewReader(html))

	page.FooterRight.Set("[page]") 

	pdfg.AddPage(page)

	if err := pdfg.Create(); err != nil {
		log.Fatal(err)
	}

	pdfBytes := pdfg.Bytes()
    return pdfBytes, nil
}
