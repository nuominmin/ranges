package ranges

import (
	"errors"
	"sort"
)

// Range 范围
type Range[T any] struct {
	Start uint64 // 包含
	Data  T
}

// ErrOverlappingRanges 用于表示区间重叠错误
var ErrOverlappingRanges = errors.New("config have overlapping ranges")

// ErrInvalidStart 用于表示起始块无效错误
var ErrInvalidStart = errors.New("start must be greater than 0")

// ErrRangeIsEmpty range 是空的
var ErrRangeIsEmpty = errors.New("range is empty")

// ErrDataNotFound 数据不存在
var ErrDataNotFound = errors.New("data not found for the given number")

type RangeProcessor[T any] struct {
	ranges []Range[T]
}

func NewRangeProcessor[T any]() *RangeProcessor[T] {
	return &RangeProcessor[T]{}
}

// AddRange 添加一个新的范围到 RangeProcessor 中
func (pb *RangeProcessor[T]) AddRange(r Range[T]) error {
	if r.Start == 0 {
		return ErrInvalidStart
	}

	pb.ranges = append(pb.ranges, r)

	sort.Slice(pb.ranges, func(i, j int) bool {
		// 根据 Start 降序
		return pb.ranges[i].Start > pb.ranges[j].Start
	})

	for i := 0; i < len(pb.ranges); i++ {
		if pb.ranges[i].Start == 0 {
			return ErrInvalidStart
		}
		if i > 0 && pb.ranges[i].Start == pb.ranges[i-1].Start {
			return ErrOverlappingRanges
		}
	}

	return nil
}

func (pb *RangeProcessor[T]) Handle(number uint64, handler func(data T) error) error {
	data, ok := pb.GetData(number)
	if !ok {
		return ErrDataNotFound
	}
	return handler(data)
}

func (pb *RangeProcessor[T]) GetData(number uint64) (data T, ok bool) {
	for i := 0; i < len(pb.ranges); i++ {
		if number >= pb.ranges[i].Start {
			return pb.ranges[i].Data, true
		}
	}
	return data, false
}
