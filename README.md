# ranges
一个通用的范围处理器。例如：可以根据某个输入值（如金额）找到对应的范围，并返回相应的数据

## 安装
你可以使用' go get '来安装这个包:

```sh
go get github.com/nuominmin/ranges
```

## 示例

``` go
type UserConfig struct {
    MaxItems int
    Level    string
}
 *
processor := ranges.NewProcessor[UserConfig]()
processor.AddRange(ranges.Range[UserConfig]{
    Start: 100,
    Data:  UserConfig{MaxItems: 10, Level: "VIP"},
})
processor.AddRange(ranges.Range[UserConfig]{
    Start: 50,
    Data:  UserConfig{MaxItems: 5, Level: "Premium"},
})
processor.AddRange(ranges.Range[UserConfig]{
    Start: 1,
    Data:  UserConfig{MaxItems: 3, Level: "Basic"},
})
 *
// 查找积分为 75 的用户配置
config, ok := processor.GetData(75)
if ok {
    fmt.Printf("用户等级: %s, 最大物品数: %d", config.Level, config.MaxItems)
    // 输出: 用户等级: Premium, 最大物品数: 5
}
```

``` go
  type Conf struct {
		WalletAddr   string
		ContractAddr string
		CodeVersion  string
	}

	ranges := []Range[Conf]{}
	ranges = append(ranges, Range[Conf]{Start: 3000, Data: Conf{
		WalletAddr: "BBBB", ContractAddr: "BBBB", CodeVersion: "v2",
	}})
	ranges = append(ranges, Range[Conf]{Start: 2000, Data: Conf{
		WalletAddr: "AAAA", ContractAddr: "AAAA", CodeVersion: "v1",
	}})
	ranges = append(ranges, Range[Conf]{Start: 4500, Data: Conf{
		WalletAddr: "CCCC", ContractAddr: "CCCC", CodeVersion: "v4",
	}})
	ranges = append(ranges, Range[Conf]{Start: 4000, Data: Conf{
		WalletAddr: "CCCC", ContractAddr: "CCCC", CodeVersion: "v3",
	}})

	p := NewRangeProcessor[Conf]()
	for i := 0; i < len(ranges); i++ {
		err := p.AddRange(ranges[i])
		if err != nil {
			t.Errorf("NewRangeProcessor[Conf] error: %s", err)
			return
		}
	}

	_ = p.Handle(1000, func(data Conf) error {
		fmt.Println(data)
		return nil
	})

	_ = p.Handle(2500, func(data Conf) error {
		fmt.Println(data)
		return nil
	})

	_ = p.Handle(3000, func(data Conf) error {
		fmt.Println(data)
		return nil
	})

	_ = p.Handle(5000, func(data Conf) error {
		fmt.Println(data)
		return nil
	})

```

```go
	// TimeToMinutes 将 HH:MM 转换为一天中的分钟数 (0~1439)
	func TimeToMinutes(hour, minute int) int64 {
		return int64(hour*60 + minute)
	}


	// 定义工作日时间段
	p := ranges.NewProcessor[string]()
	p.AddRange(ranges.Range[string]{Start: TimeToMinutes(6, 0), Data: "早晨"})   // 06:00
	p.AddRange(ranges.Range[string]{Start: TimeToMinutes(9, 0), Data: "上午"})   // 09:00
	p.AddRange(ranges.Range[string]{Start: TimeToMinutes(12, 0), Data: "中午"})  // 12:00
	p.AddRange(ranges.Range[string]{Start: TimeToMinutes(13, 0), Data: "下午"})  // 13:00
	p.AddRange(ranges.Range[string]{Start: TimeToMinutes(18, 0), Data: "傍晚"})  // 18:00
	p.AddRange(ranges.Range[string]{Start: TimeToMinutes(21, 0), Data: "深夜"})  // 21:00
	p.AddRange(ranges.Range[string]{Start: TimeToMinutes(0, 0), Data: "深夜"})   // 00:00–06:00

	// 测试几个时间点
	tests := []string{"06:00", "08:30", "09:00", "12:00", "12:30", "19:00", "23:30", "02:00"}
	for _, ts := range tests {
		t, _ := time.Parse("15:04", ts)
		minutes := TimeToMinutes(t.Hour(), t.Minute())

		if _, _, data, ok := p.GetDataWithRange(minutes); ok {
			fmt.Printf("%s → %s\n", ts, data)
		} else {
			fmt.Printf("%s → 未找到时段\n", ts)
		}
	}


输出结果
---
06:00 → 上午
08:30 → 上午
09:00 → 中午
12:00 → 中午
12:30 → 下午
19:00 → 深夜
23:30 → 深夜
02:00 → 深夜
```

