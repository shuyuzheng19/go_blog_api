package service

import (
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/pkg/common"
	"blog/pkg/configs"
	"fmt"
	"runtime"
	"time"

	"github.com/go-redis/redis"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"gorm.io/gorm"
)

// AddEyeCount 增加访问量记录
func AddEyeCount(count int64) error {
	today := time.Now().Format("2006-01-02") // 格式化为 "YYYY-MM-DD"

	var db = configs.DB

	var eyeView models.EyeView
	result := db.First(&eyeView, "id = ?", today)

	if result.Error != nil {
		// 如果找不到记录，创建新记录
		eyeView = models.EyeView{
			ID:    today,
			Count: count,
		}
		db.Create(&eyeView)
	} else {
		eyeView.Count = count
		db.Save(&eyeView)
	}
	return nil
}

// GetTodayTotalViews 计算今日所有博客的访问量总和
func GetTodayTotalViews() (int64, error) {
	var rdb = configs.REDIS
	// 从 Redis 获取当天总访问量
	totalViews, err := rdb.Get(common.EyeView).Int64()
	if err != nil {
		if err == redis.Nil {
			// 如果键不存在，说明今天还没有访问量，返回 0
			return 0, nil
		}
		// 其他错误返回
		return 0, err
	}

	return totalViews, nil
}

// func GetTodayTotalViews() (int64, error) {
// 	rdb := configs.REDIS

// 	// 使用管道获取 Redis 哈希中所有值，提升性能
// 	pipe := rdb.Pipeline()
// 	cmd := pipe.HVals(common.BlogEyeCountMapKey)
// 	_, err := pipe.Exec()
// 	if err != nil {
// 		return 0, fmt.Errorf("redis pipeline error: %v", err)
// 	}

// 	values, err := cmd.Result()
// 	if err != nil {
// 		return 0, fmt.Errorf("redis HVals error: %v", err)
// 	}

// 	var totalViews int64
// 	for _, value := range values {
// 		views, err := strconv.ParseInt(value, 10, 64)
// 		if err != nil {
// 			continue // 跳过无效值
// 		}
// 		totalViews += views
// 	}

// 	return totalViews, nil
// }

// Console 包含统计信息和系统信息
type Console struct {
	StatisticsInfo *response.Statistics `json:"statistics"`
	SystemInfo     *response.SystemInfo `json:"systemInfo"`
}

var systemInfoCache *response.SystemInfo

// GetConsole 获取控制台信息
func GetConsole() Console {
	statistics, _ := GetSystemStatistics()

	// 缓存系统信息，避免重复获取
	if systemInfoCache == nil {
		systemInfoCache, _ = GetSystemInfo()
	}

	return Console{
		StatisticsInfo: statistics,
		SystemInfo:     systemInfoCache,
	}
}

// GetSystemStatistics 获取系统统计信息
func GetSystemStatistics() (*response.Statistics, error) {
	db := configs.DB
	stats := &response.Statistics{}

	// 时间点
	now := time.Now()
	todayStart := getStartOfDay(now)
	weekAgo := now.AddDate(0, 0, -7).Unix()
	monthAgo := now.AddDate(0, -1, 0).Unix()

	// 获取统计数据
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Blog{}).Count(&stats.TotalBlogs).Error; err != nil {
			return err
		}

		// 总访问量（历史 + 今日）
		var historyViews int64
		if err := tx.Model(&models.Blog{}).Select("COALESCE(SUM(eye_count), 0)").Scan(&historyViews).Error; err != nil {
			return err
		}

		todayViews, err := GetTodayTotalViews()
		if err != nil {
			return err
		}

		stats.TotalViews = historyViews + todayViews
		stats.TodayViews = todayViews

		// 获取用户、分类、标签等统计
		if err := tx.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Category{}).Count(&stats.TotalCategories).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Tag{}).Count(&stats.TotalTags).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Topic{}).Count(&stats.TotalTopics).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.FileInfo{}).Count(&stats.TotalFiles).Error; err != nil {
			return err
		}

		// 今日新增用户和博客
		if err := tx.Model(&models.User{}).
			Where("created_at >= ?", todayStart).
			Count(&stats.TodayUsers).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Blog{}).
			Where("created_at >= ?", todayStart).
			Count(&stats.TodayBlogs).Error; err != nil {
			return err
		}

		// 周/月访问量
		if err := tx.Model(&models.EyeView{}).
			Where("updated_at >= ?", weekAgo).
			Select("COALESCE(SUM(count), 0)").
			Scan(&stats.WeeklyViews).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.EyeView{}).
			Where("updated_at >= ?", monthAgo).
			Select("COALESCE(SUM(count), 0)").
			Scan(&stats.MonthlyViews).Error; err != nil {
			return err
		}

		return nil
	})

	return stats, err
}

var upTime = time.Now().Unix()

// GetSystemInfo 获取系统信息
func GetSystemInfo() (*response.SystemInfo, error) {
	info := &response.SystemInfo{}

	// 主机信息
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("get host info error: %v", err)
	}

	info.OS = runtime.GOOS
	info.Platform = hostInfo.Platform
	info.PlatformVersion = hostInfo.PlatformVersion
	info.KernelVersion = hostInfo.KernelVersion
	info.Arch = runtime.GOARCH
	info.Hostname = hostInfo.Hostname
	info.BootTime = int64(hostInfo.BootTime)
	info.UpTime = upTime

	// CPU 信息
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("get cpu info error: %v", err)
	}

	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("get cpu percent error: %v", err)
	}

	info.CPUCores = runtime.NumCPU()
	if len(cpuPercent) > 0 {
		info.CPUUsage = cpuPercent[0]
	}
	if len(cpuInfo) > 0 {
		info.CPUModel = cpuInfo[0].ModelName
	}

	// 内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("get memory info error: %v", err)
	}

	info.TotalMemory = memInfo.Total
	info.UsedMemory = memInfo.Used
	info.FreeMemory = memInfo.Free
	info.MemoryUsage = memInfo.UsedPercent

	// Go 运行时信息
	info.GoVersion = runtime.Version()
	info.GoRoutines = runtime.NumGoroutine()

	return info, nil
}

// FormatBytes 格式化字节大小
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// getStartOfDay 获取某天的起始时间
func getStartOfDay(t time.Time) int64 {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
}
