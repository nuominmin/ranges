/*
 * @Description: 范围区间处理工具包
 * @Purpose: 用于处理基于数值范围的配置查找，支持泛型数据类型
 *
 * 主要功能：
 * 1. 根据起始值定义数据范围区间
 * 2. 支持降序排列的范围查找
 * 3. 根据输入数值快速定位对应的配置数据
 * 4. 提供完整的错误处理机制
 *
 * 使用场景：
 * - 用户等级配置（根据积分/经验值确定等级）
 * - 价格阶梯配置（根据数量确定单价）
 * - 权限配置（根据信用分确定权限等级）
 * - 合约地址配置（根据区块确定合约地址，钱包地址，去中心化代码版本）
 *
 * 使用示例参考：./processor_test.go
 */
package ranges

import (
	"errors"
	"sort"
)

// Range 范围区间定义
// T 为泛型类型，可以是任意数据类型
type Range[T any] struct {
	Start int64 // 区间起始值（包含该值），用于范围匹配
	Data  T     // 该区间对应的数据配置
}

// 错误定义

// ErrOverlappingRanges 区间重叠错误
// 当添加的范围与已存在的范围起始值相同时抛出
var ErrOverlappingRanges = errors.New("config have overlapping ranges")

// ErrInvalidStart 起始值无效错误
// 当范围起始值为负数时抛出，因为起始值必须大于0
var ErrInvalidStart = errors.New("start must be greater 0")

// ErrRangeIsEmpty 范围为空错误
// 当处理器中没有任何范围配置时抛出
var ErrRangeIsEmpty = errors.New("range is empty")

// ErrDataNotFound 数据未找到错误
// 当输入的数值无法匹配到任何范围时抛出
var ErrDataNotFound = errors.New("data not found for the given number")

// Processor 范围处理器接口
// 提供范围管理和数据查找的核心功能
type Processor[T any] interface {
	// AddRange 添加一个新的范围配置
	// 参数: r - 要添加的范围配置
	// 返回: error - 如果添加失败则返回错误信息
	AddRange(r Range[T]) error

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
	// - upperBound: 匹配的范围上界值（下一个区间的start-1，如果是最高区间则为-1表示无上界）
	// - data: 匹配范围的数据配置
	// - ok: 是否找到匹配的范围
	GetDataWithRange(number int64) (start int64, upperBound int64, data T, ok bool)
}

// processorImpl 范围处理器的具体实现
type processorImpl[T any] struct {
	ranges []Range[T] // 范围配置列表，按Start值降序排列
}

// NewProcessor 创建一个新的范围处理器实例
// 返回: Processor[T] - 范围处理器接口实例
func NewProcessor[T any]() Processor[T] {
	return &processorImpl[T]{}
}

// AddRange 添加一个新的范围到处理器中
// 实现逻辑：
// 1. 验证起始值必须大于0
// 2. 将新范围添加到列表中
// 3. 按Start值降序重新排序
// 4. 检查是否存在重复的起始值
func (pb *processorImpl[T]) AddRange(r Range[T]) error {
	// 验证起始值，必须大于0
	if r.Start < 0 {
		return ErrInvalidStart
	}

	// 添加新范围到列表
	pb.ranges = append(pb.ranges, r)

	// 按Start值降序排序，确保查找时从高到低匹配
	sort.Slice(pb.ranges, func(i, j int) bool {
		return pb.ranges[i].Start > pb.ranges[j].Start
	})

	// 检查重复的起始值（重叠范围）
	for i := 0; i < len(pb.ranges); i++ {
		if pb.ranges[i].Start < 0 {
			return ErrInvalidStart
		}
		// 检查相邻元素是否有相同的起始值
		if i > 0 && pb.ranges[i].Start == pb.ranges[i-1].Start {
			return ErrOverlappingRanges
		}
	}

	return nil
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
//
// 示例：ranges = [{Start:100}, {Start:50}, {Start:1}, {Start:0}]
// - 输入 150 -> 匹配 Start:100
// - 输入 75  -> 匹配 Start:50
// - 输入 25  -> 匹配 Start:1
// - 输入 0   -> 匹配 Start:0
func (pb *processorImpl[T]) GetData(number int64) (data T, ok bool) {
	// 遍历已排序的范围列表（降序）
	for i := 0; i < len(pb.ranges); i++ {
		// 找到第一个满足条件的范围：输入值 >= 起始值
		if number >= pb.ranges[i].Start {
			return pb.ranges[i].Data, true
		}
	}

	// 如果没有找到匹配的范围，返回零值和false
	return data, false
}

// GetDataWithRange 根据数值获取对应的数据配置，同时返回匹配的范围区间
// 返回值：
// - start: 匹配的范围起始值
// - upperBound: 匹配的范围上界值（最高区间则为-1表示无上界）
// - data: 匹配范围的数据配置
// - ok: 是否找到匹配的范围
func (pb *processorImpl[T]) GetDataWithRange(number int64) (start int64, upperBound int64, data T, ok bool) {
	// 遍历已排序的范围列表（降序）
	for i := 0; i < len(pb.ranges); i++ {
		// 找到第一个满足条件的范围：输入值 >= 起始值
		if number >= pb.ranges[i].Start {
			// 确定结束值
			var upperBound int64
			if i == len(pb.ranges)-1 {
				upperBound = -1 // 表示没有上界
			} else {
				upperBound = pb.ranges[i-1].Start
			}
			return pb.ranges[i].Start, upperBound, pb.ranges[i].Data, true
		}
	}

	// 如果没有找到匹配的范围，返回零值和false
	return 0, 0, data, false
}
