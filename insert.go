package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // init()
)

type Resp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type Auth struct {
	Username int `json:"username"`
	Pwd      int `json:"password"`
	Msg      int `json:"msg"`
}

//go连接mysql示例
var db *sql.DB //是一个连接池对象

type search struct {
	jkcount     int
	gitcodeline int
	gitcommit   int
}

func initDB() (err error) {
	//数据库信息
	dsn := "dyx:123456@tcp(192.168.38.248:3306)/mvw"
	//连接数据库
	db, err = sql.Open("mysql", dsn) //不会校验用户名密码是否正确
	if err != nil {                  //dsn格式不正确的时候会报错
		//fmt.Printf("dsn:%v invalid, err:5v\n", dsn, err)
		return
	}
	err = db.Ping()
	if err != nil {
		//fmt.Printf("open %v failed, err:5v\n", dsn, err)
		return
	}
	//设置数据库连接池的最大连接数
	db.SetMaxIdleConns(10)
	return
}

//查询多行记录
func queryRowDemo() {
	sqlStr := `select jkcount,gitcodeline,gitcommit from mvw.test`
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:5v\n", err)
		return
	}
	//非常重要：关闭rows释放持有的数据库连接
	defer rows.Close()

	//循环读取结果集中的数据
	for rows.Next() {
		var u search
		err := rows.Scan(&u.jkcount, &u.gitcodeline, &u.gitcommit)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Println(u)
	}

}

//插入数据
func insertRowDemo(a, b, c int) {
	sqlStr := `insert  into test(jkcount,gitcodeline,gitcommit) values(?,?,?);`
	ret, err := db.Exec(sqlStr, a, b, c)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	//如果是插入数据的操作，能够拿到插入数据的id
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get id failed, err:%v", err)
		return
	}
	fmt.Println("id:", id)

}

//post接口接收json数据
func login(writer http.ResponseWriter, request *http.Request) {
	var auth Auth
	if err := json.NewDecoder(request.Body).Decode(&auth); err != nil {
		request.Body.Close()
		log.Fatal(err)
	}
	fmt.Println(auth.Username, auth.Pwd, auth.Msg)
	//插入数据
	err := initDB()
	if err != nil {
		fmt.Printf("init DB failed, err:%v\n", err)
	}
	fmt.Println("连接数据库成功")

	queryRowDemo()
	insertRowDemo(auth.Username, auth.Pwd, auth.Msg)
	queryRowDemo()
	// var result  Resp
	// if auth.Username == "admin" && auth.Pwd == "123456" {
	//     result.Code = "200"
	//     result.Msg = "登录成功"
	// } else {
	//     result.Code = "401"
	//     result.Msg = "账户名或密码错误"
	// }
	// if err := json.NewEncoder(writer).Encode(result); err != nil {
	//     log.Fatal(err)
	// }
}

func main() {
	http.HandleFunc("/login", login)
	http.ListenAndServe("0.0.0.0:11111", nil)
	// err := initDB()
	// if err != nil {
	// 	fmt.Printf("init DB failed, err:%v\n", err)
	// }
	// fmt.Println("连接数据库成功")

	// queryRowDemo()
	// insertRowDemo(msg.Username, msg.Pwd, msg.Msg)
	// queryRowDemo()

}

