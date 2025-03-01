package repository

import (
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/domain/model"
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	EditCategoryID(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *gorm.DB
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepository) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	var countSlug int64
	err = c.db.Table("categories").Where("slug = ?", req.Slug).Count(&countSlug).Error
	if err != nil {
		code = "[REPOSITORY] CreateCategory - 2"
		log.Errorw(code, err)
		return err
	}

	countSlug = countSlug + 1
	slug := fmt.Sprintf("%s-%d", req.Slug, countSlug)
	modelCategory := model.Category{
		Title:       req.Title,
		Slug:        slug,
		CreatedByID: req.User.ID,
	}

	err = c.db.Create(&modelCategory).Error
	if err != nil {
		code = "[REPOSITORY] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// DeleteCategory implements CategoryRepository.
func (c *categoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	var count int64
	err = c.db.Table("contents").Where("id = ?", id).Count(&count).Error
	if err != nil {
		code = "[REPOSITORY] DeleteCategory - 1"
		log.Errorw(code, err)
		return err
	}

	if count > 0 {
		return errors.New("Cannot delete a category that has associated content")
	}

	err = c.db.Where("id = ?", id).Delete(&model.Category{}).Error
	if err != nil {
		code = "[REPOSITORY] DeleteCategory - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// EditCategoryID implements CategoryRepository.
func (c *categoryRepository) EditCategoryID(ctx context.Context, req entity.CategoryEntity) error {
	var countSlug int64
	err = c.db.Table("categories").Where("slug = ?", req.Slug).Count(&countSlug).Error
	if err != nil {
		code = "[REPOSITORY] EditCategoryID - 1"
		log.Errorw(code, err)
		return err
	}

	countSlug = countSlug + 1
	slug := req.Slug
	if countSlug == 0 {
		slug = fmt.Sprintf("%s-%d", req.Slug, countSlug)
	}
	modelCategory := model.Category{
		Title:       req.Title,
		Slug:        slug,
		CreatedByID: req.User.ID,
	}

	err = c.db.Where("id = ?", req.ID).Updates(&modelCategory).Error
	if err != nil {
		code = "[REPOSITORY] EditCategoryID - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetCategories implements CategoryRepository.
func (c *categoryRepository) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	var modelCategories []model.Category

	err = c.db.Order("created_at DESC").Preload("User").Find(&modelCategories).Error
	if err != nil {
		code = "[REPOSITORY] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if len(modelCategories) == 0 {
		code = "[REPOSITORY] GetCategories - 2"
		err = errors.New("data not found")
		log.Errorw(code, err)
		return nil, err
	}

	var resps []entity.CategoryEntity
	for _, val := range modelCategories {
		resps = append(resps, entity.CategoryEntity{
			ID:    val.ID,
			Title: val.Title,
			Slug:  val.Slug,
			User: entity.UserEntity{
				ID:       val.User.ID,
				Nama:     val.User.Nama,
				Email:    val.User.Email,
				Password: val.User.Password,
			},
		})
	}

	return resps, nil
}

// GetCategoryByID implements CategoryRepository.
func (c *categoryRepository) GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	var modelCategory model.Category
	err = c.db.Where("id = ?", id).Preload("User").First(&modelCategory).Error
	if err != nil {
		code = "[REPOSITORY] GetCategoryByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.CategoryEntity{
		ID:    modelCategory.ID,
		Title: modelCategory.Title,
		Slug:  modelCategory.Slug,
		User: entity.UserEntity{
			ID:    modelCategory.User.ID,
			Nama:  modelCategory.User.Nama,
			Email: modelCategory.User.Email,
		},
	}, nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}
