package repository

import (
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/domain/model"
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContentRepository interface {
	GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, error)
	GetContentWithQuery(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
}

type contentRepository struct {
	db *gorm.DB
}

// CreateContent implements ContentRepository.
func (c *contentRepository) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		Excerp:      req.Excerp,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: req.CreatedByID,
	}

	err = c.db.Create(&modelContent).Error
	if err != nil {
		code = "[REPOSITORY] CreateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// DeleteContent implements ContentRepository.
func (c *contentRepository) DeleteContent(ctx context.Context, id int64) error {
	err = c.db.Where("id = ?", id).Delete(&model.Content{}).Error
	if err != nil {
		code = "[REPOSITORY] DeleteContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetContentByID implements ContentRepository.
func (c *contentRepository) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	var modelContent model.Content

	err = c.db.Where("id = ?", id).Preload(clause.Associations).First(&modelContent).Error
	if err != nil {
		code = "[REPOSITORY] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	tags := strings.Split(modelContent.Tags, ",")
	resp := entity.ContentEntity{
		ID:          modelContent.ID,
		Title:       modelContent.Title,
		Excerp:      modelContent.Excerp,
		Description: modelContent.Description,
		Image:       modelContent.Image,
		Tags:        tags,
		Status:      modelContent.Status,
		CategoryID:  modelContent.CategoryID,
		CreatedByID: modelContent.CreatedByID,
		CreatedAt:   modelContent.CreatedAt,
		Category: entity.CategoryEntity{
			ID:    modelContent.Category.ID,
			Title: modelContent.Category.Title,
			Slug:  modelContent.Category.Slug,
		},
		User: entity.UserEntity{
			ID:   modelContent.User.ID,
			Nama: modelContent.User.Nama,
		},
	}

	return &resp, nil
}

func (c *contentRepository) GetContentWithQuery(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int, error) {
	var modelContents []model.Content
	var totalCount int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit
	status := ""
	if query.Status != "" {
		status = query.Status
	}

	sqlMain := c.db.Model(&modelContents).Preload(clause.Associations).
		Where("title ilike ? OR excerp ilike ? OR description ilike ?", "%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Search+"%").
		Where("status LIKE ?", "%"+status+"%")

	if query.CategoryID != 0 {
		sqlMain = sqlMain.Where("category_id = ?", query.CategoryID)
	}

	err := sqlMain.Count(&totalCount).Error
	if err != nil {
		code := "[REPOSITORY] GetContentWithQuery - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	err = sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&modelContents).Error
	if err != nil {
		code = "[REPOSITORY] GetContentWithQuery - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	resps := []entity.ContentEntity{}
	for _, val := range modelContents {
		tags := strings.Split(val.Tags, ",")
		resp := entity.ContentEntity{
			ID:          val.ID,
			Title:       val.Title,
			Excerp:      val.Excerp,
			Description: val.Description,
			Image:       val.Image,
			Tags:        tags,
			Status:      val.Status,
			CategoryID:  val.CategoryID,
			CreatedByID: val.CreatedByID,
			CreatedAt:   val.CreatedAt,
			Category: entity.CategoryEntity{
				ID:    val.Category.ID,
				Title: val.Category.Title,
				Slug:  val.Category.Slug,
			},
			User: entity.UserEntity{
				ID:   val.User.ID,
				Nama: val.User.Nama,
			},
		}

		resps = append(resps, resp)
	}

	totalPages := int(totalCount / int64(query.Limit))
	if totalCount%int64(query.Limit) > 0 {
		totalPages++
	}

	return resps, totalCount, totalPages, nil
}

// GetContents implements ContentRepository.
func (c *contentRepository) GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, error) {
	var modelContents []model.Content

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit
	status := ""
	if query.Status != "" {
		status = query.Status
	}

	sqlMain := c.db.Preload(clause.Associations).
		Where("title ilike ? OR excerp ilike ? OR description ilike ?", "%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Search+"%").
		Where("status LIKE ?", "%"+status+"%")

	if query.CategoryID != 0 {
		sqlMain = sqlMain.Where("category_id = ?", query.CategoryID)
	}

	err = sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&modelContents).Error
	if err != nil {
		code = "[REPOSITORY] GetContents - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resps := []entity.ContentEntity{}
	for _, val := range modelContents {
		tags := strings.Split(val.Tags, ",")
		resp := entity.ContentEntity{
			ID:          val.ID,
			Title:       val.Title,
			Excerp:      val.Excerp,
			Description: val.Description,
			Image:       val.Image,
			Tags:        tags,
			Status:      val.Status,
			CategoryID:  val.CategoryID,
			CreatedByID: val.CreatedByID,
			CreatedAt:   val.CreatedAt,
			Category: entity.CategoryEntity{
				ID:    val.Category.ID,
				Title: val.Category.Title,
				Slug:  val.Category.Slug,
			},
			User: entity.UserEntity{
				ID:   val.User.ID,
				Nama: val.User.Nama,
			},
		}

		resps = append(resps, resp)
	}

	return resps, nil
}

// UpdateContent implements ContentRepository.
func (c *contentRepository) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		Excerp:      req.Excerp,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: req.CreatedByID,
	}

	err = c.db.Where("id = ?", req.ID).Updates(&modelContent).Error
	if err != nil {
		code = "[REPOSITORY] UpdateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}
}
