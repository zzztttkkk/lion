# lion

reflection is useful, but I do't like complex tags.

[`lion`](https://www.dota2.com/hero/lion) is a character of Dota2, kind of evil and cute. Just like this package.

```go
type User struct {
    Name string
    Age int
}

// 1. define metainfo type
type Options struct {
    Info string
}

// 3. set metainfo
func init() {
    lion.UpdateMetaScope(func(mptr *User, update func(fptr any, meta *Options)){
        update(
            &mptr.Name,
            &Options{
                Info: "balabalaxiaomoxian"
            },
        )
    })
}

// get typeinfo
lion.TypeinfoOf[User]()

// get fieldinfo
lion.FieldOf[User](&(lion.Ptr[User]().Name))
```

# !!!!

You can noly set metainfo in `init` function, because there is no mutex.
