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