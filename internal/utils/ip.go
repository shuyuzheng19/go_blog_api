package utils

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// GetIPAddress 从请求中获取客户端的IP地址
func GetIPAddress(request fiber.Ctx) string {
	ipAddress := request.Get("X-Forwarded-For") // 尝试获取X-Forwarded-For头部

	// 如果没有获取到IP地址，尝试其他头部
	if ipAddress == "" || strings.ToLower(ipAddress) == "unknown" {
		ipAddress = request.Get("Proxy-Client-IP")
	}
	if ipAddress == "" || strings.ToLower(ipAddress) == "unknown" {
		ipAddress = request.Get("WL-Proxy-Client-IP")
	}
	if ipAddress == "" || strings.ToLower(ipAddress) == "unknown" {
		ipAddress = request.IP() // 最后尝试获取直接的IP地址
	}
	return ipAddress
}

func GetIpAndCitp(ctx fiber.Ctx) (string, string) {
	var ip = GetIPAddress(ctx)
	var city = GetIpCity(ip)
	return ip, city
}

var searcher *xdb.Searcher // IP数据库搜索器

// LoadIpDB 加载IP数据库
func LoadIpDB(dbPath string) {
	var err error
	searcher, err = xdb.NewWithFileOnly(dbPath) // 使用文件路径创建搜索器

	if err != nil {
		panic("加载IP数据库失败: " + err.Error()) // 加载失败时抛出异常
	}
}

// GetIpCity 根据IP地址获取城市信息
func GetIpCity(ip string) string {
	region, err := searcher.SearchByStr(ip) // 根据IP地址查询地区信息
	if err != nil {
		return "未知" // 查询失败时返回"未知"
	}

	// 解析地区信息
	var split = strings.Split(region, "|")
	return strings.ReplaceAll(split[0]+" "+split[2]+" "+split[3], "0", "") // 返回城市信息
}

// GetClientPlatformInfo 获取客户端平台信息
func GetClientPlatformInfo(userAgent string) string {
	if userAgent == "" {
		return "" // 如果用户代理为空，返回空字符串
	}

	userAgent = strings.ToLower(userAgent) // 转为小写以便于匹配

	var os, browser string
	// 匹配操作系统
	switch {
	case strings.Contains(userAgent, "windows"):
		os = "Windows"
	case strings.Contains(userAgent, "mac"):
		os = "Mac"
	case strings.Contains(userAgent, "android"):
		os = "Android"
	case strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad"):
		os = "iOS"
	}
	// 匹配浏览器
	switch {
	case strings.Contains(userAgent, "micromessenger"):
		browser = "微信客户端"
	case strings.Contains(userAgent, "edg"):
		browser = "Edge"
	case strings.Contains(userAgent, "chrome"):
		browser = "Chrome"
	case strings.Contains(userAgent, "firefox"):
		browser = "Firefox"
	case strings.Contains(userAgent, "safari"):
		browser = "Safari"
	}

	// 返回操作系统和浏览器信息
	if os != "" && browser != "" {
		return fmt.Sprintf("%s %s", os, browser)
	} else {
		return userAgent // 如果未匹配到，返回原始用户代理
	}
}
