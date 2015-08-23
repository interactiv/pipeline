#Pipeline

author mparaiso <mparaiso@online.fr>

copyrights 2014

license GPL-3.0

version 0.1

Martin Fowler pipelines for Go 

http://martinfowler.com/articles/collection-pipeline/

##Examples:

###Counting words

    // Counting words
    const words = `Lorem ipsum nascetur,
    nascetur adipiscing. Aenean commodo nascetur.
    Aenean nascetur commodo ridiculus nascetur,
    commodo ,nascetur consequat.`
    var result map[string]int
    err := pipeline.In(strings.Split(words, " ")).Map(func(el interface{}, i int)         interface{} {
        return strings.Trim(strings.Trim(el.(string), " \r\n\t"), ".,!")
    }).GroupBy(func(el interface{}, i int) interface{} {
    return el.(string)
    }).ToMap(func(v interface{}, k interface{}) (interface{}, interface{}) {
    return []interface{}{len(v.([]interface{})), k}, k
    }).Out(&result)
    
    // =>  map[ridiculus:1 ipsum:1 :9 Aenean:2 commodo:3 Lorem:1 nascetur:6 adipiscing:1 consequat:1]
    fmt.Print(err)     

