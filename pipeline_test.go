//    pipeline is a functional programming library for go
//    Copyright (C) 2015 mparaiso <mparaiso@online.fr>

//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.

//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.

//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <http://www.gnu.org/licenses/>.

package pipeline_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/interactiv/expect"
	"github.com/interactiv/pipeline"
)

type Person struct {
	Age  int
	Name string
}

type Product struct {
	ID         int
	Name       string
	CategoryID int
	Price      int
}

type Products []Product

func TestIntersection(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{1, 2, 4}).Intersection([]int{3, 2, 1}, []int{2, 5, 6}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(fmt.Sprint(result)).ToEqual(fmt.Sprint([]int{2}))
}

func TestGroupBy(t *testing.T) {

	e := expect.New(t)

	var result map[int][]Person
	err := pipeline.In([]Person{{12, "John"}, {12, "Jane"}, {20, "Joe"}}).
		GroupBy(func(el interface{}, i int) interface{} {
		return el.(Person).Age
	}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(2)
	e.Expect(len(result[12])).ToEqual(2)
	e.Expect(result[12][0].Age).ToEqual(12)

	var result1 map[string][]map[string]string
	err = pipeline.In([]map[string]string{
		{"product": "trousers", "category": "clothes"},
		{"product": "beer", "category": "drinks"},
		{"product": "coat", "category": "clothes"},
	}).
		GroupBy(func(el interface{}, i int) interface{} {
		return el.(map[string]string)["category"]
	}).Out(&result1)
	e.Expect(err).ToBeNil()
	e.Expect(len(result1)).ToEqual(2)
	e.Expect(len(result1["clothes"])).ToEqual(2)
	e.Expect(len(result1["drinks"])).ToEqual(1)
}

func ExamplePipeline_Xor() {
	var result []int
	err := pipeline.In([]int{1, 2}).Xor([]int{2, 3}).Out(&result)
	fmt.Print(result, " ", err)
	// Output: [1 3] <nil>
}

func TestMap(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{1, 2, 3}).Map(func(element interface{}, index int) interface{} {
		return element.(int) * 2
	}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(3)
	e.Expect(result[0]).ToEqual(2)

	var result2 *[]int
	err = pipeline.In(&[]int{1, 2, 3}).Out(&result2)
	e.Expect(err).ToBeNil()
	e.Expect(fmt.Sprint(result2)).ToEqual(fmt.Sprint(&[]int{1, 2, 3}))

	var result3 string
	err = pipeline.In("foo").Out(&result3)
	e.Expect(err).ToBeNil()
	e.Expect(result3).ToBe("foo")

	var result4 []string
	err = pipeline.In([]int{1, 2, 3}).Out(&result4)
	e.Expect(err).Not().ToBeNil()
	t.Log(err)
}

func TestMapReduce(t *testing.T) {
	e := expect.New(t)
	var result int
	err := pipeline.In([]int{1, 2, 3}).Map(func(element interface{}, index int) interface{} {
		return element.(int) * 2
	}).Reduce(func(result interface{}, element interface{}, index int) interface{} {
		return result.(int) + element.(int)
	}, 0).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToEqual(12)
}

func TestReduceRight(t *testing.T) {
	e := expect.New(t)
	var result string
	err := pipeline.In("kayak").ReduceRight(func(r interface{}, e interface{}, i int) interface{} {
		return r.(string) + string(e.(rune))
	}, "").Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToEqual("kayak")
	var result2 int
	pipeline.In([]int{1, 2, 3}).ReduceRight(func(r interface{}, e interface{}, i int) interface{} {
		return r.(int) - e.(int)
	}, nil).Out(&result2)
	e.Expect(result2).ToEqual(0)
}

func TestFilter(t *testing.T) {
	e := expect.New(t)
	var result []string
	err := pipeline.In([]string{"a", "b", "c"}).Filter(func(element interface{}, index int) bool {
		return element.(string) != "a"
	}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToBe(2)

	products := Products{{0, "Iphone 6", 0, 500}, {1, "HTC one", 0, 300}, {2, "Apple Watch", 1, 600}, {3, "ThinkPad", 2, 250}}
	var sample []Product
	Error := pipeline.In(products).
		Filter(func(el interface{}, i int) bool { return el.(Product).Price < 499 }).
		Out(&sample)
	e.Expect(Error).ToBeNil()
	e.Expect(len(sample)).ToEqual(2)
	e.Expect(sample[0].Price).ToEqual(300)
}

func TestIndexOf(t *testing.T) {
	e := expect.New(t)
	var result int = -1
	err := pipeline.In([]string{"i", "j", "k", "l"}).IndexOf("k", 0).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToEqual(2)
	var res int = -1
	err = pipeline.In("foobar").IndexOf('a', 0).Out(&res)
	e.Expect(err).ToBeNil()
	e.Expect(res).ToEqual(4)
}

func TestLastIndexOf(t *testing.T) {
	e := expect.New(t)
	var result int = -1
	err := pipeline.In("abba").LastIndexOf('a', 0).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToEqual(3)
}

func TestConcat(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{}).Concat([]int{1, 2, 3}, []int{4, 5, 6}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(6)
}

func TestReverse(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{1, 2, 3}).Reverse().Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result[0]).ToEqual(3)
}

func TestEvery(t *testing.T) {
	e := expect.New(t)
	var result bool
	err := pipeline.In([]int{2, 4, 6}).Every(func(element interface{}, index int) bool {
		return element.(int)%2 == 0
	}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToBeTrue()
	err = pipeline.In([]int{2, 4, 5}).Every(func(element interface{}, index int) bool {
		return element.(int)%2 == 0
	}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToBeFalse()
}

func TestSome(t *testing.T) {
	e := expect.New(t)
	var result bool
	err := pipeline.In([]int{1, 3, 6}).Some(func(element interface{}, index int) bool {
		return element.(int)%2 == 0
	}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToBeTrue()
}

func TestFirst(t *testing.T) {
	e := expect.New(t)
	var result int
	err := pipeline.In([]int{1, 3, 6}).First().Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToEqual(1)
}

func TestLast(t *testing.T) {
	e := expect.New(t)
	var result int
	err := pipeline.In([]int{1, 3, 6}).Last().Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToEqual(6)
}

func TestUnique(t *testing.T) {
	e := expect.New(t)
	var result []string
	err := pipeline.In([]string{"a", "b", "b", "a"}).Unique().Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(2)
	e.Expect(result[1]).ToEqual("b")
}

func TestDifference(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{1, 2, 3, 4}).Difference([]int{1, 3}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(2)
	for i, val := range []int{2, 4} {
		e.Expect(result[i]).ToEqual(val)
	}
}

func TestPush(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{2, 3, 4}).Push(5, 6, 7).Out(&result)
	e.Expect(err).ToBeNil()
	for i, val := range []int{2, 3, 4, 5, 6, 7} {
		e.Expect(result[i]).ToEqual(val)
	}
}

func TestWithout(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{1, 2, 3, 4}).Without(1, 3).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(2)
	for i, val := range []int{2, 4} {
		e.Expect(result[i]).ToEqual(val)
	}
}

func TestHead(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{1, 2, 3, 4}).Head(1).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(2)
	for i, val := range []int{1, 2} {
		e.Expect(result[i]).ToEqual(val)
	}
}

func TestTail(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{1, 2, 3, 4}).Tail(2).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(2)
	for i, val := range []int{3, 4} {
		e.Expect(result[i]).ToEqual(val)
	}
}

func TestSlice(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{1, 2, 3, 4}).Slice(0, 2).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(3)
	for i, val := range []int{1, 2, 3} {
		e.Expect(result[i]).ToEqual(val)
	}
}

func TestUnion(t *testing.T) {
	e := expect.New(t)
	var res []int
	err := pipeline.In([]int{1, 2}).Union([]int{2, 3}, []int{3, 4}).Out(&res)
	e.Expect(err).ToBeNil()
	e.Expect(fmt.Sprint(res)).ToEqual(fmt.Sprint([]int{1, 2, 3, 4}))
}

func ExamplePipeline_Splice() {
	var result []int
	err := pipeline.In([]int{1, 2, 3, 4, 5}).
		Splice(1, 2, []interface{}{6, 7, 8}...).
		Out(&result)

	fmt.Print(result, " ", err)
	// Output: [1 6 7 8 4 5] <nil>
}

func ExamplePipeline_Reduce() {
	// Using Map reduce to compile the total cost of an invoice
	type Order struct {
		ProductName string
		Quantity    int
		UnitPrice   int
	}
	var totalCost int
	command := []Order{{"Iphone", 2, 500}, {"Graphic card", 1, 250}, {"Flat screen", 3, 600}, {"Ipad air", 5, 200}}
	err := pipeline.In(command).Map(func(el interface{}, index int) interface{} {
		return el.(Order).Quantity * el.(Order).UnitPrice
	}).Reduce(func(result, el interface{}, index int) interface{} {
		return result.(int) + el.(int)
	}, 0).Out(&totalCost)

	fmt.Print(err, " ", totalCost)
	// Output: <nil> 4050
}

func ExamplePipeline_Sort() {
	var result []int
	err := pipeline.In([]int{2, 1, 6, 3, 5, 4}).Sort(func(a interface{}, b interface{}) bool {
		return a.(int) <= b.(int)
	}).Out(&result)
	fmt.Print(result, " ", err)
	// Output: [1 2 3 4 5 6] <nil>
}

func TestUnshift(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{3, 4}).Unshift(1, 2).Out(&result)
	e.Expect(err).ToBeNil()
	for i, val := range []int{1, 2, 3, 4} {
		e.Expect(result[i]).ToEqual(val)
	}

}

func TestChunk(t *testing.T) {
	e := expect.New(t)
	var result [][]int
	err := pipeline.In([]int{1, 2, 3, 4, 5}).Chunk(2).Out(&result)
	e.Expect(err).ToBeNil()
	t.Log(result)
	e.Expect(fmt.Sprint(result)).ToEqual(fmt.Sprint([][]int{{1, 2}, {3, 4}, {5}}))

}

func TestZip(t *testing.T) {
	e := expect.New(t)
	var result1 [][]interface{}
	err := pipeline.In([][]int{{1, 2, 3}}).Zip().Out(&result1)
	e.Expect(err).ToBeNil()
	e.Expect(fmt.Sprint(result1)).ToEqual(fmt.Sprint([][]interface{}{{1}, {2}, {3}}))
	err = pipeline.In([][]interface{}{{1, 2, 3}, {"John", "Jane", "David"}}).Zip().Out(&result1)
	e.Expect(err).ToBeNil()
	e.Expect(fmt.Sprint(result1)).
		ToEqual(fmt.Sprint([][]interface{}{{1, "John"}, {2, "Jane"}, {3, "David"}}))
	err = pipeline.In([][]interface{}{{"US", "FR"}, {"John", "Jane", "David"}, {true}}).Zip().Out(&result1)
	e.Expect(fmt.Sprint(result1)).
		ToEqual(fmt.Sprint([][]interface{}{{"US", "John", true}, {"FR", "Jane", nil}, {nil, "David", nil}}))
}

func TestCompact(t *testing.T) {
	e := expect.New(t)
	var result1 []interface{}
	err := pipeline.In([]interface{}{1, nil, 2, 'a', nil}).Compact().Out(&result1)
	e.Expect(err).ToBeNil()
	e.Expect(fmt.Sprint(result1)).ToEqual(fmt.Sprint([]interface{}{1, 2, 'a'}))
}

func TestEqual(t *testing.T) {
	e := expect.New(t)
	var result bool
	err := pipeline.In([]int{1, 2, 3}).Equals([]int{1, 2, 3}, []int{1, 2, 3}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToBeTrue()
	err = pipeline.In([]int{1, 2, 3}).Equals([]int{1, 2, 3}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToBeTrue()
	err = pipeline.In([]int{1, 2, 3}).Equals().Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToBeTrue()
	err = pipeline.In([]int{1, 2, 3}).Equals([]int{1, 2}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(result).ToBeFalse()
}

func TestIndexOutOfBoundsError(t *testing.T) {
	e := expect.New(t)
	err := pipeline.In([]int{1, 2}).Head(6).Out(&[]int{})
	e.Expect(err).Not()
}

func TestIterable(t *testing.T) {
	const (
		monday    = "monday"
		tuesday   = "tuesday"
		wednesday = "wednesday"
	)
	m := map[string]string{monday: "studies", tuesday: "date", wednesday: "training"}
	e := expect.New(t)
	e.Expect(m).Not().ToBeNil()
	e.Expect(pipeline.IsIterable(m)).ToBeTrue()
	i := pipeline.NewIterable(m)
	e.Expect(i.Length()).ToEqual(3)
}

func TestFlatten(t *testing.T) {
	e := expect.New(t)
	var result []int
	Error := pipeline.In([]interface{}{[]int{1, 2}, 3, []int{4, 5}}).Flatten().Out(&result)
	e.Expect(Error).ToBeNil()
	e.Expect(len(result)).ToEqual(5)
	t.Log(result)
}

func TestToMap(t *testing.T) {
	e := expect.New(t)
	in := map[string]string{"a": "angel", "b": "bookmark", "c": "card"}
	result := pipeline.In(in).ToMap(func(val interface{}, key interface{}) (interface{}, interface{}) {
		return key, val
	}).MustOut()
	e.Expect(result.(map[interface{}]interface{})["angel"]).ToEqual("a")
}

func ExamplePipeline_GroupBy() {
	// Counting words
	const words = `Lorem ipsum nascetur,
    nascetur adipiscing. Aenean commodo nascetur.
    Aenean nascetur commodo ridiculus nascetur,
    commodo ,nascetur consequat.`
	var result map[string]int
	err := pipeline.In(strings.Split(words, " ")).Map(func(el interface{}, i int) interface{} {
		return strings.Trim(strings.Trim(el.(string), " \r\n\t"), ".,!")
	}).GroupBy(func(el interface{}, i int) interface{} {
		return el.(string)
	}).ToMap(func(v interface{}, k interface{}) (interface{}, interface{}) {
		return len(v.([]interface{})), k
	}).Out(&result)

	// =>  map[ridiculus:1 ipsum:1 :9 Aenean:2 commodo:3 Lorem:1 nascetur:6 adipiscing:1 consequat:1]
	fmt.Print(err)
	// Output: <nil>
}
