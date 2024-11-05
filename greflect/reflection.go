package greflect

import (
	"reflect"

	"github.com/jinzhu/copier"
)

// Instance 创建一个指定类型 T 的零值实例。
// 参数:
//
//	T - 目标类型
//
// 返回值:
//
//	T - 指定类型 T 的零值实例
func Instance[T any]() T {
	resultType := reflect.TypeFor[T]()
	return InstanceOf(resultType).(T)
}

// InstanceOf 根据给定的 reflect.Type 创建一个实例。
// 参数:
//
//	resultType - 反射类型
//
// 返回值:
//
//	any - 根据 resultType 创建的实例
func InstanceOf(resultType reflect.Type) any {
	switch resultType.Kind() {
	case reflect.Ptr:
		return reflect.New(resultType.Elem()).Interface()
	case reflect.Slice:
		return reflect.MakeSlice(resultType, 0, 0).Interface()
	case reflect.Map:
		return reflect.MakeMap(resultType).Interface()
	case reflect.Chan:
		return reflect.MakeChan(resultType, 0).Interface()
	default:
		return reflect.New(resultType).Elem().Interface()
	}
}

// NilOrEmpty 创建一个指定类型 T 的零值实例。
// 参数:
//
//	T - 目标类型
//
// 返回值:
//
//	T - 指定类型 T 的零值实例
func NilOrEmpty[T any]() T {
	var v T
	return v
}

// Adapt 将源对象 src 复制到目标类型 T 的新实例中，并返回该实例和可能发生的错误。
// 参数:
//
//	src - 源对象
//
// 返回值:
//
//	T - 目标类型 T 的新实例
//	error - 复制过程中可能发生的错误
func Adapt[T any](src any) (T, error) {
	result := Instance[T]()
	err := copier.Copy(result, src)
	return result, err
}

// MustAdapt 将源对象 src 复制到目标类型 T 的新实例中，忽略可能发生的错误。
// 参数:
//
//	src - 源对象
//
// 返回值:
//
//	T - 目标类型 T 的新实例
func MustAdapt[T any](src any) T {
	result := Instance[T]()
	_ = copier.Copy(result, src)
	return result
}

// CanCallIsNil Can reflect.Value call reflect.Value.IsNil.
// It can avoid reflect.Value.IsNil panics.
func CanCallIsNil(v interface{}) bool {
	rv, ok := v.(reflect.Value)
	if !ok {
		return false
	}
	switch rv.Kind() {
	case reflect.Interface, reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}
