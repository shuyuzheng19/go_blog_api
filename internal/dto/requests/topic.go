package requests

import "blog/internal/models"

func (sort Sort) GetTopicOrderString(prefix string) string {
	switch sort {
	case CREATE:
		return prefix + "created_at desc"
	case UPDATE:
		return prefix + "updated_at  desc"
	case BACK:
		return prefix + "created_at asc"
	default:
		return prefix + "created_at desc"
	}
}

type TopicRequest struct {
	ID          int    `json:"id"`                               //专题id
	Name        string `json:"name" validate:"required,max=50"`  //专题名
	Description string `json:"desc" validate:"required,max=200"` //专题描述
	CoverImage  string `json:"cover" validate:"required"`
}

func (t *TopicRequest) ToModel(uid int) models.Topic {
	return models.Topic{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		CoverImage:  t.CoverImage,
		UserID:      uid,
	}
}
