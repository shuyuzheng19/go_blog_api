package service

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/internal/repository"
	"blog/internal/utils"
	"blog/pkg/common"
	"blog/pkg/configs"
	"blog/pkg/logger"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// CategoryService 分类服务
type CategoryService struct {
	repository *repository.CategoryRepository
	cache      *CategoryCache
}

// GetAdminCategoryList 获取管理员分类列表
func (c *CategoryService) GetAdminCategoryList(req requests.AdminFilterRequest, page *response.Page) error {
	list, err := c.repository.GetCategoryAdminList(req, &page.Count)
	page.Data = list
	return err
}

// SaveCategory 保存分类
func (c *CategoryService) SaveCategory(name string) error {
	if err := c.repository.SaveCategory(models.Category{Name: name}); err != nil {
		return err
	}

	go c.clearCategoryCache("保存分类")
	return nil
}

// UpdateCategory 更新分类
func (c *CategoryService) UpdateCategory(req requests.CategoryRequest) error {
	if err := c.repository.UpdateCategory(req.ID, req.Name); err != nil {
		return err
	}

	go c.clearCategoryCache("更新分类")
	return nil
}

// DeleteByIDs 批量删除分类
func (c *CategoryService) DeleteByIDs(ids []int64) error {
	if err := configs.DeleteData(models.CategoryTable, nil, ids); err != nil {
		return err
	}

	go func() {
		if err := c.repository.DeleteCategoryBlogs(ids); err != nil {
			logger.Info("删除分类博客失败", zap.String("error", err.Error()))
		}
		c.clearCategoryCache("删除分类")
	}()

	return nil
}

// UnDeleteByIDs 批量恢复分类
func (c *CategoryService) UnDeleteByIDs(ids []int64) error {
	if err := configs.UnDeleteData(models.CategoryTable, nil, ids); err != nil {
		return err
	}

	go func() {
		if err := c.repository.UndeleteCategoryBlogs(ids); err != nil {
			logger.Info("恢复分类博客失败", zap.String("error", err.Error()))
		}
		c.clearCategoryCache("恢复分类")
	}()

	return nil
}

// GetCategoryList 获取分类列表
func (c *CategoryService) GetCategoryList() []response.SimpleCategoryResponse {
	categories, err := c.cache.GetCategoryList()
	if err != nil {
		categories = c.repository.GetCategoryList()
		logger.Info("缓存分类列表")
		go func() {
			if err := c.cache.SetCategoryList(categories); err != nil {
				logger.Info("设置分类列表缓存失败", zap.String("error", err.Error()))
			}
		}()
	}
	return categories
}

// clearCategoryCache 清除分类缓存
func (c *CategoryService) clearCategoryCache(action string) {
	if err := c.cache.ClearCategoryKeys(); err != nil {
		logger.Info(action+"时清除分类缓存失败", zap.String("error", err.Error()))
	}
}

// NewCategoryService 创建新的 CategoryService 实例
func NewCategoryService() *CategoryService {
	return &CategoryService{
		repository: repository.NewCategoryRepository(),
		cache:      NewCategoryCache(),
	}
}

// CategoryCache 分类缓存
type CategoryCache struct {
	redis *redis.Client
}

// SetCategoryList 缓存分类列表
func (c *CategoryCache) SetCategoryList(categories []response.SimpleCategoryResponse) error {
	str := utils.Serialize(categories)
	return c.redis.Set(common.CategoryListKey, str, common.CategoryListExpire).Err()
}

// ClearCategoryKeys 清除分类相关的缓存
func (c *CategoryCache) ClearCategoryKeys() error {
	pages := c.redis.Keys(common.PageInfoPrefixKey + "*").Val()
	return c.redis.Del(append(pages, common.BlogMapKey, common.CategoryListKey)...).Err()
}

// GetCategoryList 从缓存获取分类列表
func (c *CategoryCache) GetCategoryList() ([]response.SimpleCategoryResponse, error) {
	r, err := c.redis.Get(common.CategoryListKey).Result()
	if err != nil {
		return nil, err
	}
	return utils.Deserialize[[]response.SimpleCategoryResponse](r), nil
}

// NewCategoryCache 创建新的 CategoryCache 实例
func NewCategoryCache() *CategoryCache {
	return &CategoryCache{redis: configs.REDIS}
}
