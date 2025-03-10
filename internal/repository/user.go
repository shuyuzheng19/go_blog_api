package repository

import (
	"blog/internal/dto/dtos"
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/pkg/common"
	"blog/pkg/configs"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db    *gorm.DB
	table string
}

// Save 保存用户信息
func (u *UserRepository) Save(user *models.User) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}
		return nil
	})
}

// UpdateLoginStatus 更新用户登录状态
func (u *UserRepository) UpdateLoginStatus(dto dtos.UserLoginStatus) error {
	result := u.db.Model(&models.User{}).
		Where("id = ?", dto.ID).
		Updates(&dto)

	if result.Error != nil {
		return fmt.Errorf("更新登录状态失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}

	return nil
}

func (u *UserRepository) UpdatePassword(id int, password string) error {
	return u.db.Model(&models.User{}).Where("id = ?", id).Update("password", password).Error
}

// FindByUsername 根据用户名查找用户
func (u *UserRepository) FindByUsername(username string) (models.User, error) {
	var user models.User
	err := u.db.Model(&models.User{}).Preload("Role").
		First(&user, "username = ?", username).Error
	if err != nil {
		return user, fmt.Errorf("查找用户失败: %w", err)
	}
	return user, nil
}

// FindById 根据用户ID查找用户
func (u *UserRepository) FindById(id int) (models.User, error) {
	var user models.User
	err := u.db.Model(&models.User{}).Preload("Role").
		First(&user, "id = ?", id).Error
	if err != nil {
		return user, fmt.Errorf("查找用户失败: %w", err)
	}
	return user, nil
}

// GetAdminUsers 获取管理员用户列表
func (u *UserRepository) GetAdminUsers(req requests.UserAdminFilter, page *response.Page) error {
	query := u.db.Model(&models.User{}).Preload("Role")

	// 应用过滤条件
	if req.Start != "" && req.End != "" {
		query = query.Where("created_at BETWEEN ? AND ?", req.Start, req.End)
	}

	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		query = query.Where("username LIKE ? OR nick_name LIKE ?", like, like)
	}

	if req.Pub != nil {
		query = query.Where("status = ?", *req.Pub)
	}

	// 获取总数
	if err := query.Count(&page.Count).Error; err != nil {
		return fmt.Errorf("计算用户数量失败: %w", err)
	}

	if page.Count == 0 {
		page.Data = make([]models.User, 0)
		return nil
	}

	// 查询列表
	var users = make([]models.User, 0)
	err := query.Offset((req.Page - 1) * page.Size).
		Limit(page.Size).
		Order(req.Sort.GetUserOrderString("")).
		Find(&users).Error

	if err != nil {
		return fmt.Errorf("查询用户列表失败: %w", err)
	}

	page.Data = users
	return nil
}

// UpdateUser 更新用户信息
func (u *UserRepository) UpdateUser(user *models.User, roleId common.RoleId) error {
	return u.db.Transaction(func(tx *gorm.DB) error {

		query := tx.Model(user).Where("id = ?", user.ID)

		result := query.Update("nick_name", user.NickName).
			Update("avatar", user.Avatar).
			Update("username", user.Username).
			Update("email", user.Email).
			Update("status", user.Status).
			Update("role_id", user.RoleID)

		if user.Password != "" {
			result.Update("password", user.Password)
		}

		if result.Error != nil {
			return fmt.Errorf("更新用户信息失败: %w", result.Error)
		}

		return nil
	})
}

func (u *UserRepository) UpdateUserRole(uid int, roleId uint) error {

	return u.db.Transaction(func(tx *gorm.DB) error {

		query := tx.Model(&models.User{}).Where("id = ?", uid)

		result := query.Update("role_id", roleId)

		if result.Error != nil {
			return fmt.Errorf("更新用户角色失败: %w", result.Error)
		}

		return nil
	})
}

// GetUserBlogList 获取用户博客列表
func (u *UserRepository) GetUserBlogList(req requests.RequestQuery, count *int64) ([]response.BlogResponse, error) {
	var list = make([]response.BlogResponse, 0)
	query := u.db.Model(&models.Blog{}).Table(models.BlogTable+" b").
		Where("b.user_id = ? and category_id is not null", req.Uid)

	// 计算总数
	if err := query.Count(count).Error; err != nil {
		return nil, fmt.Errorf("计算博客数量失败: %w", err)
	}

	if *count == 0 {
		return list, nil
	}

	// 查询列表
	err := query.Joins(fmt.Sprintf("INNER JOIN %s u ON u.id = b.user_id", models.UserTable)).
		Joins(fmt.Sprintf("LEFT JOIN %s c ON c.id = b.category_id", models.CategoryTable)).
		Select([]string{
			"b.id", "b.title", "b.description", "b.cover_image", "b.created_at",
			`u.id AS "User__id"`, `u.nick_name AS "User__nick_name"`,
			`c.id AS "Category__id"`, `c.name AS "Category__name"`,
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

// GetUserBlogTop10 获取用户前10篇博客
func (u *UserRepository) GetUserBlogTop10(uid int) ([]response.SimpleBlogResponse, error) {
	var list = make([]response.SimpleBlogResponse, 0)
	err := u.db.Model(&models.Blog{}).Table(models.BlogTable).
		Select("id, title").
		Where("user_id = ?", uid).
		Limit(10).
		Find(&list).Error

	if err != nil {
		return nil, fmt.Errorf("查询用户Top10博客失败: %w", err)
	}
	return list, nil
}

// GetTopicByUserId 获取用户的主题列表
func (u *UserRepository) GetTopicByUserId(uid int) ([]response.TopicResponse, error) {
	var topics = make([]response.TopicResponse, 0)
	err := u.db.Model(&models.Topic{}).Table(models.TopicTable+" t").
		Where("t.user_id = ?", uid).
		Select("t.id, t.name, t.cover_image, t.created_at, t.description").
		Find(&topics).Error

	if err != nil {
		return nil, fmt.Errorf("查询用户主题列表失败: %w", err)
	}
	return topics, nil
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db:    configs.DB,
		table: models.UserTable,
	}
}
