# TwoStructs

## Summary

Convert between two structs of different underlying types.

## Example

Many applications have both a `WireModel` for the handler/http/json layer and an `InternalModel` used by the internal business logic.
```
type Epoch int64

type WireModel struct {
	Name            string
	OptionalAddress *string
	TimeInSeconds   Epoch
}

type InternalModel struct {
	FullName string    // field name changes from "Name" to "FullName"
	Address  string    // type changes from *string to string
	Time     time.Time // type changes from Epoch(int64) to time.Time
}
```

Using this package (`twostructs`) we can easily map an instance of `WireModel` to `InternalModel`.  We need to help out though... not all types are easily convertable:
```
// Create a mapper
mapper := twostructs.New()

// Register any needed field type conversions
mapper.RegisterMappingFunction(func(e Epoch) time.Time {
	return time.Unix(int64(e), 0).UTC()
})
```

Let's see it in action:
```
// Source
wireObj := WireModel{Name: "David", OptionalAddress: nil, TimeInSeconds: 1553878048}

// Destination
entityObj := InternalModel{}

// TwoStructs Map
if err := mapper.Struct(wireObj, &entityObj); err != nil {
  panic(err)
}
fmt.Printf("    %#v\n", wireObj)
fmt.Printf("%#v\n", entityObj)

```
```
WireModel{
  Name:"David",
  OptionalAddress:(*string)(nil),
  TimeInSeconds:1553878048
}
InternalModel{
  FullName:"David",
  Address:"",
  Time:time.Time{ext:63689474848}
}
```
[example source](https://play.golang.org/p/QOdbNTQaRNG)

## Performance

Example above takes 0.7ms.
```
BenchmarkReadMe-8   	 2000000	       689 ns/op
```

## Background Information

`golang`'s ref/spec on [Conversions](https://golang.org/ref/spec#Conversions) allows assignment between structs of different types provided they are smiliar.

> A (struct) x can be converted to (struct) T if...
> - ignoring struct tags...
> - x's type and T have the same sequence of fields, and if corresponding fields have the same names, and identical types.
> - non-exported field names from different packages are always different.

Concretely:
```
type Person struct {
	Name string
	Age  int
}

type User struct {
	Name string // same field name & type
	Age  int
}

var p Person = Person{"David", 22}
var u User = User(p) // Valid
```
[example source](https://play.golang.org/p/DYGeLr0djTu)

