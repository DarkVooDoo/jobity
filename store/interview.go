package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type Interview struct{
    Id string `json:"id"`
    Title string
    Addr string `json:"addr"`
    UserName string `json:"name"`
    Date string `json:"interview_date"`
    Type string `json:"type"`
    UserId string `json:"userId"`
    JobId string `json:"jobId"`
}

func (i *Interview) Create(conn *sql.Conn) error{
    interviewRow := conn.QueryRowContext(context.Background(), `INSERT INTO Interview ("date", "type", adresse, user_id, job_id) VALUES($1, $2, $3, $4, $5) RETURNING id`, i.Date, i.Type, i.Addr, i.UserId, i.JobId)
    if err := interviewRow.Scan(&i.Id); err != nil{
        log.Printf("error creating the interview: %v", err)
        return errors.New("error creating the interview")
    }
    return nil
}

func (i *Interview) Interviews(conn *sql.Conn, entrepriseId string) []Interview{
    var interviews []Interview
    var inter Interview
    var adr sql.NullString
    interviewsRows, err := conn.QueryContext(context.Background(), `SELECT j.title, i.date, i.type, i.id, CONCAT(u.firstname, ' ', u.lastname), i.adresse, j.id, u.id FROM Job AS j RIGHT JOIN Interview AS i ON j.id=i.job_id RIGHT JOIN Users AS u ON u.id=i.user_id WHERE j.entreprise_id=$1`, entrepriseId)
    if err != nil{
        log.Printf("error in the query: %v", err)
    }
    for interviewsRows.Next(){
        if err := interviewsRows.Scan(&inter.Title, &inter.Date, &inter.Type, &inter.Id, &inter.UserName, &adr, &inter.JobId, &inter.UserId); err != nil{
            log.Printf("error scan %v", err)
        }
        inter.Addr = adr.String
        inter.Date = inter.Date[:len(inter.Date)-4]
        interviews = append(interviews, inter)
    }
    return interviews
}

func (i *Interview) Modify(conn *sql.Conn)error{
    result, err := conn.ExecContext(context.Background(), `UPDATE Interview SET "date"=$1, "type"=$2, adresse=$3 WHERE id=$4`, i.Date, i.Type, i.Addr, i.Id)
    if err != nil {
        log.Printf("error executing query: %v", err)
        return errors.New("errors query")
    }
    affected, err := result.RowsAffected()
    if err != nil || affected == 0{
        log.Printf("error no rows affected: %v", err)
        return errors.New("error affected")
    }
    return nil
}
