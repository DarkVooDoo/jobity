package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Job struct{
    Id string
    Title string `json:"title"`
    Description string `json:"description"`
    Salary []float64 `json:"salary"`
    SalaryString string
    Contract string `json:"contract"`
    ContractArray string
    EntrepriseName string
    EntrepriseId string 
    Experience int `json:"exp"`
    Advantage []string `json:"advantage"`
    SkillNeeded []Skill `json:"skill"`
    Date string
    City string `json:"city"`
    Postal string `json:"postal"`
    Lat float64   `json:"lat"`
    Long float64 `json:"long"`
    FullAdresse string
    Fulltime bool
    FulltimeString string
    Applied bool
    WeeklyWorkTime float64 `json:"weeklyWorkTime"`
    StartDay string `json:"startTime"`
    IsThirdParty bool
    ApplicationCount int
    RecomendationVector []float64
    ThirdPartyLink string
    TemplateName string `json:"tname"`
    BookmarkId string
}

type Filter struct{
    Query string `json:"q"`
    Postal string `json:"postal"`
    City string `json:"city"`
    Experience string `json:"exp"`
    Contract string `json:"contract"`
    Fulltime bool `json:"fulltime"`
    ThirdParty bool `json:"thirdParty"`
    Lat float64 `json:"lat"`
    Long float64 `json:"long"`
    ExperienceMin string `json:"minExp"`
    MinSalary float64 `json:"minSalary"`
    MaxSalary float64 `json:"maxSalary"`
    Order string `json:"order"`
    Distance float64 `json:"distance"`
}

type Skill struct{
    Label string `json:"label"`
    Required bool `json:"required"`
}

type EntrepriseTemplates struct{
    Id string
    Name string
}

type Contract struct{
    Id int `json:"id"`
    Name string `json:"value"`
}

type FranceTravailToken struct{
    Token string `json:"access_token"`
    Type string `json:"token_type"`
    Scope string `json:"scope"`
    ExpireIn int `json:"expires_in"`
    ExpirationDate time.Time
}

type Payload struct{
    Result []FranceTravailJob `json:"resultats"`
}

type FranceTravailJob struct{
    Id string `json:"id"`
    Title string `json:"intitule"`
    Description string `json:"description"`
    Contract string `json:"typeContrat"`
    Date string `json:"dateCreation"`
    WeeklyWorkTime string `json:"dureeTravailLibelle"`
    Fulltime string `json:"dureeTravailLibelleConverti"`
    Salary struct{
        Label string `json:"commentaire"`
        Amount string `json:"libelle"`
    } `json:"salaire"`
    Adresse struct{
        Name string `json:"libelle"`
        Postal string `json:"codePostal"`
        Commune string `json:"commune"`
    } `json:"lieuTravail"`
    Enpreprise struct{
        Name string `json:"nom"`
    } `json:"entreprise"`
    Contact struct{
        Lien string `json:"coordonnees1"`
        Url string `json:"urlPostulation"`
    } `json:"contact"`
}

var AuthToken FranceTravailToken

var cache struct{
    Jobs []Job
    Time time.Time
    ExpireIn int
}

func GetEmplois(recomendation []float64) (ftJobs []Job, ftRecomendationJobs []Job, jobs []Job){
    jobs = GetAppJobs(recomendation)
    ftJobs, ftRecomendationJobs, _ = GetFranceTravailEmplois(len(jobs))
    jobs = append(jobs, ftJobs...)
    rand.Shuffle(len(jobs), func(i, j int) {{jobs[i], jobs[j] = jobs[j], jobs[i]}})
    return ftJobs, ftRecomendationJobs, jobs
}

func GetRecomendationJobs()[]Job{
    //Get App Recomendation jobs
    var jobs []Job   
    return jobs
}

func GetAppJobs(recomendationVector []float64)[]Job{
    var jobRows *sql.Rows
    var jobs []Job
    var sal sql.NullString
    var fulltime bool
    var title, fulltimeString, contract, date, id, salaryString, entrepriseName, fullAddr string
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return jobs
    }
    defer conn.Close()
    query := ""
    if len(recomendationVector) > 0{
        query = `SELECT j.id, j.title, j.salary, LEFT(j.postal, 2) || ' - ' || j.city, c.name, TO_CHAR(AGE(NOW(), j.created), 'Y-MM-DD-HH24-MI-SS'), e.name, j.fulltime FROM Job AS j LEFT JOIN Entreprise AS e ON j.entreprise_id=e.id LEFT JOIN Contract AS c ON c.id=j.contract WHERE Recomendation($1,vector) >= 0.8 LIMIT 10`
        rows, err := conn.QueryContext(context.Background(), query, pq.Array(recomendationVector))
        if err != nil{
            log.Println(err)
        }
        jobRows = rows
    }else{
        query = `SELECT j.id, j.title, j.salary, LEFT(j.postal, 2) || ' - ' || j.city, c.name, TO_CHAR(AGE(NOW(), j.created), 'Y-MM-DD-HH24-MI-SS'), e.name, j.fulltime FROM Job AS j LEFT JOIN Entreprise AS e ON j.entreprise_id=e.id LEFT JOIN Contract AS c ON c.id=j.contract ORDER BY RANDOM() LIMIT 10`
        rows, err := conn.QueryContext(context.Background(), query)
        if err != nil{
            log.Println(err)
        }
        jobRows = rows
    }
    if err != nil{
        log.Println(err)
        return jobs
    }
    if err != nil{
        log.Printf("Query error %v", err)
        return jobs
    }
    for jobRows.Next(){
        if err := jobRows.Scan(&id, &title, &sal, &fullAddr, &contract, &date, &entrepriseName, &fulltime); err != nil{
            log.Println(err)
        }
        formatedDate := PostgresIntervalIntoString(strings.Split(date, "-"))
        if fulltime{
            fulltimeString = "Temps plein"
        }else{
            fulltimeString = "Temps partiel"
        }

        jobs = append(jobs, Job{
            Id: id, 
            Title: title, 
            Contract: contract, 
            Date: formatedDate, 
            FullAdresse: fullAddr, 
            FulltimeString: fulltimeString,
            SalaryString: salaryString,
            EntrepriseName: entrepriseName,
        })
    }
    return jobs
}

func (j *Job) CreateJob()(string, error){
    conn, err := GetDBConn()   
    if err != nil{
        log.Println(err)
        return "", errors.New("error db conn")
    }
    defer conn.Close()
    salary, _ := json.Marshal(j.Salary)
    h := fmt.Sprintf("{%v}", string(salary)[1:len(salary)-1])
    skill, _ := json.Marshal(j.SkillNeeded)
    if j.WeeklyWorkTime >= 35{
        j.Fulltime = true
    }
    jobRow := conn.QueryRowContext(context.Background(), `INSERT INTO Job(title, description, salary, city, postal, contract, worktime, advantage, skill, fulltime, lat, long, experience, entreprise_id) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING id`, j.Title, j.Description, h, j.City, j.Postal, j.Contract, j.WeeklyWorkTime, pq.Array(j.Advantage), string(skill), j.Fulltime, j.Lat, j.Long, j.Experience, j.EntrepriseId)
    if err = jobRow.Scan(&j.Id); err != nil{
        log.Println(err)
        return "", errors.New("error inserting job to the db")
    }
    //if j.TemplateName != ""{
    //    //Save as template
    //    if _, err = conn.ExecContext(context.Background(), `INSERT INTO JobTemplate (name, entreprise_id, job_id) VALUES($1,$2,$3)`, j.TemplateName, j.EntrepriseName, j.Id); err != nil{
    //        log.Println(err)
    //        return "", errors.New("error inserting to jobtemplate")
    //    }
    //}
    return j.Id, nil
}

func (j *Job) ModifyJob()error{
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("error conn to the db")
    }
    defer conn.Close()
    salaryMarshal, _ := json.Marshal(j.Salary)
    advantage, _ := json.Marshal(j.Advantage)
    skill, _ := json.Marshal(j.SkillNeeded)
    if j.WeeklyWorkTime >= 35{
        j.Fulltime = true
    }
    salary := fmt.Sprintf("{%v}", string(salaryMarshal)[1:len(salaryMarshal)-1])
    _, err = conn.ExecContext(context.Background(), `UPDATE Job SET title=$1, city=$2, postal=$3, contract=$4, description=$5, salary=$6, worktime=$7, advantage=$8, skill=$9, fulltime=$10 WHERE id=$11 AND entreprise_id=$12`,j.Title, j.City, j.Postal, j.Contract, j.Description, salary, j.WeeklyWorkTime, string(advantage), string(skill), j.Fulltime, j.Id, j.EntrepriseId)
    if err != nil{
        log.Println(err)
        return errors.New("error updating job")
    }
    return nil
}

func (j *Job) DeleteJob()error{
    conn, err := GetDBConn()
    if err != nil{
        return errors.New("error conn to the db")
    }
    _, err = conn.ExecContext(context.Background(), `DELETE FROM Job WHERE id=$1 AND entreprise_id=$2`, j.Id, j.EntrepriseId)
    if err != nil{
        return errors.New("error deleting job")
    }
    return nil
}

func (j *Job) GetJobById()error{
    var skill, applicationId, bookmarkId  sql.NullString
    var sal []float64
    var dateAge string
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("error conn to the db")
    }
    defer conn.Close()
    if err != nil{
        log.Println(err)
        return errors.New("error transction")
    }
    jobRow := conn.QueryRowContext(context.Background(), `SELECT UPPER(e.name) , TO_CHAR(AGE(NOW(), j.created), 'Y-MM-DD-HH24-MI-SS'), j.title, j.description, c.name, j.city, j.postal,  CONCAT(LEFT(j.postal, 2), ' - ', J.city), j.salary, j.advantage, j.skill, j.worktime, j.fulltime, j.lat, j.long, j.experience, j.vector, ja.id, b.id FROM Job AS j LEFT JOIN Entreprise AS e ON j.entreprise_id=e.id LEFT JOIN JobApplication AS ja ON ja.job_id=j.id LEFT JOIN Bookmark AS b ON j.id=b.job_id LEFT JOIN Contract AS c ON c.id=j.contract WHERE j.id=$1`, j.Id)
    if err := jobRow.Scan(&j.EntrepriseName, &dateAge, &j.Title, &j.Description, &j.Contract, &j.City, &j.Postal, &j.FullAdresse, pq.Array(&sal), pq.Array(&j.Advantage), &skill, &j.WeeklyWorkTime, &j.Fulltime, &j.Lat, &j.Long,&j.Experience, pq.Array(&j.RecomendationVector),  &applicationId, &bookmarkId); err != nil{
        log.Println(err)
        return errors.New("error scanning job")
    }
    if applicationId.String != "" {j.Applied = true}
    j.Description = strings.ReplaceAll(j.Description, `\n`, "\n")
    j.BookmarkId = bookmarkId.String
    j.Date = PostgresIntervalIntoString(strings.Split(dateAge, "-"))
    j.SalaryString, j.Salary = postgresSalaryIntoString(sal)
    j.ContractArray =  GetContracts()
    if err := json.Unmarshal([]byte(skill.String), &j.SkillNeeded); err != nil{
        log.Println(err)
        return errors.New("error making skill into json")
    }
    return nil
}

func GetEntrepriseJobs(entrepriseId string)[]Job{
    var jobs []Job
    var title, id, city, postal, contract string
    var count int
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return jobs
    }
    defer conn.Close()
    entrepriseJobRows, err := conn.QueryContext(context.Background(), `SELECT j.title, j.id, LEFT(j.postal, 2), j.city, (SELECT count(id) FROM JobApplication WHERE status <= 'Vue' AND j.id=job_id), c.name FROM Job AS j LEFT JOIN Contract AS c ON c.id=j.contract WHERE j.entreprise_id=$1`, entrepriseId)
    if err != nil{
        log.Println(err)
        return jobs
    }
    for entrepriseJobRows.Next(){
        entrepriseJobRows.Scan(&title, &id, &postal, &city, &count, &contract)
        jobs = append(jobs, Job{Title: title, Id: id, Postal: postal, City: city, ApplicationCount: count, Contract: contract})
    }
    return jobs
}

func GetEntrepriseJobCards(entrepriseId string)[]Job{
    var date, fulltimeString string
    var jobList []Job
    var job Job
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error in the db: %v", err)
        return jobList
    }
    defer conn.Close()

    jobRows, err := conn.QueryContext(context.Background(), `SELECT j.id, j.title, j.salary, LEFT(j.postal, 2) || ' - ' || j.city, c.name, TO_CHAR(AGE(NOW(), j.created), 'Y-MM-DD-HH24-MI-SS'), e.name, j.fulltime FROM Job AS j LEFT JOIN Entreprise AS e ON j.entreprise_id=e.id LEFT JOIN Contract AS c ON c.id=j.contract WHERE j.entreprise_id=$1`, entrepriseId)
    if err != nil{
        log.Printf("error in the query: %v", err)
        return jobList
    }
    for jobRows.Next(){
        if err := jobRows.Scan(&job.Id, &job.Title, pq.Array(&job.Salary), &job.FullAdresse, &job.Contract, &date, &job.EntrepriseName, &job.Fulltime); err != nil{
            log.Printf("error scan: %v", err)
        }
        formatedDate := PostgresIntervalIntoString(strings.Split(date, "-"))
        if job.Fulltime{
            fulltimeString = "Temps plein"
        }else{
            fulltimeString = "Temps partiel"
        }
        job.FulltimeString = fulltimeString
        job.Date = formatedDate
        job.SalaryString, _ = postgresSalaryIntoString(job.Salary)

        jobList = append(jobList, job)
    }
    return jobList
}

func (j *Job) GetJobByTemplateId(id string)error{
    var advantage, skill sql.NullString
    var sal []float64
    var dateAge string
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("error conn to the db")
    }
    defer conn.Close()
    jobRow := conn.QueryRowContext(context.Background(), `SELECT UPPER(e.name) , TO_CHAR(AGE(NOW(), j.created), 'Y-MM-DD-HH24-MI-SS'), j.title, j.description, j.contract, j.city, j.postal, j.salary, j.advantage, j.skill, j.worktime FROM Job AS j LEFT JOIN Entreprise AS e ON j.entreprise_id=e.id LEFT JOIN JobTemplate AS jt ON jt.job_id=j.id WHERE jt.id=$1`, id)
    if err := jobRow.Scan(&j.EntrepriseName, &dateAge, &j.Title, &j.Description, &j.Contract, &j.City, &j.Postal, pq.Array(&sal), &advantage, &skill, &j.WeeklyWorkTime); err != nil{
        log.Println(err)
        return errors.New("error scanning job")
    }
    j.SalaryString, j.Salary = postgresSalaryIntoString(sal)
    j.Date = PostgresIntervalIntoString(strings.Split(dateAge, "-"))
    if err := json.Unmarshal([]byte(advantage.String), &j.Advantage); err != nil{
        log.Println(err)
        return errors.New("error making advantage into json")
    }
    if err := json.Unmarshal([]byte(skill.String), &j.SkillNeeded); err != nil{
        log.Println(err)
        return errors.New("error making skill into json")
    }
    return nil
}

func (j *Job)SaveAsTemplate(name string)error{
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return errors.New("error conn to the db")
    }
    defer conn.Close()
    _, err = conn.ExecContext(context.Background(), `INSERT INTO JobTemplate (name, job_id, entreprise_id) VALUES($1,$2,$3)`, name, j.Id, j.EntrepriseId)
    if err != nil{
        log.Println(err)
        return errors.New("error creating job template")
    }
    return nil
}

func (j *Job)GetTemplates()([]EntrepriseTemplates, error){
    var templates []EntrepriseTemplates
    var name, id string
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return templates, errors.New("error db conn")
    }
    defer conn.Close()
    rows, err := conn.QueryContext(context.Background(), `SELECT id, name From JobTemplate WHERE entreprise_id=$1`, j.EntrepriseId)
    if err != nil{
        log.Println(err)
        return templates, errors.New("error selecting job templates")
    }
    for rows.Next(){
        if err := rows.Scan(&id, &name); err != nil{
            log.Println(err)
        }
        templates = append(templates, EntrepriseTemplates{Id: id, Name: name})
    }
    return templates, nil
}

func SearchJobWithFilter(query string, filter Filter)[]Job{
    var jobs []Job
    var name, addr, id, contract, date, entreprise, fulltimeString string
    var fulltime bool
    var sal []float64
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return jobs
    }
    defer conn.Close()
    rows, err := conn.QueryContext(context.Background(), query)
    if err != nil{
        log.Printf("error query: %v",err)
        return jobs
    }
    for rows.Next(){
        rows.Scan(&name, &addr, &entreprise, &contract, &fulltime, &date, &id, pq.Array(&sal))
        salaryString, _ := postgresSalaryIntoString(sal)
        formatedDate := PostgresIntervalIntoString(strings.Split(date, "-"))
        if fulltime{
            fulltimeString = "Temps plein"
        }else{
            fulltimeString = "Temps partiel"
        }
        jobs = append(jobs, Job{
            Title: name, 
            FullAdresse: addr, 
            EntrepriseName: entreprise, 
            Id: id, 
            FulltimeString: fulltimeString, 
            SalaryString: salaryString, 
            Contract: contract, 
            Date: formatedDate,
        })

    }

    if filter.ThirdParty{
        var franceTravailJobs Payload
        startRange := len(jobs)
        total := 9 - startRange
        req, err := http.NewRequest("GET", fmt.Sprintf(`https://api.francetravail.io/partenaire/offresdemploi/v2/offres/search?commune=%v&range=%v-%v&motsCles=%v&salaireMin=%v&tempsPlein=%v&distance=%v`, filter.Postal, startRange, total, filter.Query, filter.MinSalary, filter.Fulltime, filter.Distance), nil)
        if err != nil{
            log.Printf("error third party: %v", err)
        }
        req.Header.Add("Accept", "application/json")
        req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", AuthToken.Token))
        res, err := http.DefaultClient.Do(req)
        if err != nil{
            log.Println(err)
            return jobs        
        }
        dec := json.NewDecoder(res.Body)
        if err := dec.Decode(&franceTravailJobs); err != nil{
            log.Println(err)
            return jobs        
        }
        for _, value := range franceTravailJobs.Result{
            date := DateDifference(value.Date)
            jobs = append(jobs, Job{
                Id: value.Id,
                Title: value.Title, 
                FullAdresse: value.Adresse.Name,
                SalaryString: value.Salary.Amount,
                FulltimeString: value.Fulltime,
                EntrepriseName: value.Enpreprise.Name,
                IsThirdParty: true, 
                Contract: value.Contract,
                Date:  date,
            })
        }
    }
    return jobs
}

func GetJobBySearch(query string, postal string, startRange int, appLastPosition int)(jobs []Job, ftOffset int, appPosition int,  err error){
    var id, title, contract, date, addrValue, salaryString, entrepriseName string
    var fulltime bool
    var franceTravailJobs Payload
    var fulltimeString string
    var job Job
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
    }
    if err == nil{
        defer conn.Close()
        jobRows, err := conn.QueryContext(context.Background(), `SELECT j.id, j.title, j.salary, LEFT(j.postal, 2) || ' - ' || j.city, j.contract, TO_CHAR(AGE(NOW(), j.created), 'Y-MM-DD-HH24-MI-SS'), e.name, j.fulltime FROM Job AS j LEFT JOIN Entreprise AS e ON e.id=j.entreprise_id WHERE ts @@ websearch_to_tsquery('french', $1) AND LEFT(j.postal, 2)=$2 AND j.id > $3  ORDER BY j.created DESC LIMIT 10`, query, postal[0:2], appLastPosition)
        if err != nil{
            log.Println(err)
        }
        appPosition = appLastPosition
        for jobRows.Next(){
            if err := jobRows.Scan(&id, &title, pq.Array(&job.Salary), &addrValue, &contract, &date, &entrepriseName, &fulltime); err != nil{
                log.Println(err)
            }
            salaryString, _ = postgresSalaryIntoString(job.Salary)
            formatedDate := PostgresIntervalIntoString(strings.Split(date, "-"))
            if fulltime{
                fulltimeString = "Temps plein"
            }else{
                fulltimeString = "Temps partiel"
            }
            jobs = append(jobs, Job{
                Id: id, 
                Title: title, 
                Contract: contract, 
                Date: formatedDate, 
                FullAdresse: addrValue, 
                FulltimeString: fulltimeString,
                SalaryString: salaryString,
                EntrepriseName: entrepriseName,
            })
            appPosition, _ = strconv.Atoi(id)
        }
    }
    total := 9 - len(jobs) + startRange
    if err := GetFranceTravailToken(); err != nil{
        log.Println("we dont have token")
    }
    query = strings.ReplaceAll(query, " ", "+")
    req, err := http.NewRequest("GET", fmt.Sprintf(`https://api.francetravail.io/partenaire/offresdemploi/v2/offres/search?commune=%v&range=%v-%v&motsCles=%v&distance=20`, postal, startRange+1, total, query), nil)
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", AuthToken.Token))
    if err != nil{
        log.Println(err)
        return jobs, total, appPosition, errors.New("error requesting job query")
    }
    res, err := http.DefaultClient.Do(req)
    if err != nil{
        log.Println(err)
        return jobs, total, appPosition, errors.New("errors response to france travail")
    }
    dec := json.NewDecoder(res.Body)
    if err := dec.Decode(&franceTravailJobs); err != nil{
        log.Println(err)
        return jobs, total, appPosition, errors.New("error decoding france travail jobs")
    }
    for _, value := range franceTravailJobs.Result{
        date := DateDifference(value.Date)
        jobs = append(jobs, Job{
            Id: value.Id,
            Title: value.Title, 
            FullAdresse: value.Adresse.Name,
            SalaryString: value.Salary.Amount,
            FulltimeString: value.Fulltime,
            EntrepriseName: value.Enpreprise.Name,
            IsThirdParty: true, 
            Contract: value.Contract,
            Date:  date,
        })
    }
    return jobs, total, appPosition,  nil
}

func TemplatesIntoString(templates []EntrepriseTemplates)string{
    stringTemplate := `[{"id": "", "value": " "}`
    for _, v := range templates{
        stringTemplate += fmt.Sprintf(`,{"id": "%v", "value": "%v"}`, v.Id, v.Name)
    }
    return stringTemplate[:]+"]"
}

func GetFranceTravailJobById(id string)(Job, error){
    var job FranceTravailJob
    var convertedJob Job
    var bookmarkId sql.NullString
    if err := GetFranceTravailToken(); err != nil{
        log.Println(err)
        return convertedJob, errors.New("error request token")
    }
    req, err := http.NewRequest("GET", fmt.Sprintf("https://api.francetravail.io/partenaire/offresdemploi/v2/offres/%v", id), nil)
    if err != nil{
        log.Println(err)
        return convertedJob, errors.New("error request")
    }
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", AuthToken.Token))
    res, err := http.DefaultClient.Do(req)
    if err != nil{
        log.Println(err)
        return convertedJob, errors.New("error doing request")
    }
    decoder := json.NewDecoder(res.Body)
    if err := decoder.Decode(&job); err != nil{
        log.Println(err)
        return convertedJob, errors.New("error decoding request")
    }
    var weeklyWorkTime float64 = 0
    if job.WeeklyWorkTime != ""{
        weeklyWorkTime, err = strconv.ParseFloat(job.WeeklyWorkTime[:2], 64)
        if err != nil{
            weeklyWorkTime = 0
        }
    }
    conn, err := GetDBConn()
    if err != nil{
        log.Printf("error conn db")
    }
    defer conn.Close()
    bookmarkRow := conn.QueryRowContext(context.Background(), `SELECT id FROM Bookmark WHERE id_third_party=$1`, id)
    if err := bookmarkRow.Scan(&bookmarkId); err != nil{
        log.Println(err)
    }
    convertedJob.BookmarkId = bookmarkId.String
    convertedJob.Title = job.Title
    convertedJob.Id = job.Id
    convertedJob.Description = job.Description
    convertedJob.SalaryString = job.Salary.Amount
    convertedJob.Contract = job.Contract
    convertedJob.WeeklyWorkTime = weeklyWorkTime
    convertedJob.Fulltime = strings.HasPrefix(job.Fulltime, "Temps plein")
    convertedJob.FullAdresse = job.Adresse.Name
    convertedJob.IsThirdParty = true
    convertedJob.Postal = job.Adresse.Postal
    convertedJob.EntrepriseName = job.Enpreprise.Name
    startIndex := strings.Index(job.Contact.Lien, "http")
    if startIndex > -1{
        convertedJob.ThirdPartyLink = job.Contact.Lien[startIndex:]
    }else{
        convertedJob.ThirdPartyLink = job.Contact.Lien
    }
    convertedJob.Date = DateDifference(job.Date)
    return convertedJob, nil
}

func GetFranceTravailToken()( error){
    var token FranceTravailToken
    if AuthToken.ExpirationDate.After(time.Now()) {
        return nil
    }
    values := url.Values{}
    values.Set("grant_type", "client_credentials")
    values.Set("client_id", "PAR_ndepart_39b9519b2f103a4d00d809e6ae5d5607ea073add8a9946af3be179fa9382db2e")
    values.Set("client_secret", "6f833954cb4be784d61b2ad0b394f2a1f805f06055738488479672a7796cfe85")
   //values.Set("client_id", os.Getenv("FT_CLIENT_ID"))
   //values.Set("client_secret", os.Getenv("FT_SECRET_KEY"))
    values.Set("scope", "api_offresdemploiv2 o2dsoffre")
    res, err := http.PostForm("https://entreprise.francetravail.fr/connexion/oauth2/access_token?realm=/partenaire", values)
    if err != nil{
        return errors.New("error request")
    }
    defer res.Body.Close()
    dec := json.NewDecoder(res.Body)
    if err := dec.Decode(&token); err != nil{
        return errors.New("error decoding token")
    }
    dur, _ := time.ParseDuration(fmt.Sprintf("%vs", token.ExpireIn))
    token.ExpirationDate = time.Now().Add(dur)
    AuthToken = token
    return nil
}

func GetFranceTravailEmplois(jobLength int)(mostRecentJobs []Job, recomendationJobs []Job,  requestError error){
    mostRecentJobs, requestError = getFranceTravailFrontpageJobs(fmt.Sprintf("https://api.francetravail.io/partenaire/offresdemploi/v2/offres/search?origineOffre=1&range=0-%v", 9-jobLength))
    if requestError != nil{
        log.Println(requestError)
    }
    //recomendationJobs, requestError = getFranceTravailFrontpageJobs(fmt.Sprintf("https://api.francetravail.io/partenaire/offresdemploi/v2/offres/search?origineOffre=1&range=0-%v&departement=%v&typeContrat=%v", 9-jobLength, recomendation.Postal, recomendation.Contract))
    return mostRecentJobs, recomendationJobs, nil
}

func GetContracts()string{
    var contractArray []Contract
    var contract Contract
    conn, err := GetDBConn()
    if err != nil{
        log.Println(err)
        return ""
    }
    defer conn.Close()
    row, err := conn.QueryContext(context.Background(), `SELECT id, name FROM Contract`)
    if err != nil{
        log.Println(err)
        return ""
    }
    for row.Next(){
        if err := row.Scan(&contract.Id, &contract.Name); err != nil{
            log.Println(err)
            return ""
        }
        contractArray = append(contractArray, contract)
    }
    contractMarshal, _ := json.Marshal(contractArray)
    return string(contractMarshal)
}

func postgresSalaryIntoString(salary []float64)(string, []float64){
    if len(salary) > 1{
        return fmt.Sprintf("Salaire entre %v et %v euros", salary[0], salary[1]), salary
    }else{
        return fmt.Sprintf("Salaire %v euros", salary[0]), salary
    }
}

func DateDifference(date string)string{
    jobDate, _ := time.Parse(time.RFC3339, fmt.Sprintf("%v+02:00", date[:len(date)-1]))
    diff := time.Since(jobDate)
    if fmt.Sprintf("%0.f", diff.Hours()) != "0"{
        date = fmt.Sprintf("Il y a %0.f heures", diff.Hours())
    }else if fmt.Sprintf("%0.f", diff.Minutes()) != "0"{
        date = fmt.Sprintf("Il y a %0.f minute", diff.Minutes())
    }else{
        date = fmt.Sprintf("Il y a %.0f seconde", diff.Seconds())
    }
    return date
}

//Conver Postgres interval into string format: Y-MM-DD-HH24-MI-SS
func PostgresIntervalIntoString(interval []string) string{

    labels := []string{"ans", "mois", "jours", "heures", "minutes", "secondes"}
    age := "0"
    for i, v := range interval{
        if v == "0" || v == "00"{
            continue
        }
        ageInt, _ := strconv.Atoi(v)
        age = fmt.Sprintf("Il y a %v %v", ageInt, labels[i])
        break
    }
    return age
}

func getFranceTravailFrontpageJobs(request string)([]Job, error){
    var job Payload
    var jobs []Job
    if len(cache.Jobs) > 0 && time.Since(cache.Time).Minutes() < 5{
        return cache.Jobs, nil
    }
    if err := GetFranceTravailToken(); err != nil{
        log.Println(err)
        return jobs, errors.New("error request france travail token")
    }
    req, err := http.NewRequest("GET", request, nil)
    if err != nil{
        return nil, errors.New("error in api request")
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", AuthToken.Token))
    resp, err := http.DefaultClient.Do(req)
    if err != nil{
        log.Printf("error search offres: %v", err)
        return nil, errors.New("error on request")
    }
    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&job); err != nil{
        return nil, errors.New("error decoding france travail payload")
    }
    for _, value := range job.Result{
        date := DateDifference(value.Date)
        jobs= append(jobs, Job{
            Id: value.Id,
            Title: value.Title, 
            Description: value.Description, 
            SalaryString: value.Salary.Amount,
            EntrepriseName: value.Enpreprise.Name,
            FullAdresse: value.Adresse.Name,
            FulltimeString: value.Fulltime,
            IsThirdParty: true, 
            Contract: value.Contract,
            Date:  date,
        })
    }
    
    cache = struct{Jobs []Job; Time time.Time; ExpireIn int}{Jobs: jobs, Time: time.Now(), ExpireIn: 5*60}
    return jobs, nil
}
