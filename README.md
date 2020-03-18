# wrap

### wrap for me.

## modules 

### wraps : factories for instances
1. wraps.GetMysql 
2. wraps.GetRedis
3. wraps.GetOss
4. wraps.GetRequest
5. wraps.GetCall

### utils
1. encryptutil
2. fileutil
3. hashutil
4. imageutil
5. lru
6. mathutil
7. netutil
8. nutil for number
9. osutil
10. sqlutil
11. strngutil
12. sutil for struct
13. timeutil

### wrapiris : for iris
```go
package main
import (
	"github.com/RocksonZeta/wrap/iriswrap"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)
func main() {
	app := iris.New()
	app.ContextPool.Attach(func() context.Context {
		return &iriswrap.Context{
			Context:        context.NewContext(app),
			PageSize:       20,
			AutoIncludeCss: true,
			AutoIncludeJs:  false,
			AutoHead: true,
		}
	})
	app.Use(iriswrap.SessionFilter)

	app.Get("/", func(ctx iris.Context) {
		c := ctx.(*iriswrap.Context)
		c.Ok(c.Session.Uid())
	})
	app.Listen(":9000")
}

```

### wraplog : for logging

### errs : common errs
```go
type Err struct {
	State     int         //0:ok ,other:error
	Data      interface{} 
	Message   string      
	Pkg       string
	Module    string 
	Err       error  //source error
	UserError bool   //user error or system error
}
```

## Architecture
![Architecture](https://github.com/RocksonZeta/wrap/blob/master/arch.png)
