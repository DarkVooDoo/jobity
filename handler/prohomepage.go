package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"job/store"
	"log"
	"net/http"
)

type HomepageData struct{
    User store.ConnectedUser
    Profile store.ProUser
    Job []store.Job
}

type InterviewData struct{
    Interview []store.JobApplication
    InterviewType []string
}

type Adresse struct{
    Data []struct{
        Adresse struct{
            Name string `json:"name"`
            City string `json:"city"`
            Postal string `json:"postcode"`
        } `json:"properties"`
    } `json:"features"`
}

type Payload struct{
    Query string `jsonn:"city"`
}

var ProHomepageHandler = func(res http.ResponseWriter, req *http.Request){
    route, err := NewRoute(res, req)
    if req.URL.Path != "/"{
        route.Render("route/protemplate.htmk", "route/notfound.html")
        return
    }

    route.Get(func() {
        if err != nil{
            route.Response.Header().Add("Location", "/connexion")
            route.Response.WriteHeader(http.StatusTemporaryRedirect)
            return
        }
        q := route.Request.URL.Query().Get("city")
        if q != ""{
            //Search adresse from BAN api
            var addr Adresse
            req, _ := http.NewRequest("GET", fmt.Sprintf("https://api-adresse.data.gouv.fr/search/?q=%v&limit=5", q), nil)
            res, err := http.DefaultClient.Do(req)
            if err != nil{
                log.Println(err)
            }
            dec := json.NewDecoder(res.Body)
            if err := dec.Decode(&addr); err != nil{
                log.Println(err)
            }
            log.Println(addr)
            return
        }
        var pageData HomepageData
        entreprise := store.ProUser{Id: route.User.Id}
        entreprise.GetProfile()
        jobs := store.GetEntrepriseJobs(route.User.Id)
        pageData.Profile = entreprise
        pageData.Job = jobs
        pageData.User = route.User
        route.Render(pageData, "route/protemplate.html", "route/prohomepage.html")
    })

    route.Post(nil, func() {
        jobs := store.GetEntrepriseJobs(route.User.Id)
        temp, err := template.New("Entreprise Jobs").Parse(`
        <div class="offert-container">
            <a class="create-offert" href="/job/creer">Creer une annonce</a>
            {{range .}}
            <a class="offert-card" href="/job/{{.Id}}">
                <div class="data">
                    <h3>{{.Title}}</h3>
                    <div class="metadata">
                        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                            <circle cx="12" cy="10" r="3"></circle>
                        </svg>
                        <p>{{.City}}, {{.Postal}}</p>
                    </div>
                    <div class="metadata">
                        <img src="/static/contract.svg" class="icon"></img>
                        <p>{{.Contract}}</p>
                    </div>
                </div>
                <p class="reqAmount">{{.ApplicationCount}}</p>
            </a>
            {{end}}
        </div>
        `)
        if err != nil{
            log.Printf("Error dans entreprise jobs: %v", err)
        }
        temp.Execute(route.Response, jobs)
    })

    route.Put(nil, func() {
        application := store.JobApplication{}
        interviews := application.Interviews(route.User.Id)
        interviewData := InterviewData{Interview: interviews, InterviewType: store.GetInterviewType()}
        temp, err := template.New("Get interviews").Parse(`
            <div>
                {{$InterviewType := .InterviewType}}
                {{range .Interview}}
                    <div class="interview-card">
                        <div class="interview-card-info">
                            <img src="/static/netflix.webp" class="interview-photo" />
                            <div style="position: relative; width: 100%">
                                <img src="/static/{{.Type}}.svg" class="interview-type" id="icon-{{.Id}}"/>
                                <div class="interview-row">
                                    <svg class="interview-card-info-icon" viewBox="0 0 16 16" fill="none" >
                                        <path d="M8 7C9.65685 7 11 5.65685 11 4C11 2.34315 9.65685 1 8 1C6.34315 1 5 2.34315 5 4C5 5.65685 6.34315 7 8 7Z" fill="#000000"/>
                                        <path d="M14 12C14 10.3431 12.6569 9 11 9H5C3.34315 9 2 10.3431 2 12V15H14V12Z" fill="#000000"/>
                                    </svg>
                                    <h1 class="interview-user-name">{{.UserName}}</h1>
                                </div>
                                <div class="interview-row">
                                    <svg class="interview-card-info-icon" viewBox="0 0 24 24" fill="none" >
                                        <path d="M3 9H21M7 3V5M17 3V5M6 12H8M11 12H13M16 12H18M6 15H8M11 15H13M16 15H18M6 18H8M11 18H13M16 18H18M6.2 21H17.8C18.9201 21 19.4802 21 19.908 20.782C20.2843 20.5903 20.5903 20.2843 20.782 19.908C21 19.4802 21 18.9201 21 17.8V8.2C21 7.07989 21 6.51984 20.782 6.09202C20.5903 5.71569 20.2843 5.40973 19.908 5.21799C19.4802 5 18.9201 5 17.8 5H6.2C5.0799 5 4.51984 5 4.09202 5.21799C3.71569 5.40973 3.40973 5.71569 3.21799 6.09202C3 6.51984 3 7.07989 3 8.2V17.8C3 18.9201 3 19.4802 3.21799 19.908C3.40973 20.2843 3.71569 20.5903 4.09202 20.782C4.51984 21 5.07989 21 6.2 21Z" stroke="#000000" stroke-width="2" stroke-linecap="round"/>
                                    </svg>
                                    <p id="date-{{.Id}}">{{.InterviewDate}}</p>
                                </div>
                                {{if .Addr}}
                                    <div class="interview-row">
                                        <svg class="interview-card-info-icon" viewBox="-4 0 32 32" version="1.1">
                                            <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
                                                <g id="Icon-Set" sketch:type="MSLayerGroup" transform="translate(-104.000000, -411.000000)" fill="#000000">
                                                    <path d="M116,426 C114.343,426 113,424.657 113,423 C113,421.343 114.343,420 116,420 C117.657,420 119,421.343 119,423 C119,424.657 117.657,426 116,426 L116,426 Z M116,418 C113.239,418 111,420.238 111,423 C111,425.762 113.239,428 116,428 C118.761,428 121,425.762 121,423 C121,420.238 118.761,418 116,418 L116,418 Z M116,440 C114.337,440.009 106,427.181 106,423 C106,417.478 110.477,413 116,413 C121.523,413 126,417.478 126,423 C126,427.125 117.637,440.009 116,440 L116,440 Z M116,411 C109.373,411 104,416.373 104,423 C104,428.018 114.005,443.011 116,443 C117.964,443.011 128,427.95 128,423 C128,416.373 122.627,411 116,411 L116,411 Z" id="location" sketch:type="MSShapeGroup">

                                                    </path>
                                                </g>
                                            </g>
                                        </svg>
                                        <p id="addr-{{.Id}}">{{.Addr}}</p>
                                    </div>
                                {{end}}
                            </div>
                        </div>
                        <div class="interview-card-footer">
                            <div class="card-footer-user">
                                <img src="/static/netflix.webp" class="card-footer-photo" />
                                <div>
                                    <h1>Alice Doe</h1>
                                    <p>Responsable</p>
                                </div>
                            </div>
                            <button type="button" class="interview-footer-edit" onclick="onEditInterview(this)">
                                <svg fill="#000000" style="width: 80%;aspect-ratio: 1/1;" viewBox="-4 0 32 32" version="1.1" xmlns="http://www.w3.org/2000/svg">
                                <path d="M17.438 22.469v-4.031l2.5-2.5v7.344c0 1.469-1.219 2.688-2.656 2.688h-14.625c-1.469 0-2.656-1.219-2.656-2.688v-14.594c0-1.469 1.188-2.688 2.656-2.688h14.844v0.031l-2.5 2.469h-11.5c-0.531 0-1 0.469-1 1.031v12.938c0 0.563 0.469 1 1 1h12.938c0.531 0 1-0.438 1-1zM19.813 7.219l2.656 2.656 1.219-1.219-2.656-2.656zM10.469 16.594l2.625 2.656 8.469-8.469-2.625-2.656zM8.594 21.094l3.625-0.969-2.656-2.656z"></path>
                                </svg>
                            </button>
                        </div>

                        <dialog class="interview-modal">
                            <h1 style="margin-bottom: .5rem;text-align: center;">Modifier interview</h1>
                            <form hx-post="/job/{{.JobId}}/cv/{{.UserId}}" hx-vals='{"id": "{{.Id}}", "status": "Interview"}' 
                                hx-swap="none">
                                <div style="display: flex;justify-content: space-between;gap: .5rem;">
                                    <div>
                                        <b class="input-label">Date</b>
                                        <input type="datetime-local" name="interview_date" class="interview-input" value="{{.InterviewDate}}" />
                                    </div>
                                    <div>
                                        <b class="input-label">Type</b>
                                        <select name="type" class="interview-input" onchange="onTypeChange(this)">
                                            {{range $InterviewType}}
                                                <option value="{{.}}">{{.}}</option>
                                            {{end}}
                                        </select>
                                    </div>
                                </div>
                                <div class="interviewModal-location hidden" style="position: relative;">
                                    <b class="input-label">Location</b>
                                    <input type="text" name="location" autocomplte="off" class="interview-input" oninput="onFetchAddr(this)"/>
                                    <div class="interviewModal-location-suggest hidden"></div>
                                </div>
                                <div class="interview-buttons">
                                    <button type="button" class="interview-btn" onclick="onCloseInterview(this)">Cancel</button>
                                    <button type="submit" class="interview-btn">Modifier</button>
                                </div>
                            </form>
                        </dialog>
                    </div>
                {{end}}
            </div>
        `)   
        if err != nil{
            log.Printf("error dans la requete put %v", err)
        }
        temp.Execute(route.Response, interviewData) 
    })
}


