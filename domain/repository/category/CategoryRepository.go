package category

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
)

type categoryRepository struct {
	*gorm.DB
}

func ProvideCategoryRepository(logger *zap.SugaredLogger, conn *gorm.DB) _interface.CategoryRepository {
	logger.Info("Executing ProvideCategoryRepository.")
	return &categoryRepository{conn}
}

func (m *categoryRepository) GetByID(ctx context.Context, id uint) (entitie.Category, error) {
	var category entitie.Category
	tx := m.DB.First(&category, id)
	if tx.Error != nil {
		return entitie.Category{}, tx.Error
	}
	return category, nil
}

func (m *categoryRepository) GetAll(ctx context.Context) ([]entitie.Category, error) {
	var category []entitie.Category
	tx := m.DB.Find(&category)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return category, nil
}

func (m *categoryRepository) GetByName(ctx context.Context, name string) (entitie.Category, error) {
	var category entitie.Category
	tx := m.DB.Model(&entitie.Category{}).Where(&entitie.Category{Name: name}).First(&category)
	if tx.Error != nil {
		return entitie.Category{}, tx.Error
	}
	return category, nil
}

func (m *categoryRepository) Update(ctx context.Context, ar *entitie.Category) error {
	tx := m.DB.Save(&ar)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *categoryRepository) Store(ctx context.Context, a *entitie.Category) error {
	tx := m.DB.Create(&a)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *categoryRepository) Delete(ctx context.Context, id uint) error {
	tx := m.DB.Delete(&entitie.Category{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
