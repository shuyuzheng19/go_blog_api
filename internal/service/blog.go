package service

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/job"
	"blog/internal/models"
	"blog/internal/repository"
	"blog/internal/search"
	"blog/internal/utils"
	"blog/pkg/common"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BlogService struct {
	repository *repository.BlogRepository
	cache      *BlogCache
	search     *search.MeiliSearchClient
	index      string
}

// CreateBlog 添加博客
func (b *BlogService) CreateBlog(uid int, request requests.BlogRequest) (*models.Blog, error) {
	blog := request.ToBlogModel(uid)

	if err := b.repository.CreateBlog(&blog); err != nil {
		logger.Info("创建博客失败", zap.String("err", err.Error()))
		return nil, err
	}

	go func() {
		b.updateCacheAndSearch(blog)
		b.cache.ClearBlogKeys()
	}()

	logger.Info("博客创建成功", zap.Int("user_id", uid), zap.String("title", blog.Title))
	return &blog, nil
}

// UpdateBlog 更新博客
func (b *BlogService) UpdateBlog(bid int64, uid int, super bool, request requests.BlogRequest) (*models.Blog, error) {
	blog := request.ToBlogModel(uid)
	blog.ID = bid

	if err := b.repository.UpdateBlog(uid, super, &blog); err != nil {
		logger.Info("更新博客失败", zap.String("err", err.Error()))
		return nil, err
	}

	go func() {
		b.updateCacheAndSearch(blog)
		b.cache.DeleteByIds([]int64{blog.ID})
	}()

	logger.Info("博客更新成功", zap.Int64("id", bid), zap.Int("user_id", uid))
	return &blog, nil
}

func (b *BlogService) SaveEditBlog(uid int, content string) error {
	return b.repository.SaveEditBlog(models.EditBlog{UID: uid, Content: content})
}

func (b *BlogService) GetSaveEditBlog(uid int) (string, error) {
	return b.repository.GetEditBlog(uid)
}

// updateCacheAndSearch 更新缓存和搜索索引
func (b *BlogService) updateCacheAndSearch(blog models.Blog) {
	b.cache.ClearBlogPageInfo()
	searchBlog := response.SearchBlogResponse{
		Id:          blog.ID,
		Title:       blog.Title,
		Description: blog.Description,
	}
	str := utils.Serialize([]response.SearchBlogResponse{searchBlog})
	if err := b.search.SaveDocument(b.index, str); err != nil {
		logger.Info("更新搜索索引失败", zap.String("err", err.Error()))
	}
}

// DeleteBlogByIDs 删除博客
func (b *BlogService) DeleteBlogByIDs(uid *int, ids []int64) error {
	if err := configs.DeleteData(models.BlogTable, uid, ids); err != nil {
		logger.Info("删除博客失败", zap.String("err", err.Error()))
		return err
	}
	go func() {
		b.cache.DeleteByIds(ids)
		b.cache.ClearBlogKeys()
	}()

	logger.Info("博客删除成功", zap.Int64s("ids", ids))
	return nil
}

// UnDeleteBlogByIDs 恢复删除的博客
func (b *BlogService) UnDeleteBlogByIDs(uid *int, ids []int64) error {
	if err := configs.UnDeleteData(models.BlogTable, uid, ids); err != nil {
		logger.Info("恢复博客失败", zap.String("err", err.Error()))
		return err
	}

	go b.cache.ClearBlogKeys()

	logger.Info("博客恢复成功", zap.Int64s("ids", ids))
	return nil
}

// GetBlogByID 根据 ID 获取博客
func (b *BlogService) GetBlogByID(id int64) (*models.Blog, error) {
	blog, err := b.cache.GetBlogInfo(id)
	if err == nil {
		return blog, nil
	}

	blog, err = b.repository.GetBlogById(id)
	if err != nil || blog == nil || blog.ID == 0 {
		logger.Info("获取博客失败", zap.Int64("id", id), zap.String("err", err.Error()))
		return nil, err
	}

	go b.cache.SetBlogInfo(id, blog)
	return blog, nil
}

// GetBlogEyeCount 获取博客浏览次数
func (b *BlogService) GetBlogEyeCount(count, id int64) int64 {
	var result = b.cache.GetBlogEyeCount(count, id)
	go b.UpdateDailyTotalPv()
	return result
}

// UpdateDailyTotalPv 更新当天的全局浏览量
func (b *BlogService) UpdateDailyTotalPv() {
	err := b.cache.IncrementDailyPv()
	if err != nil {
		// 记录日志以便排查问题
		logger.Info("Failed to update daily total PV:", zap.String("error", err.Error()))
	}
}

// SaveRecommend 保存推荐博客
func (b *BlogService) SaveRecommend(ids []int) error {
	if len(ids) != common.RecommendBlogCount {
		return errors.New("推荐博客数量不正确")
	}

	blogs, err := b.repository.FindByIdInSimpleBlog(ids)
	if err != nil {
		logger.Info("获取推荐博客失败", zap.String("err", err.Error()))
		return err
	}

	go b.cache.SetRecommend(blogs)
	logger.Info("保存推荐博客成功", zap.Ints("ids", ids))
	return nil
}

// GetRecommend 获取推荐博客
func (b *BlogService) GetRecommend() ([]response.SimpleBlogResponse, error) {
	blogs, err := b.cache.GetRecommend()
	if err != nil {
		logger.Info("获取推荐博客失败", zap.String("err", err.Error()))
		return nil, err
	}
	return blogs, nil
}

// GetHotBlogs 获取热门博客
func (b *BlogService) GetHotBlogs() ([]response.SimpleBlogResponse, error) {
	blogs, err := b.cache.GetHotBlog()
	if err != nil {
		blogs, err = b.repository.GetHotBlog()
		if err != nil {
			logger.Info("获取热门博客失败", zap.String("err", err.Error()))
			return nil, err
		}
		go b.cache.SetHotBlog(blogs)
	}
	return blogs, nil
}

// GetAdminBlogList 获取管理员博客列表
func (b *BlogService) GetAdminBlogList(uid *int, req requests.AdminFilterRequest, page *response.Page) error {
	list, err := b.repository.GetBlogAdminList(uid, req, &page.Count)
	if err != nil {
		logger.Info("获取管理员博客列表失败", zap.String("err", err.Error()))
		return err
	}
	page.Data = list
	return nil
}

// GetLatestBlogs 获取最新博客
func (b *BlogService) GetLatestBlogs() ([]response.SimpleBlogResponse, error) {
	blogs, err := b.cache.GetLatestBlog()
	if err != nil {
		blogs, err = b.repository.GetLatestBlog()
		if err != nil {
			logger.Info("获取最新博客失败", zap.String("err", err.Error()))
			return nil, err
		}
		go b.cache.SetLatestBlog(blogs)
	}
	return blogs, nil
}

// InitEyeCount 初始化博客浏览量
func (b *BlogService) InitEyeCount() {
	maps := b.cache.GetAllBlogEyeCount()
	if len(maps) == 0 {
		logger.Info("没有需要初始化的浏览量数据")
		return
	}
	var (
		wg            sync.WaitGroup
		concurrentSem = make(chan struct{}, 10) // 控制最大并发数，例如 10
	)
	for id, count := range maps {
		wg.Add(1)
		// 使用带缓冲的通道限制 Goroutine 并发数
		concurrentSem <- struct{}{}
		go func(id, count string) {
			defer wg.Done()
			defer func() { <-concurrentSem }() // 释放通道占用
			// 转换 ID 和 count，并处理可能的错误
			idNumber, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return
			}
			countNumber, err := strconv.ParseInt(count, 10, 64)
			if err != nil {
				return
			}
			// 更新数据库中的浏览量
			b.repository.UpdateEyeCount(idNumber, countNumber)
		}(id, count)
	}
	wg.Wait()
	logger.Info("初始化浏览量完成")
	var count, err = GetTodayTotalViews()
	if err == nil {
		err = AddEyeCount(count)
		if err == nil {
			b.cache.DeletePvViewCount()
		}

	}
	b.cache.DeleteBlogEyeCount()
}

// GetBlogList 获取博客列表
func (b *BlogService) GetBlogList(prequest requests.RequestQuery, page *response.Page) error {
	pageInfo, err := b.cache.GetPageInfo(prequest)

	if err == nil {
		*page = *pageInfo
		return nil
	}

	list, err := b.repository.GetBlogList(prequest, &page.Count)

	if err != nil {
		logger.Info("获取博客列表失败", zap.String("err", err.Error()))
		return err
	}

	page.Data = list
	go func() {
		b.cache.SetPageInfo(prequest, page)
	}()
	return nil
}

// GetArchiveBlog 获取归档博客
func (b *BlogService) GetArchiveBlog(request requests.ArchiveBlogRequest, page *response.Page) error {
	return b.repository.GetArchiveBlog(request, page)
}

func checkTimestamp(expiryTimestamp int64) (bool, time.Duration, error) {
	// 当前时间
	currentTime := time.Now()

	// 将前端传过来的秒级时间戳转换为 time.Time
	expiryTime := time.Unix(expiryTimestamp, 0)

	// 判断是否已过期
	if currentTime.After(expiryTime) {
		return true, 0, nil // 已过期，返回 true
	}

	// 计算距离失效时间的差值
	timeLeft := expiryTime.Sub(currentTime)
	return false, timeLeft, nil
}

func (b *BlogService) SetTempBlog(req requests.TmpBlog) (string, error) {
	var expire, time, err = checkTimestamp(req.Unix)

	if expire || err != nil {
		return "", fmt.Errorf("保存失败")
	}

	var uuid = uuid.NewString()

	return uuid, b.cache.SetTempBlog(uuid, time, req)
}

func (b *BlogService) GetTempBlog(id string) *requests.TmpBlog {
	blog, err := b.cache.GetTempBlog(id)
	if err != nil {
		return nil
	}
	return blog
}

func (b *BlogService) SetPinnedBlog(req requests.PinnedBlogRequest) error {
	var err = b.repository.UpdatePinned(req)

	if err != nil {
		return err
	}

	go b.cache.ClearPinnedKey()

	return nil
}

func (b *BlogService) GetPinnedBlog() []response.BlogResponse {
	var blogs, err = b.cache.GetPinnedBlog()

	if err != nil {
		blogs = b.repository.GetPinnedBlogList()
		b.cache.SetPinnedBlog(blogs)
	}

	return blogs
}

// InitSearch 初始化搜索索引
func (b *BlogService) InitSearch() error {
	if err := b.search.DeleteAllDocument(b.index); err != nil {
		logger.Info("初始化搜索索引失败", zap.String("err", err.Error()))
		return err
	}
	blogs, err := b.repository.FindAllSearchBlog()
	if err != nil {
		logger.Info("获取所有博客失败", zap.String("err", err.Error()))
		return err
	}
	jsonStr := utils.Serialize(blogs)
	return b.search.SaveDocument(b.index, jsonStr)
}

// SimilarBlog 获取相似博客
func (b *BlogService) SimilarBlog(keyword string) ([]any, error) {
	req := getBlogSearchRequest(requests.SearchBlogRequest{Page: 1, Keyword: keyword})
	response := b.search.SearchDocument(b.index, req)
	return response.Hits, nil
}

// SearchBlog 搜索博客
func (b *BlogService) SearchBlog(req requests.SearchBlogRequest) response.Page {
	result := b.search.SearchDocument(b.index, getBlogSearchRequest(req))
	return response.Page{
		Page:  req.Page,
		Count: result.EstimatedTotalHits,
		Size:  common.SearchBlogPageCount,
		Data:  result.Hits,
	}
}

func (b *BlogService) initHotBlog() {
	b.cache.ClearHotBlog()
	b.cache.GetHotBlog()
}

func (b *BlogService) initSearch() {
	b.InitSearch()
}

// getBlogSearchRequest 构建博客搜索请求
func getBlogSearchRequest(req requests.SearchBlogRequest) search.MeiliSearchRequest {
	return search.MeiliSearchRequest{
		Q:                     req.Keyword,
		Offset:                (req.Page - 1) * common.SearchBlogPageCount,
		Limit:                 common.SearchBlogPageCount,
		AttributesToHighlight: []string{"*"},
		ShowMatchesPosition:   false,
		HighlightPreTag:       "<b>",
		HighlightPostTag:      "</b>",
	}
}

// NewBlogService 创建新的 BlogService 实例
func NewBlogService() *BlogService {
	var service = &BlogService{
		repository: repository.NewBlogRepository(),
		cache:      NewBlogCache(),
		search:     configs.SEARCH,
		index:      configs.CONFIG.Search.BlogIndex,
	}

	if configs.CONFIG.Server.Cron {
		job.AddJob(job.Job{
			Hour:        0,
			Eq:          true,
			Description: "初始化浏览量",
			Job:         service.InitEyeCount,
		})

		job.AddJob(job.Job{
			Hour:        1,
			Eq:          true,
			Description: "初始化搜索",
			Job:         service.initSearch,
		})

		job.AddJob(job.Job{
			Hour:        6,
			Eq:          false,
			Description: "更新热门博客",
			Job:         service.initHotBlog,
		})
	}

	return service
}

type BlogCache struct {
	redis *redis.Client
}

// SetBlogInfo 缓存博客详情信息
func (b *BlogCache) SetBlogInfo(id int64, blog *models.Blog) error {
	var str string
	if blog != nil {
		str = utils.Serialize(blog)
	}
	return b.redis.HSet(common.BlogMapKey, strconv.FormatInt(id, 10), str).Err()
}

// DeleteByIds 删除缓存中的博客
func (b *BlogCache) DeleteByIds(ids []int64) error {
	fields := make([]string, len(ids))
	for i, id := range ids {
		fields[i] = strconv.FormatInt(id, 10)
	}
	return b.redis.HDel(common.BlogMapKey, fields...).Err()
}

// GetAllBlogEyeCount 获取所有的浏览量
func (b *BlogCache) GetAllBlogEyeCount() map[string]string {
	return b.redis.HGetAll(common.BlogEyeCountMapKey).Val()
}

// GetBlogInfo 获取博客详情信息
func (b *BlogCache) GetBlogInfo(id int64) (*models.Blog, error) {
	r := b.redis.HGet(common.BlogMapKey, strconv.FormatInt(id, 10))
	if r.Err() != nil {
		return nil, errors.New("获取博客信息失败")
	}
	return utils.Deserialize[*models.Blog](r.Val()), nil
}

// SetRecommend 缓存推荐博客
func (b *BlogCache) SetRecommend(blogs []response.SimpleBlogResponse) error {
	blogsJson := utils.Serialize(blogs)
	return b.redis.Set(common.RecommendKey, blogsJson, -1).Err()
}

// DeleteBlogEyeCount 删除浏览量
func (b *BlogCache) DeleteBlogEyeCount() error {
	return b.redis.Del(common.BlogEyeCountMapKey).Err()
}

func (b *BlogCache) DeletePvViewCount() error {
	return b.redis.Del(common.EyeView).Err()
}

// GetRecommend 从缓存获取推荐博客
func (b *BlogCache) GetRecommend() ([]response.SimpleBlogResponse, error) {
	str := b.redis.Get(common.RecommendKey).Val()
	if str == "" {
		return nil, errors.New("推荐博客未找到")
	}
	return utils.Deserialize[[]response.SimpleBlogResponse](str), nil
}

// getSimpleBlog 获取简单博客列表缓存
func (b *BlogCache) getSimpleBlog(key string) ([]response.SimpleBlogResponse, error) {
	str := b.redis.Get(key)
	if str.Err() != nil {
		return nil, errors.New("获取失败")
	}
	return utils.Deserialize[[]response.SimpleBlogResponse](str.Val()), nil
}

// SetLatestBlog 缓存最新的10条博客
func (b *BlogCache) SetLatestBlog(blogs []response.SimpleBlogResponse) error {
	return b.setSimpleBlog(common.LatestBlogKey, common.LatestBlogExpire, blogs)
}

// GetLatestBlog 从缓存获取最新博客
func (b *BlogCache) GetLatestBlog() ([]response.SimpleBlogResponse, error) {
	return b.getSimpleBlog(common.LatestBlogKey)
}

func (b *BlogCache) SetPinnedBlog(blogs []response.BlogResponse) error {
	blogsJson := utils.Serialize(blogs)
	return b.redis.Set(common.PinnedBlog, blogsJson, -1).Err()
}

func (b *BlogCache) GetPinnedBlog() ([]response.BlogResponse, error) {
	var jsonStr, err = b.redis.Get(common.PinnedBlog).Result()
	if err != nil {
		return nil, err
	}
	return utils.Deserialize[[]response.BlogResponse](jsonStr), nil
}

// GetHotBlog 从缓存获取热门博客
func (b *BlogCache) GetHotBlog() ([]response.SimpleBlogResponse, error) {
	return b.getSimpleBlog(common.HotBlogKey)
}

func (b *BlogCache) ClearHotBlog() error {
	return b.redis.Del(common.HotBlogKey).Err()
}

// SetHotBlog 缓存热门的10条博客
func (b *BlogCache) SetHotBlog(blogs []response.SimpleBlogResponse) error {
	return b.setSimpleBlog(common.HotBlogKey, common.HotBlogExpire, blogs)
}

// setSimpleBlog 设置简单博客列表缓存
func (b *BlogCache) setSimpleBlog(key string, expire time.Duration, blogs []response.SimpleBlogResponse) error {
	blogsJson := utils.Serialize(blogs)
	return b.redis.Set(key, blogsJson, expire).Err()
}

// ClearBlogPageInfo 清空所有博客页面缓存
func (b *BlogCache) ClearBlogPageInfo() error {
	keys := b.redis.Keys(common.PageInfoPrefixKey + "*").Val()
	return b.redis.Del(keys...).Err()
}

// SetPageInfo 将博客列表存入redis
func (b *BlogCache) SetPageInfo(req requests.RequestQuery, pageInfo *response.Page) error {
	key := fmt.Sprintf("%s:page_%d_cid_%d_sort:%s", common.PageInfoPrefixKey, req.Page, req.Cid, req.Sort)
	str := utils.Serialize(pageInfo)
	return b.redis.Set(key, str, common.PageInfoExpire).Err()
}

// GetPageInfo 从redis获取博客列表
func (b *BlogCache) GetPageInfo(req requests.RequestQuery) (*response.Page, error) {
	key := fmt.Sprintf("%s:page_%d_cid_%d_sort:%s", common.PageInfoPrefixKey, req.Page, req.Cid, req.Sort)
	str := b.redis.Get(key).Val()
	if str == "" {
		return nil, errors.New("not found")
	}
	return utils.Deserialize[*response.Page](str), nil
}

func (b *BlogCache) ClearPinnedKey() error {
	return b.redis.Del(common.PinnedBlog).Err()
}

func (b *BlogCache) ClearBlogKeys() error {
	var pages = b.redis.Keys(common.PageInfoPrefixKey + "*").Val()
	return b.redis.Del(append(pages, common.BlogMapKey, common.HotBlogKey, common.LatestBlogKey)...).Err()
}

func (b *BlogCache) SetTempBlog(id string, time time.Duration, blog requests.TmpBlog) error {
	str := utils.Serialize(blog)
	return b.redis.SetNX(common.TmpBlogKey+id, str, time).Err()
}

func (b *BlogCache) GetTempBlog(id string) (*requests.TmpBlog, error) {
	str, err := b.redis.Get(common.TmpBlogKey + id).Result()
	if err != nil {
		return nil, err
	}
	var blog = utils.Deserialize[requests.TmpBlog](str)
	return &blog, nil
}

// IncrementDailyPv 更新当天的全局浏览量
func (c *BlogCache) IncrementDailyPv() error {
	// Redis 原子操作增加计数
	return c.redis.Incr(common.EyeView).Err()
}

// GetBlogEyeCount 获取博客浏览次数
func (b *BlogCache) GetBlogEyeCount(defaultCount int64, id int64) int64 {
	field := strconv.FormatInt(id, 10)
	luaScript := `
		local count = redis.call("HINCRBY", KEYS[1], ARGV[1], 1)
		if count == 1 then
			redis.call("HSET", KEYS[1], ARGV[1], ARGV[2])
			return ARGV[2]
		end
		return count
	`
	count, err := b.redis.Eval(luaScript, []string{common.BlogEyeCountMapKey}, field, defaultCount+1).Int64()
	if err != nil {
		logger.Info("获取博客浏览次数失败", zap.String("err", err.Error()))
		return defaultCount
	}
	return count
}

// NewBlogCache 创建新的 BlogCache 实例
func NewBlogCache() *BlogCache {
	return &BlogCache{redis: configs.REDIS}
}
