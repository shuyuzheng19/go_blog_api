package response

// 博客详情页返回对象
type BlogContentResponse struct {
	ID          int64                   `json:"id"`          // 博客ID
	Title       string                  `json:"title"`       // 博客标题
	Description string                  `json:"description"` // 博客描述
	CoverImage  string                  `json:"coverImage"`  // 博客封面
	SourceURL   *string                 `json:"sourceUrl"`   // 博客原文链接
	Content     string                  `json:"content"`     // 博客正文
	EyeCount    int64                   `json:"eye_count"`   // 浏览量
	User        SimpleUserResponse      `json:"user"`        // 用户信息
	Category    *SimpleCategoryResponse `json:"category"`    // 分类信息
	Topic       *SimpleTopicResponse    `json:"topic"`       // 专题信息
	Tags        []SimpleTagResponse     `json:"tags"`        // 标签信息
	CreateTime  int64                   `json:"create_time"`
	UpdateTime  int64                   `json:"update_time"`
}

// BlogResponse 博客概要。通常是博客列表信息
type BlogResponse struct {
	ID          int64                   `json:"id"`                 //博客ID
	Title       string                  `json:"title"`              //博客标题
	Description string                  `json:"desc"`               //博客描述
	CoverImage  string                  `json:"coverImage"`         //博客封面图片
	CreatedAt   int64                   `json:"timestamp"`          //博客发布时间戳
	Category    *SimpleCategoryResponse `json:"category,omitempty"` //博客的分类概要
	User        SimpleUserResponse      `json:"user"`               // 博客用户概要
	Topic       *SimpleTopicResponse    `json:"topic,omitempty"`
	CategoryId  int                     `json:"-"` //分类ID
	TopicId     int                     `json:"-"`
	UserId      int                     `json:"-"` //用户ID
}

// SimpleBlogResponse 推荐博客
type SimpleBlogResponse struct {
	Id         int64  `json:"id"`         //博客ID
	Title      string `json:"title"`      //博客标题
	CoverImage string `json:"coverImage"` //博客封面
}

// ArchiveBlogResponse 归档博客
// @Description 归档博客概要
type ArchiveBlogResponse struct {
	Id          int64  `json:"id"`     //博客ID
	Title       string `json:"title"`  //博客标题
	Description string `json:"desc"`   //博客描述
	CreateTime  int64  `json:"create"` //博客创建日期
}

type SearchBlogResponse struct {
	Id          int64  `json:"id"`          //博客ID
	Title       string `json:"title"`       //博客标题
	Description string `json:"description"` //博客描述
}

// AdminBlogResponse 后台管理博客列表
// @Description 后台管理博客列表
type AdminBlogResponse struct {
	Id          int64                   `json:"id"`          //博客id
	Title       string                  `json:"title"`       //博客标题
	Description string                  `json:"description"` //博客描述
	CoverImage  string                  `json:"coverImage"`  //博客封面
	EyeCount    int64                   `json:"eyeCount"`    //博客浏览量
	Category    *SimpleCategoryResponse `json:"category"`    //博客分类
	Topic       *SimpleTopicResponse    `json:"topic"`       //博客专题
	CreatedAt   int64                   `json:"createAt"`    //创建时间
	UpdatedAt   int64                   `json:"updateAt"`    //修改时间
	User        SimpleUserResponse      `json:"user"`        //博客用户信息
	SourceURL   *string                 `json:"sourceUrl"`
	Pinned      bool                    `json:"pinned"`
	Order       *int64                  `json:"order"`
	UserId      int                     `json:"-"`
	CategoryId  int                     `json:"-"`
	TopicId     int                     `json:"-"`
}

type HomeData struct {
	Pinned []BlogResponse `json:"pinned"`
	Page   Page           `json:"page"`
}
