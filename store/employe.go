package store

import (
	"context"
	"errors"
	"log"
)

type Employe struct{
    Id string `json:"id"`
    Name string
    ShortName string
    Email string
    Age string
    ShopName string
    ShopId string `json:"shopId"`
    UserId string `json:"userId"`
    EntrepriseId string `json:"entrepriseId"`
}

func (e *Employe) Search()[]Employe{
    var employeList []Employe
    var employe Employe
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error conn db %v", err)
        return employeList
    }
    defer conn.Close()
    usersRow, err := conn.QueryContext(context.Background(), `SELECT email, id, LEFT(firstname, 1) || LEFT(lastname, 1) FROM Users WHERE email=$1 LIMIT 5`,  e.Email)
    if err != nil{
        log.Printf("error in query %v", err)
        return employeList
    }

    for usersRow.Next(){
        usersRow.Scan(&employe.Email, &employe.UserId, &employe.ShortName)
        employeList = append(employeList, employe)
    }
    log.Printf("employes: %v", e.Email)
    return employeList
}

func (e *Employe) New() error{
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error conn db %v", err)
        return errors.New("error conn db")
    }
    defer conn.Close()

    employeRow := conn.QueryRowContext(context.Background(), `INSERT INTO Employe (user_id, entreprise_id) VALUES($1, $2) RETURNING id`, e.UserId, e.EntrepriseId)
    if err := employeRow.Scan(&e.Id); err != nil{
        log.Printf("error new employe query")
        return errors.New("error query")
    }
    return nil
}

func (e *Employe) Delete()error{
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error conn db %v", err)
        return errors.New("error conn db")
    }
    defer conn.Close()
    _, err = conn.ExecContext(context.Background(), `DELETE FROM Employe WHERE id=$1`, e.Id)
    if err != nil{
        log.Printf("error deleting employe")
        return errors.New("query error deleting employe")
    }
    return nil
}

func (e *Employe) ShopAssign()error{
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error conn db %v", err)
        return errors.New("error conn db")
    }
    defer conn.Close()
    _, err = conn.ExecContext(context.Background(), `UPDATE Employe SET shop_id=$1 WHERE id=$2 AND entreprise_id=$3`, e.ShopId, e.Id, e.EntrepriseId)
    if err != nil{
        log.Printf("error in the query %v", err)
        return errors.New("error in the query")
    }
    return nil
}
