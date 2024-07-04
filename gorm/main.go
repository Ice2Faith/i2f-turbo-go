package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*

create table sys_user
(
 id bigint auto_increment primary key comment 'ID',
 username varchar(300) not null comment '用户名',
 age tinyint comment '年龄',
 high decimal(5,2) comment '升高'
) comment '用户表';

*/

type SysUser struct {
	Id       int     `gorm:"column:id;primaryKey"`
	Username string  `gorm:"column:username"`
	Age      int     `gorm:"column:age"`
	High     float32 `gorm:"column:high"`
}

func (su *SysUser) TableName() string {
	return "sys_user"
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:ltb12315@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 根据条件查询
	eu := SysUser{}
	db.Where("username = ?", "root").Find(&eu)
	fmt.Println("是否存在root:", eu)

	if eu.Id > 0 {
		// 删除
		db.Where("id = ?", eu.Id).Delete(&SysUser{})
		fmt.Println("删除已存在的root:", eu)
		eu.Id = 0
	}

	if eu.Id == 0 {
		// 插入
		su := SysUser{
			Username: "root",
			Age:      22,
			High:     1.73,
		}
		db.Save(&su)
		fmt.Println("不存在root，插入：", su)
	}

	// 查询所有
	list := []SysUser{}
	db.Find(&list)

	fmt.Println("查询所有：")
	fmt.Println(list)

	// 更新
	// 单列更新
	db.Model(&SysUser{}).Where("username = ?", "root").Update("age", 23)
	// 多列更新
	db.Model(&SysUser{}).Where("username = ?", "root").Updates(map[string]interface{}{
		"age":  23,
		"high": 1.75,
	})

	// 查询所有
	list = []SysUser{}
	db.Find(&list)

	fmt.Println("查询所有：")
	fmt.Println(list)
}
