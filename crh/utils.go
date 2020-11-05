package main

import (
	"fmt"
	"strings"
)
//角色信息组合
type Role struct {
    Name        string      //名字
    FlagCode    int         //角色位标志
    RaceCode    int         //种族标志
}
//获取全部角色列表
func getRole() (roles []Role) {
    roleSql:=`
        select sprmd.name,sprmd.model,race 
        from SPRoleManDefault sprmd join SPKindle1 spk1 on sprmd.model=spk1.id
    `
    rows,_ := Db.Query(roleSql)
    for rows.Next() {
        role := Role{}
        rows.Scan(&role.Name,&role.FlagCode,&role.RaceCode)
        roles = append(roles, role)
    }
    rows.Close()
    return
}
//利用角色标志位，获得适用角色的名称列表，并予以浓缩精简
func getValidRole(flagCode int) (names string) {
    //先处理极端情况
    if flagCode==0 {
        names="无适用角色"
        return
    }
    roles:=getRole()
    races:=make(map[int]int)
    for _,role:= range roles {
        races[role.RaceCode]++
    }
    validRoles:=[]Role{}
    validRaces:=make(map[int]int)
    for _,role := range roles {
        if flagCode & (1<<(role.FlagCode-1)) != 0 {
            validRoles=append(validRoles,role)
            validRaces[role.RaceCode]++
            names+=fmt.Sprintf("%s ",role.Name)
        }
    }
    if validRaces[0]==races[0] && validRaces[5]==0 {
        names="男性"
    }
    if validRaces[5]==races[5] && validRaces[0]==0 {
        names="女性"
    }
    if len(validRoles)==len(roles) {
        names="全体"
    }
    names=strings.TrimSuffix(names," ")
    return
}
//以更好看的方式，显示全部的百分比数值
func perDisp(f float32) (fs string) {
	fs = fmt.Sprintf("%.2f", f)
	for {
		hasDot, TrailZero := strings.Contains(fs, "."), strings.HasSuffix(fs, "0")
		if !TrailZero || !hasDot {
			break
		} else {
			fs = strings.TrimSuffix(fs, "0")
			fs = strings.TrimSuffix(fs, ".")
		}
	}
	return
}

//两个函数，一个通过ID查名称，另一个反过来
func getName(table string, id int) (name string) {
	sql := fmt.Sprintf("select name from %s where id=%d", table, id)
	Db.QueryRow(sql).Scan(&name)
	return
}
func getId(table, name string) (id int) {
	sql := fmt.Sprintf(`select id from %s where name="%s"`, table, name)
	Db.QueryRow(sql).Scan(&id)
	return
}