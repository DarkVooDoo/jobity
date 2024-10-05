package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type Shop struct{ 
    Id string `json:""` 
    Name string `json:"name"` 
    City string `json:"city"` 
    Postal string `json:"postal"` 
    Adresse string `json:"adresse"`
    Lat float64 `json:"lat"`
    Long float64 `json:"long"`
    EmployeCount int 
    EntrepriseId string `json:"entrepriseId"`
}

func (s *Shop) GetAll()[]Shop{
    var allShop []Shop
    var shop Shop
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error connecting to the db %v", err)
        return allShop
    }
    defer conn.Close()
    shopRows, err := conn.QueryContext(context.Background(), `SELECT id, name, addr, postal FROM Shop WHERE entreprise_id=$1`, s.EntrepriseId)
    if err != nil{
        log.Printf("error in the query %v", err)
        return allShop
    }
    for shopRows.Next(){
        if err := shopRows.Scan(&shop.Id, &shop.Name, &shop.Adresse, &shop.Postal); err != nil{
            log.Printf("scan error: %v", err)
        }
        allShop = append(allShop, shop)
    }
    return allShop
}

func (s *Shop) Employes()[]Employe{
    var employeList []Employe
    var employe Employe
    var id, name sql.NullString
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error connecting to the db %v", err)
        return employeList
    }
    defer conn.Close()
    employeRow, err := conn.QueryContext(context.Background(), `SELECT e.id, u.firstname || ' ' || u.lastname, '4 mois', s.name FROM Shop AS s LEFT JOIN Employe AS e ON s.id=e.shop_id LEFT JOIN Users AS u ON e.user_id=u.id  WHERE s.id=$1`, s.Id)
    if err != nil{
        log.Printf("error in the query %v", err)
        return employeList
    }
    for employeRow.Next(){
        if err := employeRow.Scan(&id, &name, &employe.Age, &s.Name); err != nil{
            log.Printf("Error scan %v", err)
        }
        if id.Valid{
            employe = Employe{Id: id.String, Name: name.String}
            employeList = append(employeList, employe)
        }
    }
    s.EmployeCount = len(employeList)
    return employeList
}

func (s *Shop) UnassignEmploye()[]Employe{
    var unassignEmployeList []Employe
    var employe Employe
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error connecting to the db %v", err)
        return unassignEmployeList
    }
    defer conn.Close()
    unassignEmployesRows, err := conn.QueryContext(context.Background(), `SELECT e.id, u.firstname || ' ' || u.lastname, '4 mois' FROM Employe AS e LEFT JOIN Users AS u ON u.id=e.user_id WHERE e.entreprise_id=$1 AND shop_id IS NULL `, s.EntrepriseId)
    if err != nil{
        log.Println(err)
        return unassignEmployeList
    }
    for unassignEmployesRows.Next(){
        if err := unassignEmployesRows.Scan(&employe.Id, &employe.Name, &employe.Age); err != nil{
            log.Printf("error scan %v", err)
        }
        unassignEmployeList = append(unassignEmployeList, employe)
    }
    return unassignEmployeList
}

func (s *Shop) Create() error{
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error connecting to the db %v", err)
        return errors.New("error creating shop")
    }
    defer conn.Close()
    shopRow := conn.QueryRowContext(context.Background(), `INSERT INTO Shop (name, addr, city, postal, lat, long, entreprise_id) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`, s.Name, s.Adresse, s.City, s.Postal, s.Lat, s.Long, s.EntrepriseId)
    if err := shopRow.Scan(&s.Id); err != nil{
        log.Printf("error query %v", err)
        return errors.New("error in the query")
    }
    return nil
}

func (s *Shop) Delete()error{
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error connecting to the db %v", err)
        return errors.New("error creating shop")
    }
    defer conn.Close()
    if _, err = conn.ExecContext(context.Background(), `DELETE FROM Shop WHERE id=$1`, s.Id); err != nil{
        return errors.New("error in the query delete")
    }
    return nil
}
