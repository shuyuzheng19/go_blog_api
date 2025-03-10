package search

// MeiliSearchRequest 表示 MeiliSearch 的搜索请求结构体
type MeiliSearchRequest struct {
	Q                     string   `json:"q"`                     // 搜索查询字符串
	Offset                int      `json:"offset,omitempty"`      // 偏移量，用于分页
	Limit                 int      `json:"limit,omitempty"`       // 每页返回的最大结果数
	HighlightPreTag       string   `json:"highlightPreTag"`       // 高亮前标签
	HighlightPostTag      string   `json:"highlightPostTag"`      // 高亮后标签
	ShowMatchesPosition   bool     `json:"showMatchesPosition"`   // 是否显示匹配位置
	Sort                  []string `json:"sort"`                  // 排序字段
	AttributesToHighlight []string `json:"attributesToHighlight"` // 需要高亮的属性
}

// NewSearchRequest 初始化一个新的 MeiliSearchRequest
func NewSearchRequest() *MeiliSearchRequest {
	return &MeiliSearchRequest{}
}

// SetAttributesToHighlight 设置需要高亮的属性
func (s *MeiliSearchRequest) SetAttributesToHighlight(highlight []string) *MeiliSearchRequest {
	s.AttributesToHighlight = highlight
	return s
}

// SetShowMatchesPosition 设置是否显示匹配位置
func (s *MeiliSearchRequest) SetShowMatchesPosition(show bool) *MeiliSearchRequest {
	s.ShowMatchesPosition = show
	return s
}

// SetQ 设置搜索查询字符串
func (s *MeiliSearchRequest) SetQ(q string) *MeiliSearchRequest {
	s.Q = q
	return s
}

// SetOffset 设置偏移量
func (s *MeiliSearchRequest) SetOffset(offset int) *MeiliSearchRequest {
	s.Offset = offset
	return s
}

// SetLimit 设置每页返回的最大结果数
func (s *MeiliSearchRequest) SetLimit(limit int) *MeiliSearchRequest {
	s.Limit = limit
	return s
}

// SetHighlightPreTag 设置高亮前标签
func (s *MeiliSearchRequest) SetHighlightPreTag(highlightPreTag string) *MeiliSearchRequest {
	s.HighlightPreTag = highlightPreTag
	return s
}

// SetHighlightPostTag 设置高亮后标签
func (s *MeiliSearchRequest) SetHighlightPostTag(highlightPostTag string) *MeiliSearchRequest {
	s.HighlightPostTag = highlightPostTag
	return s
}

// SetSort 设置排序字段
func (s *MeiliSearchRequest) SetSort(sort []string) *MeiliSearchRequest {
	s.Sort = sort
	return s
}

// Build 返回构建好的 MeiliSearchRequest
func (s *MeiliSearchRequest) Build() MeiliSearchRequest {
	return *s
}
