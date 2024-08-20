package ranges

import (
	"fmt"
	"testing"
)

func TestRangeProcessor_Handle(t *testing.T) {
	// 2000 A
	// 3000 B
	// 4000 C

	// f(2500) B

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

}
