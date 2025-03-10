package common

import (
	"blog/internal/models"
	"time"
)

const EnableEmailCode = false

const PageRequest = "page_request"

const AdminRequest = "admin_request"

const (
	PageInfoPrefixKey = "PAGE_INFO" //缓存博客列表的key
)

// 用户缓存键集合
const (
	UserTokenKey       = "USER_TOKEN:"    //缓存用户Token的Key
	EmailCodeKey       = "EMAIL_CODE:"    //缓存注册邮箱验证码的key
	EmailCodeKeyExpire = time.Minute      //邮箱验证码过期实际
	UserInfoKey        = "USER_INFO:"     //缓存用户信息的key
	UserInfoKeyExpire  = time.Minute * 30 //用户信息过期时间
)

// 博客相关缓存
const (
	BlogMapKey         = "BLOG_MAP"       //缓存博客详情的key
	RecommendKey       = "RECOMMEND_BLOG" //缓存推荐博客的key
	HotBlogKey         = "HOT_BLOG"       //缓存热门博客的key
	HotBlogExpire      = time.Hour * 24   //热门博客过期时间
	LatestBlogKey      = "LATEST_BLOG"    //缓存最新博客的key
	PageInfoExpire     = time.Hour * 3    //博客列表过期时间
	LatestBlogExpire   = time.Hour * 10   //最新博客过期时间
	BlogEyeCountMapKey = "EYE_MAP"        //缓存博客的浏览量
	EyeView            = "EYE_VIEW"       //统计今日浏览量
	PinnedBlog         = "PINNED_BLOG"    //置顶博客
)

// 分类缓存键集合
const (
	CategoryListKey    = "CATEGORY_LIST"    //缓存分类列表的key
	CategoryListExpire = time.Hour * 24 * 3 //分类列表过期时间
)

// 标签缓存键集合
const (
	RandomTagListKey    = "RANDOM_TAG"       //随机获取标签的数量
	RandomTagCount      = 25                 //随机获取标签的数量
	RandomTagListExpire = time.Hour * 24 * 3 //随机标签的过期时间
	TagMapKey           = "TAG_MAP"          //标签简要信息的key
)

// 专题缓存集合
const (
	TopicPageKey    = "TOPIC_PAGE:"      //缓存专题页的key
	TopicPageExpire = time.Hour * 24 * 3 //专题页过期时间
	TopicMapKey     = "TOPIC_MAP"        //专题简要信息的key
)

// Count
const (
	RecommendBlogCount  = 4
	TopicPageCount      = 20
	ArchivePageCount    = 15
	SearchBlogPageCount = 10
	FileListPageCount   = 15
)

// JWT
const (
	TokenEncrypted = "shuyuYuice----"
	TokenExpire    = time.Hour * 24 * 7
)

type RoleId uint

// Role 角色
const (
	UserRoleId       RoleId = 1
	AdminRoleId      RoleId = 2
	SuperAdminRoleId RoleId = 3
)

// GetJwtUser 根据用户ID获取用户信息的函数变量
var GetJwtUser = func(id int) *models.User {
	return &models.User{}
}

// GetToken 根据用户ID获取token的函数变量
var GetToken = func(id int) string {
	return ""
}

const (
	WebSiteConfigKey = "WEB_SITE_CONFIG" //缓存网站配置信息

	TmpBlogKey = "TMP_BLOG:" //缓存临时的博客

	TarKey = "TAR:"
)
