//    pipeline is a functional programming library for go
//    Copyright (C) 2015 mparaiso <mparaiso@online.fr>
//
//    pipeline program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    pipeline program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with pipeline program.  If not, see <http://www.gnu.org/licenses/>.

package pipeline

import (
	"fmt"
	"log"
	"reflect"
	"sort"
)

// Array is a place holder for interface{}
type Array interface{}

// VERSION of pipeline
const VERSION = "0.1"

/*********************************/
/*            PIPELINE           */
/*********************************/

// Pipeline allow sequential operations on slices, arrays or strings
type Pipeline struct {
	in       Array
	commands []func() (interface{}, error)
	current  interface{}
}

// Map send each element of a iterable through a function and return an array of results
func (pipeline *Pipeline) Map(callback func(interface{}, int) interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Map(pipeline.in, callback)
	})
	return pipeline
}

// Reduce folds the array into a single value
func (pipeline *Pipeline) Reduce(callback func(result interface{}, element interface{}, index int) interface{}, initialOrNil interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Reduce(pipeline.in, callback, initialOrNil)
	})
	return pipeline
}

// ReduceRight folds the array from end into a single value
func (pipeline *Pipeline) ReduceRight(callback func(result interface{}, element interface{}, index int) interface{}, initialOrNil interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return ReduceRight(pipeline.in, callback, initialOrNil)
	})
	return pipeline
}

// Sort sorts an array given a compare function
func (pipeline *Pipeline) Sort(compareFunc func(a, b interface{}) bool) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Sort(pipeline.in, compareFunc)
	})
	return pipeline
}

// Filter Iterates over elements of collection, returning a collection of all elements the predicate returns truthy for
func (pipeline *Pipeline) Filter(predicate func(element interface{}, index int) bool) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Filter(pipeline.in, predicate)
	})
	return pipeline
}

func (pipeline *Pipeline) Flatten() *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Flatten(pipeline.in)
	})
	return pipeline
}

// Compact remove nil values from array
func (pipeline *Pipeline) Compact() *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Compact(pipeline.in)
	})
	return pipeline
}

// Intersection creates a collection of unique values that are included in all
// of the provided collections.
func (pipeline *Pipeline) Intersection(arrays ...interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Intersection(append(append([]interface{}{}, pipeline.in), arrays...)...)
	})
	return pipeline
}

// IndexOf returns the index at which the first occurrence of element is found in array
// or -1 if the element is not found
func (pipeline *Pipeline) IndexOf(value interface{}, fromIndex int) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return IndexOf(pipeline.in, value, fromIndex)
	})
	return pipeline
}

// LastIndexOf method returns the last index at which a given element
// can be found in the array, or -1 if it is not present. The array is searched backwards, starting at fromIndex.
func (pipeline *Pipeline) LastIndexOf(value interface{}, fromIndex int) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return LastIndexOf(pipeline.in, value, fromIndex)
	})
	return pipeline
}

// Concat adds arrays to the end of the array and returns an new array
func (pipeline *Pipeline) Concat(arrays ...interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Concat(pipeline.in, arrays...)
	})
	return pipeline
}

// Zip creates an array of grouped elements,
// the first of which contains the first elements of the given arrays,
// the second of which contains the second elements of the given arrays, and so on.
func (pipeline *Pipeline) Zip() *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Zip(pipeline.in)
	})
	return pipeline
}

// Chunk Creates an array of elements split into groups the length of size. If collection can’t be split evenly, the final chunk will be the remaining elements.
func (pipeline *Pipeline) Chunk(length int) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Chunk(pipeline.in, length)
	})
	return pipeline
}

// Reverse reverse the order of the elements of the array and returns a new one
func (pipeline *Pipeline) Reverse() *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Reverse(pipeline.in)
	})
	return pipeline
}

// Some returns true if the callback predicate is satisfied
func (pipeline *Pipeline) Some(predicate func(element interface{}, index int) bool) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Some(pipeline.in, predicate)
	})
	return pipeline
}

// Push adds an element at the  end of the array
func (pipeline *Pipeline) Push(values ...interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Push(pipeline.in, values...)
	})
	return pipeline
}

// Unshift add an element at the beginning of a collection
func (pipeline *Pipeline) Unshift(values ...interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Unshift(pipeline.in, values...)
	})
	return pipeline
}

// Every returns true if the callback predicate is true for every element of the array
func (pipeline *Pipeline) Every(predicate func(element interface{}, index int) bool) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Every(pipeline.in, predicate)
	})
	return pipeline
}

// First returns the first element
func (pipeline *Pipeline) First() *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return First(pipeline.in)
	})
	return pipeline
}

// GroupBy Creates a map composed of keys generated
// from the results of running each element of collection through iteratee
func (pipeline *Pipeline) GroupBy(iteratee func(element interface{}, index int) interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return GroupBy(pipeline.in, iteratee)
	})
	return pipeline
}

// Op insert a custom operation in the pipeline
func (pipeline *Pipeline) Op(callback func(in interface{}) (interface{}, error)) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return callback(pipeline.in)
	})
	return pipeline
}

// Last returns the last element
func (pipeline *Pipeline) Last() *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Last(pipeline.in)
	})
	return pipeline
}

// Head returns the head until end
func (pipeline *Pipeline) Head(end int) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Head(pipeline.in, end)
	})
	return pipeline
}

// Tail returns the tail starting from start
func (pipeline *Pipeline) Tail(start int) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Tail(pipeline.in, start)
	})
	return pipeline
}

// ToMap takes a collection or a map and a callback, and returns a map[interface{}]interface{}
func (pipeline *Pipeline) ToMap(callback func(value interface{}, key interface{}) (resultValue interface{}, resultKey interface{})) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return ToMap(pipeline.in, callback)
	})
	return pipeline
}

// Slice returns a slice of an array
func (pipeline *Pipeline) Slice(start int, end int) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Slice(pipeline.in, start, end)
	})
	return pipeline
}

// Unique returns all the unique elements in a collection
func (pipeline *Pipeline) Unique() *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Unique(pipeline.in)
	})
	return pipeline
}

// Splice  deletes 'deleteCount' elements of an array from 'start' index
// and optionally inserts 'items'
func (pipeline *Pipeline) Splice(start int, deleteCount int, items ...interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Splice(pipeline.in, start, deleteCount, items...)
	})
	return pipeline
}

// Union returns an array filled by all unique values of the arrays
func (pipeline *Pipeline) Union(arrays ...interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Union(append(append([]interface{}{}, pipeline.in), arrays...)...)
	})
	return pipeline
}

// Difference returns a collection of the differences between 2 collections
func (pipeline *Pipeline) Difference(array interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Difference(pipeline.in, array)
	})
	return pipeline
}

// Without returns a collection without the values
func (pipeline *Pipeline) Without(values ...interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Without(pipeline.in, values...)
	})
	return pipeline
}

// Xor creates an array of unique values that is the symmetric difference of the provided arrays.
func (pipeline *Pipeline) Xor(arrays ...interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Xor(append(append([]interface{}{}, pipeline.in), arrays...)...)
	})
	return pipeline
}

// Out sets the output for the pipeline or return an error if an operation has failed
// output must be a pointer.
func (pipeline *Pipeline) Out(output interface{}) error {
	// output must be a pointer !
	if !isPointer(output) {
		return NotAPointerError{output}
	}
	// execute pipeline
	for i, command := range pipeline.commands {
		current, err := command()
		if err != nil {
			return StepError{i + 1, err}
		}
		pipeline.in = current

	}
	// first try
	if canAssignTo(pipeline.in, output) {
		valueOf(output).Elem().Set(valueOf(pipeline.in))
		return nil
	}
	// if in and out are maps let's try to match them
	if isMap(pipeline.in) && isMap(output) {
		maP := makeMapFrom(output)
		candidate := valueOf(pipeline.in)
		for _, key := range candidate.MapKeys() {
			if (candidate.MapIndex(key).Kind() == reflect.Ptr || candidate.MapIndex(key).Kind() == reflect.Interface) && canAssignTo(candidate.MapIndex(key).Elem(), maP) {
				maP.SetMapIndex(key.Elem(), candidate.MapIndex(key).Elem())
			} else if isSlice(candidate.MapIndex(key).Interface()) {
				val, err := convertSliceOfInterfaceToTypedSlice(candidate.MapIndex(key).Interface(), maP.Type().Elem())
				if err != nil {
					return err
				}
				maP.SetMapIndex(key.Elem(), val)
			}
		}
		pipeline.in = maP.Interface()
	}
	// If in and out are slices , let's try to match them
	if isSlice(pipeline.in) && isSlice(output) {
		arr := reflect.MakeSlice(reflect.TypeOf(output).Elem(), 0, 0)
		it := NewIterable(pipeline.in)
		for i := 0; i < it.Length(); i++ {
			val := it.At(i)
			// if arr is not of type []interface{} and []val cannot be assigned to arr
			if reflect.TypeOf(arr.Interface()) != reflect.TypeOf([]interface{}{}) && !reflect.SliceOf(reflect.TypeOf(val)).AssignableTo(reflect.TypeOf(arr.Interface())) {
				return CannotAppendError{arr.Interface(), val}
			}
			arr = reflect.Append(arr, valueOf(val))
		}
		pipeline.in = arr.Interface()
	}
	// if in can be assigned to out , do it
	if canAssignTo(pipeline.in, output) {
		valueOf(output).Elem().Set(valueOf(pipeline.in))
	} else {
		return CannotAssignError{pipeline.in, output}
	}
	return nil
}

// MustOut panics on error or returns the result of the pipeline
func (pipeline *Pipeline) MustOut() interface{} {
	var res interface{}
	err := pipeline.Out(&res)
	if err != nil {
		panic(err)
	}
	return res
}

// Equals returns true if all arrays are of equal length and Equal content
func (pipeline *Pipeline) Equals(arrays ...interface{}) *Pipeline {
	pipeline.commands = append(pipeline.commands, func() (interface{}, error) {
		return Equals(append(append([]interface{}{}, pipeline.in), arrays...)...)
	})
	return pipeline
}

// In Returns a new Pipeline
func In(sliceOrStringOrMap Array) *Pipeline {
	return &Pipeline{in: sliceOrStringOrMap, commands: []func() (interface{}, error){}}
}

// Must returns value or panics if err is not nil
func Must(value interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	} else {
		return value
	}
}

func IsString(value interface{}) bool {
	if _, ok := value.(string); ok {
		return true
	} else if _, ok2 := value.(*string); ok2 {
		return true
	}
	return false
}

// IsIterable returns true if value is iterable
func IsIterable(value interface{}) bool {
	arrayValue := reflect.ValueOf(value)
	switch arrayValue.Kind() {
	case reflect.Array, reflect.Slice, reflect.String, reflect.Map:
		return true
	default:
		switch value.(type) {
		case IterableInterface:
			return true
		default:
			return false
		}
	}

}

// Map is a map function
func Map(value interface{}, callback func(interface{}, int) interface{}) (interface{}, error) {
	if !IsIterable(value) {
		return nil, NotIterableError{value}
	}
	iterable := NewIterable(value)
	var result []interface{}

	for i := 0; i < iterable.Length(); i++ {
		val := iterable.At(i)
		result = append(result, callback(val, i))
	}
	return result, nil
}

// Reduce folds the array into a single value
func Reduce(value interface{}, callback func(result interface{}, element interface{}, index int) interface{}, initialOrNil interface{}) (interface{}, error) {
	if !IsIterable(value) {
		return nil, NotIterableError{value}
	}
	iterable := NewIterable(value)
	var result interface{}
	if initialOrNil != nil {
		result = initialOrNil
		for i := 0; i < iterable.Length(); i++ {
			result = callback(result, iterable.At(i), i)
		}
	} else {
		result = iterable.At(0)
		for i := 1; i < iterable.Length(); i++ {
			result = callback(result, iterable.At(i), i)
		}
	}

	return result, nil
}

// ReduceRight folds the array into a single value,from the last value to the first value
func ReduceRight(array interface{}, callback func(result interface{}, element interface{}, index int) interface{}, initialOrnil interface{}) (interface{}, error) {
	array, Error := Reverse(array)
	if Error != nil {
		return nil, Error
	}
	return Reduce(array, callback, initialOrnil)
}

func Flatten(array interface{}) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iterable := NewIterable(array)
	result := []interface{}{}
	for i := 0; i < iterable.Length(); i++ {
		if !IsString(iterable.At(i)) && IsIterable(iterable.At(i)) {
			candidate := NewIterable(iterable.At(i))
			for j := 0; j < candidate.Length(); j++ {
				result = append(result, candidate.At(j))
			}
		} else {
			result = append(result, iterable.At(i))
		}
	}
	return result, nil
}

// Filter Iterates over elements of collection, returning a collection of all elements the predicate returns truthy for
func Filter(array interface{}, predicate func(element interface{}, index int) bool) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	result := []interface{}{}
	iterable := NewIterable(array)
	for i := 0; i < iterable.Length(); i++ {
		if predicate(iterable.At(i), i) {
			result = append(result, iterable.At(i))
		}
	}
	return result, nil
}

// IndexOf returns the index at which the first occurrence of element is found in array
// or -1 if the element is not found
func IndexOf(array interface{}, searchedElement interface{}, fromIndex int) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iterable := NewIterable(array)
	for i := fromIndex; i < iterable.Length(); i++ {
		if iterable.At(i) == searchedElement {
			return i, nil
		}
	}
	return -1, nil
}

// LastIndexOf method returns the last index at which a given element
// can be found in the array, or -1 if it is not present. The array is searched backwards, starting at fromIndex.
func LastIndexOf(array interface{}, searchElement interface{}, fromIndex int) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iterable := NewIterable(array)
	for i := iterable.Length() - 1; i >= 0; i-- {
		if iterable.At(i) == searchElement {
			return i, nil
		}
	}

	return -1, nil
}

// Concat adds arrays to the end of the array and returns an new array
func Concat(array interface{}, arrays ...interface{}) (interface{}, error) {
	if len(arrays) == 0 {
		return array, nil
	}
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}

	result := []interface{}{}
	initial := NewIterable(array)
	for i := 0; i < initial.Length(); i++ {
		result = append(result, initial.At(i))
	}
	for _, array := range arrays {
		if !IsIterable(array) {
			return nil, NotIterableError{array}
		}
		iter := NewIterable(array)
		for i := 0; i < iter.Length(); i++ {
			result = append(result, iter.At(i))
		}
	}
	return result, nil
}

// Reverse reverse the order of the elements of the array and returns a new one
func Reverse(array interface{}) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iterable := NewIterable(array)
	result := []interface{}{}
	for i := iterable.Length() - 1; i >= 0; i-- {
		result = append(result, iterable.At(i))
	}
	return result, nil
}

// Some returns true if the callback predicate is satisfied
func Some(array interface{}, predicate func(v interface{}, index int) bool) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iterable := NewIterable(array)
	for i := 0; i < iterable.Length(); i++ {
		if predicate(iterable.At(i), i) {
			return true, nil
		}
	}
	return false, nil
}

// Every returns true if the callback predicate is true for every element of the array
func Every(array interface{}, predicate func(v interface{}, index int) bool) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iterable := NewIterable(array)
	for i := 0; i < iterable.Length(); i++ {
		if !predicate(iterable.At(i), i) {
			return false, nil
		}
	}
	return true, nil
}

// First returns the first element of an array
func First(array interface{}) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iter := NewIterable(array)
	if iter.Length() < 1 {
		return nil, nil
	}
	return iter.At(0), nil
}

// Last returns the last element of an array
func Last(array interface{}) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iter := NewIterable(array)
	if iter.Length() < 1 {
		return nil, nil
	}
	return iter.At(iter.Length() - 1), nil
}

// Head returns the head until end
func Head(array interface{}, endIndex int) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iterable := NewIterable(array)
	if endIndex >= iterable.Length() {
		return nil, IndexOutOfBoundsError{endIndex}
	}
	if endIndex < 0 {
		return nil, IndexOutOfBoundsError{endIndex}
	}
	result := []interface{}{}

	for i := 0; i <= endIndex; i++ {
		result = append(result, iterable.At(i))
	}
	return result, nil

}

// Tail returns the tail starting from start
func Tail(array interface{}, startIndex int) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iterable := NewIterable(array)
	if startIndex >= iterable.Length() {
		return nil, IndexOutOfBoundsError{startIndex}
	}
	if startIndex < 0 {
		return nil, IndexOutOfBoundsError{startIndex}
	}
	result := []interface{}{}
	for i := startIndex; i < iterable.Length(); i++ {
		result = append(result, iterable.At(i))
	}
	return result, nil

}

// Slice returns a slice of an array
func Slice(array interface{}, start int, end int) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	result, err := Head(array, end)
	if err != nil {
		return nil, err
	}

	result, err = Tail(result, start)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Splice  deletes 'deleteCount' elements of an array from 'start' index
// and optionally inserts 'items'
func Splice(array interface{}, start int, deleteCount int, items ...interface{}) (interface{}, error) {
	head, err := Head(array, start-1)
	if err != nil {
		return nil, err
	}
	tail, err := Tail(array, start+deleteCount)
	if err != nil {
		return nil, err
	}
	return Concat(head, items, tail)
}

// Unique filters remove duplicate values from an array
func Unique(array interface{}) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	result := []interface{}{}
	iter := NewIterable(array)
	for i := 0; i < iter.Length(); i++ {
		if index, err := IndexOf(result, iter.At(i), 0); err == nil && index.(int) < 0 {
			result = append(result, iter.At(i))
		} else if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// Without creates an array excluding all provided values
func Without(array interface{}, values ...interface{}) (interface{}, error) {
	return Difference(array, values)
}

// Difference creates an array excluding all provided values
func Difference(array interface{}, values interface{}) (interface{}, error) {
	var Error error
	val, err := Filter(array, func(element interface{}, index int) bool {
		val, err := IndexOf(values, element, 0)
		if err != nil {
			Error = err
		}
		return val.(int) < 0
	})
	if Error == nil && err != nil {
		return nil, err
	}
	return val, nil
}

// Push adds an element at the  end of the array
func Push(array interface{}, values ...interface{}) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iter := NewIterable(array)
	result := []interface{}{}
	for i := 0; i < iter.Length(); i++ {
		result = append(result, iter.At(i))
	}
	result = append(result, values...)
	return result, nil
}

// Unshift add an element at the beginning of a collection
func Unshift(array interface{}, values ...interface{}) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iterable := NewIterable(array)
	return append(append([]interface{}{}, values...), iterable.ToArrayOfInterface()...), nil
}

// Sort sorts an array given a compare function
func Sort(array interface{}, compareFunc func(a, b interface{}) bool) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	iterable := NewIterable(array)
	sorter := &sorter{iterable, compareFunc}
	sort.Sort(sorter)
	return sorter.array.ToArrayOfInterface(), nil
}

// Chunk Creates an array of elements split into groups the length of size. If collection can’t be split evenly, the final chunk will be the remaining elements.
func Chunk(array interface{}, length int) (interface{}, error) {
	if !IsIterable(array) {
		return nil, NotIterableError{array}
	}
	if (length) < 0 {
		return nil, IndexOutOfBoundsError{length}
	}
	arrayType := reflect.TypeOf(array)

	iterable := NewIterable(array)
	result := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(array)), 0, 0)
	for i := 0; i < iterable.Length(); i++ {
		if (i % (length)) == 0 {
			result = reflect.Append(result, reflect.MakeSlice(arrayType, 0, 0))
		}
		result.Index(result.Len() - 1).Set(reflect.Append(result.Index(result.Len()-1), reflect.ValueOf(iterable.At(i))))
	}
	return result.Interface(), nil
}

// Union returns an array filled by all unique values of the arrays
func Union(arrays ...interface{}) (interface{}, error) {
	for _, array := range arrays {
		if !IsIterable(array) {
			return nil, NotIterableError{array}
		}
	}
	result := []interface{}{}
	for _, array := range arrays {
		iterable := NewIterable(array)
		for i := 0; i < iterable.Length(); i++ {
			val, err := IndexOf(result, iterable.At(i), 0)
			if err != nil {
				return nil, err
			}
			if val.(int) == -1 {
				result = append(result, iterable.At(i))
			}
		}
	}
	return result, nil
}

// Xor creates an array of unique values that is the symmetric difference of the provided arrays.
func Xor(arrays ...interface{}) (interface{}, error) {
	for _, array := range arrays {
		if !IsIterable(array) {
			return nil, NotIterableError{array}
		}
	}
	switch len(arrays) {
	case 0:
		return nil, nil
	case 1:
		return []interface{}{}, nil
	case 2:
		return Must(Difference(Must(Union(arrays...)), Must(Intersection(arrays...)))), nil
	default:
		return Must(Xor(append(append([]interface{}{}, Must(Xor(arrays[:2]...)).([]interface{})), arrays[2:]...)...)), nil
	}
}

// Zip Creates an array of grouped elements,
// the first of which contains the first elements of the given arrays,
// the second of which contains the second elements of the given arrays, and so on.
func Zip(arrayS interface{}) (interface{}, error) {
	if !IsIterable(arrayS) {
		return nil, NotIterableError{arrayS}
	}
	arrays := NewIterable(arrayS)
	if arrays.Length() == 0 {
		return []interface{}{}, nil
	}
	for i := 0; i < arrays.Length(); i++ {
		if !IsIterable(arrays.At(i)) {
			return nil, NotIterableError{arrays.At(i)}
		}
	}
	if arrays.Length() == 1 {
		result := [][]interface{}{}
		iterable := NewIterable(arrays.At(0))
		for i := 0; i < iterable.Length(); i++ {
			result = append(result, []interface{}{})
			result[i] = append(result[i], iterable.At(i))
		}
		return result, nil
	}
	result := [][]interface{}{}

	maxLength := func() int {
		l := 0
		for i := 0; i < arrays.Length(); i++ {
			iterable := NewIterable(arrays.At(i))
			if iterable.Length() > l {
				l = iterable.Length()
			}
		}
		return l
	}()
	for i := 0; i < maxLength; i++ {
		result = append(result, []interface{}{})
	}
	for i := 0; i < maxLength; i++ {
		for j := 0; j < arrays.Length(); j++ {
			iterable := NewIterable(arrays.At(j))
			if i < iterable.Length() {
				result[i] = append(result[i], iterable.At(i))
			} else {
				result[i] = append(result[i], nil)
			}
		}
	}
	return result, nil
}

// Compact remove nil values from array
func Compact(array interface{}) (interface{}, error) {
	return Filter(array, func(el interface{}, i int) bool {
		return el != nil
	})

}

// Equals returns true if all arrays are of equal length and equal content
func Equals(arrays ...interface{}) (interface{}, error) {
	iterables := []IterableInterface{}
	for _, array := range arrays {
		if !IsIterable(array) {
			return nil, NotIterableError{array}
		}
	}
	for _, array := range arrays {
		iterables = append(iterables, NewIterable(array))
	}
	return equal(iterables...), nil
}

// Equal returns true if all arrays are of equal length and equal content
func equal(arrays ...IterableInterface) bool {

	switch len(arrays) {
	case 0, 1:
		return true
	case 2:
		if arrays[0].Length() != arrays[1].Length() {
			return false
		}
		for i := 0; i < arrays[0].Length(); i++ {
			if arrays[0].At(i) != arrays[1].At(i) {
				return false
			}
		}
		return true
	default:
		return arrays[0].Length() == arrays[1].Length() && arrays[0].Length() == arrays[2].Length() && equal(arrays[0], arrays[1]) && equal(arrays[:2]...)
	}
}

// Intersection creates a collection of unique values that are included in all
// of the provided collections.
func Intersection(arrays ...interface{}) (interface{}, error) {
	switch len(arrays) {
	case 0:
		return nil, nil
	case 1:
		return arrays[0], nil
	default:
		iterables := []IterableInterface{}
		for _, array := range arrays {
			if !IsIterable(array) {
				return nil, NotIterableError{array}
			}
			iterables = append(iterables, NewIterable(Must(Unique(array))))
		}
		first := iterables[0]
		iterables = iterables[1:]
		return Filter(first, func(element interface{}, i int) bool {
			return Must(Every(iterables, func(iterable interface{}, i int) bool {
				return Must(IndexOf(iterable, element, 0)).(int) >= 0
			})).(bool)
		})
	}
}

// GroupBy creates an object composed of keys generated from the results of running each element of collection through iteratee
func GroupBy(collection interface{}, iteratee func(interface{}, int) interface{}) (interface{}, error) {
	result := map[interface{}][]interface{}{}
	if !IsIterable(collection) {
		return nil, NotIterableError{collection}
	}
	iterable := NewIterable(collection)
	for i := 0; i < iterable.Length(); i++ {
		group := iteratee(iterable.At(i), i)
		if result[group] == nil {
			result[group] = []interface{}{}
		}
		result[group] = append(result[group], iterable.At(i))
	}
	return result, nil
}

// ToMap takes a collection or a map and a callback, and returns a map[interface{}]interface{}
func ToMap(mapOrSlice interface{}, mapper func(value interface{}, key interface{}) (valueResult interface{}, keyResult interface{})) (interface{}, error) {
	if !IsIterable(mapOrSlice) {
		return nil, NotIterableError{mapOrSlice}
	}
	m := map[interface{}]interface{}{}
	if !isMap(mapOrSlice) {
		iterable := NewIterable(mapOrSlice)
		for i := 0; i < iterable.Length(); i++ {
			m[i] = iterable.At(i)
		}
	} else {
		v := reflect.ValueOf(mapOrSlice)
		k := v.MapKeys()
		for i := 0; i < v.Len(); i++ {
			m[k[i].Interface()] = v.MapIndex(k[i]).Interface()
		}
	}
	result := map[interface{}]interface{}{}
	for k, v := range m {
		valueResult, keyResult := mapper(v, k)
		result[keyResult] = valueResult
	}
	return result, nil
}

/*********************************/
/*           ITERABLE            */
/*********************************/

// IterableInterface represents an value that can be iterated on
type IterableInterface interface {
	Length() int
	At(index int) interface{}
	ToArrayOfInterface() []interface{}
}

// Iterable implements IterableInterface
type Iterable struct {
	array  reflect.Value
	length int
	isMap  bool
	keys   []reflect.Value
}

// NewIterable returns a new iterable
func NewIterable(array interface{}) IterableInterface {
	if a, ok := array.(IterableInterface); ok == true {
		return a
	}
	switch t := array.(type) {
	case string:
		res := []rune{}
		for _, char := range t {
			res = append(res, char)
		}
		return &Iterable{array: reflect.ValueOf(res), length: len(res)}
	default:
		arr := reflect.ValueOf(array)

		return &Iterable{array: arr, length: arr.Len(), isMap: arr.Kind() == reflect.Map}
	}

}

// String returns a string representation
func (iterable Iterable) String() string {
	return fmt.Sprintf("<Iterable %#v>", iterable.array.Interface())
}

// Length is the length
func (iterable Iterable) Length() int {
	return iterable.length
}

// At is the value at index
func (iterable *Iterable) At(index int) interface{} {
	if iterable.isMap {
		if iterable.keys == nil {
			iterable.keys = iterable.array.MapKeys()
			log.Print(iterable.keys)
		}
		return iterable.array.MapIndex(iterable.keys[index]).Interface()
	}
	return iterable.array.Index(index).Interface()
}

// ToArrayOfInterface returns []interface{}
func (iterable Iterable) ToArrayOfInterface() []interface{} {
	result := []interface{}{}
	for i := 0; i < iterable.Length(); i++ {
		result = append(result, iterable.At(i))
	}
	return result
}

/*********************************/
/*            ERRORS             */
/*********************************/

// NotIterableError discriminates a value that cannot be iterated on
type NotIterableError struct {
	value interface{}
}

// Error returns a string
func (err NotIterableError) Error() string {
	return fmt.Sprintf("%#v is not iterable", err.value)
}

// IndexOutOfBoundsError discriminates an index out of bounds error
type IndexOutOfBoundsError struct {
	index int
}

// Error returns a string
func (indexOutOfBoundsError IndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("Index out of bounds %d", indexOutOfBoundsError.index)
}

// CannotAppendError discriminates a not appendable error
type CannotAppendError struct {
	array interface{}
	value interface{}
}

// Error returns a string
func (cannotAppendError CannotAppendError) Error() string {
	return fmt.Sprintf("Cannot append value %#v to %#v .", cannotAppendError.value, cannotAppendError.array)
}

// CannotAssignError discriminates a non assignable value error
type CannotAssignError struct {
	from interface{}
	to   interface{}
}

// Error returns a string
func (cannotAssignError CannotAssignError) Error() string {
	return fmt.Sprintf("Cannot assign the result of the pipeline %#v to output %#v .", cannotAssignError.from, cannotAssignError.to)
}

// StepError discriminates a step error
type StepError struct {
	step   int
	reason interface{}
}

// Error returns a string
func (stepError StepError) Error() string {
	return fmt.Sprintf("Error at step %d : %#v ", stepError.step, stepError.reason)
}

// NotAPointerError discriminate pointer errors
type NotAPointerError struct {
	value interface{}
}

// Error returns a string
func (notAPointerError NotAPointerError) Error() string {
	return fmt.Sprintf(" %#v should be a pointer ", notAPointerError.value)
}

// sorter is used for array.Sort
type sorter struct {
	array       IterableInterface
	compareFunc func(a, b interface{}) bool
}

// Len returns the length of the array
func (s *sorter) Len() int {
	return s.array.Length()
}

// Less compare 2 elements. if i is less than j return true ,else return false
func (s *sorter) Less(i, j int) bool {
	return s.compareFunc(s.array.At(i), s.array.At(j))
}

// Swap swaps 2 elements
func (s *sorter) Swap(i, j int) {
	a := s.array.ToArrayOfInterface()
	a[i], a[j] = a[j], a[i]

	s.array = NewIterable(a)
}

/*********************************/
/*             HELPERS           */
/*********************************/

func convertSliceOfInterfaceToTypedSlice(from interface{}, to reflect.Type) (reflect.Value, error) {
	arr := reflect.MakeSlice(to, 0, 0)
	it := NewIterable(from)
	for i := 0; i < it.Length(); i++ {
		val := it.At(i)
		// if arr is not of type []interface{} and []val cannot be assigned to arr
		if reflect.TypeOf(arr.Interface()) != reflect.TypeOf([]interface{}{}) && !reflect.SliceOf(reflect.TypeOf(val)).AssignableTo(reflect.TypeOf(arr.Interface())) {
			return valueOf(nil), CannotAppendError{arr.Interface(), val}
		}
		arr = reflect.Append(arr, valueOf(val))
	}
	return arr, nil
}

func valueOf(in interface{}) reflect.Value {
	return reflect.ValueOf(in)
}

func isMap(in interface{}) bool {
	return reflect.TypeOf(in).Kind() == reflect.Map || (reflect.TypeOf(in).Kind() == reflect.Ptr && reflect.TypeOf(in).Elem().Kind() == reflect.Map)
}

func isSlice(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Slice || (reflect.TypeOf(value).Kind() == reflect.Ptr && reflect.TypeOf(value).Elem().Kind() == reflect.Slice)
}

func isPointer(in interface{}) bool {
	return reflect.TypeOf(in).Kind() == reflect.Ptr
}

func canAssignTo(in, out interface{}) bool {
	return reflect.TypeOf(in).AssignableTo(reflect.TypeOf(out)) || reflect.TypeOf(in).AssignableTo(reflect.TypeOf(out).Elem())
}

func makeMapFrom(from interface{}) reflect.Value {
	if isPointer(from) {
		return reflect.MakeMap(reflect.TypeOf(from).Elem())
	}
	return reflect.MakeMap(reflect.TypeOf(from))

}
