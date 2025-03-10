package handler

import (
	"blog/internal/models"
	"blog/pkg/common"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// ResultSuccessToResponse 返回成功响应
func ResultSuccessToResponse(data interface{}, c fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(common.SUCCESS.DoData(data)) // 返回状态200和数据
}

// ResultErrorToResponse 返回错误响应
func ResultErrorToResponse(code common.Code, c fiber.Ctx, err string) error {
	var data = common.R{
		Code:    code,
		Message: err,
	} // 获取错误代码信息
	return c.Status(http.StatusOK).JSON(data) // 返回状态200和错误信息
}

// ResultValidatorErrorToResponse 返回验证错误响应
func ResultValidatorErrorToResponse(c fiber.Ctx, errs interface{}) error {
	return c.Status(http.StatusOK).JSON(common.R{
		Code:    common.BAD_REQUEST,
		Message: "请求参数验证失败",
		Data:    errs,
	}) // 返回状态200和验证错误信息
}

// ErrorResponse 验证错误响应结构体
type ErrorResponse struct {
	Error       bool   `json:"error"`   // 是否错误
	FailedField string `json:"field"`   // 失败字段
	Tag         string `json:"tag"`     // 验证标签
	Message     string `json:"message"` // 错误信息
}

var validate = validator.New() // 创建验证器实例

// Validate 验证数据结构
func Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{} // 存储验证错误

	errs := validate.Struct(data) // 验证数据结构
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse // 创建错误响应实例

			fieldName := err.Field()                      // 获取失败字段名
			elem.FailedField = strings.ToLower(fieldName) // 转为小写
			elem.Tag = err.Tag()                          // 获取验证标签
			elem.Error = true                             // 标记为错误

			// 获取字段的标签信息
			typeOf := reflect.TypeOf(data)
			if typeOf.Kind() == reflect.Ptr {
				typeOf = typeOf.Elem() // 获取指针指向的类型
			}
			field, ok := typeOf.FieldByName(fieldName) // 查找字段
			if ok {
				elem.Message = field.Tag.Get("error") // 获取字段的错误信息
			}
			validationErrors = append(validationErrors, elem) // 添加到错误列表
		}
	}

	return validationErrors // 返回所有验证错误
}

// GetUserIdIfSuper 获取用户ID，如果是超级管理员则返回-1
func GetUserIdIfSuper(user models.User) int {
	if user.Role.Name == string(common.SuperAdminRoleId) {
		return -1 // 超级管理员返回-1
	}
	return user.ID // 返回用户ID
}

// GetUserInfo 从上下文中获取用户信息
func GetUserInfo(ctx fiber.Ctx) *models.User {
	user := ctx.Locals("user") // 从上下文中获取用户信息

	if user == nil {
		return nil // 用户信息不存在，返回nil
	}

	return user.(*models.User) // 返回用户信息
}
