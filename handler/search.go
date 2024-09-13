package handler

import (
	"fmt"
	"job/store"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"text/template"
)

type SearchPayload struct{
    RequireData   
    Jobs []store.Job
    Contract string
}


var SearcHandler = func(res http.ResponseWriter, req *http.Request){
    route, _ := NewRoute(res, req)
    var filter store.Filter
    route.Get(func() {
        query := route.Request.URL.Query().Get("q")
        lat := route.Request.URL.Query().Get("la")
        long := route.Request.URL.Query().Get("lo")
        city := route.Request.URL.Query().Get("city")
        postal := route.Request.URL.Query().Get("postal")
        appLastPosition, _ := strconv.Atoi(route.Request.URL.Query().Get("lp"))
        startRange, _ := strconv.Atoi(route.Request.URL.Query().Get("sr"))
        jobs,ftOffset, lastPosition, err := store.GetJobBySearch(query, postal, startRange, appLastPosition)
        if err != nil{
            log.Printf("Q paso %v", err)
        }
        contract := store.GetContracts()
        searchData := SearchQuery{
            Query: query, 
            Postal: postal, 
            FrancetravailPosition: ftOffset, 
            Lastposition: lastPosition, 
            Lat: lat, 
            Long: long, 
            City: city,
        }
        payload := SearchPayload{
            RequireData: RequireData{
                route.User,
                searchData,
            },
            Jobs: jobs,
            Contract: contract,
        }

        route.Render(payload, "route/template.html", "route/search.html")
    })
    
    route.Post(nil, func() {
        query := route.UrlEncoded["q"]
        postal := route.UrlEncoded["postal"]
        appLastPosition, _ := strconv.Atoi(route.UrlEncoded["lp"])
        startRange, _ := strconv.Atoi(route.UrlEncoded["sr"])
        jobs,ftOffset, lastPosition, err := store.GetJobBySearch(query, postal, startRange, appLastPosition)
        if err != nil{
            log.Println(err)
        }
        payload := SearchPayload{
            RequireData: RequireData{
                store.ConnectedUser{},
                SearchQuery{
                    Query: query, Postal: postal, FrancetravailPosition: ftOffset, Lastposition: lastPosition,
                },
            },
            Jobs: jobs,
        }
        templ, err := template.New("MoreJobs").Parse(`
        {{range .Jobs}}
            <job-card id="{{.Id}}" salary="{{.SalaryString}}" entreprise="{{.EntrepriseName}}"  isThirdParty="{{.IsThirdParty}}" contract="{{.Contract}}" title="{{.Title}}" fulladdr="{{.FullAdresse}}" date="{{.Date}}" ></job-card>
        {{end}}
        <button type="button" hx-post="/search" hx-vals='{"postal": "{{.Search.Postal}}", "q": "{{.Search.Query}}", "city": "{{.Search.City}}", "sr": "{{.Search.FrancetravailPosition}}", "lp": "{{.Search.Lastposition}}"}' hx-swap="beforeend" hx-swap-oob="true" hx-target="#search_jobs" id="search_more">More</button>
        `)
        if err != nil{
            log.Println(err)
            return
        }
        if err := templ.Execute(route.Response, payload); err != nil{
            log.Println(err)
        }
        //https://api-adresse.data.gouv.fr/search/?q=8+bd+du+port
    })

    route.Put(&filter, func() {
        if filter.MinSalary > filter.MaxSalary{
            route.Notification("error", "salaire minimum est majeur que salaire maximum")
            return
        }

        filterString := "SELECT j.title, LEFT(j.postal, 2) || ' - ' || j.city, e.name, j.contract, j.fulltime, TO_CHAR(AGE(NOW(), j.created), 'Y-MM-DD-HH24-MI-SS'), j.id, j.salary FROM Job AS j LEFT JOIN Entreprise AS e ON j.entreprise_id=e.id WHERE"
        var minSalary float64
        reflectValue := reflect.ValueOf(filter)
        reflectType := reflect.TypeOf(filter)
        counter := 0
        for i := 0; i < reflect.TypeOf(filter).NumField(); i++{
            if reflectValue.Field(i).IsZero() && reflectValue.Field(i).Kind().String() != "bool" {continue}
            //if typeValue.Field(i).Name
            prefix := ""
            if counter > 0{prefix = "AND"}
            counter++
            switch reflectValue.Field(i).Kind().String() {
                case "string":
                    value := reflectValue.Field(i).String()
                    switch reflectType.Field(i).Name{
                        case "Postal":
                            filterString += fmt.Sprintf(" %v LEFT(j.postal, 2)='%v'", prefix, value[0:2])
                        case "Contract":
                            filterString += fmt.Sprintf(" %v j.contract='%v'", prefix, value)
                        case "Query":
                            filterString += fmt.Sprintf(" %v ts @@ websearch_to_tsquery('french', '%v')", prefix, value)
                        case "ExperienceMin":
                            filterString += fmt.Sprintf(" %v j.require_exp <= %v", prefix, value)
                        default:
                            log.Println("default")
                    }
                case "bool":
                    value := reflectValue.Field(i).Bool()
                    switch reflectType.Field(i).Name{
                        case "Fulltime":
                            filterString += fmt.Sprintf(" %v j.fulltime=%t", prefix, value)
                        default:
                            log.Println("default")
                    }
                default:
                    value := reflectValue.Field(i).Float()
                    if filter.MinSalary > 0 && filter.MaxSalary > 0{
                        switch reflectType.Field(i).Name{
                            case "MinSalary":
                                minSalary = value
                            case "MaxSalary":
                                filterString += fmt.Sprintf(" %v j.salary[1] BETWEEN %v AND %v", prefix, minSalary, value)
                            default:
                                log.Println("default")
                        }
                    }
            }

        }
        filterString += fmt.Sprintf(" AND Distance(j.lat, j.long, %v, %v) < %v ORDER BY j.created %v LIMIT 2", filter.Lat, filter.Long, filter.Distance, filter.Order)
        
        jobs := store.SearchJobWithFilter(filterString, filter)
        temp, _ := template.New("Filter Job").Parse(`
            {{range .}}
            <job-card id="{{.Id}}" salary="{{.SalaryString}}" entreprise="{{.EntrepriseName}}"  isThirdParty="{{.IsThirdParty}}" contract="{{.Contract}}" title="{{.Title}}" fulladdr="{{.FullAdresse}}" date="{{.Date}}" fulltime="{{.FulltimeString}}" ></job-card>
            {{end}}
        `)
        temp.Execute(route.Response, jobs)
        
    })
}
