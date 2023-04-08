package tag

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"microservice/domain/entitie"
	_interface "microservice/domain/interface"
)

type tagRepository struct {
	*gorm.DB
}

func ProvideTagRepository(logger *zap.SugaredLogger, conn *gorm.DB) _interface.TagRepository {
	logger.Info("Executing ProvideTagRepository.")
	return &tagRepository{conn}
}

func (m *tagRepository) GetByID(ctx context.Context, id uint) (entitie.Tag, error) {
	var tag entitie.Tag
	tx := m.DB.First(&tag, id)
	if tx.Error != nil {
		return entitie.Tag{}, tx.Error
	}
	return tag, nil
}

func (m *tagRepository) GetAll(ctx context.Context) ([]entitie.Tag, error) {
	var tag []entitie.Tag
	tx := m.DB.Find(&tag)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tag, nil
}

func (m *tagRepository) GetByName(ctx context.Context, title string) (entitie.Tag, error) {
	var tag entitie.Tag
	tx := m.DB.Model(&entitie.Tag{}).Where(&entitie.Tag{Name: title}).First(&tag)
	if tx.Error != nil {
		return entitie.Tag{}, tx.Error
	}
	return tag, nil
}

func (m *tagRepository) Update(ctx context.Context, ar *entitie.Tag) error {
	tx := m.DB.Save(&ar)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *tagRepository) Store(ctx context.Context, a *entitie.Tag) error {
	tx := m.DB.Create(&a)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *tagRepository) Delete(ctx context.Context, id uint) error {
	tx := m.DB.Delete(&entitie.Tag{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
