package models

import "blog/internal/dto/response"

// Blog 博客模型
type Blog struct {
	Model
	ID          int64     `gorm:"primary_key;comment:博客ID"`
	Description string    `gorm:"size:255;not null;comment:博客描述"`
	Title       string    `gorm:"size:255;not null;comment:博客标题"`
	CoverImage  string    `gorm:"not null;comment:博客封面"`
	SourceURL   *string   `gorm:"default:null;comment:博客原文链接"`
	Content     string    `gorm:"type:text;comment:博客正文"`
	EyeCount    int64     `gorm:"default:0;comment:浏览量"`
	CategoryID  *int      `gorm:"column:category_id;type:integer;comment:博客分类ID"`
	UserID      int       `gorm:"column:user_id;type:integer;comment:创建的用户ID"`
	TopicID     *int      `gorm:"column:topic_id;type:integer;comment:博客专题ID"`
	Tags        []Tag     `gorm:"many2many:blogs_tags"`
	Category    *Category `gorm:"foreignKey:CategoryID"`
	User        User      `gorm:"foreignKey:UserID"`
	Topic       *Topic    `gorm:"foreignKey:TopicID"`
	IsPrivate   bool      `gorm:"comment:是否为私有博客"`
	Pinned      bool      `gorm:"comment:是否为置顶博客"`
	Order       *int64    `gorm:"default:null;comment:置顶博客排序，只有开启置顶的时候才有用"`
	Password    *string   `gorm:"size:255;comment:访问私有博客的密码"` // 使用指针以便于处理空值
}

func (*Blog) TableName() string {
	return BlogTable
}

func (b *Blog) ToBlogContentResponse() response.BlogContentResponse {
	var category *response.SimpleCategoryResponse
	if b.Category != nil {
		category = &response.SimpleCategoryResponse{
			ID:   b.Category.ID,
			Name: b.Category.Name,
		}
	}

	tags := make([]response.SimpleTagResponse, len(b.Tags))
	for i, tag := range b.Tags {
		tags[i] = response.SimpleTagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	var topic *response.SimpleTopicResponse
	if b.Topic != nil {
		topic = &response.SimpleTopicResponse{
			ID:   b.Topic.ID,
			Name: b.Topic.Name,
		}
	}

	return response.BlogContentResponse{
		ID:          b.ID,
		Title:       b.Title,
		Content:     b.Content,
		Description: b.Description,
		CoverImage:  b.CoverImage,
		EyeCount:    b.EyeCount,
		Category:    category,
		Topic:       topic,
		CreateTime:  b.CreatedAt,
		UpdateTime:  b.UpdatedAt,
		SourceURL:   b.SourceURL,
		User:        b.User.ToSimpleUser(),
		Tags:        tags,
	}
}

type EditBlog struct {
	Model
	UID     int    `gorm:"primaryKey;type:int;column:uid;comment:保存博客的用户ID"`
	Content string `gorm:"type:text;comment:保存博客正文"`
}

func (*EditBlog) TableName() string {
	return EditBlogTable
}
