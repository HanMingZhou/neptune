package order

import (
	"fmt"
	"time"
)

// 根据时间+随机数生成订单号
func GenerateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf(
		"%s%06d",
		now.Format("20060102150405"),
		now.Nanosecond()/1000, // 微秒
	)
}
