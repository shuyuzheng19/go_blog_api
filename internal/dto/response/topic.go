package response

type TopicResponse struct {
	ID          int                `json:"id"`          //专题id
	Name        string             `json:"name"`        //专题名
	Description string             `json:"description"` //专题描述
	CoverImage  string             `json:"cover"`       //专题封面
	CreatedAt   int64              `json:"timestamp"`   //专题创建时间戳
	User        SimpleUserResponse `json:"user"`
	UserId      int                `json:"-"`
}

type SimpleTopicResponse struct {
	ID   int    `json:"id"`   //专题id
	Name string `json:"name"` //专题名子
}

// AdminTopicResponse 后台管理专题模型
// @Description 后台管理专题模型
type AdminTopicResponse struct {
	Id          int                `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	CreatedAt   int64              `json:"created_at"`
	UpdatedAt   int64              `json:"updated_at"`
	CoverImage  string             `json:"coverImage"`
	User        SimpleUserResponse `json:"user"`
	UserId      int                `json:"-"`
}
