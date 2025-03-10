package service

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/internal/repository"
	"blog/internal/store"
	"blog/internal/utils"
	"blog/pkg/common"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

type FileService struct {
	repository *repository.FileRepository
	config     configs.UploadConfig
}

// NewFileService 创建一个新的 FileService 实例
func NewFileService() *FileService {
	return &FileService{
		repository: repository.NewFileRepository(),
		config:     configs.CONFIG.Upload,
	}
}

const mb = 1024 * 1024 // 1 MB 转换单位

func (f *FileService) GetConfig() configs.UploadConfig {
	return f.config
}

func (f *FileService) SetConfig(config configs.UploadConfig) {
	f.config = config
}

// SaveFile 处理文件上传并返回上传的文件信息
func (f *FileService) SaveFile(form *multipart.Form, uid *int, isPub, isImg bool) []response.SimpleFileResponse {
	var wg sync.WaitGroup
	fileList := make([]models.FileInfo, 0)
	var fileMap sync.Map

	for _, file := range form.File["files"] {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			if file.Size > int64(f.config.MaxFileSize)*mb {
				logger.Info("文件大小超过限制", zap.String("filename", file.Filename))
				return
			}
			if err := f.processFile(file, uid, isPub, isImg, &fileMap); err != nil {
				logger.Info("处理文件时出错", zap.String("error", err.Error()))
			}
		}(file)
	}

	wg.Wait()

	fileMap.Range(func(_, value interface{}) bool {
		fileList = append(fileList, value.(models.FileInfo))
		return true
	})

	simpleList := make([]response.SimpleFileResponse, len(fileList))
	for i, file := range fileList {
		simpleList[i] = response.SimpleFileResponse{
			Name:    file.NewName,
			OldName: file.OldName,
			Url:     file.FileMd5Info.Url,
		}
	}

	f.repository.BatchSave(fileList)

	return simpleList
}

// processFile 处理单个文件的上传逻辑
func (f *FileService) processFile(file *multipart.FileHeader, userId *int, isPub, isImg bool, fileMap *sync.Map) error {
	md5Value := utils.CalculateMD5(file)
	if md5Value == "" {
		return fmt.Errorf("计算MD5失败")
	}

	ext := filepath.Ext(file.Filename)
	newFileName := md5Value + ext
	url := f.repository.FindByMd5(md5Value)
	var path string
	if url == "" {
		dstPath := filepath.Join(f.config.Path, newFileName)
		if isImg {
			if image, err := store.UploadImageToVeyme(file, f.config.VeymeToken); err != nil {
				if err := store.SaveToFile(file, dstPath); err != nil {
					return err
				}
				url = fmt.Sprintf("%s/%s", f.config.Uri, newFileName)
				logger.Info("Veyme上传图片成功,使用本地文件上传")
				path = filepath.Join(f.config.Path, newFileName)
			} else {
				path = image.Url
				url = image.Image
				logger.Info("Veyme上传图片成功")
			}
		} else if f.config.Store == "github" {
			var name = "blog/" + md5Value + ext
			if urll, err := store.UploadImageToGitHub(file, name, f.config.Github); err != nil {
				return err
			} else {
				url = urll
			}
			logger.Info("Github上传文件成功")
			path = name
		} else {
			if err := store.SaveToFile(file, dstPath); err != nil {
				return err
			}
			url = fmt.Sprintf("%s/%s", f.config.Uri, newFileName)
			logger.Info("本地文件上传成功")
			path = filepath.Join(f.config.Path, newFileName)
		}
		newMd5Info := models.FileMd5Info{
			Md5:          md5Value,
			Url:          url,
			AbsolutePath: path,
		}
		if err := f.repository.SaveFileMd5(&newMd5Info); err != nil {
			return err
		}
	}

	newFile := models.FileInfo{
		OldName: file.Filename,
		NewName: newFileName,
		UserID:  userId,
		Suffix:  ext,
		Size:    file.Size,
		FileMd5: md5Value,
		FileMd5Info: models.FileMd5Info{
			Md5: md5Value,
			Url: url,
		},
		IsPub: isPub,
	}

	fileMap.Store(md5Value, newFile)

	return nil
}

func (f *FileService) GetCurrentFiles(uid *int, req requests.FileRequest, page *response.Page) error {
	if req.Page <= 0 {
		req.Page = 1
	}
	list, err := f.repository.GetFileList(uid, req, &page.Count)
	page.Page = req.Page
	page.Size = common.FileListPageCount
	page.Data = list
	return err
}

// GetAdminFileList 获取管理员文件列表
func (f *FileService) GetAdminFileList(uid *int, req requests.AdminFilterRequest, page *response.Page) error {
	files, err := f.repository.GetAdminFile(uid, req, &page.Count)
	if err != nil {
		logger.Info("获取管理员文件列表失败", zap.Int("uid", *uid), zap.String("error", err.Error()))
	}
	page.Data = files
	return err
}

func (f *FileService) UpdateFileInfo(uid *int, req requests.FileUpdateRequest) error {
	return f.repository.UpdateFileInfo(uid, req)
}

func (f *FileService) DeleteMd5(md5 string) error {
	var info, err = f.repository.FindByMd5Info(md5)
	if err != nil {
		return err
	}

	err = f.repository.DeleteMd5Infos(info.Md5)

	if err != nil {
		return err
	}

	err = os.Remove(info.AbsolutePath)

	if err != nil {
		logger.Warn("删除本地文件MD5失败", zap.String("md5", info.Md5))
	}

	return nil
}

func (f *FileService) DeleteFileByIDs(uid *int, ids []int64) error {
	return configs.DeleteData(models.FileInfoTable, uid, ids)
}
