package usecase

import (
	"context"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

type TableService struct {
	tableRepo ports.TableRepository
}

func NewTableService(tableRepo ports.TableRepository) *TableService {
	return &TableService{tableRepo: tableRepo}
}

func (s *TableService) GetAll(ctx context.Context) ([]domain.Table, error) {
	return s.tableRepo.GetAll(ctx)
}

func (s *TableService) GetByID(ctx context.Context, id int) (*domain.Table, error) {
	return s.tableRepo.GetByID(ctx, id)
}

func (s *TableService) Create(ctx context.Context, table *domain.Table) error {
	table.Status = domain.TableFree
	return s.tableRepo.Create(ctx, table)
}

func (s *TableService) UpdateStatus(ctx context.Context, id int, status domain.TableStatus) error {
	return s.tableRepo.UpdateStatus(ctx, id, status)
}

func (s *TableService) Delete(ctx context.Context, id int) error {
	return s.tableRepo.Delete(ctx, id)
}
