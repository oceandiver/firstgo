package session

import (
    "testing"
    "time"
)

func TestInit(t *testing.T) {
    var lastInsertId int
    err := mydb.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
    checkErr(err)
    t.Log("last inserted id =", lastInsertId)

    t.Log("# Updating")
    stmt, err := mydb.Prepare("update userinfo set username=$1 where uid=$2")
    checkErr(err)

    res, err := stmt.Exec("astaxieupdate", lastInsertId)
    checkErr(err)

    affect, err := res.RowsAffected()
    checkErr(err)

    t.Log(affect, "rows changed")

    t.Log("# Querying")
    rows, err := mydb.Query("SELECT * FROM userinfo")
    checkErr(err)

    for rows.Next() {
        var uid int
        var username string
        var department string
        var created time.Time
        err = rows.Scan(&uid, &username, &department, &created)
        checkErr(err)
        t.Log("uid | username | department | created ")
        t.Log("%3v | %8v | %6v | %6v\n", uid, username, department, created)
    }

    t.Log("# Deleting")
    stmt, err = mydb.Prepare("delete from userinfo where uid=$1")
    checkErr(err)

    res, err = stmt.Exec(lastInsertId)
    checkErr(err)

    affect, err = res.RowsAffected()
    checkErr(err)

    t.Log(affect, "rows changed")
    
}

