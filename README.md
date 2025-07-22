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

