package repository

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TagRepository 标签仓储
type TagRepository struct {
	db *gorm.DB
}

// GetTagList 获取标签列表
func (t *TagRepository) GetTagList() []response.SimpleTagResponse {
	var tags = make([]response.SimpleTagResponse, 0)
	t.db.Model(&models.Tag{}).
		Select("id, name").
		Order(requests.UPDATE.GetTagOrderString("")).
		Scan(&tags)
	return tags
}

// GetTagBlogList 获取标签相关的博客列表
func (t *TagRepository) GetTagBlogList(req requests.RequestQuery, count *int64) ([]response.BlogResponse, error) {
	var list = make([]response.BlogResponse, 0)

	if req.Tid == nil {
		return nil, errors.New("标签ID不能为空")
	}

	query := t.db.Model(&models.Blog{}).
		Table(models.BlogTable+" b").
		Joins(fmt.Sprintf("INNER JOIN %s tb ON tb.blog_id = b.id", models.BlogTagTable)).
		Where("tb.tag_id = ?", *req.Tid)

	// 计算总数
	if err := query.Count(count).Error; err != nil {
		return nil, fmt.Errorf("计算博客数量失败: %w", err)
	}

	if *count == 0 {
		return list, nil
	}

	// 查询博客列表
	err := query.Joins(fmt.Sprintf("INNER JOIN %s u ON u.id = b.user_id", models.UserTable)).
		Joins(fmt.Sprintf("LEFT JOIN %s c ON c.id = b.category_id", models.CategoryTable)).
		Select([]string{
			"b.id, b.title, b.description, b.cover_image, b.created_at",
			`u.id AS "User__id", u.nick_name AS "User__nick_name"`,
			`c.id AS "Category__id", c.name AS "Category__name"`,
		}).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Order(req.Sort.GetBlogOrderString("b.")).
		Find(&list).Error

	if err != nil {
		return nil, fmt.Errorf("查询博客列表失败: %w", err)
	}

	return list, nil
}

// GetTagByID 根据ID获取标签
func (t *TagRepository) GetTagByID(id int) *response.SimpleTagResponse {
	var tag response.SimpleTagResponse
	err := t.db.Model(&models.Tag{}).
		Select("id, name").
		Where("id = ?", id).
		First(&tag).Error

	if err != nil {
		logger.Error("获取标签失败", zap.Error(err), zap.Int("id", id))
		return nil
	}
	return &tag
}

// GetTagAdminList 获取管理员标签列表
func (t *TagRepository) GetTagAdminList(req requests.AdminFilterRequest, count *int64) ([]models.Tag, error) {
	var tags = make([]models.Tag, 0)
	query := t.db.Model(&models.Tag{}).Table(models.TagTable)

	// 构建查询条件
	if req.Deleted {
		query = query.Unscoped().Where("deleted_at IS NOT NULL")
	}

	if req.Keyword != nil {
		query = query.Where("name LIKE ?", "%"+*req.Keyword+"%")
	}

	if req.Start != nil && req.End != nil {
		query = query.Where("created_at BETWEEN ? AND ?", req.Start, req.End)
	}

	// 计算总数
	if err := query.Count(count).Error; err != nil {
		return nil, fmt.Errorf("计算标签数量失败: %w", err)
	}

	if *count == 0 {
		return tags, nil
	}

	// 查询标签列表
	err := query.Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Order(req.Sort.GetTagOrderString("")).
		Find(&tags).Error

	if err != nil {
		return nil, fmt.Errorf("查询标签列表失败: %w", err)
	}

	return tags, nil
}

// UpdateTag 更新标签
func (t *TagRepository) UpdateTag(id int, name string) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Tag{}).
			Where("id = ?", id).
			Update("name", name)

		if result.Error != nil {
			return fmt.Errorf("更新标签失败: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return errors.New("标签不存在")
		}

		return nil
	})
}

// SaveTag 保存标签
func (t *TagRepository) SaveTag(tag models.Tag) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&tag).Error; err != nil {
			return fmt.Errorf("创建标签失败: %w", err)
		}
		return nil
	})
}

// DeleteTags 批量删除标签
func (t *TagRepository) DeleteTags(ids []int64) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		// 删除标签
		if err := tx.Where("id IN ?", ids).Delete(&models.Tag{}).Error; err != nil {
			return fmt.Errorf("删除标签失败: %w", err)
		}
		return nil
	})
}

// UnDeleteTags 批量恢复标签
func (t *TagRepository) UnDeleteTags(ids []int64) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Tag{}).
			Unscoped().
			Where("id IN ?", ids).
			Update("deleted_at", nil)

		if result.Error != nil {
			return fmt.Errorf("恢复标签失败: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return errors.New("没有找到要恢复的标签")
		}

		return nil
	})
}

// NewTagRepository 创建标签仓储实例
func NewTagRepository() *TagRepository {
	return &TagRepository{db: configs.DB}
}
