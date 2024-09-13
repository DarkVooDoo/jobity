package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/lib/pq"
)

type User struct{
    Id        string 
    Firstname string  `json:"firstname"`
    Lastname  string  `json:"lastname"`
    Age       string  `json:"age"` 
    Description string `json:"description"`
    City      string  `json:"city"`
    Postal    string  `json:"postal"`
    Lat       float64 `json:"lat"`
    Long      float64 `json:"long"`
    Adresse   string  `json:"adresse"`
    Gender    string  `json:"gender"`
    BirthDate string  `json:"birthDate"`
    Email     string
    Cv        string
}

type ProUser struct{
    Id string
    Email string `json:"email"`
    City string 
    Postal string
    Name string
    Siren string `json:"siren"`
    Description string `json:"description"`
    Picture string
    
}

type EntrepriseFetchInfo struct{
    Results []struct{
        Name string `json:"nom_complet"`
        Siren string `json:"siren"`
        Hq struct{
            Naf string `json:"activite_principale"`
            City string `json:"libelle_commune"`
            Postal string `json:"code_postal"`
            Street string `json:"libelle_voie"`
            StreetType string `json:"type_voie"`
            StreetNumber string `json:"numero_voie"`
            Convention []string `json:"liste_idcc"`
        } `json:"siege"`

    } `json:"results"`
}

func (u *User) Create(password string)error{
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("db conn error")
    }
    defer conn.Close()
    cryptPassword, salt := cryptPassword(password)
    _, err = conn.ExecContext(context.Background(), `INSERT INTO Users(email, firstname, lastname, password, salt) VALUES($1,$2,$3,$4,$5)`, u.Email, u.Firstname, u.Lastname, cryptPassword, salt)
    if err != nil{
        log.Println(err)
        return errors.New("error inserting user")
    }
    return nil
}

func (u *User) UploadPhoto(file io.Reader){
    //INsert user photo
}

func (u *User) Modify()error{
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("error conn to the db")
    }
    defer conn.Close()
    _, err = conn.ExecContext(context.Background(), `UPDATE Users SET firstname=INITCAP($1), lastname=INITCAP($2), gender=$3, birthdate=$4, city=$5, postal=$6, description=$7, lat=$8, long=$9 WHERE id=$10`, u.Firstname, u.Lastname, u.Gender, u.BirthDate, u.City, u.Postal, u.Description, u.Lat, u.Long,  u.Id)
    if err != nil{
        log.Println(err)
        return errors.New("error updating users")
    }
    return nil
}

func (u *User) GetProfile()error{
    var city, postal, birthdate, gender, description sql.NullString
    var lat, long sql.NullFloat64
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("error conn to the db")
    }
    defer conn.Close()
    userRow := conn.QueryRowContext(context.Background(), `SELECT firstname, lastname, city, postal, TO_CHAR(birthdate, 'YYYY-MM-DD'), gender, description, lat, long FROM Users WHERE id=$1`, u.Id)
    if err := userRow.Scan(&u.Firstname, &u.Lastname, &city, &postal, &birthdate, &gender, &description, &lat, &long); err != nil{
        log.Println(err)
        return errors.New("error getting user profile")
    }
    u.City = city.String
    u.Postal = postal.String
    u.BirthDate = birthdate.String
    u.Gender = gender.String
    u.Description = description.String
    u.Lat = lat.Float64
    u.Long = long.Float64
    return nil
}

func (p *ProUser) Create(password string)error{
    var company  EntrepriseFetchInfo
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("error conn")
    }
    defer conn.Close()
    req, err := http.NewRequest("GET", fmt.Sprintf("https://recherche-entreprises.api.gouv.fr/search?q=%v", p.Siren), nil)
    req.Header.Add("Accept", "application/json")
    if err != nil{
        log.Println(err)
        return errors.New("error fetch company request")
    }
    res, err := http.DefaultClient.Do(req)
    if err != nil{
        log.Println(err)
        return errors.New("error fetch company response")
    }
    dec := json.NewDecoder(res.Body)
    if err := dec.Decode(&company); err != nil{
        log.Println(err)
        return errors.New("error decoding company data")
    }
    var  myCompany = company.Results[0]
    adresse := fmt.Sprintf("%v %v %v", myCompany.Hq.StreetNumber, myCompany.Hq.StreetType, myCompany.Hq.Street)
    cryptPassword, salt := cryptPassword(password)
    if _, err = conn.ExecContext(context.Background(), `INSERT INTO Entreprise (email, password, salt, siren, naf, name, adresse, city, postal, convention) VALUES($1,$2,$3,$4,$5,UPPER($6),$7,$8,$9,$10)`, p.Email, cryptPassword, salt, myCompany.Siren, myCompany.Hq.Naf, myCompany.Name, adresse, myCompany.Hq.City, myCompany.Hq.Postal, pq.Array(myCompany.Hq.Convention)); err != nil{
        log.Println(err)
        return errors.New("error creating new pro user")
    }
    return nil
}

func (p *ProUser) GetProfile()(error){
    var picture, description sql.NullString
    var name, siren string
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("error db conn")
    }
    defer conn.Close()
    profileRows, err := conn.QueryContext(context.Background(), `SELECT e.name, e.siren, e.description, e.picture FROM Entreprise AS e WHERE e.id=$1`, p.Id)
    if err != nil{
        log.Println(err)
        return errors.New("error postgres query")
    }
    for profileRows.Next(){
        if err := profileRows.Scan(&name, &siren, &description, &picture); err != nil{
            log.Println(err)
        }
    }
    p.Name = name
    p.Siren = siren
    p.Description = description.String
    p.Picture = picture.String
    return nil
}

func cryptPassword(password string)(cryptPassword string, salt int){
    salt = rand.Intn(9999-1000) + 1000
    cypher := sha256.Sum256([]byte(fmt.Sprintf("%v%v%v", password,salt,os.Getenv("PASSWORD_SECRET_KEY"))))
    cryptPassword = fmt.Sprintf("%x", cypher)
    return cryptPassword, salt
}

