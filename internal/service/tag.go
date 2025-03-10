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
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// TagService 标签服务
type TagService struct {
	repository *repository.TagRepository
	cache      *TagCache
}

// NewTagService 创建标签服务实例
func NewTagService() *TagService {
	return &TagService{
		repository: repository.NewTagRepository(),
		cache:      NewTagCache(),
	}
}

// RandomTags 获取随机标签
func (t *TagService) RandomTags() []response.SimpleTagResponse {
	tags, err := t.cache.GetRandomTags()
	if err == nil {
		return tags
	}

	tagList := t.repository.GetTagList()
	go func() {
		if err := t.cache.SetTags(tagList); err != nil {
			logger.Error("缓存随机标签失败", zap.Error(err))
		}
	}()

	logger.Info("添加随机标签缓存")
	return tagList
}

// GetTagList 获取标签列表
func (t *TagService) GetTagList() []response.SimpleTagResponse {
	return t.repository.GetTagList()
}

// GetAdminTagList 获取管理员标签列表
func (t *TagService) GetAdminTagList(req requests.AdminFilterRequest, page *response.Page) error {
	list, err := t.repository.GetTagAdminList(req, &page.Count)
	if err != nil {
		return fmt.Errorf("获取管理员标签列表失败: %w", err)
	}
	page.Data = list
	return nil
}

// SaveTag 保存标签
func (t *TagService) SaveTag(name string) error {
	if err := t.repository.SaveTag(models.Tag{Name: name}); err != nil {
		return fmt.Errorf("保存标签失败: %w", err)
	}

	go func() {
		if err := t.cache.ClearTagKeys(); err != nil {
			logger.Error("清除标签缓存失败", zap.Error(err))
		}
	}()

	return nil
}

// UpdateTag 更新标签
func (t *TagService) UpdateTag(req requests.TagRequest) error {
	if err := t.repository.UpdateTag(req.ID, req.Name); err != nil {
		return fmt.Errorf("更新标签失败: %w", err)
	}

	go func() {
		if err := t.cache.ClearTagKeys(); err != nil {
			logger.Error("清除标签缓存失败", zap.Error(err))
		}
	}()

	return nil
}

// DeleteByIDs 批量删除标签
func (t *TagService) DeleteByIDs(ids []int64) error {
	if err := t.repository.DeleteTags(ids); err != nil {
		return fmt.Errorf("删除标签失败: %w", err)
	}

	go func() {
		if err := t.cache.ClearTagKeys(); err != nil {
			logger.Error("清除标签缓存失败", zap.Error(err))
		}
	}()

	return nil
}

// UnDeleteByIDs 批量恢复标签
func (t *TagService) UnDeleteByIDs(ids []int64) error {
	if err := t.repository.UnDeleteTags(ids); err != nil {
		return fmt.Errorf("恢复标签失败: %w", err)
	}

	go func() {
		if err := t.cache.ClearTagKeys(); err != nil {
			logger.Error("清除标签缓存失败", zap.Error(err))
		}
	}()

	return nil
}

// GetTagByID 根据ID获取标签
func (t *TagService) GetTagByID(id int) *response.SimpleTagResponse {
	result, err := t.cache.GetTagFromMap(id)
	if err != nil {
		result = t.repository.GetTagByID(id)
		go func() {
			if err := t.cache.SetTagToMap(id, result); err != nil {
				logger.Error("缓存标签信息失败", zap.Error(err), zap.Int("tagID", id))
			}
		}()
		logger.Info("缓存标签信息", zap.Any("标签信息", result))
	}
	return result
}

// GetTagIdBlogs 获取标签相关的博客列表
func (t *TagService) GetTagIdBlogs(req requests.RequestQuery, page *response.Page) error {
	blogs, err := t.repository.GetTagBlogList(req, &page.Count)
	if err != nil {
		return fmt.Errorf("获取标签博客列表失败: %w", err)
	}
	page.Data = blogs
	return nil
}

// TagCache 标签缓存
type TagCache struct {
	redis *redis.Client
}

// ClearTagKeys 清除标签相关的所有缓存
func (t *TagCache) ClearTagKeys() error {
	return t.redis.Del(common.TagMapKey, common.RandomTagListKey).Err()
}

// GetRandomTags 获取随机标签
func (t *TagCache) GetRandomTags() ([]response.SimpleTagResponse, error) {
	r := t.redis.SRandMemberN(common.RandomTagListKey, common.RandomTagCount).Val()
	if len(r) == 0 {
		return nil, errors.New("获取随机标签失败")
	}

	result := make([]response.SimpleTagResponse, len(r))
	for i, str := range r {
		result[i] = utils.Deserialize[response.SimpleTagResponse](str)
	}
	return result, nil
}

// SetTags 缓存标签列表
func (t *TagCache) SetTags(tags []response.SimpleTagResponse) error {
	if len(tags) == 0 {
		return errors.New("标签列表为空")
	}

	jsons := make([]interface{}, len(tags))
	for i, tag := range tags {
		jsons[i] = utils.Serialize(tag)
	}

	pipe := t.redis.Pipeline()
	pipe.SAdd(common.RandomTagListKey, jsons...)
	pipe.Expire(common.RandomTagListKey, common.RandomTagListExpire)
	_, err := pipe.Exec()
	return err
}

// SetTagToMap 将标签信息存入缓存map
func (t *TagCache) SetTagToMap(id int, tag *response.SimpleTagResponse) error {
	if tag == nil {
		return errors.New("标签信息为空")
	}
	return t.redis.HSet(common.TagMapKey, strconv.Itoa(id), utils.Serialize(tag)).Err()
}

// GetTagFromMap 从缓存map中获取标签信息
func (t *TagCache) GetTagFromMap(tagID int) (*response.SimpleTagResponse, error) {
	val, err := t.redis.HGet(common.TagMapKey, strconv.Itoa(tagID)).Result()
	if err != nil {
		return nil, fmt.Errorf("从缓存获取标签信息失败: %w", err)
	}
	return utils.Deserialize[*response.SimpleTagResponse](val), nil
}

// NewTagCache 创建标签缓存实例
func NewTagCache() *TagCache {
	return &TagCache{redis: configs.REDIS}
}
