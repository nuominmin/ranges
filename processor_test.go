package ranges_test

import (
	"fmt"
	"github.com/nuominmin/ranges"
	"testing"
)

func TestRangeProcessor_Handle(t *testing.T) {
	// 测试基本范围功能
	// 2000 A
	// 3000 B
	// 4000 C

	// f(2500) B

	type Conf struct {
		WalletAddr   string
		ContractAddr string
		CodeVersion  string
	}

	// 使用新的 Builder API
	builder := ranges.NewRangeBuilder[Conf]()

	// 链式调用添加范围
	processor, err := builder.
		AddRange(3000, Conf{
			WalletAddr: "BBBB", ContractAddr: "BBBB", CodeVersion: "v2",
		}).
		AddRange(2000, Conf{
			WalletAddr: "AAAA", ContractAddr: "AAAA", CodeVersion: "v1",
		}).
		AddRange(4500, Conf{
			WalletAddr: "CCCC", ContractAddr: "CCCC", CodeVersion: "v4",
		}).
		AddRange(4000, Conf{
			WalletAddr: "CCCC", ContractAddr: "CCCC", CodeVersion: "v3",
		}).
		Build()

	if err != nil {
		t.Errorf("Build error: %s", err)
		return
	}

	_ = processor.Handle(1000, func(data Conf) error {
		fmt.Println("1000:", data)
		return nil
	})

	_ = processor.Handle(2500, func(data Conf) error {
		fmt.Println("2500:", data)
		return nil
	})

	_ = processor.Handle(3000, func(data Conf) error {
		fmt.Println("3000:", data)
		return nil
	})

	_ = processor.Handle(5000, func(data Conf) error {
		fmt.Println("5000:", data)
		return nil
	})
}

// TestTimeRangeProcessor 测试时间段范围处理器
// 测试场景：
// 早晨：06:00 - 09:00 → start 范围：360 - 540
// 上午：09:00 - 12:00 → start 范围：540 - 720
// 中午：12:00 - 13:00 → start 范围：720 - 780
// 下午：13:00 - 18:00 → start 范围：780 - 1080
// 傍晚：18:00 - 21:00 → start 范围：1080 - 1260
// 深夜：21:00 - 06:00 → start 范围：1260 - 1440 或 0 - 360
func TestTimeRangeProcessor(t *testing.T) {
	type TimeConfig struct {
		Period      string
		Description string
	}

	// 使用新的 Builder API
	builder := ranges.NewRangeBuilder[TimeConfig]()

	// 链式调用添加各个时间段，深夜时段使用多个起始值
	p, err := builder.
		AddRange(360, TimeConfig{Period: "早晨", Description: "06:00 - 09:00"}).
		AddRange(540, TimeConfig{Period: "上午", Description: "09:00 - 12:00"}).
		AddRange(720, TimeConfig{Period: "中午", Description: "12:00 - 13:00"}).
		AddRange(780, TimeConfig{Period: "下午", Description: "13:00 - 18:00"}).
		AddRange(1080, TimeConfig{Period: "傍晚", Description: "18:00 - 21:00"}).
		AddRange(0, TimeConfig{Period: "深夜", Description: "21:00 - 06:00"}).    // 深夜时段有两个起始值
		AddRange(1260, TimeConfig{Period: "深夜", Description: "21:00 - 06:00"}). // 深夜时段有两个起始值
		Build()

	if err != nil {
		t.Errorf("Build error: %s", err)
		return
	}

	// 测试各个时间点
	testCases := []struct {
		minute   int64
		expected string
	}{
		{100, "深夜"},  // 01:40 (深夜时段，匹配起始值0)
		{400, "早晨"},  // 06:40 (早晨时段)
		{600, "上午"},  // 10:00 (上午时段)
		{750, "中午"},  // 12:30 (中午时段)
		{900, "下午"},  // 15:00 (下午时段)
		{1200, "傍晚"}, // 20:00 (傍晚时段)
		{1300, "深夜"}, // 21:40 (深夜时段，匹配起始值1260)
		{30, "深夜"},   // 00:30 (深夜时段，匹配起始值0)
	}

	for _, tc := range testCases {
		data, ok := p.GetData(tc.minute)
		if !ok {
			t.Errorf("GetData(%d) failed: no data found", tc.minute)
			continue
		}
		if data.Period != tc.expected {
			t.Errorf("GetData(%d) = %s, expected %s", tc.minute, data.Period, tc.expected)
		} else {
			fmt.Printf("✓ 时间 %d 分钟 -> %s (%s)\n", tc.minute, data.Period, data.Description)
		}
	}

	// 测试 GetDataWithRange
	fmt.Println("\n测试 GetDataWithRange:")
	testMinute := int64(100) // 01:40
	start, upperBound, data, ok := p.GetDataWithRange(testMinute)
	if ok {
		fmt.Printf("时间 %d 分钟 -> 起始值: %d, 上界: %d, 时段: %s\n", testMinute, start, upperBound, data.Period)
	}

	testMinute = int64(1300) // 21:40
	start, upperBound, data, ok = p.GetDataWithRange(testMinute)
	if ok {
		fmt.Printf("时间 %d 分钟 -> 起始值: %d, 上界: %d, 时段: %s\n", testMinute, start, upperBound, data.Period)
	}
}
