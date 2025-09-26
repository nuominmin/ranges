# ranges
一个通用的范围处理器。例如：根据某个输入值（如金额、时间、区块、等级...）找到对应的范围，并返回相应的数据

## 安装
你可以使用' go get '来安装这个包:

```sh
go get github.com/nuominmin/ranges
```

## 示例1 - 用户等级处理

``` go
type UserConfig struct {
    MaxItems int
    Level    string
}

// 使用 Builder 模式创建处理器
processor, err := ranges.NewRangeBuilder[UserConfig]().
    AddRange(100, UserConfig{MaxItems: 10, Level: "VIP"}).
    AddRange(50, UserConfig{MaxItems: 5, Level: "Premium"}).
    AddRange(1, UserConfig{MaxItems: 3, Level: "Basic"}).
    Build()

if err != nil {
    panic(err)
}

// 查找积分为 75 的用户配置
config, ok := processor.GetData(75)
if ok {
    fmt.Printf("用户等级: %s, 最大物品数: %d", config.Level, config.MaxItems)
    // 输出: 用户等级: Premium, 最大物品数: 5
}
```

## 示例2 - 区块处理
``` go
type Conf struct {
    WalletAddr   string
    ContractAddr string
    CodeVersion  string
}

// 使用 Builder 模式链式调用创建处理器
processor, err := ranges.NewRangeBuilder[Conf]().
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
    panic(err)
}

// 使用 Handle 方法处理不同的输入值
_ = processor.Handle(1000, func(data Conf) error {
    fmt.Printf("区块 1000: %+v\n", data)
    return nil
})

_ = processor.Handle(2500, func(data Conf) error {
    fmt.Printf("区块 2500: %+v\n", data)
    return nil
})

_ = processor.Handle(3000, func(data Conf) error {
    fmt.Printf("区块 3000: %+v\n", data)
    return nil
})

_ = processor.Handle(5000, func(data Conf) error {
    fmt.Printf("区块 5000: %+v\n", data)
    return nil
})
```

## 示例3 - 时间段处理
```go
// TimeToMinutes 将 HH:MM 转换为一天中的分钟数 (0~1439)
func TimeToMinutes(hour, minute int) int64 {
    return int64(hour*60 + minute)
}

// 使用 Builder 模式定义工作日时间段
processor, err := ranges.NewRangeBuilder[string]().
    AddRange(TimeToMinutes(6, 0), "早晨").   // 06:00
    AddRange(TimeToMinutes(9, 0), "上午").   // 09:00
    AddRange(TimeToMinutes(12, 0), "中午").  // 12:00
    AddRange(TimeToMinutes(13, 0), "下午").  // 13:00
    AddRange(TimeToMinutes(18, 0), "傍晚").  // 18:00
    AddRange(TimeToMinutes(21, 0), "深夜").  // 21:00
    WithCircular().
    Build()

if err != nil {
    panic(err)
}

// 测试几个时间点
tests := []string{"06:00", "08:30", "09:00", "12:00", "12:30", "19:00", "23:30", "02:00"}
for _, ts := range tests {
    t, _ := time.Parse("15:04", ts)
    minutes := TimeToMinutes(t.Hour(), t.Minute())

    if start, upperBound, data, ok := processor.GetDataWithRange(minutes); ok {
        fmt.Printf("%s → %s (范围: %d-%d)\n", ts, data, start, upperBound)
    } else {
        fmt.Printf("%s → 未找到时段\n", ts)
    }
}

// 输出结果:
// 06:00 → 早晨 (范围: 360-540)
// 08:30 → 早晨 (范围: 360-540)
// 09:00 → 上午 (范围: 540-720)
// 12:00 → 中午 (范围: 720-780)
// 12:30 → 中午 (范围: 720-780)
// 19:00 → 傍晚 (范围: 1080-1260)
// 23:30 → 深夜 (范围: 1260-360)
// 02:00 → 深夜 (范围: 1260-360)
```

