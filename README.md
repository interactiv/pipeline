#Pipeline

Pipeline is a functionnal programming library for the Go langague. With Pipeline developpers can use 
functionnal principles on their collection types such as slices,maps and string. Pipeline is written in go and inspired by underscore.js , lodash.js and Martin Fowler's pipelines :

http://martinfowler.com/articles/collection-pipeline/

author mparaiso <mparaiso@online.fr>

copyrights 2014

license GPL-3.0

version 0.1



##Installating:

* Install the Go language

* Use 'go get' with a command line interface

    go get github.com/interactiv/pipeline

##Examples:

###Counting words

```go
    // Counting words
    const words = `Lorem ipsum nascetur,
    nascetur adipiscing. Aenean commodo nascetur.
    Aenean nascetur commodo ridiculus nascetur,
    commodo ,nascetur consequat.`
    var result map[string]int
    err := pipeline.In(strings.Split(words, " ")).Map(
		func(el interface{}, i int) interface{} {
        	return strings.Trim(strings.Trim(el.(string), " \r\n\t"), ".,!")
    	}).GroupBy(func(el interface{}, i int) interface{} {
    		return el.(string)
    	}).ToMap(func(v interface{}, k interface{}) (interface{}, interface{}) {
    		return []interface{}{len(v.([]interface{})), k}, k
    	}).Out(&result)
    
    // =>  map[ridiculus:1 ipsum:1 :9 Aenean:2 commodo:3 Lorem:1 nascetur:6 adipiscing:1 consequat:1]
    fmt.Print(err)     
```

###Calculating the total cost of an customer order

```go
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
```

##Implemented pipelines 

* Chunk
* Compact
* Concat
* Difference
* Equals
* Every
* Filter
* First
* Flatten
* GroupBy
* Head
* IndexOf
* Intersection
* Last
* LastIndexOf
* Map
* Must
* Push
* Reduce
* ReduceRight
* Reverse
* Slice
* Some
* Sort
* Splice
* Tail
* ToMap
* Union
* Unique
* Unshift
* Without
* Xor
* Zip

