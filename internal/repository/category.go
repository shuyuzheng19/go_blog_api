package repository

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/pkg/configs"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

// GetCategoryList 获取分类列表
func (c *CategoryRepository) GetCategoryList() []response.SimpleCategoryResponse {
	var categories = make([]response.SimpleCategoryResponse, 0)
	c.db.Model(&models.Category{}).Select("id, name").Order(requests.UPDATE.GetCategoryOrderString("")).Scan(&categories)
	return categories
}

// UpdateCategory 更新分类名称
func (c *CategoryRepository) UpdateCategory(id int, name string) error {
	// 重要操作，建议使用事务
	return c.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&models.Category{}).Where("id = ?", id).Update("name", name).Error
	})
}

// SaveCategory 保存新分类
func (c *CategoryRepository) SaveCategory(category models.Category) error {
	// 重要操作，建议使用事务
	return c.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&category).Error
	})
}

// DeleteCategoryBlogs 批量软删除指定分类下的所有博客
func (c *CategoryRepository) DeleteCategoryBlogs(categoryIds []int64) error {
	if len(categoryIds) == 0 {
		return errors.New("分类ID列表不能为空")
	}

	return c.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Blog{}).
			Where("category_id IN ?", categoryIds).
			Update("deleted_at", time.Now())

		if result.Error != nil {
			return fmt.Errorf("删除分类博客失败: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return nil // 或者返回一个特定的错误，表示没有博客被删除
		}

		return nil
	})
}

// UndeleteCategoryBlogs 批量恢复指定分类下的所有博客
func (c *CategoryRepository) UndeleteCategoryBlogs(categoryIds []int64) error {
	if len(categoryIds) == 0 {
		return errors.New("分类ID列表不能为空")
	}

	return c.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Blog{}).
			Unscoped().
			Where("category_id IN ?", categoryIds).
			Update("deleted_at", nil)

		if result.Error != nil {
			return fmt.Errorf("恢复分类博客失败: %w", result.Error)
		}

		return nil
	})
}

// GetCategoryAdminList 获取管理员分类列表
func (c *CategoryRepository) GetCategoryAdminList(req requests.AdminFilterRequest, count *int64) ([]models.Category, error) {
	var categories = make([]models.Category, 0)
	query := c.db.Model(&models.Category{}).Table(models.CategoryTable)

	if req.Deleted {
		query = query.Unscoped().Where("deleted_at IS NOT NULL")
	}

	if req.Keyword != nil {
		like := "%" + *req.Keyword + "%"
		query = query.Where("name LIKE ?", like)
	}

	if req.Start != nil && req.End != nil {
		query = query.Where("created_at BETWEEN ? AND ?", req.Start, req.End)
	}

	if err := query.Count(count).Error; err != nil {
		return nil, err
	}

	if *count == 0 {
		return categories, nil
	}

	err := query.Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Order(req.Sort.GetCategoryOrderString("")).
		Find(&categories).Error

	return categories, err
}

// NewCategoryRepository 创建新的 CategoryRepository 实例
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{db: configs.DB}
}
