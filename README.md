# lion

reflection is useful, but I do't like complex tags.

[`lion`](https://www.dota2.com/hero/lion) is a character of Dota2, kind of evil and cute. Just like this package.

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
    lion.RegisterOf[Options]().TagNames("opts")
}

// 3. set metainfo
func init() {
    ptr := lion.Ptr[User]()

    lion.FieldOf[User, Options](&ptr.Name).Meta = &Options{Info : "balabala"}
}

// 4. read metainfo

// get typeinfo
lion.TypeinfoOf[User, Options]()

// get fieldinfo
lion.FieldOf[User, Options](&(reflectx.Ptr[User]().Name))
```

# !!!!

You can noly set metainfo in `init` function, because there is no mutex.
