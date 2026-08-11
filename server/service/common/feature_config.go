package common

import (
	"strings"

	"fresh-shop/server/global"
	"fresh-shop/server/model/system"
)

// Feature flag keys stored in sys_config.name（M0-1 / industry-config）
const (
	KeyUserAudit   = "feature.user_audit"   // 1=启用用户审核；0=关闭（默认 B2C）
	KeySettleMonth = "feature.settle_month" // 1=允许月结；0=关闭
	KeyCourierMode = "feature.courier_mode" // delivery=城配；courier=快递
)

// FeatureEnabled 读取功能开关。缺键、禁用、非法值时返回 defaultEnabled。
func FeatureEnabled(key string, defaultEnabled bool) bool {
	raw, ok := readFeatureValue(key)
	if !ok {
		return defaultEnabled
	}
	enabled, parsed := parseFeatureBool(raw)
	if !parsed {
		return defaultEnabled
	}
	return enabled
}

// FeatureString 读取功能配置字符串。缺键或禁用时返回 defaultVal。
func FeatureString(key string, defaultVal string) string {
	raw, ok := readFeatureValue(key)
	if !ok {
		return defaultVal
	}
	return raw
}

func readFeatureValue(key string) (value string, ok bool) {
	if global.DB == nil {
		return "", false
	}
	var config system.SysConfig
	if err := global.DB.Where("name = ?", key).First(&config).Error; err != nil {
		return "", false
	}
	if config.Status == nil || *config.Status != 1 {
		return "", false
	}
	return strings.TrimSpace(config.Value), true
}

// parseFeatureBool 解析开关字符串。ok=false 表示无法识别。
func parseFeatureBool(value string) (enabled bool, ok bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "on":
		return true, true
	case "0", "false", "off":
		return false, true
	default:
		return false, false
	}
}
