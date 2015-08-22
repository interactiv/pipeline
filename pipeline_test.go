package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/interactiv/expect"
	"github.com/interactiv/pipeline"
)

type Val interface{}

func TestMap(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{1, 2, 3}).Map(func(element interface{}, index int) interface{} {
		return element.(int) * 2
	}).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(3)
	e.Expect(result[0]).ToEqual(2)

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

func TestSplice(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{1, 2, 3, 4, 5}).Splice(1, 2, []interface{}{6, 7, 8}...).Out(&result)
	e.Expect(err).ToBeNil()
	e.Expect(len(result)).ToEqual(6)
	for i, val := range []int{1, 6, 7, 8, 4, 5} {
		e.Expect(result[i]).ToEqual(val)
	}
}

func TestSort(t *testing.T) {
	e := expect.New(t)
	var result []int
	err := pipeline.In([]int{2, 1, 6, 3, 5, 4}).Sort(func(a interface{}, b interface{}) bool {
		return a.(int) <= b.(int)
	}).Out(&result)
	e.Expect(err).ToBeNil()
	for i, val := range []int{1, 2, 3, 4, 5, 6} {
		e.Expect(result[i]).ToEqual(val)
	}

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

func TestNotIterableError(t *testing.T) {
	e := expect.New(t)
	err := pipeline.NotIterable("foo")
	_, ok := err.(pipeline.NotIterableError)
	e.Expect(ok).ToBeTrue()
}

func TestIndexOutOfBoundsError(t *testing.T) {
	e := expect.New(t)
	err := pipeline.In([]int{1, 2}).Head(6).Out(&[]int{})
	e.Expect(err).Not()
}
