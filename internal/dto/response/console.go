package response

type Statistics struct {
	TotalBlogs      int64 `json:"totalBlogs"`
	TotalViews      int64 `json:"totalViews"`
	TotalUsers      int64 `json:"totalUsers"`
	TotalCategories int64 `json:"totalCategories"`
	TotalTags       int64 `json:"totalTags"`
	TotalTopics     int64 `json:"totalTopics"`
	TotalFiles      int64 `json:"totalFiles"`
	TodayViews      int64 `json:"todayViews"`
	TodayUsers      int64 `json:"todayUsers"`
	TodayBlogs      int64 `json:"todayBlogs"`
	WeeklyViews     int64 `json:"weeklyViews"`
	MonthlyViews    int64 `json:"monthlyViews"`
}

type SystemInfo struct {
	// 操作系统信息
	OS              string `json:"os"`              // 操作系统
	Platform        string `json:"platform"`        // 平台
	PlatformVersion string `json:"platformVersion"` // 平台版本
	KernelVersion   string `json:"kernelVersion"`   // 内核版本
	Arch            string `json:"arch"`            // 系统架构
	Hostname        string `json:"hostname"`        // 主机名
	BootTime        int64  `json:"bootTime"`        // 启动时间
	UpTime          int64  `json:"upTime"`          //程序启动时间

	// CPU信息
	CPUCores int     `json:"cpuCores"` // CPU核心数
	CPUUsage float64 `json:"cpuUsage"` // CPU使用率
	CPUModel string  `json:"cpuModel"` // CPU型号

	// 内存信息
	TotalMemory uint64  `json:"totalMemory"` // 总内存
	UsedMemory  uint64  `json:"usedMemory"`  // 已用内存
	FreeMemory  uint64  `json:"freeMemory"`  // 空闲内存
	MemoryUsage float64 `json:"memoryUsage"` // 内存使用率
	// Go运行时信息
	GoVersion  string `json:"goVersion"`  // Go版本
	GoRoutines int    `json:"goRoutines"` // 协程数量
}
