package main

import (
	"blog/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbData struct {
	server *gorm.DB
	local  *gorm.DB
}

func (d *DbData) CopyTag() {
	var tableName = "tags"

	var tags = make([]models.Tag, 0)

	d.local.Table(tableName).Find(&tags)

	var err = d.server.Model(&models.Tag{}).Save(&tags)

	fmt.Println("err=", err)

	fmt.Println("tags=", tags)

}

func (d *DbData) CopyCategory() {
	var tableName = "categories"

	var categories = make([]models.Category, 0)

	d.local.Table(tableName).Find(&categories)

	d.server.Model(&models.Category{}).Table(tableName).Save(&categories)

	fmt.Println("categories=", categories)
}

func (d *DbData) CopyFileMD5() {
	var tableName = "file_md5_infos"

	var md5s = make([]models.FileMd5Info, 0)

	d.local.Table(tableName).Find(&md5s)

	d.server.Table(tableName).Save(&md5s)

	fmt.Println("md5s=", md5s)
}

func (d *DbData) CopyTopics() {
	var tableName = "topics"

	var topics = make([]models.Topic, 0)

	d.local.Table(tableName).Find(&topics)

	d.server.Table(tableName).Save(&topics)

	fmt.Println("topics=", topics)

}

/**
INSERT INTO file_infos (id, new_name, suffix, created_at, updated_at, deleted_at, md5, is_pub, old_name, user_id, size)
*/

func (d *DbData) CopyFile() {
	var tableName = "file_infos"

	var files = make([]models.FileInfo, 0)

	d.local.Table(tableName).Find(&files)

	d.server.Table(tableName).Save(&files)

	fmt.Println("files=", files)

}

func (d *DbData) CopyUser() {
	var tableName = "users"

	var users = make([]models.User, 0)

	d.local.Table(tableName).Find(&users)

	d.server.Table(tableName).Save(&users)

	fmt.Println("users=", users)

}

func (d *DbData) CopyBlog() {
	var tableName = "blogs"

	var blogs = make([]models.Blog, 0)

	d.local.Table(tableName).Find(&blogs)

	var err = d.server.Table(tableName).Save(&blogs).Error

	fmt.Println("err=", err)

	// fmt.Println("tags=", tags)

}

func (d *DbData) CopyBlogTags() {
	var tableName = "blogs_tags"

	var maps = make([]map[string]interface{}, 0)

	d.local.Table(tableName).Find(&maps)

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

	d.server.Table(tableName).Create(&maps2)

	fmt.Println("tags=", maps2)

}

func main() {
	var service = NewDbData()
	service.CopyUser()
	// service.CopyBlogTags()
	// service.CopyBlog()
	// service.CopyFile()
	// service.CopyFileMD5()
	// service.CopyTopics()
	// service.CopyCategory()
	// service.CopyTag()
}

func NewDbData() *DbData {
	var localdsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable Timezone=%s",
		"127.0.0.1", 5432, "root", "go_zmc_blog", "123456", "Asia/Shanghai")
	var serverdsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable Timezone=%s",
		"45.77.1.30", 5432, "zsy", "enming_blog", "xiaoyu2528959216 ", "Asia/Shanghai")

	var local, _ = gorm.Open(postgres.Open(localdsn), &gorm.Config{})

	var server, _ = gorm.Open(postgres.Open(serverdsn), &gorm.Config{})

	return &DbData{local: local, server: server}
}
