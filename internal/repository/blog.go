package repository

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/pkg/common"
	"blog/pkg/configs"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BlogRepository struct {
	db *gorm.DB
}

// NewBlogRepository 创建新的 BlogRepository 实例
func NewBlogRepository() *BlogRepository {
	return &BlogRepository{db: configs.DB}
}

// CreateBlog 保存博客到数据库
func (b *BlogRepository) CreateBlog(blog *models.Blog) error {
	return b.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(blog).Error; err != nil {
			return fmt.Errorf("无法创建博客: %w", err)
		}
		return b.handleTags(tx, blog)
	})
}

func (b *BlogRepository) UpdateBlog(uid int, super bool, blog *models.Blog) error {
	return b.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&models.Blog{}).Select(
			"Description", "Title", "CoverImage", "SourceURL", "Content",
			"CategoryID", "TopicID", "IsPrivate", "Password",
		).Where("id = ?", blog.ID)

		// 如果不是超级用户, 限制只能更新自己的博客
		if !super {
			query = query.Where("user_id = ?", uid)
		}

		// 更新博客的其他字段
		if err := query.Updates(blog).Error; err != nil {
			return fmt.Errorf("无法更新博客: %w", err)
		}

		// 处理标签
		return b.handleTags(tx, blog)
	})
}

// handleTags 处理标签关联
func (b *BlogRepository) handleTags(tx *gorm.DB, blog *models.Blog) error {
	if len(blog.Tags) > 0 {
		if err := tx.Model(blog).Association("Tags").Replace(blog.Tags); err != nil {
			return fmt.Errorf("无法处理标签关联: %w", err)
		}
	} else {
		if err := tx.Model(blog).Association("Tags").Clear(); err != nil {
			return fmt.Errorf("无法清除标签关联: %w", err)
		}
	}
	return nil
}

// GetBlogDetail 获取博客详情
func (b *BlogRepository) GetBlogDetail(id int64) (*models.Blog, error) {
	var blog models.Blog
	err := b.db.Preload("Tags").Preload("Category").Preload("User").Preload("Topic").First(&blog, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("博客不存在")
		}
		return nil, fmt.Errorf("无法获取博客详情: %w", err)
	}
	return &blog, nil
}

// FindAllSearchBlog 查找所有可搜索的博客
func (b *BlogRepository) FindAllSearchBlog() ([]response.SearchBlogResponse, error) {
	var blogs []response.SearchBlogResponse
	err := b.db.Model(&models.Blog{}).Select("id, title, description, created_at").Scan(&blogs).Error
	return blogs, err
}

// GetBlogById 根据 ID 获取博客
func (b *BlogRepository) GetBlogById(id int64) (*models.Blog, error) {
	return b.getBlogWithJoins("b.id = ?", id)
}

// getBlogWithJoins 通过条件获取博客并连接相关表
func (b *BlogRepository) getBlogWithJoins(condition string, args ...interface{}) (*models.Blog, error) {
	var blog models.Blog
	err := b.db.Table(models.BlogTable+" b").
		Scopes(joinUser, joinCategory, joinTopic).
		Select(blogFields).
		Where(condition, args...).Find(&blog).Error
	if err == nil {
		err = b.loadBlogTags(&blog)
	}
	return &blog, err
}

// loadBlogTags 加载博客的标签
func (b *BlogRepository) loadBlogTags(blog *models.Blog) error {
	var tags []models.Tag
	err := b.db.Model(&models.Tag{}).
		Raw(fmt.Sprintf(`SELECT id, name FROM %s WHERE id IN (SELECT tag_id FROM %s WHERE blog_id = ?)`, models.TagTable, models.BlogTagTable), blog.ID).
		Find(&tags).Error
	blog.Tags = tags
	return err
}

// UpdateEyeCount 更新博客浏览次数
func (b *BlogRepository) UpdateEyeCount(id int64, count int64) error {
	return b.db.Exec("UPDATE "+models.BlogTable+" SET eye_count = ? WHERE id = ?", count, id).Error
}

// GetBlogList 获取博客列表
func (b *BlogRepository) GetBlogList(prequest requests.RequestQuery, count *int64) ([]response.BlogResponse, error) {
	var list = make([]response.BlogResponse, 0)

	db := b.db.Model(&models.Blog{}).Table(models.BlogTable + " b")

	if prequest.Cid != nil && *prequest.Cid > 0 {
		db = db.Where("b.category_id = ?", prequest.Cid)
	} else {
		db = db.Where("b.category_id IS NOT NULL")
	}

	if err := db.Count(count).Error; err != nil || *count == 0 {
		return list, nil
	}

	err := db.Scopes(joinUser, joinCategory).
		Select(blogListFields).
		Offset((prequest.Page - 1) * prequest.Size).
		Limit(prequest.Size).
		Order(prequest.Sort.GetBlogOrderString("b.")).
		Find(&list).Error

	return list, err
}

func (b *BlogRepository) UpdatePinned(req requests.PinnedBlogRequest) error {
	return b.db.Model(&models.Blog{}).Where("id = ?", req.Id).Update("pinned", req.Pinned).Update("order", req.Order).Error
}

// 获取所有置顶博客
func (b *BlogRepository) GetPinnedBlogList() []response.BlogResponse {
	var list = make([]response.BlogResponse, 0)

	var fields = append(blogListFields, `t.id AS "Topic__id", t.name AS "Topic__name"`)

	db := b.db.Model(&models.Blog{}).Table(models.BlogTable + " b").Where("b.pinned = true")

	db.Scopes(joinUser, joinCategory, joinTopic).
		Select(fields).
		Order("b.order asc").
		Find(&list)

	return list
}

// GetArchiveBlog 获取归档博客
func (b *BlogRepository) GetArchiveBlog(req requests.ArchiveBlogRequest, page *response.Page) error {
	var blogs []response.ArchiveBlogResponse
	var archiveCount = common.ArchivePageCount

	build := b.db.Table(models.BlogTable).
		Select("id, title, created_at AS create_time, description").
		Where("created_at BETWEEN ? AND ?", req.Start, req.End)

	if err := build.Count(&page.Count).Error; err != nil {
		return err
	}

	page.Data = blogs
	page.Page = req.Page
	page.Size = archiveCount

	if page.Count == 0 {
		return nil
	}

	if err := build.Offset((req.Page - 1) * archiveCount).Limit(archiveCount).
		Order(requests.CREATE.GetBlogOrderString("")).Scan(&blogs).Error; err != nil {
		return err
	}

	page.Data = blogs
	return nil
}

// FindByIdInSimpleBlog 根据博客 ID 列表查找简单博客信息
func (b *BlogRepository) FindByIdInSimpleBlog(ids []int) ([]response.SimpleBlogResponse, error) {
	var blogs []response.SimpleBlogResponse
	err := b.db.Table(models.BlogTable).
		Select("id, title, cover_image").
		Where("id IN ?", ids).
		Scan(&blogs).Error
	return blogs, err
}

// GetHotBlog 获取热门博客
func (b *BlogRepository) GetHotBlog() ([]response.SimpleBlogResponse, error) {
	return b.getSimpleBlogs("id, title", requests.EYE.GetBlogOrderString(""), 10)
}

// GetLatestBlog 获取最新博客
func (b *BlogRepository) GetLatestBlog() ([]response.SimpleBlogResponse, error) {
	return b.getSimpleBlogs("id, title", requests.CREATE.GetBlogOrderString(""), 10)
}

// getSimpleBlogs 获取简单博客信息
func (b *BlogRepository) getSimpleBlogs(selectFields, order string, limit int) ([]response.SimpleBlogResponse, error) {
	var blogs []response.SimpleBlogResponse
	err := b.db.Table(models.BlogTable).
		Select(selectFields).
		Limit(limit).
		Order(order).
		Scan(&blogs).Error
	return blogs, err
}

func (b *BlogRepository) SaveEditBlog(edit models.EditBlog) error {
	var existingEdit models.EditBlog

	// Try to find the existing record
	err := b.db.Where("uid = ?", edit.UID).First(&existingEdit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Record not found, create a new one
			return b.db.Create(&edit).Error
		}
		// Some other error occurred
		return err
	}

	// Record exists, update it
	return b.db.Model(&existingEdit).Updates(models.EditBlog{Content: edit.Content}).Error
}

func (b *BlogRepository) GetEditBlog(uid int) (string, error) {
	var content string

	var err = b.db.Model(&models.EditBlog{}).Select("content").Where("uid = ?", uid).Scan(&content).Error

	return content, err
}

// GetBlogAdminList 获取博客管理列表
func (b *BlogRepository) GetBlogAdminList(uid *int, req requests.AdminFilterRequest, count *int64) ([]response.AdminBlogResponse, error) {
	var blogs []response.AdminBlogResponse
	build := b.db.Model(&models.Blog{}).Table(models.BlogTable + " b")

	if req.Deleted {
		build = build.Unscoped().Where("b.deleted_at IS NOT NULL")
	}

	if uid != nil {
		build = build.Where("b.user_id = ?", uid)
	}

	if req.Pub != nil {
		build = build.Where("is_private = ?", !*req.Pub)
	}

	if req.Category != nil {
		build = build.Where("b.category_id = ?", *req.Category)
	} else if req.Topic != nil {
		build = build.Where("b.topic_id = ?", *req.Topic)
	}

	if req.Keyword != nil {
		like := "%" + *req.Keyword + "%"
		build = build.Where("b.title LIKE ? OR b.description LIKE ?", like, like)
	}

	if req.Start != nil && req.End != nil {
		build = build.Where("b.created_at BETWEEN ? AND ?", req.Start, req.End)
	}

	if err := build.Count(count).Error; err != nil || *count == 0 {
		return blogs, nil
	}

	build = build.Scopes(joinUser, joinCategory, joinTopic).Select(blogFields, "pinned", "order")

	err := build.Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Order(req.Sort.GetBlogOrderString("b.")).Find(&blogs).Error

	return blogs, err
}

// Helper functions and constants
var blogFields = []string{
	"b.id, b.title, b.description, b.cover_image, b.created_at, b.content, b.updated_at, b.eye_count, b.source_url, b.is_private, b.password",
	`u.id AS "User__id", u.nick_name AS "User__nick_name"`,
	`c.id AS "Category__id", c.name AS "Category__name"`,
	`t.id AS "Topic__id", t.name AS "Topic__name"`,
}

var blogListFields = []string{
	"b.id, b.title, b.description, b.cover_image, b.created_at",
	`u.id AS "User__id", u.nick_name AS "User__nick_name"`,
	`c.id AS "Category__id", c.name AS "Category__name"`,
}

func joinUser(db *gorm.DB) *gorm.DB {
	return db.Joins(fmt.Sprintf("INNER JOIN %s u ON u.id = b.user_id", models.UserTable))
}

func joinCategory(db *gorm.DB) *gorm.DB {
	return db.Joins(fmt.Sprintf("LEFT JOIN %s c ON c.id = b.category_id", models.CategoryTable))
}

func joinTopic(db *gorm.DB) *gorm.DB {
	return db.Joins(fmt.Sprintf("LEFT JOIN %s t ON t.id = b.topic_id", models.TopicTable))
}
