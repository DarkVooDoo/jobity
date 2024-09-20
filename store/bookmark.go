package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Bookmark struct{
    Id string
    Title string
    Addr string
    Metadata []string
    IsThirdParty bool
    ThirdPartyId string
    UserId string
    JobId string
}

func (b *Bookmark) Get()[]Bookmark{
    var jobBookmark []Bookmark
    var job FranceTravailJob
    var  id, title, addr, contract, jobId, thirdPartyId sql.NullString
    var fulltime, thirdParty sql.NullBool
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("erro conn db %v", err)
        return jobBookmark
    }
    defer conn.Close()
    bookmarkRows, err := conn.QueryContext(context.Background(), `SELECT b.id, j.id, j.title, j.postal || ' - ' || j.city, j.contract, j.fulltime, b.third_party, b.id_third_party FROM Bookmark AS b LEFT JOIN Job AS j ON j.id=b.job_id WHERE b.user_id=$1`, b.UserId)
    if err !=  nil{
        log.Printf("error query bookmark %v", err)
        return jobBookmark
    }
    for bookmarkRows.Next(){
        if err := bookmarkRows.Scan(&id, &jobId,  &title, &addr, &contract, &fulltime, &thirdParty, &thirdPartyId); err != nil{
            log.Println(err)
        }
        if thirdParty.Bool{
            if err := GetFranceTravailToken(); err != nil{
                log.Println(err)
                continue
            }
            req, err := http.NewRequest("GET", fmt.Sprintf("https://api.francetravail.io/partenaire/offresdemploi/v2/offres/%v", thirdPartyId.String), nil)
            if err != nil{
                continue
            }
            req.Header.Add("Accept", "application/json")
            req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", AuthToken.Token))
            resp, err := http.DefaultClient.Do(req)
            if err != nil{
                log.Printf("error search offres: %v", err)
                continue
            }
            decoder := json.NewDecoder(resp.Body)
            if err := decoder.Decode(&job); err != nil{
                continue
            }
            jobBookmark = append(jobBookmark, Bookmark{
                Id: id.String,
                Title: job.Title,
                Addr: job.Adresse.Name,
                Metadata: []string{job.Contract, job.Fulltime},
                IsThirdParty: true,
                ThirdPartyId: job.Id,
            })
        }else{
            var fulltimeString string
            if fulltime.Bool == true{
                fulltimeString = "Temps plein"
            }else{
                fulltimeString = "Temps partiel"
            }
            jobBookmark = append(jobBookmark, Bookmark{
                Id: id.String, 
                JobId: jobId.String, 
                Title: title.String, 
                Addr: addr.String, 
                Metadata: []string{fulltimeString, contract.String},
                IsThirdParty: thirdParty.Bool,
                ThirdPartyId: thirdPartyId.String,
            })
        }
    }
    return jobBookmark
}

func (b *Bookmark) Create()error{
    conn, err := GetDBConn()
    if err != nil{
        return errors.New("error getting db conn")
    }
    defer conn.Close()
    var bookmarkRow *sql.Row
    if b.IsThirdParty{
        bookmarkRow = conn.QueryRowContext(context.Background(), `INSERT INTO Bookmark (user_id, third_party, id_third_party) VALUES($1, $2, $3) RETURNING id`, b.UserId, b.IsThirdParty, b.JobId)
    }else{
        bookmarkRow = conn.QueryRowContext(context.Background(), `INSERT INTO Bookmark (job_id, user_id, third_party) VALUES($1, $2, $3) RETURNING id`, b.JobId, b.UserId, b.IsThirdParty)
    }
    if err := bookmarkRow.Scan(&b.Id); err != nil{
        return errors.New("error query creating bookmark")
    }
    return nil
}

func (b *Bookmark) Delete()error{
    conn, err := GetDBConn()
    if err != nil{
        return errors.New("error conn db")
    }
    defer conn.Close()
    _, err = conn.ExecContext(context.Background(), `DELETE FROM Bookmark WHERE id=$1 AND user_id=$2`, b.Id, b.UserId)
    if err != nil{
        return errors.New("error query bookmark delete")
    }
    return nil
}
