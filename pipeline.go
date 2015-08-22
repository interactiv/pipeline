package pipeline

import (
	"fmt"
	"reflect"
	"sort"
)

type Array interface{}
type Result interface{}
type Value interface{}

/*********************************/
/*            ITERATOR           */
/*********************************/

type IterableInterface interface {
	Length() int
	At(index int) interface{}
	ToArrayOfInterface() []interface{}
}

type Iterable struct {
	array  reflect.Value
	length int
}

func NewIterable(array Value) IterableInterface {
	if a, ok := array.(IterableInterface); ok == true {
		return a
	} else {
		switch t := array.(type) {
		case string:
			res := []rune{}
			for _, char := range t {
				res = append(res, char)
			}
			return &Iterable{array: reflect.ValueOf(res), length: len(res)}
		default:
			arr := reflect.ValueOf(array)
			return &Iterable{array: arr, length: arr.Len()}
		}
	}
}

func (this Iterable) Length() int {
	return this.length
}

func (this Iterable) At(index int) interface{} {
	return this.array.Index(index).Interface()
}

func (this Iterable) ToArrayOfInterface() []interface{} {
	result := []interface{}{}
	for i := 0; i < this.Length(); i++ {
		result = append(result, this.At(i))
	}
	return result
}

/*********************************/
/*            PIPELINE           */
/*********************************/
type Pipeline struct {
	in       Array
	commands []func() (Value, error)
	current  Value
}

func (this *Pipeline) Map(callback func(interface{}, int) interface{}) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Map(this.in, callback)
	})
	return this
}

// Reduce folds the array into a single value
func (this *Pipeline) Reduce(callback func(result interface{}, element interface{}, index int) interface{}, initialOrNil interface{}) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Reduce(this.in, callback, initialOrNil)
	})
	return this
}

func (this *Pipeline) ReduceRight(callback func(result interface{}, element interface{}, index int) interface{}, initialOrNil interface{}) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return ReduceRight(this.in, callback, initialOrNil)
	})
	return this
}

// Sort sorts an array given a compare function
func (this *Pipeline) Sort(compareFunc func(a, b interface{}) bool) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Sort(this.in, compareFunc)
	})
	return this
}

func (this *Pipeline) Filter(predicate func(element interface{}, index int) bool) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Filter(this.in, predicate)
	})
	return this
}

func (this *Pipeline) IndexOf(value Value, fromIndex int) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return IndexOf(this.in, value, fromIndex)
	})
	return this
}

// The lastIndexOf() method returns the last index at which a given element
// can be found in the array, or -1 if it is not present. The array is searched backwards, starting at fromIndex.
func (this *Pipeline) LastIndexOf(value Value, fromIndex int) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return LastIndexOf(this.in, value, fromIndex)
	})
	return this
}

// Concat adds arrays to the end of the array and returns an new array
func (this *Pipeline) Concat(arrays ...Value) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Concat(this.in, arrays...)
	})
	return this
}

// Chunk Creates an array of elements split into groups the length of size. If collection can’t be split evenly, the final chunk will be the remaining elements.
func (this *Pipeline) Chunk(length int) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Chunk(this.in, length)
	})
	return this
}

// Reverse reverse the order of the elements of the array and returns a new one
func (this *Pipeline) Reverse() *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Reverse(this.in)
	})
	return this
}

// Some returns true if the callback predicate is satisfied
func (this *Pipeline) Some(predicate func(element interface{}, index int) bool) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Some(this.in, predicate)
	})
	return this
}

func (this *Pipeline) Push(values ...interface{}) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Push(this.in, values...)
	})
	return this
}

func (this *Pipeline) Unshift(values ...interface{}) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Unshift(this.in, values...)
	})
	return this
}

// Every returns true if the callback predicate is true for every element of the array
func (this *Pipeline) Every(predicate func(element interface{}, index int) bool) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Every(this.in, predicate)
	})
	return this
}

func (this *Pipeline) First() *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return First(this.in)
	})
	return this
}

func (this *Pipeline) Last() *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Last(this.in)
	})
	return this
}

func (this *Pipeline) Head(end int) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Head(this.in, end)
	})
	return this
}
func (this *Pipeline) Tail(start int) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Tail(this.in, start)
	})
	return this
}

// Slice returns a slice of an array
func (this *Pipeline) Slice(start int, end int) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Slice(this.in, start, end)
	})
	return this
}

func (this *Pipeline) Unique() *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Unique(this.in)
	})
	return this
}

// Splice  deletes 'deleteCount' elements of an array from 'start' index
// and optionally inserts 'items'
func (this *Pipeline) Splice(start int, deleteCount int, items ...interface{}) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Splice(this.in, start, deleteCount, items...)
	})
	return this
}

func (this *Pipeline) Difference(array Value) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Difference(this.in, array)
	})
	return this
}

func (this *Pipeline) Without(values ...Value) *Pipeline {
	this.commands = append(this.commands, func() (Value, error) {
		return Without(this.in, values...)
	})
	return this
}

// Out sets the output for the pipeline or return an error if an operation has failed
func (this *Pipeline) Out(output interface{}) error {
	// output must be a pointer to something

	if reflect.TypeOf(output).Kind() != reflect.Ptr {
		return fmt.Errorf("%v should be a pointer ", output)
	}
	for i, command := range this.commands {
		current, err := command()
		if err != nil {
			return fmt.Errorf("Error at step %d : %s ", i+1, err)
		} else {
			this.in = current
		}
	}
	if reflect.TypeOf(this.in).Kind() == reflect.Slice && reflect.TypeOf(output).Elem().Kind() == reflect.Slice {
		// if output is a pointer of a slice and if this.in is also a pointer of a slice then set the output accordingly
		arr := reflect.MakeSlice(reflect.TypeOf(output).Elem(), 0, 0)
		it := NewIterable(this.in)
		for i := 0; i < it.Length(); i++ {
			val := it.At(i)
			arr = reflect.Append(arr, reflect.ValueOf(val))
		}
		reflect.ValueOf(output).Elem().Set(arr)
	} else if reflect.TypeOf(this.in).AssignableTo(reflect.TypeOf(output).Elem()) {
		// if ouput points to the same type as this.in , set output to this.in
		reflect.ValueOf(output).Elem().Set(reflect.ValueOf(this.in))
	} else {
		return fmt.Errorf("Cannot assign the result of the pipeline %v to output %v", this.in, output)
	}

	return nil
}

// In Returns a new Pipeline
func In(array Array) *Pipeline {
	return &Pipeline{in: array, commands: []func() (Value, error){}}
}

// IsArray returns true if value is iterable
func IsArray(value Value) bool {
	arrayValue := reflect.ValueOf(value)
	switch arrayValue.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
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

// map_ is a map function
func Map(value Value, callback func(interface{}, int) interface{}) (Value, error) {
	if !IsArray(value) {
		return nil, NotIterable(value)
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
func Reduce(value Value, callback func(result interface{}, element interface{}, index int) interface{}, initialOrNil interface{}) (Value, error) {
	if !IsArray(value) {
		return nil, NotIterable(value)
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

// Reduce folds the array into a single value,from the last value to the first value
func ReduceRight(array Value, callback func(result interface{}, element interface{}, index int) interface{}, initialOrnil interface{}) (Value, error) {
	array, Error := Reverse(array)
	if Error != nil {
		return nil, Error
	}
	return Reduce(array, callback, initialOrnil)
}

func Filter(array Value, predicate func(element interface{}, index int) bool) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
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

func IndexOf(array Value, searchedElement Value, fromIndex int) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
	}
	iterable := NewIterable(array)
	for i := fromIndex; i < iterable.Length(); i++ {
		if iterable.At(i) == searchedElement {
			return i, nil
		}
	}
	return -1, nil
}

// The lastIndexOf() method returns the last index at which a given element
// can be found in the array, or -1 if it is not present. The array is searched backwards, starting at fromIndex.
func LastIndexOf(array Value, searchElement interface{}, fromIndex int) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
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
func Concat(array Value, arrays ...Value) (Value, error) {
	if len(arrays) == 0 {
		return array, nil
	}
	if !IsArray(array) {
		return nil, NotIterable(array)
	}

	result := []interface{}{}
	initial := NewIterable(array)
	for i := 0; i < initial.Length(); i++ {
		result = append(result, initial.At(i))
	}
	for _, array := range arrays {
		if !IsArray(array) {
			return nil, NotIterable(array)
		}
		iter := NewIterable(array)
		for i := 0; i < iter.Length(); i++ {
			result = append(result, iter.At(i))
		}
	}
	return result, nil
}

// Reverse reverse the order of the elements of the array and returns a new one
func Reverse(array Value) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
	}
	iterable := NewIterable(array)
	result := []interface{}{}
	for i := iterable.Length() - 1; i >= 0; i-- {
		result = append(result, iterable.At(i))
	}
	return result, nil
}

// Some returns true if the callback predicate is satisfied
func Some(array Value, predicate func(v interface{}, index int) bool) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
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
func Every(array Value, predicate func(v interface{}, index int) bool) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
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
func First(array Value) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
	}
	iter := NewIterable(array)
	if iter.Length() < 1 {
		return nil, nil
	}
	return iter.At(0), nil
}

// Last returns the last element of an array
func Last(array Value) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
	}
	iter := NewIterable(array)
	if iter.Length() < 1 {
		return nil, nil
	}
	return iter.At(iter.Length() - 1), nil
}

func Head(array Value, endIndex int) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
	}
	iterable := NewIterable(array)
	if endIndex >= iterable.Length() {
		return nil, IndexOutOfBounds(endIndex)
	}
	if endIndex < 0 {
		return nil, IndexOutOfBounds(endIndex)
	}
	result := []interface{}{}

	for i := 0; i <= endIndex; i++ {
		result = append(result, iterable.At(i))
	}
	return result, nil

}

func Tail(array Value, startIndex int) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
	}
	iterable := NewIterable(array)
	if startIndex >= iterable.Length() {
		return nil, IndexOutOfBounds(startIndex)
	}
	if startIndex < 0 {
		return nil, IndexOutOfBounds(startIndex)
	}
	result := []interface{}{}
	for i := startIndex; i < iterable.Length(); i++ {
		result = append(result, iterable.At(i))
	}
	return result, nil

}

// Slice returns a slice of an array
func Slice(array Value, start int, end int) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
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
func Splice(array Value, start int, deleteCount int, items ...interface{}) (Value, error) {
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
func Unique(array Value) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
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
func Without(array Value, values ...Value) (Value, error) {
	return Difference(array, values)
}

// Difference creates an array excluding all provided values
func Difference(array Value, values Value) (Value, error) {
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

func Push(array Value, values ...interface{}) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
	}
	iter := NewIterable(array)
	result := []interface{}{}
	for i := 0; i < iter.Length(); i++ {
		result = append(result, iter.At(i))
	}
	result = append(result, values...)
	return result, nil
}

func Unshift(array Value, values ...interface{}) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
	}
	iterable := NewIterable(array)
	return append(append([]interface{}{}, values...), iterable.ToArrayOfInterface()...), nil
}

// Sort sorts an array given a compare function
func Sort(array Value, compareFunc func(a, b interface{}) bool) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
	}
	iterable := NewIterable(array)
	sorter := &sorter{iterable, compareFunc}
	sort.Sort(sorter)
	return sorter.array.ToArrayOfInterface(), nil
}

// Chunk Creates an array of elements split into groups the length of size. If collection can’t be split evenly, the final chunk will be the remaining elements.
func Chunk(array Value, length int) (Value, error) {
	if !IsArray(array) {
		return nil, NotIterable(array)
	}
	if (length) < 0 {
		return nil, IndexOutOfBounds(length)
	}
	arrayType := reflect.TypeOf(array)

	iterable := NewIterable(array)
	result := reflect.ValueOf([]reflect.Value{})
	for i := 0; i < iterable.Length(); i++ {
		if (i % (length)) == 0 {
			result = reflect.Append(result, reflect.MakeSlice(arrayType, 0, 0))
		}
		result.Index(result.Len() - 1).Set(reflect.Append(result.Index(result.Len()-1), reflect.ValueOf(iterable.At(i))))
	}
	return result.Interface(), nil
}

/*********************************/
/*            ERRORS             */
/*********************************/

// NotIterable returns a NotIterableError
func NotIterable(value Value) error {
	return NotIterableError{fmt.Errorf("%v is not iterable", value)}
}

type NotIterableError struct {
	e error
}

func (err NotIterableError) Error() string {
	return err.e.Error()
}

func IndexOutOfBounds(index int) error {
	return IndexOutOfBoundsError{index}
}

type IndexOutOfBoundsError struct {
	index int
}

func (this IndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("Index out of bounds %d", this.index)
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
