package service

import (
	"blog/internal/models"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DataBaseService 定义数据库服务接口
type DataBaseService interface {
	// 执行数据库SQL
	ExecDataBaseSQL(sql string) int64

	// 获取各种表的插入SQL
	GetBlogInsertSQL(page int) []string
	GetTagInsertSQL(page int) []string
	GetCategoryInsertSQL(page int) []string
	GetTopicInsertSQL(page int) []string
	GetUserInsertSQL(page int) []string
	GetRoleInsertSQL(page int) []string
	GetFileInsertSQL(page int) []string
	GetFileMd5InsertSQL(page int) []string
	GetBlogTagInsertSQL(page int) []string

	// 通用的获取插入SQL方法
	GetInsertSQL(tableName string, page int) []string
}

// ValueConverter 定义类型转换接口
type ValueConverter interface {
	Convert(value interface{}) string
}

// 各种类型转换器实现
type (
	IntConverter    struct{}
	FloatConverter  struct{}
	TimeConverter   struct{}
	StringConverter struct{}
	BoolConverter   struct{}
	NullConverter   struct{}
)

func (c IntConverter) Convert(number interface{}) string {
	if number == nil {
		return "null"
	}
	return fmt.Sprintf("%d", number)
}

func (c FloatConverter) Convert(number interface{}) string {
	if number == nil {
		return "null"
	}
	return fmt.Sprintf("%f", number)
}

func (c TimeConverter) Convert(t interface{}) string {
	if t == nil {
		return "null"
	}
	return fmt.Sprintf("'%s'", t.(time.Time).Format("2006-01-02 15:04:05"))
}

func (c StringConverter) Convert(t interface{}) string {
	if t == nil {
		return "null"
	}
	return fmt.Sprintf("'%s'", strings.ReplaceAll(t.(string), "'", "''"))
}

func (c BoolConverter) Convert(t interface{}) string {
	if t == nil {
		return "null"
	}
	return fmt.Sprintf("%t", t.(bool))
}

func (c NullConverter) Convert(_ interface{}) string {
	return "null"
}

// DataBaseServiceImpl 实现了 DataBaseService 接口
type DataBaseServiceImpl struct {
	converters map[string]ValueConverter
	db         *gorm.DB
}

// ExecDataBaseSQL 执行数据库 SQL 语句
func (d *DataBaseServiceImpl) ExecDataBaseSQL(sql string) int64 {
	result := d.db.Exec(sql)
	if result.Error != nil {
		logger.Error("执行 SQL 失败", zap.Error(result.Error))
		return 0
	}
	return result.RowsAffected
}

// getGlobalResult 获取全局查询结果并生成插入 SQL
func (d *DataBaseServiceImpl) getGlobalResult(tableName string, page int) []string {
	var records []map[string]interface{}

	result := d.db.Table(tableName).
		Offset((page - 1) * 50).
		Limit(50).
		Scan(&records)

	if result.Error != nil {
		logger.Error("查询数据失败",
			zap.String("table", tableName),
			zap.Error(result.Error))
		return nil
	}

	if len(records) == 0 {
		return nil
	}

	typeMap := d.createTypeMap(records[0])

	var arrays = make([]string, 0)

	for _, record := range records {
		columns, values := d.processRecord(record, typeMap)
		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
			tableName,
			strings.Join(columns, ", "),
			strings.Join(values, ", "))
		arrays = append(arrays, sql)
	}

	return arrays
}

func (d *DataBaseServiceImpl) createTypeMap(record map[string]interface{}) map[string]string {
	typeMap := make(map[string]string)
	for column, value := range record {
		typeMap[column] = fmt.Sprintf("%v", reflect.TypeOf(value))
	}
	return typeMap
}

func (d *DataBaseServiceImpl) processRecord(
	record map[string]interface{},
	typeMap map[string]string,
) ([]string, []string) {
	columns := make([]string, 0, len(record))
	values := make([]string, 0, len(record))

	for column, value := range record {
		columns = append(columns, column)

		converter, ok := d.converters[typeMap[column]]
		if !ok {
			converter = d.converters["<nil>"]
		}

		values = append(values, converter.Convert(value))
	}

	return columns, values
}

// GetInsertSQL 通用的获取插入 SQL 方法
func (d *DataBaseServiceImpl) GetInsertSQL(tableName string, page int) []string {
	return d.getGlobalResult(tableName, page)
}

// 为所有表格方法添加具体实现
func (d *DataBaseServiceImpl) GetBlogInsertSQL(page int) []string {
	return d.GetInsertSQL(models.BlogTable, page)
}

func (d *DataBaseServiceImpl) GetTagInsertSQL(page int) []string {
	return d.GetInsertSQL(models.TagTable, page)
}

func (d *DataBaseServiceImpl) GetCategoryInsertSQL(page int) []string {
	return d.GetInsertSQL(models.CategoryTable, page)
}

func (d *DataBaseServiceImpl) GetTopicInsertSQL(page int) []string {
	return d.GetInsertSQL(models.TopicTable, page)
}

func (d *DataBaseServiceImpl) GetUserInsertSQL(page int) []string {
	return d.GetInsertSQL(models.UserTable, page)
}

func (d *DataBaseServiceImpl) GetRoleInsertSQL(page int) []string {
	return d.GetInsertSQL(models.RoleTable, page)
}

func (d *DataBaseServiceImpl) GetFileInsertSQL(page int) []string {
	return d.GetInsertSQL(models.FileInfoTable, page)
}

func (d *DataBaseServiceImpl) GetFileMd5InsertSQL(page int) []string {
	return d.GetInsertSQL(models.FileInfoMd5Table, page)
}

func (d *DataBaseServiceImpl) GetBlogTagInsertSQL(page int) []string {
	return d.GetInsertSQL(models.BlogTagTable, page)
}

// NewDataBaseService 创建一个新的 DataBaseService 实例
func NewDataBaseService() DataBaseService {
	var once sync.Once
	var service *DataBaseServiceImpl

	once.Do(func() {
		converters := map[string]ValueConverter{
			"int":       IntConverter{},
			"int8":      IntConverter{},
			"int16":     IntConverter{},
			"int32":     IntConverter{},
			"int64":     IntConverter{},
			"uint":      IntConverter{},
			"uint8":     IntConverter{},
			"uint16":    IntConverter{},
			"uint32":    IntConverter{},
			"uint64":    IntConverter{},
			"float32":   FloatConverter{},
			"float64":   FloatConverter{},
			"bool":      BoolConverter{},
			"time.Time": TimeConverter{},
			"<nil>":     NullConverter{},
			"string":    StringConverter{},
		}

		service = &DataBaseServiceImpl{
			converters: converters,
			db:         configs.DB,
		}
	})
	return service
}
