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
    Addr string
    Status string
    JobId string
    UserId string
}

var statusMap = map[string]int{"Non vue":  0, "Vue": 1, "Reject": 2, "Possible": 3, "Interview": 4}

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

func (a *JobApplication) GetUserApplications()(applications []JobApplication){
    var id, title, status, addr string
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return 
    }
    defer conn.Close()
    applicationRow, err := conn.QueryContext(context.Background(), `SELECT a.id, j.title, a.status, LEFT(j.postal, 2) || ' - ' || j.city FROM JobApplication AS a LEFT JOIN Job AS j ON j.id=a.job_id WHERE a.user_id=$1`, a.UserId)
    if err != nil{
        log.Println(err)
        return 
    }
    for applicationRow.Next(){
        applicationRow.Scan(&id, &title, &status, &addr)
        applications = append(applications, JobApplication{Id: id, Title: title, Status: status, Addr: addr})
    }
    return 
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

func (a JobApplication) UpdateStatus(status string)(error){
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error conn db %v", err)
        return errors.New("error db conn")
    }
    _, err = conn.ExecContext(context.Background(), `UPDATE JobApplication SET status=$1 WHERE job_id=$2 AND user_id=$3 AND status < $4`, status, a.JobId, a.UserId, status)
    if err != nil{
        log.Printf("error updating application status \n%v", err)
        return errors.New("error updating application")
    }  
    return nil
}

