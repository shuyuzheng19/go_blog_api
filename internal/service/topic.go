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

// TopicService 专题服务
type TopicService struct {
	repository *repository.TopicRepository
	cache      *TopicCache
}

// GetTopicByPage 分页获取专题列表
func (t *TopicService) GetTopicByPage(pindex int) *response.Page {

	pageInfo, err := t.cache.GetPageInfo(pindex)

	if err == nil {
		return pageInfo
	}

	pageInfo = &response.Page{}

	pageInfo.Size = common.TopicPageCount

	pageInfo.Page = pindex

	list, count := t.repository.FindTopicByPage(pindex)

	pageInfo.Count = count

	pageInfo.Data = list

	go t.cache.SetPageInfo(pindex, pageInfo)

	return pageInfo
}

// GetAdminTopicList 获取管理员专题列表
func (t *TopicService) GetAdminTopicList(req requests.AdminFilterRequest, page *response.Page) error {
	list, err := t.repository.GetAdminTopic(req, &page.Count)
	if err != nil {
		return fmt.Errorf("获取管理员专题列表失败: %w", err)
	}
	page.Data = list
	return nil
}

// SaveTopic 保存专题
func (t *TopicService) SaveTopic(uid int, req requests.TopicRequest) error {
	if err := t.repository.Create(req.ToModel(uid)); err != nil {
		return fmt.Errorf("创建专题失败: %w", err)
	}

	go func() {
		if err := t.cache.ClearTopicKeys(); err != nil {
			logger.Error("清除专题缓存失败", zap.Error(err))
		}
	}()

	return nil
}

// UpdateTopic 更新专题
func (t *TopicService) UpdateTopic(uid int, req requests.TopicRequest) error {
	if err := t.repository.Update(req.ToModel(uid)); err != nil {
		return fmt.Errorf("更新专题失败: %w", err)
	}

	go func() {
		if err := t.cache.ClearTopicKeys(); err != nil {
			logger.Error("清除专题缓存失败", zap.Error(err))
		}
	}()

	return nil
}

// DeleteByIDs 批量删除专题
func (t *TopicService) DeleteByIDs(ids []int64) error {
	if err := configs.DeleteData(models.TopicTable, nil, ids); err != nil {
		return fmt.Errorf("删除专题失败: %w", err)
	}

	go func() {
		if err := t.cache.ClearTopicKeys(); err != nil {
			logger.Error("清除专题缓存失败", zap.Error(err))
		}
		if err := t.repository.DeleteTopicBlogs(ids); err != nil {
			logger.Error("删除专题博客失败", zap.Error(err))
		}
	}()

	return nil
}

// UnDeleteByIDs 批量恢复专题
func (t *TopicService) UnDeleteByIDs(ids []int64) error {
	if err := configs.UnDeleteData(models.TopicTable, nil, ids); err != nil {
		return fmt.Errorf("恢复专题失败: %w", err)
	}

	go func() {
		if err := t.cache.ClearTopicKeys(); err != nil {
			logger.Error("清除专题缓存失败", zap.Error(err))
		}
		if err := t.repository.UndeleteTopicBlogs(ids); err != nil {
			logger.Error("恢复专题博客失败", zap.Error(err))
		}
	}()

	return nil
}

// GetAllTopicList 获取所有专题列表
func (t *TopicService) GetAllTopicList() ([]response.SimpleTopicResponse, error) {
	return t.repository.GetAllTopicList()
}

// GetTopicBlogs 获取专题博客列表
func (t *TopicService) GetTopicBlogs(tid int) []response.SimpleBlogResponse {
	return t.repository.GetTopicBlogs(tid)
}

// GetTopicBlogList 获取专题博客分页列表
func (t *TopicService) GetTopicBlogList(req requests.RequestQuery, page *response.Page) error {
	list, err := t.repository.GetTopicBlogList(req, &page.Count)

	if err != nil {
		page.Data = make([]response.BlogResponse, 0)
		return fmt.Errorf("获取专题博客列表失败: %w", err)
	}

	page.Data = list
	return nil
}

// GetTopicInfo 获取专题信息
func (t *TopicService) GetTopicInfo(id int) *response.SimpleTopicResponse {
	result, err := t.cache.GetTopicFromMap(id)
	if err != nil {
		result = t.repository.FindById(id)
		go func() {
			if err := t.cache.SetTopicToMap(id, result); err != nil {
				logger.Error("缓存专题信息失败", zap.Error(err), zap.Int("topicID", id))
			}
		}()
		logger.Info("缓存专题信息", zap.Any("专题信息", result))
	}
	return result
}

// NewTopicService 创建专题服务实例
func NewTopicService() *TopicService {
	return &TopicService{
		repository: repository.NewTopicRepository(),
		cache:      NewTopicCache(),
	}
}

// TopicCache 专题缓存
type TopicCache struct {
	redis *redis.Client
}

// SetTopicToMap 将专题信息存入缓存
func (t *TopicCache) SetTopicToMap(id int, topic *response.SimpleTopicResponse) error {
	if topic == nil {
		return errors.New("专题信息为空")
	}
	return t.redis.HSet(common.TopicMapKey, strconv.Itoa(id), utils.Serialize(topic)).Err()
}

// ClearTopicKeys 清除专题相关的所有缓存
func (t *TopicCache) ClearTopicKeys() error {
	var keys = t.redis.Keys(common.TopicPageKey + "*").Val()
	return t.redis.Del(append(keys, common.TopicMapKey)...).Err()
}

// GetTopicFromMap 从缓存获取专题信息
func (t *TopicCache) GetTopicFromMap(topicID int) (*response.SimpleTopicResponse, error) {
	val, err := t.redis.HGet(common.TopicMapKey, strconv.Itoa(topicID)).Result()
	if err != nil {
		return nil, fmt.Errorf("从缓存获取专题信息失败: %w", err)
	}
	return utils.Deserialize[*response.SimpleTopicResponse](val), nil
}

// SetPageInfo 将博客列表存入redis
func (b *TopicCache) SetPageInfo(page int, pageInfo *response.Page) error {
	key := fmt.Sprintf(common.TopicPageKey+"%d", page)
	str := utils.Serialize(pageInfo)
	return b.redis.Set(key, str, common.PageInfoExpire).Err()
}

// GetPageInfo 从redis获取博客列表
func (b *TopicCache) GetPageInfo(page int) (*response.Page, error) {
	key := fmt.Sprintf(common.TopicPageKey+"%d", page)
	str := b.redis.Get(key).Val()
	if str == "" {
		return nil, errors.New("not found")
	}
	return utils.Deserialize[*response.Page](str), nil
}

// NewTopicCache 创建专题缓存实例
func NewTopicCache() *TopicCache {
	return &TopicCache{redis: configs.REDIS}
}
