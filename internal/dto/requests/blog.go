package requests

import (
	"blog/internal/dto/response"
	"blog/internal/models"
)

// BlogRequest 用于添加博客的请求结构体
type BlogRequest struct {
	Description string  `json:"description" validate:"required,max=255" error:"博客描述为必填项，且长度不能超过255个字符"`
	Title       string  `json:"title" validate:"required,max=255" error:"博客标题为必填项，且长度不能超过255个字符"`
	CoverImage  string  `json:"coverImage" validate:"required" error:"博客封面为必填项"`
	SourceURL   *string `json:"source_url" validate:"omitempty,url" error:"博客原文链接格式不正确"`
	Content     string  `json:"content" validate:"required" error:"博客正文为必填项"`
	IsPrivate   bool    `json:"isPrivate" error:"是否为私有博客"`
	Password    *string `json:"password" validate:"omitempty,min=6" error:"私有博客密码长度至少为6个字符"` // 可选，且长度至少为6个字符
	CategoryID  *int    `json:"category" validate:"omitempty" error:"博客分类ID必须为有效值"`
	TopicID     *int    `json:"topic" validate:"omitempty" error:"博客专题ID必须为有效值"`
	Tags        []int   `json:"tags" validate:"omitempty" error:"标签ID列表"` // 标签ID列表
}

// ToBlogDo 将请求模型转为数据库模型
func (b BlogRequest) ToBlogModel(uid int) models.Blog {
	var tags []models.Tag
	if b.CategoryID != nil {
		for _, tag := range b.Tags {
			tags = append(tags, models.Tag{
				ID: tag,
			})
		}
		b.TopicID = nil
	}
	if b.TopicID != nil {
		b.CategoryID = nil
	}
	return models.Blog{
		Description: b.Description,
		Title:       b.Title,
		CoverImage:  b.CoverImage,
		SourceURL:   b.SourceURL,
		Password:    b.Password,
		IsPrivate:   b.IsPrivate,
		Content:     b.Content,
		CategoryID:  b.CategoryID,
		UserID:      uid,
		TopicID:     b.TopicID,
		Tags:        tags,
	}
}

// GetOrderString 博客列表排序方式
func (sort Sort) GetBlogOrderString(prefix string) string {
	switch sort {
	case CREATE:
		return prefix + "created_at desc"
	case UPDATE:
		return prefix + "updated_at  desc"
	case EYE:
		return prefix + "eye_count desc"
	case BACK:
		return prefix + "created_at asc"
	default:
		return prefix + "created_at desc"
	}
}

type ArchiveBlogRequest struct {
	Start int64 `form:"start"` //开始日期的时间戳
	End   int64 `form:"end"`   //结束日期的时间戳
	Page  int   `form:"page"`  //第几页
}

// SearchBlogRequest 搜索博客
// @Description 搜索博客
type SearchBlogRequest struct {
	Keyword string `form:"keyword"` //搜索关键字
	Page    int    `form:"page"`    //第几页
}

type TmpBlog struct {
	Title   string                       `json:"title"  validate:"required" error:"标题为必填项"`
	Desc    string                       `json:"desc" validate:"required" error:"描述为必填项"`
	Unix    int64                        `json:"unix" validate:"required" error:"失效日期为必填项"`
	Content string                       `json:"content" validate:"required" error:"内容为必填项"`
	Create  int64                        `json:"create"`
	User    *response.SimpleUserResponse `json:"user"`
}

type PinnedBlogRequest struct {
	Id     int64  `json:"id"`
	Order  *int64 `json:"order"`
	Pinned bool   `json:"pinned"`
}
