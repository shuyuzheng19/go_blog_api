package handler

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/internal/service"
	"blog/internal/utils"
	"blog/pkg/common"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"blog/pkg/smail"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type FileController struct {
	service *service.FileService
	system  *service.SystemService
}

// UploadFile 处理文件上传
func (f *FileController) UploadFile(ctx fiber.Ctx) error {
	files, err := ctx.MultipartForm()
	if err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "未检测到任何文件，请选择文件后重试")
	}

	uid := ctx.Locals("uid").(int)
	isPub, _ := strconv.ParseBool(ctx.Query("is_pub", "false"))

	list := f.service.SaveFile(files, &uid, isPub, false)
	return ResultSuccessToResponse(list, ctx)
}

// UploadAvatar 处理头像上传
func (f *FileController) UploadAvatar(ctx fiber.Ctx) error {
	return f.uploadSingleFile(ctx, "files", "只能上传一个头像文件", "文件类型不支持，请上传图片格式", "头像文件不能超过5MB")
}

func (f *FileController) DownloadTar(ctx fiber.Ctx) error {
	var uuid = ctx.Query("id")

	str, err := configs.REDIS.Get(common.TarKey + uuid).Result()

	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "不允许访问")
	}

	// 获取文件名
	fileName := path.Base(str)

	// 发送文件，自动触发下载
	return ctx.Download(str, fileName)
}

func (f *FileController) TarDockerComposeData(ctx fiber.Ctx) error {

	logger.Info("重要！！！重要！！！！重要！！！执行打包下载功能")

	var req requests.TarRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数绑定失败!")
	}

	if errs := Validate(req); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	go func() {
		if uuid, err := service.CreateTarGzAndStoreRedis(req.Path, req.Min); err == nil {
			var downloadUrl = "https://blog.shuyuz.com/api/v1/file/admin/system_file/tar?id=" + uuid
			smail.SendEmail(configs.CONFIG.MyEmail, "打包完毕", true, fmt.Sprintf(`链接<br/> <a style="color:blue" href="%s">%s</a>`, downloadUrl, downloadUrl))
			logger.Info("打包成功,已发送至你的邮箱", zap.String("下载链接", downloadUrl))
		}
	}()

	return ResultSuccessToResponse(nil, ctx)
}

// UploadImage 处理图片上传
func (f *FileController) UploadImage(ctx fiber.Ctx) error {
	return f.uploadSingleFile(ctx, "files", "只能上传一个图片文件", "文件类型不支持，请上传图片格式", "图片文件不能超过5MB")
}

func (f *FileController) GetAdminFileList(ctx fiber.Ctx) error {
	req := ctx.Locals(common.AdminRequest).(requests.AdminFilterRequest)

	page := response.Page{
		Page: req.Page,
		Size: req.Size,
	}

	uid := ctx.Locals("uid").(int)
	var userId *int
	rid := ctx.Locals("rid").(uint)

	if rid != uint(common.SuperAdminRoleId) {
		userId = &uid
	}

	err := f.service.GetAdminFileList(userId, req, &page)
	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取文件列表，请稍后重试")
	}

	return ResultSuccessToResponse(&page, ctx)
}

func (f *FileController) DeleteSystemFile(ctx fiber.Ctx) error {
	user := ctx.Locals("user").(*models.User)

	if user.Email != configs.CONFIG.MyEmail {
		return ResultErrorToResponse(common.Unauthorized, ctx, "您没有权限执行此操作")
	}

	var paths []string
	if err := ctx.Bind().Body(&paths); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请求参数无效，请检查输入")
	}

	result := f.system.DeleteSystemFile(paths)
	return ResultSuccessToResponse(result, ctx)
}

func getSystemFilePath(c fiber.Ctx) string {
	token := c.Query("token")
	_, rid := utils.ParseTokenUserIdAndRoleId(token)

	if rid != int(common.SuperAdminRoleId) {
		ResultErrorToResponse(common.Unauthorized, c, "权限验证失败")
		return ""
	}

	path := c.Query("path")
	if path == "" {
		ResultErrorToResponse(common.Unauthorized, c, "文件路径不能为空")
		return ""
	}

	return path
}

// DownloadSystemFile 下载本地文件
func (f FileController) DownloadSystemFile(c fiber.Ctx) error {
	path := getSystemFilePath(c)

	if path != "" {
		name := filepath.Base(path)

		return c.Download(path, name)
	}

	return nil
}

func (f *FileController) GetSystemFile(ctx fiber.Ctx) error {
	user := ctx.Locals("user").(*models.User)

	if user.Email != configs.CONFIG.MyEmail {
		return ResultErrorToResponse(common.Unauthorized, ctx, "您没有权限执行此操作")
	}

	var req requests.SystemFileRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请求参数无效，请检查输入")
	}

	result := f.system.GetSystemFile(req)
	return ResultSuccessToResponse(result, ctx)
}

func (f FileController) ClearSystemFileContent(ctx fiber.Ctx) error {
	user := ctx.Locals("user").(*models.User)

	if user.Email != configs.CONFIG.MyEmail {
		return ResultErrorToResponse(common.Unauthorized, ctx, "您没有权限执行此操作")
	}

	path := ctx.Query("path")
	if path != "" {
		if err := f.system.ClearFileContent(path); err != nil {
			return ResultErrorToResponse(common.FAIL, ctx, "无法清除文件内容，请稍后重试")
		}
	}

	return ResultSuccessToResponse(nil, ctx)
}

func (f *FileController) GetCurrentLog(c fiber.Ctx) error {
	logger := configs.CONFIG.Logger
	path := filepath.Join(logger.LoggerDir, logger.DefaultName)
	name := filepath.Base(path)

	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Transfer-Encoding", "binary")

	return c.SendFile(path)
}

func (f *FileController) GetLogFileList(ctx fiber.Ctx) error {
	keyword := ctx.Query("keyword")
	result := f.system.GetSystemFile(requests.SystemFileRequest{
		Path:    configs.CONFIG.Logger.LoggerDir,
		Keyword: keyword,
	})

	return ResultSuccessToResponse(result, ctx)
}

func (f *FileController) UpdateFileInfo(ctx fiber.Ctx) error {
	var req requests.FileUpdateRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请求参数无效，请检查输入")
	}

	if errs := Validate(&req); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	uid := ctx.Locals("uid").(int)
	var userId *int
	rid := ctx.Locals("rid").(uint)

	if rid != uint(common.SuperAdminRoleId) {
		userId = &uid
	}

	if err := f.service.UpdateFileInfo(userId, req); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "文件信息更新失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// DeleteByIds 删除文件
func (f *FileController) DeleteByIDs(ctx fiber.Ctx) error {
	var ids []int64

	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请求参数无效，请检查输入")
	}

	uid := ctx.Locals("uid").(int)
	var userId *int
	rid := ctx.Locals("rid").(uint)

	if rid != uint(common.SuperAdminRoleId) {
		userId = &uid
	}

	if err := f.service.DeleteFileByIDs(userId, ids); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "文件删除失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// uploadSingleFile 处理单文件上传的公共逻辑
func (f *FileController) uploadSingleFile(ctx fiber.Ctx, fieldName, errMsgFileCount, errMsgContentType, errMsgFileSize string) error {
	files, err := ctx.MultipartForm()
	if err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "未检测到任何文件，请选择文件后重试")
	}

	filesList := files.File[fieldName]
	if len(filesList) != 1 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, errMsgFileCount)
	}

	file := filesList[0]
	if !strings.Contains(file.Header.Get("Content-Type"), "image/") {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, errMsgContentType)
	}

	if file.Size > 1024*1024*5 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, errMsgFileSize)
	}

	list := f.service.SaveFile(files, nil, false, true)
	return ResultSuccessToResponse(list, ctx)
}

// GetPublicFileList 获取公开文件列表
func (f *FileController) GetPublicFileList(ctx fiber.Ctx) error {
	return f.getFileList(ctx, nil)
}

func (f *FileController) GetUploadConfig(ctx fiber.Ctx) error {
	return ResultSuccessToResponse(f.service.GetConfig(), ctx)
}

func (f *FileController) SetUploadConfig(ctx fiber.Ctx) error {
	var req configs.UploadConfig

	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数错误")
	}

	f.service.SetConfig(req)

	return ResultSuccessToResponse(f.service.GetConfig(), ctx)
}

func (f *FileController) DeleteFileByMd5(ctx fiber.Ctx) error {
	var md5 = ctx.Query("md5")

	if md5 == "" {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "md5为空")
	}

	var err = f.service.DeleteMd5(md5)

	if err != nil {
		logger.Warn("删除文件MD5失败", zap.String("md5", md5))
		return ResultErrorToResponse(common.ERROR, ctx, "删除失败")
	}

	logger.Info("删除MD5文件成功", zap.String("md5", md5))

	return ResultSuccessToResponse(nil, ctx)
}

// GetCurrentFileFileList 获取当前用户文件列表
func (f *FileController) GetCurrentFileFileList(ctx fiber.Ctx) error {
	uid := ctx.Locals("uid").(int)
	return f.getFileList(ctx, &uid)
}

// getFileList 获取文件列表的公共逻辑
func (f *FileController) getFileList(ctx fiber.Ctx, uid *int) error {
	var req requests.FileRequest

	if err := ctx.Bind().Query(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "文件列表请求参数无效，请检查输入")
	}

	var page response.Page

	err := f.service.GetCurrentFiles(uid, req, &page)
	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取文件列表，请稍后重试")
	}

	return ResultSuccessToResponse(&page, ctx)
}

// NewFileController 创建新的文件控制器
func NewFileController() *FileController {
	return &FileController{service: service.NewFileService(), system: service.NewSystemService()}
}
