package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/core"
	"xorm.io/xorm"
)

// DbModel2 数据库表db_model2对应的模型
type DbModel2 struct {
	ID        int64     // 不添加标签，由转换器自动转换为数据库字段，默认id为自增主键
	Title     string    `xorm:"varchar(25) notnull 'db_title'"` // 通过标签指定属性及字段名
	Created   time.Time `xorm:"created"`                        // 这个Field将在Insert时自动赋值为当前时间
	Updated   time.Time `xorm:"updated"`                        // 这个Field将在Insert或Update时自动赋值为当前时间
	DeletedAt time.Time `xorm:"deleted"`                        // 如果带DeletedAt这个字段和标签，xorm删除时自动软删除
}

func main() {
	const (
		dbType   string = "mysql"
		dbHost   string = "***"
		dbPort   string = "3306"
		dbName   string = "demo"
		dbUser   string = "***"
		dbPasswd string = "***"
		dbParams string = "charset=utf8&parseTime=true"
	)
	var dbURL = fmt.Sprintf("%s:%s@(%s:%s)/%s?%s", dbUser, dbPasswd, dbHost, dbPort, dbName, dbParams)

	// 创建一个引擎，对于xorm，一个引擎对应一个数据库
	// 但是创建引擎时不会检查连接有效性，只是简单讲配置文件解析为xorm数据结构
	// 当执行的操作需要用到连接时，xorm才会发现错误并返回错误
	engine, err := xorm.NewEngine(dbType, dbURL)
	if err != nil {
		log.Printf("打开mysql失败:%v\n", err)
		panic(err)
	}
	defer engine.Close()

	Config(engine)
	SyncTable(engine)
	C(engine)
	R(engine)
	U(engine)
	D(engine)

	// 获取表结构信息，通过调用engine.DBMetas()可以获取到数据库中所有的表，字段，索引的信息
	// dbMetas, err := engine.DBMetas()
	// if err != nil {
	// 	log.Printf("获取数据库表信息失败:%v\n", err)
	// 	panic(err)
	// }
	// log.Printf("获取数据库表信息成功")
	// jsonSring, _ := json.Marshal(dbMetas)
	// log.Println(string(jsonSring)
}

// SyncTable 自动同步数据库字段和结构体字段（根据结构体字段创建表或添加表的字段）
func SyncTable(engine *xorm.Engine) {
	err := engine.Sync2(new(DbModel2))
	if err != nil {
		log.Printf("同步数据库和结构体字段失败:%v\n", err)
		panic(err)
	}
}

// Config xorm一些可选设置
func Config(engine *xorm.Engine) {
	// 设置日志等级，设置显示sql，设置显示执行时间
	engine.SetLogLevel(xorm.DEFAULT_LOG_LEVEL)
	engine.ShowSQL(true)
	engine.ShowExecTime(true)

	// 指定结构体字段到数据库字段的转换器
	// 默认为core.SnakeMapper
	// 但是我们通常在struct中使用"ID"
	// 而SnakeMapper将"ID"转换为"i_d"
	// 因此我们需要手动指定转换器为core.GonicMapper{}
	engine.SetMapper(core.GonicMapper{})
}

// C 增
func C(engine *xorm.Engine) {
	// 增（无法获取插入主键等信息）
	newMsg := DbModel2{
		Title: "new message",
	}
	affected, err := engine.Insert(newMsg)
	if err != nil {
		log.Printf("插入数据失败:%v\n", err)
		panic(err)
	}
	log.Printf("插入数据成功，影响行数:%v\n", affected)

	// 增（利用指针可获取插入后的主键）
	newMsg2 := DbModel2{
		Title: "new message 2",
	}
	affected, err = engine.Insert(&newMsg2)
	if err != nil {
		log.Printf("插入数据失败:%v\n", err)
		panic(err)
	}
	log.Printf("插入数据成功，影响行数:%v\n，插入数据为%v\n", affected, newMsg2)

	// 增，使用数组插入多条记录，且每条记录都可以获取到主键
	// 注意需要传入的是数组解构后的指针，如果传入的不是指针将无法获取主键
	msgList := make([]*DbModel2, 2)
	msgList[0] = new(DbModel2)
	msgList[0].Title = "list 1"
	msgList[1] = new(DbModel2)
	msgList[1].Title = "list 2"

	// 此处官方有误，这样无法把msgList...解构传递过去，会报错
	// cannot use msgList (variable of type []*DbModel2) as []interface{}
	// affected, err = engine.Insert(msgList...)
	// 只能这样
	affected, err = engine.Insert(msgList[0], msgList[1])
	if err != nil {
		log.Printf("插入数据失败:%v\n", err)
		panic(err)
	}
	log.Printf("插入数据成功，影响行数:%v\n，插入数据为%v, %v\n",
		affected, msgList[0], msgList[1])
}

// R 查
func R(engine *xorm.Engine) {
	// Get 获得单条数据
	var selectedMsg DbModel2
	// Alias给表起一个别名
	has, err := engine.Alias("o").Where("o.db_title = ?", "list 1").Get(&selectedMsg)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	if has {
		log.Println(selectedMsg)
	} else {
		log.Println("查询结果为空")
	}
	// Find 获取多条数据
	var selectedMsgList []DbModel2
	// Alias给表起一个别名
	err = engine.Alias("o").Where("o.db_title = ?", "list 2").Find(&selectedMsgList)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println(selectedMsgList)
}

// U 改
func U(engine *xorm.Engine) {
	updateMsg1 := DbModel2{
		Title: "updated by xorm",
	}
	// Update方法将返回两个参数，第一个为 更新的记录数
	// 需要注意的是 SQLITE 数据库返回的是根据更新条件查询的记录数
	// 而不是真正受更新的记录数

	// Update会自动从user结构体中提取非零值和非nil得值作为需要更新的内容
	// 因此，如果需要更新一个值为零值，则此种方法将无法实现
	// 并且当没有修改任何值的时候，xorm也会自动更新updated的时间
	// 从而使得返回的影响行数不为零值
	affected, err := engine.Id(3).Update(&updateMsg1)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Printf("更新成功，影响%v行", affected)

	// 如果需要更新一个值为零值

	// 方法1：通过指定Cols，指定的Cols一定会在结构体中取值
	updateMsg2 := DbModel2{
		Title: "",
	}
	affected, err = engine.Id(1).Cols("db_title").Update(&updateMsg2)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Printf("更新成功，影响%v行", affected)

	// 方法2：通过传入map[string]interface{}来进行更新，但这时需要额外指定更新到哪个表，因为通过map是无法自动检测更新哪个表的。
	affected, err = engine.Table(new(DbModel2)).Id(1).Update(map[string]interface{}{"db_title": ""})
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Printf("更新成功，影响%v行", affected)
}

// D 删
func D(engine *xorm.Engine) {
	// 由于设置了DeletedAt，删除会自动变成软删除

	// 删除，用Id显式指定主键
	// 用传入的结构体指针代表需要操作的表
	// 并读取非零值，非nil的字段进行组合限制删除
	// 相当于WHERE `db_title`=? AND `id`=?
	deleteMsg1 := DbModel2{
		Title: "new message 2",
	}
	affected, err := engine.Id(2).Delete(&deleteMsg1)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Printf("删除成功，影响%v行", affected)

	// 也可以直接在结构体中给主键赋值达到指定主键
	// 同时指定其他字段的组合限制效果
	deleteMsg2 := DbModel2{
		ID:    2,
		Title: "new message",
	}
	affected, err = engine.Delete(&deleteMsg2)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Printf("删除成功，影响%v行", affected)

	// // 但注意！！当没有使用Id进行指定，且结构体中没有非零非nil的主键时
	// // xorm会删除全部数据！！
	// // 删库跑路一气呵成~~
	// deleteAll := DbModel2{}
	// affected, err = engine.Delete(&deleteAll)
	// if err != nil {
	// 	log.Println(err)
	// 	panic(err)
	// }
	// log.Printf("删除成功，影响%v行", affected)
}
