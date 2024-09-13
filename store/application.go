package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type JobApplication struct{
    Id string
    Title string
    Status string
    JobId string
    UserId string
}

func (a *JobApplication) CreateJobApplication()error{
    var id string
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("error conn to the db")
    }
    defer conn.Close()
    tx, _ := conn.BeginTx(context.Background(), &sql.TxOptions{})
    row := tx.QueryRow(`SELECT id FROM JobApplication WHERE user_id=$1 AND job_id=$2`, a.UserId, a.JobId)
    if err := row.Scan(&id); err == nil{
        tx.Rollback()
        return errors.New("Application already exist")
    }
    _, err = tx.Exec(`INSERT INTO JobApplication (user_id, job_id) VALUES($1,$2)`, a.UserId, a.JobId)
    if err != nil{
        tx.Rollback()
        log.Println(err)
        return errors.New("error inserting job application")
    }
    tx.Commit()
    return nil
}

func (a *JobApplication) GetUserApplications()[]JobApplication{
    var applications []JobApplication
    var id, title string
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return applications
    }
    defer conn.Close()
    applicationRow, err := conn.QueryContext(context.Background(), `SELECT a.id, j.title FROM JobApplication AS a LEFT JOIN Job AS j ON j.id=a.job_id WHERE a.user_id=$1`, a.UserId)
    if err != nil{
        log.Println(err)
        return applications
    }
    for applicationRow.Next(){
        applicationRow.Scan(&id, &title)
        applications = append(applications, JobApplication{Id: id, Title: title})
    }
    return applications
}

func (a *JobApplication) Delete()error{
    conn, err := GetDBConn()
    if err !=  nil{
        log.Println(err)
        return errors.New("error db conn")
    }
    defer conn.Close()
    _, err = conn.ExecContext(context.Background(), `DELETE FROM JobApplication WHERE id=$1`, a.Id)
    if err != nil{
        log.Println(err)
        return errors.New("error deleting job application")
    }
    return nil
}

