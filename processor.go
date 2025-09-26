package ranges

import (
	"errors"
	"sort"
)

// rangeItem 范围区间定义
type rangeItem[T any] struct {
	start int64 // 区间起始值（包含该值），用于范围匹配
	data  T     // 该区间对应的数据配置
}

// 错误定义

var (
	// ErrOverlappingRanges 区间重叠错误
	// 当添加的范围与已存在的范围起始值相同时抛出
	ErrOverlappingRanges = errors.New("config have overlapping ranges")

	// ErrInvalidStart 起始值无效错误
	// 当范围起始值为负数时抛出，因为起始值必须大于等于0
	ErrInvalidStart = errors.New("start must be greater than or equal to 0")

	// ErrRangeIsEmpty 范围为空错误
	// 当处理器中没有任何范围配置时抛出
	ErrRangeIsEmpty = errors.New("range is empty")

	// ErrDataNotFound 数据未找到错误
	// 当输入的数值无法匹配到任何范围时抛出
	ErrDataNotFound = errors.New("data not found for the given number")
)

// RangeBuilder 范围构建器接口
// 用于添加范围配置，采用Builder模式
type RangeBuilder[T any] interface {
	// AddRange 添加一个新的范围配置
	// 参数: start - 要添加的起始值
	//       data - 该起始值对应的数据配置
	AddRange(start int64, data T) RangeBuilder[T]

	// WithCircular 设置为循环模式
	WithCircular() RangeBuilder[T]

	// Build 构建最终的范围处理器
	// 执行验证、排序等操作，返回可用的处理器
	// 返回: processor - 范围处理器实例
	//       error - 如果构建失败则返回错误信息
	Build() (Processor[T], error)
}

// Processor 范围处理器接口
// 提供范围管理和数据查找的核心功能
type Processor[T any] interface {
	// Handle 处理指定数值，找到对应数据后执行处理函数
	// 参数: number - 要查找的数值
	//       handler - 数据处理函数
	// 返回: error - 如果查找失败或处理函数执行失败则返回错误
	Handle(number int64, handler func(data T) error) error

	// GetData 根据数值获取对应的数据配置
	// 参数: number - 要查找的数值
	// 返回: data - 匹配的数据配置
	//       ok - 是否找到匹配的配置
	GetData(number int64) (data T, ok bool)

	// GetDataWithRange 根据数值获取对应的数据配置，同时返回匹配的范围区间
	// 返回值：
	// - start: 匹配的范围起始值
	// - upperBound: 匹配的范围上界值（最高区间则为-1表示无上界）
	// - data: 匹配范围的数据配置
	// - ok: 是否找到匹配的范围
	GetDataWithRange(number int64) (start int64, upperBound int64, data T, ok bool)
}

// rangeBuilderImpl 范围构建器的具体实现
type rangeBuilderImpl[T any] struct {
	circular bool
	ranges   []rangeItem[T] // 范围配置列表，构建阶段未排序
}

// processorImpl 范围处理器的具体实现
type processorImpl[T any] struct {
	circular bool           // 是否是循环模式
	ranges   []rangeItem[T] // 范围配置列表，已排序且验证过
}

/*
循环模式说明：
当为循环模式时，范围的结束值为范围的起始值

开启循环模式示例：
- ranges = [{Start:100}, {Start:50}, {Start:1}, {Start:0}], number = -1, 则返回 Start:100
- ranges = [{Start:100}, {Start:50}, {Start:2}, {Start:1}], number = 0, 则返回 Start:100

未开启循环模式示例：
- ranges = [{Start:100}, {Start:50}, {Start:1}, {Start:0}]
  - 输入 150 -> 匹配 Start:100
  - 输入 75  -> 匹配 Start:50
  - 输入 25  -> 匹配 Start:1
  - 输入 0   -> 匹配 Start:0
  - 输入 -1  -> 未匹配
*/

// NewRangeBuilder 创建一个新的范围构建器实例
func NewRangeBuilder[T any]() RangeBuilder[T] {
	return &rangeBuilderImpl[T]{}
}

// AddRange 添加一个新的范围到处理器中
func (rb *rangeBuilderImpl[T]) AddRange(start int64, data T) RangeBuilder[T] {
	rb.ranges = append(rb.ranges, rangeItem[T]{
		start: start,
		data:  data,
	})
	return rb
}

// WithCircular 设置为循环模式
func (rb *rangeBuilderImpl[T]) WithCircular() RangeBuilder[T] {
	rb.circular = true
	return rb
}

// Build 构建最终的范围处理器
// 执行验证、排序等操作，返回可用的处理器
func (rb *rangeBuilderImpl[T]) Build() (Processor[T], error) {
	// 检查是否有范围数据
	if len(rb.ranges) == 0 {
		return nil, ErrRangeIsEmpty
	}

	// 按start值降序排序，确保查找时从高到低匹配
	sort.Slice(rb.ranges, func(i, j int) bool {
		return rb.ranges[i].start > rb.ranges[j].start
	})

	// 检查重复的起始值（重叠范围）
	for i := 0; i < len(rb.ranges); i++ {
		if rb.ranges[i].start < 0 {
			return nil, ErrInvalidStart
		}
		// 检查相邻元素是否有相同的起始值
		if i > 0 && rb.ranges[i].start == rb.ranges[i-1].start {
			return nil, ErrOverlappingRanges
		}
	}

	// 创建处理器实例并复制范围
	processor := &processorImpl[T]{
		circular: rb.circular,
		ranges:   make([]rangeItem[T], len(rb.ranges)),
	}
	copy(processor.ranges, rb.ranges)
	return processor, nil
}

// Handle 查找指定数值对应的数据并执行处理函数
// 这是GetData方法的高级封装，支持链式处理
func (pb *processorImpl[T]) Handle(number int64, handler func(data T) error) error {
	data, ok := pb.GetData(number)
	if !ok {
		return ErrDataNotFound
	}
	return handler(data)
}

// GetData 根据输入数值查找对应的范围配置数据
// 查找逻辑：
// 1. 遍历降序排列的范围列表
// 2. 找到第一个 number >= Start 的范围
// 3. 返回该范围的数据配置
func (pb *processorImpl[T]) GetData(number int64) (data T, ok bool) {
	_, _, data, ok = pb.GetDataWithRange(number)
	return data, ok
}

// GetDataWithRange 根据数值获取对应的数据配置，同时返回匹配的范围区间
// 返回值：
// - start: 匹配的范围起始值
// - upperBound: 匹配的范围上界值（最高区间则为-1表示无上界）
// - data: 匹配范围的数据配置
// - ok: 是否找到匹配的范围
func (pb *processorImpl[T]) GetDataWithRange(number int64) (start int64, upperBound int64, data T, ok bool) {
	// 没有配置，直接返回
	if len(pb.ranges) == 0 {
		return 0, 0, data, false
	}

	// 遍历已排序的范围列表（降序）
	for i := 0; i < len(pb.ranges); i++ {
		// 找到第一个满足条件的范围：输入值 >= 起始值
		if number >= pb.ranges[i].start {
			if i == 0 {
				upperBound = -1 // 最高区间，表示没有上界
			} else {
				upperBound = pb.ranges[i-1].start // 前一个区间的起始值作为当前区间的上界
			}
			return pb.ranges[i].start, upperBound, pb.ranges[i].data, true
		}
	}

	// 如果没有找到匹配的范围
	if pb.circular {
		// 循环模式下
		last := pb.ranges[0]
		first := pb.ranges[len(pb.ranges)-1]
		return last.start, first.start, last.data, true
	}
	return 0, 0, data, false
}
