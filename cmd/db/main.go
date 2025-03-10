package main

import (
	"blog/internal/models"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbData struct {
	server *gorm.DB
	local  *gorm.DB
}

func (d *DbData) CopyTag() {
	var tableName = "tags"

	var maps = make([]map[string]interface{}, 0)

	d.server.Table(tableName).Find(&maps)

	var tags = make([]models.Tag, 0)

	for _, value := range maps {
		tags = append(tags, models.Tag{
			ID:   int(value["id"].(int32)),
			Name: value["name"].(string),
			Model: models.Model{
				CreatedAt: value["created_at"].(time.Time).Unix(),
				UpdatedAt: value["updated_at"].(time.Time).Unix(),
			},
		})
	}

	d.local.Table(tableName).Save(&tags)

	fmt.Println("tags=", tags)

}

func (d *DbData) CopyTag2() {
	var tableName = "tags"

	var maps = make([]map[string]interface{}, 0)

	d.local.Table(tableName).Find(&maps)

	var tags = make([]models.Tag, 0)

	for _, value := range maps {
		tags = append(tags, models.Tag{
			ID:   int(value["id"].(int64)),
			Name: value["name"].(string),
			Model: models.Model{
				CreatedAt: value["created_at"].(time.Time).Unix(),
				UpdatedAt: value["updated_at"].(time.Time).Unix(),
			},
		})
	}

	d.server.Table(tableName).Save(&tags)

	fmt.Println("tags=", tags)

}

func (d *DbData) CopyCategory() {
	var tableName = "categories"

	var maps = make([]map[string]interface{}, 0)

	d.server.Table(tableName).Find(&maps)

	fmt.Println(maps)

	var categories = make([]models.Category, 0)

	for _, value := range maps {
		categories = append(categories, models.Category{
			ID:   int(value["id"].(int32)),
			Name: value["name"].(string),
			Model: models.Model{
				CreatedAt: value["created_at"].(time.Time).Unix(),
				UpdatedAt: value["updated_at"].(time.Time).Unix(),
			},
		})
	}

	d.local.Table(tableName).Save(&categories)

	fmt.Println("categories=", categories)
}

func (d *DbData) CopyCategory2() {
	var tableName = "categories"

	var maps = make([]map[string]interface{}, 0)

	d.local.Table(tableName).Find(&maps)

	fmt.Println(maps)

	var categories = make([]models.Category, 0)

	for _, value := range maps {
		categories = append(categories, models.Category{
			ID:   int(value["id"].(int64)),
			Name: value["name"].(string),
			Model: models.Model{
				CreatedAt: value["created_at"].(time.Time).Unix(),
				UpdatedAt: value["updated_at"].(time.Time).Unix(),
			},
		})
	}

	d.server.Table(tableName).Save(&categories)

	fmt.Println("categories=", categories)
}

func (d *DbData) CopyFileMD5() {
	var tableName = "file_md5_infos"

	var maps = make([]map[string]interface{}, 0)

	d.server.Table(tableName).Find(&maps)

	var tags = make([]models.FileMd5Info, 0)

	for _, value := range maps {
		tags = append(tags, models.FileMd5Info{
			Md5:          value["md5"].(string),
			Url:          value["url"].(string),
			AbsolutePath: value["absolute_path"].(string),
		})
	}

	d.local.Table(tableName).Save(&tags)

	fmt.Println("tags=", tags)
}

func (d *DbData) CopyTopics() {
	var tableName = "topics"

	var maps = make([]map[string]interface{}, 0)

	d.server.Table(tableName).Find(&maps)

	var tags = make([]models.Topic, 0)

	for _, value := range maps {
		tags = append(tags, models.Topic{
			ID:          int(value["id"].(int32)),
			Name:        value["name"].(string),
			Description: value["description"].(string),
			CoverImage:  value["cover_image"].(string),
			UserID:      int(value["user_id"].(int64)),
			Model: models.Model{
				CreatedAt: value["created_at"].(time.Time).Unix(),
				UpdatedAt: value["updated_at"].(time.Time).Unix(),
			},
		})
	}

	d.local.Table(tableName).Save(&tags)

	fmt.Println("tags=", tags)

}

func (d *DbData) CopyTopics2() {
	var tableName = "topics"

	var maps = make([]map[string]interface{}, 0)

	d.local.Table(tableName).Find(&maps)

	var tags = make([]models.Topic, 0)

	for _, value := range maps {
		tags = append(tags, models.Topic{
			ID:          int(value["id"].(int64)),
			Name:        value["name"].(string),
			Description: value["description"].(string),
			CoverImage:  value["cover_image"].(string),
			UserID:      int(value["user_id"].(int64)),
			Model: models.Model{
				CreatedAt: value["created_at"].(time.Time).Unix(),
				UpdatedAt: value["updated_at"].(time.Time).Unix(),
			},
		})
	}

	d.server.Table(tableName).Save(&tags)

	fmt.Println("tags=", tags)

}

/**
INSERT INTO file_infos (id, new_name, suffix, created_at, updated_at, deleted_at, md5, is_pub, old_name, user_id, size)
*/

func (d *DbData) CopyFile() {
	var tableName = "file_infos"

	var maps = make([]map[string]interface{}, 0)

	d.server.Table(tableName).Find(&maps)

	var tags = make([]models.FileInfo, 0)

	for _, value := range maps {
		var uid = value["user_id"].(int64)
		var userID = int(uid)
		tags = append(tags, models.FileInfo{
			ID:      int(value["id"].(int64)),
			NewName: value["new_name"].(string),
			Suffix:  value["suffix"].(string),
			FileMd5: value["md5"].(string),
			IsPub:   value["is_pub"].(bool),
			UserID:  &userID,
			OldName: value["old_name"].(string),
			Size:    value["size"].(int64),
			Model: models.Model{
				CreatedAt: value["created_at"].(time.Time).Unix(),
				UpdatedAt: value["updated_at"].(time.Time).Unix(),
			},
		})
	}

	d.local.Table(tableName).Save(&tags)

	fmt.Println("tags=", tags)

}

func (d *DbData) CopyBlog() {
	var tableName = "blogs"

	var maps = make([]map[string]interface{}, 0)

	d.server.Table(tableName).Find(&maps)

	var tags = make([]models.Blog, 0)

	for _, value := range maps {
		var uid = value["user_id"].(int64)
		var userID = int(uid)
		/*
			*
			INSERT INTO blogs (created_at, deleted_at, source_url, cover_image, updated_at, id, title, like_count, topic_id, description, content, eye_count, category_id, user_id)
		*/
		var source_url = value["source_url"]

		var source *string = nil

		if source_url != nil {
			var s = source_url.(string)
			source = &s
		}

		var topic_id = value["topic_id"]

		var tid *int = nil

		if topic_id != nil {
			var i = int(topic_id.(int64))
			tid = &i
		}

		var category_id = value["category_id"]

		var cid *int

		if category_id != nil {
			var c = int(category_id.(int64))
			cid = &c
		}

		tags = append(tags, models.Blog{
			ID:          int64(value["id"].(int64)),
			UserID:      userID,
			SourceURL:   source,
			CoverImage:  value["cover_image"].(string),
			Title:       value["title"].(string),
			EyeCount:    value["eye_count"].(int64),
			TopicID:     tid,
			Description: value["description"].(string),
			Content:     value["content"].(string),
			CategoryID:  cid,
			Model: models.Model{
				CreatedAt: value["created_at"].(time.Time).Unix(),
				UpdatedAt: value["updated_at"].(time.Time).Unix(),
			},
		})
	}

	var err = d.local.Table(tableName).Save(&tags).Error

	fmt.Println("err=", err)

	// fmt.Println("tags=", tags)

}

func (d *DbData) CopyBlog2() {
	var tableName = "blogs"

	var maps = make([]map[string]interface{}, 0)

	d.local.Table(tableName).Where("user_id = ?", 4).Find(&maps)

	var tags = make([]models.Blog, 0)

	for _, value := range maps {
		var uid = value["user_id"].(int64)
		var userID = int(uid)
		/*
			*
			INSERT INTO blogs (created_at, deleted_at, source_url, cover_image, updated_at, id, title, like_count, topic_id, description, content, eye_count, category_id, user_id)
		*/
		var source_url = value["source_url"]

		var source *string = nil

		if source_url != nil {
			var s = source_url.(string)
			source = &s
		}

		var topic_id = value["topic_id"]

		var tid *int = nil

		if topic_id != nil {
			var i = int(topic_id.(int64))
			tid = &i
		}

		var category_id = value["category_id"]

		var cid *int

		if category_id != nil {
			var c = int(category_id.(int64))
			cid = &c
		}

		tags = append(tags, models.Blog{
			ID:          int64(value["id"].(int64)),
			UserID:      userID,
			SourceURL:   source,
			CoverImage:  value["cover_image"].(string),
			Title:       value["title"].(string),
			EyeCount:    value["eye_count"].(int64),
			TopicID:     tid,
			Description: value["description"].(string),
			Content:     value["content"].(string),
			CategoryID:  cid,
			Model: models.Model{
				CreatedAt: value["created_at"].(time.Time).Unix(),
				UpdatedAt: value["updated_at"].(time.Time).Unix(),
			},
		})
	}

	var err = d.server.Table(tableName).Save(&tags).Error

	fmt.Println("err=", err)

	// fmt.Println("tags=", tags)

}

func (d *DbData) CopyBlogTags() {
	var tableName = "blogs_tags"

	var maps = make([]map[string]interface{}, 0)

	d.server.Table(tableName).Find(&maps)

	var maps2 = make([]map[string]interface{}, 0)

	// INSERT INTO blogs_tags (blog_id, tag_id) VALUES (414, 55);
	for _, value := range maps {
		var bid = value["blog_id"].(int64)
		var tid = value["tag_id"].(int64)
		maps2 = append(maps2, map[string]interface{}{
			"blog_id": bid,
			"tag_id":  tid,
		})
	}

	d.local.Table(tableName).Create(&maps2)

	fmt.Println("tags=", maps2)

}

func (d *DbData) CopyBlogTags2() {
	var tableName = "blogs_tags"

	var maps = make([]map[string]interface{}, 0)

	d.local.Table(tableName).Find(&maps)

	var maps2 = make([]map[string]interface{}, 0)

	// INSERT INTO blogs_tags (blog_id, tag_id) VALUES (414, 55);
	for _, value := range maps {
		var bid = value["blog_id"].(int64)
		var tid = value["tag_id"].(int64)
		d.server.Table(tableName).Create(&map[string]interface{}{
			"blog_id": bid,
			"tag_id":  tid,
		})
		// maps2 = append(maps2, map[string]interface{}{
		// 	"blog_id": bid,
		// 	"tag_id":  tid,
		// })
	}

	// d.server.Table(tableName).Create(&maps2)

	fmt.Println("tags=", maps2)

}

func main() {
	var service = NewDbData()
	service.CopyBlogTags2()
	// service.CopyBlog2()
	// service.CopyFile()
	// service.CopyFileMD5()
	// service.CopyTopics2()
	// service.CopyCategory2()
	// service.CopyTag2()
}

func NewDbData() *DbData {
	var localdsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable Timezone=%s",
		"127.0.0.1", 5432, "root", "go_zmc_blog", "123456", "Asia/Shanghai")
	var serverdsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable Timezone=%s",
		"45.77.1.30", 5432, "zsy", "enming_blog", "xiaoyu2528959216 ", "Asia/Shanghai")

	// var localdsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable Timezone=%s",
	// 	"127.0.0.1", 5432, "root", "blog", "123456", "Asia/Shanghai")
	// var serverdsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable Timezone=%s",
	// 	"144.202.96.150", 5432, "zsy", "go-shuyu-blog", "xiaoyu2528959216 ", "Asia/Shanghai")

	var local, _ = gorm.Open(postgres.Open(localdsn), &gorm.Config{})

	var server, _ = gorm.Open(postgres.Open(serverdsn), &gorm.Config{})

	return &DbData{local: local, server: server}
}
