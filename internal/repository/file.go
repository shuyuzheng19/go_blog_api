package repository

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/pkg/common"
	"blog/pkg/configs"
	"fmt"

	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func (u *FileRepository) FindByMd5(md5 string) string {
	var r string
	u.db.Model(&models.FileMd5Info{}).Select("url").Where("md5=?", md5).Scan(&r)
	return r
}

func (u *FileRepository) FindByMd5Info(md5 string) (models.FileMd5Info, error) {
	var md5Info models.FileMd5Info
	var err = u.db.Model(&models.FileMd5Info{}).First(&md5Info, "md5 = ?", md5).Error
	return md5Info, err
}

func (u *FileRepository) DeleteMd5Infos(md5 string) error {
	// 删除 FileInfo 表中的记录
	var info = &models.FileInfo{}
	if err := u.db.Unscoped().Model(info).Delete(info, "md5 = ?", md5).Error; err != nil {
		return err
	}

	// 删除 FileMd5Info 表中的记录
	var info2 = &models.FileMd5Info{}
	if err := u.db.Unscoped().Model(info2).Delete(info2, "md5 = ?", md5).Error; err != nil {
		return err
	}

	return nil
}

func (u *FileRepository) SaveFileMd5(md5Info *models.FileMd5Info) error {
	return u.db.Model(&models.FileMd5Info{}).Create(md5Info).Error
}

func (u FileRepository) BatchSave(files []models.FileInfo) error {
	return u.db.Model(&models.FileInfo{}).Create(&files).Error
}

func (u FileRepository) GetFileList(uid *int, req requests.FileRequest, count *int64) ([]response.FileResponse, error) {

	var files = make([]response.FileResponse, 0)

	var build = u.db.Model(&models.FileInfo{}).Table(models.FileInfoTable + " f")

	if uid != nil {
		build.Where("f.user_id = ?", uid)
	} else {
		build.Where("f.is_pub = ?", true)
	}

	if req.Keyword != "" {
		build.Where("old_name like ?", "%"+req.Keyword+"%")
	}

	if build.Count(count); *count == 0 {
		return files, nil
	}

	build.Joins(fmt.Sprintf("join %s fm on fm.md5 = f.md5", models.FileInfoMd5Table))

	build.Select("f.id,f.old_name as name,f.created_at,f.suffix,f.size",
		"fm.md5 as md5,fm.url as url").Offset((req.Page - 1) * req.Page).Limit(common.FileListPageCount)

	if req.Sort == "size" {
		build.Order("f.size desc")
	} else {
		build.Order("f.created_at desc")
	}

	var err = build.Scan(&files).Error

	return files, err
}

func (u *FileRepository) UpdateFileInfo(uid *int, req requests.FileUpdateRequest) error {
	var db = u.db.Model(&models.FileInfo{}).Where("id = ?", req.ID)

	if uid != nil {
		db.Where("user_id = ?", *uid)
	}

	if req.Name != nil {
		db.Update("old_name", *req.Name)
	}

	return db.Update("is_pub", req.IsPublic).Error
}

func (u *FileRepository) GetAdminFile(uid *int, req requests.AdminFilterRequest, count *int64) ([]response.FileAdminResponse, error) {

	var db = u.db.Model(&models.FileInfo{}).Table(models.FileInfoTable + " as f")

	var result = make([]response.FileAdminResponse, 0)

	if uid != nil {
		db.Where("f.user_id = ?", uid)
	}

	if req.Start != nil && req.End != nil {
		db.Where("f.created_at BETWEEN ? AND ?", req.Start, req.End)
	}

	if req.Pub != nil {
		db.Where("is_pub = ?", *req.Pub)
	}

	if req.Keyword != nil {
		db.Where("f.old_name like ?", "%"+*req.Keyword+"%")
	}

	if db.Count(count); *count == 0 {
		return result, nil
	}

	var pageCount = req.Size

	var err = db.Joins(fmt.Sprintf("join %s fm on f.md5 = fm.md5", models.FileInfoMd5Table)).
		Joins(fmt.Sprintf("left join %s u on u.id = f.user_id", models.UserTable)).
		Select("f.id as id,f.old_name as name,f.created_at,f.is_pub as public,f.size as size",
			"u.id as uid,u.nick_name as nickname",
			"fm.url as url,fm.md5 as md5,fm.absolute_path as path").
		Offset((req.Page - 1) * pageCount).
		Limit(pageCount).
		Order(req.Sort.GetFilegOrderString("f.")).
		Scan(&result).Error

	return result, err
}

func NewFileRepository() *FileRepository {
	return &FileRepository{db: configs.DB}
}
