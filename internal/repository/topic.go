package repository

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/pkg/common"
	"blog/pkg/configs"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TopicRepository 专题仓储
type TopicRepository struct {
	db *gorm.DB
}

// FindTopicByPage 分页获取专题列表
func (t *TopicRepository) FindTopicByPage(page int) ([]response.TopicResponse, int64) {
	var topics = make([]response.TopicResponse, 0)
	var count int64

	query := t.db.Model(&models.Topic{}).Table(models.TopicTable + " t").
		Joins(fmt.Sprintf("INNER JOIN %s u ON u.id = t.user_id", models.UserTable))

	if err := query.Count(&count).Error; err != nil {
		return nil, 0
	}

	if count == 0 {
		return topics, 0
	}

	fields := []string{
		"t.id", "t.name", "t.cover_image", "t.created_at", "t.description",
		`u.id AS "User__id"`, `u.nick_name AS "User__nick_name"`,
	}

	query.Select(fields).
		Offset((page - 1) * common.TopicPageCount).
		Limit(common.TopicPageCount).
		Find(&topics)

	return topics, count
}

// GetAllTopicList 获取所有专题列表
func (t *TopicRepository) GetAllTopicList() ([]response.SimpleTopicResponse, error) {
	var list = make([]response.SimpleTopicResponse, 0)
	err := t.db.Model(&models.Topic{}).
		Select("id, name").
		Find(&list).Error
	return list, err
}

// GetTopicBlogList 获取专题相关的博客列表
func (t *TopicRepository) GetTopicBlogList(req requests.RequestQuery, count *int64) ([]response.BlogResponse, error) {
	var list = make([]response.BlogResponse, 0)

	query := t.db.Model(&models.Blog{}).
		Table(models.BlogTable + " b")

	if req.Tid != nil && *req.Tid > 0 {
		query = query.Where("b.topic_id = ?", req.Tid)
	}

	if err := query.Count(count).Error; err != nil {
		return nil, fmt.Errorf("计算博客数量失败: %w", err)
	}

	if *count == 0 {
		return list, nil
	}

	err := query.Joins(fmt.Sprintf("INNER JOIN %s u ON u.id = b.user_id", models.UserTable)).
		Joins(fmt.Sprintf("INNER JOIN %s t ON t.id = b.topic_id", models.TopicTable)).
		Select([]string{
			"b.id", "b.title", "b.description", "b.cover_image", "b.created_at",
			`u.id AS "User__id"`, `u.nick_name AS "User__nick_name"`,
		}).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Order(requests.BACK.GetBlogOrderString("b.")).
		Find(&list).Error

	if err != nil {
		return nil, fmt.Errorf("查询博客列表失败: %w", err)
	}

	return list, nil
}

// FindById 根据ID查找专题
func (t *TopicRepository) FindById(id int) *response.SimpleTopicResponse {
	var topic response.SimpleTopicResponse
	if err := t.db.Model(&models.Topic{}).
		Select("id, name").
		First(&topic, "id = ?", id).Error; err != nil {
		return nil
	}
	return &topic
}

// GetTopicBlogs 获取专题下的所有博客
func (t *TopicRepository) GetTopicBlogs(tid int) []response.SimpleBlogResponse {
	var list = make([]response.SimpleBlogResponse, 0)
	t.db.Model(&models.Blog{}).
		Select("id, title").
		Where("topic_id = ?", tid).
		Order(requests.BACK.GetBlogOrderString("")).
		Find(&list)
	return list
}

// Create 创建专题
func (t *TopicRepository) Create(topic models.Topic) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&topic).Error; err != nil {
			return fmt.Errorf("创建专题失败: %w", err)
		}
		return nil
	})
}

// Update 更新专题
func (t *TopicRepository) Update(topic models.Topic) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Topic{}).
			Where("id = ?", topic.ID).
			Updates(topic)

		if result.Error != nil {
			return fmt.Errorf("更新专题失败: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return errors.New("专题不存在")
		}

		return nil
	})
}

// GetAdminTopic 获取管理员专题列表
func (t *TopicRepository) GetAdminTopic(req requests.AdminFilterRequest, count *int64) ([]response.AdminTopicResponse, error) {
	var result = make([]response.AdminTopicResponse, 0)
	query := t.db.Model(&models.Topic{}).
		Table(models.TopicTable + " t")

	// 构建查询条件
	if req.Deleted {
		query = query.Unscoped().Where("t.deleted_at IS NOT NULL")
	} else {
		query = query.Where("t.deleted_at IS NULL")
	}

	if req.Start != nil && req.End != nil {
		query = query.Where("t.created_at BETWEEN ? AND ?", req.Start, req.End)
	}

	if req.Keyword != nil {
		query = query.Where("t.name LIKE ?", "%"+*req.Keyword+"%")
	}

	if err := query.Count(count).Error; err != nil {
		return nil, fmt.Errorf("计算专题数量失败: %w", err)
	}

	if *count == 0 {
		return result, nil
	}

	err := query.Joins(fmt.Sprintf("INNER JOIN %s u ON u.id = t.user_id", models.UserTable)).
		Select([]string{
			"t.id", "t.name", "t.description", "t.cover_image", "t.user_id",
			"t.created_at", "t.updated_at",
			`u.id AS "User__id"`, `u.nick_name AS "User__nick_name"`,
		}).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Order(req.Sort.GetTopicOrderString("t.")).
		Find(&result).Error

	if err != nil {
		return nil, fmt.Errorf("查询专题列表失败: %w", err)
	}

	return result, nil
}

// DeleteTopicBlogs 批量软删除指定专题下的所有博客
func (t *TopicRepository) DeleteTopicBlogs(topicIds []int64) error {
	if len(topicIds) == 0 {
		return errors.New("专题ID列表不能为空")
	}

	return t.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Blog{}).
			Where("topic_id IN ?", topicIds).
			Update("deleted_at", time.Now())

		if result.Error != nil {
			return fmt.Errorf("删除专题博客失败: %w", result.Error)
		}

		return nil
	})
}

// UndeleteTopicBlogs 批量恢复指定专题下的所有博客
func (t *TopicRepository) UndeleteTopicBlogs(topicIds []int64) error {
	if len(topicIds) == 0 {
		return errors.New("专题ID列表不能为空")
	}

	return t.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Blog{}).
			Unscoped().
			Where("topic_id IN ?", topicIds).
			Update("deleted_at", nil)

		if result.Error != nil {
			return fmt.Errorf("恢复专题博客失败: %w", result.Error)
		}

		return nil
	})
}

// NewTopicRepository 创建专题仓储实例
func NewTopicRepository() *TopicRepository {
	return &TopicRepository{db: configs.DB}
}
