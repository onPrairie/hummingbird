package utilsEx

import (
	"fmt"
	"strings"
	"utils/ulog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MysqlDb *sqlx.DB

var constr string //clear

func OpenMysql(conn string) (err error) {
	constr = conn
	MysqlDb, err = sqlx.Open("mysql", conn)
	for true {
		err := MysqlDb.Ping()
		if err != nil {
			ulog.Warnln("数据库链接失败", err)
		} else {
			break
		}
	}
	MysqlDb.SetMaxOpenConns(10) // 设置连接池最大连接数
	MysqlDb.SetMaxIdleConns(10) // 设置连接池最大空闲连接数
	return nil
}
func CeateOrm(tableaname string) {
	createsql2go()
	execsql2go(tableaname)
	deletsql2go()
	//删除
	constr = ""
}

//创建mysql函数
func createsql2go() {
	sqlstr := "CREATE PROCEDURE sql2go ( IN s_t_name CHAR(30) ) BEGIN DECLARE t_name, s_name CHAR(20); SELECT SUBSTRING_INDEX(s_t_name, '.', 1) INTO s_name; SELECT SUBSTRING_INDEX(s_t_name, '.', -1) INTO t_name; SELECT concat(concat(UPPER(LEFT(column_name, 1)), RIGHT(column_name, LENGTH(column_name) - 1)), '  ', CASE  WHEN data_type IN ('varchar', 'char', 'text') THEN 'string' WHEN data_type IN ('int', 'tinyint') THEN 'int' WHEN data_type IN ('bigint') THEN 'int64' WHEN data_type IN ('datetime') THEN 'time.Time' WHEN data_type IN ('bit', 'boolean') THEN 'bool' ELSE '类型不确定' END, ' ', concat('`db:\"', column_name, '\"`')) AS golang FROM information_schema.COLUMNS WHERE table_name = t_name AND TABLE_SCHEMA = s_name; END"
	// mysqlstr := readConfig("mysql")
	r := MysqlDb.MustExec(sqlstr)
	_, err := r.LastInsertId()
	if err != nil {
		panic(err)
	}
	_, err = r.RowsAffected()
	if err != nil {
		panic(err)

	}
}

//执行sql2go函数
func execsql2go(tableaname string) {
	sqlstr := "CALL sql2go('%s.%s')"
	i := strings.LastIndex(constr, "/")
	sqlstr = fmt.Sprintf(sqlstr, constr[i+1:], tableaname)
	rows, err := MysqlDb.Query(sqlstr)
	if err != nil {
		panic("errr CeateOrm")
	}
	structs := "type  struct {\n\t"
	for rows.Next() { //循环结果
		var args string
		err = rows.Scan(&args)
		if err != nil {
			panic(err)
		}
		structs += args + "\n\t"
	}
	structs = structs[:len(structs)-1]
	structs += "}"
	fmt.Println(structs)
}
func deletsql2go() {
	sqlstr := "drop procedure if exists sql2go;"
	r := MysqlDb.MustExec(sqlstr)
	_, err := r.LastInsertId()
	if err != nil {
		panic(err)
	}
	_, err = r.RowsAffected()
	if err != nil {
		panic(err)
	}
}
