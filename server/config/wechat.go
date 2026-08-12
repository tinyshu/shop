package config

type Wechat struct {
	Appid          string `mapstructure:"appid" json:"appid" yaml:"appid"`
	Secret         string `mapstructure:"secret" json:"secret" yaml:"secret"`
	RandomPassword string `mapstructure:"randomPassword" json:"randomPassword" yaml:"randomPassword"`
}

type WechatPay struct {
	MchId      string                 `mapstructure:"mchId" json:"mchId" yaml:"mchId"`             // 商户号
	ApiV2Key   string                 `mapstructure:"apiV2Key" json:"apiV2Key" yaml:"apiV2Key"`    // API v2 密钥
	NotifyURL  string                 `mapstructure:"notifyUrl" json:"notifyUrl" yaml:"notifyUrl"` // 微信支付通知地址
	Debug      bool                   `mapstructure:"debug" json:"debug" yaml:"debug"`
	CertPath   string                 `mapstructure:"certPath" json:"certPath" yaml:"certPath"`
	KeyPath    string                 `mapstructure:"keyPath" json:"keyPath" yaml:"keyPath"`
	Compensate WechatPayCompensate    `mapstructure:"compensate" json:"compensate" yaml:"compensate"`
}

// WechatPayCompensate 掉单查单 / 定时补偿配置（PAY-01）
type WechatPayCompensate struct {
	Enable           bool   `mapstructure:"enable" json:"enable" yaml:"enable"`                               // 定时扫描总开关；测试建议 false
	Spec             string `mapstructure:"spec" json:"spec" yaml:"spec"`                                     // cron，如 @every 5m
	MinAgeMinutes    int    `mapstructure:"minAgeMinutes" json:"minAgeMinutes" yaml:"minAgeMinutes"`           // 创建超过 N 分钟的待付单才扫
	BatchSize        int    `mapstructure:"batchSize" json:"batchSize" yaml:"batchSize"`                       // 每轮最多查单数
	MaxQueryPerOrder int    `mapstructure:"maxQueryPerOrder" json:"maxQueryPerOrder" yaml:"maxQueryPerOrder"` // 单笔进程内最多查微信次数，0=不限制
	Mock             bool   `mapstructure:"mock" json:"mock" yaml:"mock"`                                     // true 不调微信，视为未支付
	AdminSyncEnable  bool   `mapstructure:"adminSyncEnable" json:"adminSyncEnable" yaml:"adminSyncEnable"`     // 管理端一键同步开关
}
