package category

import (
	"context"
	"gorm.io/gorm"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
)

type categoryRepository struct {
	*gorm.DB
}

func NewCategoryRepository(conn *gorm.DB) _interface.CategoryRepository {
	// FIXME the best WAY?
	conn.AutoMigrate(&entitie.Category{})
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

func (m *categoryRepository) GetByName(ctx context.Context, title string) (entitie.Category, error) {
	return entitie.Category{}, nil
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
