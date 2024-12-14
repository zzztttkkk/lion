# reflectx

reflection is useful, but I do't like complex tag.

```go

type User struct {
    Name string `opts:"name"`
    Age int `opts:"age"`
}

// 1. define metainfo type
type Options struct {
    Info string
}

// 2. set tag name
func init(){
    reflectx.RegisterOf[Options]().TagNames("opts")
}

// 3. set metainfo
func init() {
    ptr := reflectx.Ptr[User]()

    reflectx.FieldOf[User, Options](&ptr.Name).Meta = &Options{Info : "balabala"}
}

// 4. read metainfo

// get typeinfo
reflectx.TypeinfoOf[User, Options]()

// get fieldinfo
reflectx.FieldOf[User, Options](&(reflectx.Ptr[User]().Name))
```

# !!!!

You can noly set metainfo in `init` function, because there is no mutex.
