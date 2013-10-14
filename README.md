gofig
=============

Gofig provides methods to load configuration values from a given JSON file.

When inflating structs, key name underscores are dropped & the first letter of the next
word is capitalized.  For example, a json key of 'my_config_value' will be
populated into a struct field called MyConfigValue.

Documentation available at [godoc](http://godoc.org/github.com/JKallhoff/gofig).

Example
-------
    if conf, err := gofig.Load("./application.conf"); err == nil {
        // Basic types
        s, err := conf.Str("first_name")
        b, err := conf.Bool("is_active")
        i, err := conf.Int("item_count")
        f, err := conf.Float("ranking")
        a, err := conf.Array("some_array")
        m, err := conf.Map("complex_object")
        
        // Inflating structs
        obj := &SomeType{}
        if err := conf.Struct("complex_object", obj); err == nil {
            // do something with inflated obj struct
        }

        var items []Item
        if err := conf.StructArray("lots_of_items", &items); err == nil {
            // do something with populated items array  
        }
    }