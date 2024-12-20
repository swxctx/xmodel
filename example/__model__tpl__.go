// Command example is the xmodel tools.
// The framework reference: https://github.com/swxctx/xmodel
package __TPL__

// __MYSQL_MODEL__ create mysql model
type __MYSQL_MODEL__ struct {
	User
}

// __MONGO_MODEL__ create mongodb model
type __MONGO_MODEL__ struct {
	Meta
}

// User user info
type User struct {
	Id   int64  `key:"pri"`
	Name string `key:"uni"`
	Age  int32
}

type Meta struct {
	Hobby []string
	Tags  []string
}
