package requests

func (sort Sort) GetTagOrderString(prefix string) string {
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

type TagRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
