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
    Type string
    InterviewDate string
    UserName string
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

func (a *JobApplication) Interviews(entrepriseId string) []JobApplication{
    var interviews []JobApplication
    var inter JobApplication
    var date sql.NullString
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error in the conn %v", err)
        return interviews
    }
    defer conn.Close()
    interviewsRows, err := conn.QueryContext(context.Background(), `SELECT DISTINCT j.title, ja.location, ja.interview_date, ja.interview_type, ja.id, CONCAT(u.firstname, ' ', u.lastname), j.id, ja.user_id FROM Job AS j RIGHT JOIN JobApplication AS ja ON j.id=ja.job_id AND ja.status='Interview' RIGHT JOIN Users AS u ON u.id=ja.user_id WHERE j.entreprise_id=$1`, entrepriseId)
    if err != nil{
        log.Printf("error in the query: %v", err)
    }
    for interviewsRows.Next(){
        if err := interviewsRows.Scan(&inter.Title, &inter.Addr, &date, &inter.Type, &inter.Id, &inter.UserName, &inter.JobId, &inter.UserId); err != nil{
            log.Printf("error scan %v", err)
        }
        inter.InterviewDate = date.String[:len(date.String)-4]
        interviews = append(interviews, inter)
    }
    return interviews
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
    if status == "Interview"{
        _, err = conn.ExecContext(context.Background(), `UPDATE JobApplication SET status=$1, interview_date=$2, location=$3, interview_type=$4 WHERE id=$5`, status, a.InterviewDate, a.Addr, a.Type, a.Id)
    }else{
        _, err = conn.ExecContext(context.Background(), `UPDATE JobApplication SET status=$1 WHERE id=$2 AND status < $3`, status, a.Id, status)
    }
    if err != nil{
        log.Printf("error updating application status \n%v", err)
        return errors.New("error updating application")
    }  
    return nil
}









