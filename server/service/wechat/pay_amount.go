package wechat

import (
	"errors"
	"math"
)

// ExpectedPayFen 计算回调应核对的金额（单位：分）。
// debug=true 时与 JSAPIPay 一致，期望为 1 分。
func ExpectedPayFen(orderTotal float64, debug bool) (int64, error) {
	if debug {
		return 1, nil
	}
	if orderTotal < 0 || math.IsNaN(orderTotal) || math.IsInf(orderTotal, 0) {
		return 0, errors.New("订单金额异常")
	}
	// 与 JSAPIPay 中 fmt.Sprintf("%.0f", amount*100) 对齐
	return int64(math.Round(orderTotal * 100)), nil
}

// NotifyAmountMatches 判断微信回调 TotalFee（分）是否与订单应付一致。
func NotifyAmountMatches(orderTotal float64, totalFeeFen int64, debug bool) bool {
	expected, err := ExpectedPayFen(orderTotal, debug)
	if err != nil {
		return false
	}
	return expected == totalFeeFen
}
