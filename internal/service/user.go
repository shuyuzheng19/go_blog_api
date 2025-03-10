package service

import (
	"blog/internal/dto/dtos"
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/internal/repository"
	"blog/internal/utils"
	"blog/pkg/common"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"blog/pkg/smail"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// UserService 用户服务
type UserService struct {
	dao   *repository.UserRepository
	cache *UserCache
}

// SendCodeToEmail 发送注册验证码
func (u *UserService) SendCodeToEmail(email string) error {
	code := utils.RandomNumberCode()

	if err := smail.SendEmail(email, "注册验证码", false, code); err != nil {
		logger.Info("发送邮箱验证码失败", zap.String("error", err.Error()), zap.String("email", email))
		return fmt.Errorf("发送验证码失败: %w", err)
	}

	if err := u.cache.SetEmailCode(code, email); err != nil {
		logger.Info("缓存邮箱验证码失败", zap.String("error", err.Error()), zap.String("email", email))
		return fmt.Errorf("保存验证码失败: %w", err)
	}

	logger.Info("发送邮箱验证码成功", zap.String("email", email))
	return nil
}

// GetWebSiteConfig 获取网站配置
func (u *UserService) GetWebSiteConfig() response.BlogConfigInfo {
	config := u.cache.GetWebSiteConfig()
	return config
}

// Contact 用户反馈
func (u *UserService) Contact(req requests.ContactRequest) error {
	text := fmt.Sprintf("<h3>%s</h3><p>对方名字: %s</p><p>对方邮箱: %s</p>留言内容:<p>%s</p>",
		req.Subject, req.Name, req.Email, req.Content)

	if err := smail.SendEmail(configs.CONFIG.MyEmail, req.Subject, true, text); err != nil {
		logger.Info("发送反馈邮件失败", zap.String("error", err.Error()), zap.Any("request", req))
		return fmt.Errorf("发送反馈邮件失败: %w", err)
	}

	logger.Info("用户反馈成功", zap.Any("request", req))
	return nil
}

func (u *UserService) SetWebSiteConfig(c response.BlogConfigInfo) error {
	logger.Info("更新网站配置")
	return u.cache.SetWebSiteConfig(c)
}

// GetAdminUserList 获取管理员用户列表
func (u *UserService) GetAdminUserList(req requests.UserAdminFilter, page *response.Page) error {
	if err := u.dao.GetAdminUsers(req, page); err != nil {
		logger.Info("获取管理员用户列表失败", zap.String("error", err.Error()), zap.Any("request", req))
		return fmt.Errorf("获取管理员用户列表失败: %w", err)
	}
	return nil
}

// GetRedisKeys 获取所有Redis里的key
func (u *UserService) GetRedisKeys() ([]string, error) {
	keys := u.cache.GetAllKey()
	if len(keys) == 0 {
		return nil, errors.New("未找到任何键")
	}
	return keys, nil
}

// DeleteRedisKeys 删除指定的Redis键
func (u *UserService) DeleteRedisKeys(keys []string) error {
	if err := u.cache.DeleteKeys(keys); err != nil {
		logger.Info("删除Redis键失败", zap.String("error", err.Error()), zap.Strings("keys", keys))
		return fmt.Errorf("删除Redis键失败: %w", err)
	}
	return nil
}

// DeleteMatchKeys 删除匹配的Redis键
func (u *UserService) DeleteMatchKeys(key string) error {
	if err := u.cache.MatchDelete(key); err != nil {
		logger.Info("删除匹配的Redis键失败", zap.String("error", err.Error()), zap.String("key", key))
		return fmt.Errorf("删除匹配的Redis键失败: %w", err)
	}
	return nil
}

// RegisteredUser 用户注册
func (u *UserService) RegisteredUser(ip string, req requests.UserRequest) (models.User, error) {
	if common.EnableEmailCode {
		cacheCode := u.cache.GetEmailCode(req.Email)
		if cacheCode == "" || cacheCode != req.Code {
			return models.User{}, errors.New("邮箱验证码错误")
		}
	}

	user := req.ToUserModel(ip)

	if err := u.dao.Save(&user); err != nil {
		logger.Info("用户注册失败", zap.String("error", err.Error()), zap.Any("user", user))
		return models.User{}, fmt.Errorf("用户注册失败: %w", err)
	}

	user.RoleID = uint(common.UserRoleId)

	logger.Info("用户注册成功", zap.Int("UserID", user.ID), zap.String("username", user.Username))
	return user, nil
}

// UpdateUserStatus 更新用户登录状态
func (u *UserService) UpdateUserStatus(dto dtos.UserLoginStatus) error {
	if err := u.dao.UpdateLoginStatus(dto); err != nil {
		logger.Info("更新用户状态失败", zap.String("error", err.Error()), zap.Int("UserID", dto.ID))
		return fmt.Errorf("更新登录状态失败: %w", err)
	}
	logger.Info("更新用户状态成功", zap.Int("UserID", dto.ID))
	return nil
}

func (u *UserService) ResetPassword(id int, password string) error {

	var hashPassword = utils.EncryptPassword(password)

	var err = u.dao.UpdatePassword(id, hashPassword)

	if err == nil {
		go u.cache.RemoveToken(id)
	}

	return err
}

func (u *UserService) UpdateRoleID(uid int, rid uint) error {
	err := u.dao.UpdateUserRole(uid, rid)
	if err != nil {
		logger.Info("更新用户角色失败", zap.String("error", err.Error()), zap.Int("UserID", uid))
		return err
	}

	go u.cache.ClearUserInfoByID(uid)

	logger.Info("更新用户角色成功", zap.Int("UserID", uid), zap.Uint("RoleID", rid))
	return nil
}

// Login 用户登录
func (u *UserService) Login(request requests.LoginRequest) (response.TokenResponse, error) {
	user, err := u.dao.FindByUsername(request.Username)
	if err != nil {
		logger.Info("用户登录失败", zap.String("username", request.Username), zap.String("error", err.Error()))
		return response.TokenResponse{}, errors.New("用户未找到")
	}

	if !utils.VerifyPassword(user.Password, request.Password) {
		return response.TokenResponse{}, errors.New("密码错误")
	}

	if !user.Status {
		return response.TokenResponse{}, errors.New("用户被禁用")
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return response.TokenResponse{}, errors.New("生成令牌失败")
	}

	if err := u.cache.SetToken(user.ID, token.Token); err != nil {
		logger.Info("缓存令牌失败", zap.String("error", err.Error()), zap.Int("UserID", user.ID))
		return response.TokenResponse{}, errors.New("缓存令牌失败")
	}

	logger.Info("用户登录成功", zap.String("username", user.Username))
	return response.TokenResponse{Token: token.Token, User: user.ToVo()}, nil
}

// GetUsetBlogByPage 获取用户博客分页
func (u *UserService) GetUsetBlogByPage(prequest requests.RequestQuery, page *response.Page) error {
	list, err := u.dao.GetUserBlogList(prequest, &page.Count)
	if err != nil {
		return fmt.Errorf("获取博客列表失败: %w", err)
	}

	page.Data = list
	return nil
}

// GetUsetBlogTop10 获取用户前10篇博客
func (u *UserService) GetUsetBlogTop10(uid int) ([]response.SimpleBlogResponse, error) {
	blogs, err := u.dao.GetUserBlogTop10(uid)
	if err != nil {
		return nil, fmt.Errorf("获取Top10博客失败: %w", err)
	}
	return blogs, nil
}

// GetUsetTopic 获取用户的主题列表
func (u *UserService) GetUsetTopic(uid int) ([]response.TopicResponse, error) {
	topics, err := u.dao.GetTopicByUserId(uid)
	if err != nil {
		return nil, fmt.Errorf("获取用户主题失败: %w", err)
	}
	return topics, nil
}

// GetUser 获取用户信息
func (u *UserService) GetUser(id int) *models.User {
	user := u.cache.GetUser(id)
	if user == nil {
		dbUser, err := u.dao.FindById(id)
		if err != nil {
			logger.Info("获取用户信息失败", zap.Int("UserID", id), zap.String("error", err.Error()))
			return nil
		}
		u.cache.SetUser(id, &dbUser)
		return &dbUser
	}
	return user
}

// GetToken 获取用户令牌
func (u *UserService) GetToken(id int) string {
	return u.cache.GetToken(id)
}

// Logout 用户登出
func (u *UserService) Logout(uid int) error {
	if err := u.cache.RemoveToken(uid); err != nil {
		logger.Info("用户登出失败", zap.Int("UserID", uid), zap.String("error", err.Error()))
		return fmt.Errorf("删除令牌失败: %w", err)
	}
	logger.Info("用户登出成功", zap.Int("UserID", uid))
	return nil
}

// UpdateUser 修改用户信息
func (u *UserService) UpdateUser(userRequest *requests.UpdateUserRequest) (*models.User, error) {
	userModel := userRequest.ToModel()

	if err := u.dao.UpdateUser(&userModel, common.RoleId(userRequest.CurrentRoleId)); err != nil {
		logger.Info("更新用户信息失败", zap.Int("UserID", userModel.ID), zap.String("error", err.Error()))
		return nil, fmt.Errorf("更新用户信息失败: %w", err)
	}

	go u.cache.ClearUserInfoByID(userModel.ID)

	logger.Info("更新用户信息成功", zap.Int("UserID", userModel.ID))
	return &userModel, nil
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	service := &UserService{
		dao:   repository.NewUserRepository(),
		cache: NewUserCache(),
	}

	common.GetJwtUser = service.GetUser
	common.GetToken = service.GetToken

	return service
}

// UserCache 用户缓存
type UserCache struct {
	redis *redis.Client
}

// SetEmailCode 缓存邮箱验证码
func (u *UserCache) SetEmailCode(code, email string) error {
	key := fmt.Sprintf(common.EmailCodeKey+"%s", email)
	return u.redis.Set(key, code, common.EmailCodeKeyExpire).Err()
}

// GetEmailCode 获取邮箱验证码
func (u *UserCache) GetEmailCode(email string) string {
	key := fmt.Sprintf(common.EmailCodeKey+"%s", email)
	return u.redis.Get(key).Val()
}

func (u *UserCache) ClearUserInfoByID(id int) error {
	return u.redis.Del(fmt.Sprintf(common.UserInfoKey+"%d", id)).Err()
}

// GetWebSiteConfig 获取网站配置
func (u *UserCache) GetWebSiteConfig() response.BlogConfigInfo {
	str := u.redis.Get(common.WebSiteConfigKey).Val()
	if str == "" {
		return response.GetDefaultBlogConfigInfo()
	}
	return utils.Deserialize[response.BlogConfigInfo](str)
}

// SetWebSiteConfig 缓存网站信息
func (u *UserCache) SetWebSiteConfig(siteConfig response.BlogConfigInfo) error {
	str := utils.Serialize(siteConfig)
	return u.redis.Set(common.WebSiteConfigKey, str, -1).Err()
}

// SetUser 缓存用户信息
func (u *UserCache) SetUser(id int, user *models.User) error {
	if user == nil {
		return errors.New("用户信息为空")
	}

	key := common.UserInfoKey + strconv.Itoa(id)
	return u.redis.Set(key, utils.Serialize(user), common.UserInfoKeyExpire).Err()
}

// GetUser 获取用户缓存
func (u *UserCache) GetUser(id int) *models.User {
	key := common.UserInfoKey + strconv.Itoa(id)
	val := u.redis.Get(key).Val()
	if val == "" {
		return nil
	}
	return utils.Deserialize[*models.User](val)
}

// SetToken 缓存用户token
func (u *UserCache) SetToken(id int, token string) error {
	if token == "" {
		return errors.New("token不能为空")
	}

	key := fmt.Sprintf(common.UserTokenKey+"%d", id)
	return u.redis.Set(key, token, common.TokenExpire).Err()
}

// GetToken 获取用户token
func (u *UserCache) GetToken(id int) string {
	key := fmt.Sprintf(common.UserTokenKey+"%d", id)
	return u.redis.Get(key).Val()
}

// RemoveToken 删除用户token
func (u *UserCache) RemoveToken(uid int) error {
	var uidStr = strconv.Itoa(uid)
	key := common.UserTokenKey + uidStr
	key2 := common.UserInfoKey + uidStr
	return u.redis.Del(key, key2).Err()
}

// GetAllKey 获取所有Redis键
func (u *UserCache) GetAllKey() []string {
	return u.redis.Keys("*").Val()
}

// DeleteKeys 删除Redis键
func (u *UserCache) DeleteKeys(keys []string) error {
	if len(keys) == 0 {
		return errors.New("键数组为空")
	}

	return u.redis.Del(keys...).Err()
}

// MatchDelete 删除匹配的Redis键
func (u *UserCache) MatchDelete(key string) error {
	arrays := u.redis.Keys(key).Val()
	return u.DeleteKeys(arrays)
}

// NewUserCache 创建用户缓存实例
func NewUserCache() *UserCache {
	return &UserCache{redis: configs.REDIS}
}
