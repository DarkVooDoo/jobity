package handler

import (
	"fmt"
	"html/template"
	"job/store"
	"log"
	"net/http"
)

type NewJobPage struct{
    RequireData
    Templates []store.EntrepriseTemplates
    Contract []store.Contract
    Category []store.Category
}

type TemplateLoad struct{
    Contract []store.Contract
    Category []store.Category
    Job store.Job
}

type CategoryUpdate struct{
    Category store.Category
    SelectSubcategory store.Category
    Subcategory []store.Category
}

var CreateJobHandler = func(res http.ResponseWriter, req *http.Request){
    route, err := NewRoute(res, req)
    job := store.Job{EntrepriseId: route.User.Id}
    route.Get(func() {
        if err != nil{
            route.Response.Header().Add("Location", "/")
            route.Response.WriteHeader(http.StatusTemporaryRedirect)
        }
        conn := store.GetDBPoolConn()
        defer conn.Close()
        templ, err := job.GetTemplates(conn)
        if err != nil{
            log.Println(err)
            //route.Notification("warning", "impossible dans les modéles")
        }
        contract := store.GetContracts(conn)
        category := store.GetCategorys(conn)
        data := NewJobPage{
            Contract: contract,
            Templates: templ,
            RequireData: RequireData{Search: SearchQuery{Query: ""}},
            Category: category,
        }
        route.Render(data, "route/protemplate.html", "route/create_job.html")
    })

    route.Patch(nil, func() {
        title := route.UrlEncoded["title"]
        category := route.UrlEncoded["category"]

        conn, _ := store.GetDBConn()
        defer conn.Close()
        if title != ""{
            category, subcategory := store.GetCategoryByTitle(title)
            if category.Id == ""{
                return
            }
            subcategoryList := store.GetSubcategory(conn, category.Id)
            updateCategory := CategoryUpdate{category, subcategory,  subcategoryList}
            temp, err := template.New("category").Parse(`
                {{$selectSubcategory := .SelectSubcategory}}
                <option value="{{.Category.Id}}" id="category-{{.Category.Id}}" hx-swap-oob="true" selected>{{.Category.Name}}</option>
                <div>
                    <h2 class="subheader">Subcategorie</h2>
                    <select name="subcategory" class="select-element input">
                        {{range .Subcategory}}
                            <option value="{{.Id}}" {{if eq $selectSubcategory.Id .Id}}selected{{end}}>{{.Name}}</option>
                        {{end}}
                    </select>
                </div>
            `)
            if err != nil{
                log.Printf("error parsing template: %v", err)
            }
            if err := temp.Execute(route.Response, updateCategory); err != nil{
                log.Printf("execute temp error: %v", err)
            }
        }else{
            subcategory := store.GetSubcategory(conn, category)
            temp, err := template.New("Subcategory").Parse(`
                <div>
                    <h2 class="subheader">Subcategorie</h2>
                    <select name="subcategory" class="select-element inpu">
                        {{range .}}
                            <option value="{{.Id}}">{{.Name}}</option>
                        {{end}}
                    </select>
                </div>
            `)
            if err != nil{
                log.Printf("error parsing template: %v", err)
            }
            if err := temp.Execute(route.Response, subcategory); err != nil{
                log.Printf("error executing template: %v", err)
            }
        }
    })

    route.Put(nil, func() {
        templateId := route.UrlEncoded["template"]
        conn, _ := store.GetDBConn() 
        defer conn.Close()
        job.GetJobByTemplateId(templateId)
        data := TemplateLoad{store.GetContracts(conn), store.GetCategorys(conn), job}
        templ, _ := template.New("form").Parse(`
            <div class="newjob-section">
                <h2 class="newjob-section-name">General</h2>
                <div class="field">
                    <label for="title" class="label">Titre</label>
                    <input type="text" value="{{.Job.Title}}" autocomplete="off" required id="title" onchange="onSaveSnapshot(this)"  name="title" class="input" hx-patch="/job/creer"
                    hx-params="title" hx-trigger="input delay:500ms" hx-ext="ignore:json-enc" hx-target="#subcategory" hx-swap="innerHTML" />
                </div>
                <div class="newjob-flex-form">
                    <div class="field" style="flex: 1;">
                        <label for="city" class="label">Departement</label>
                        <input type="text" value="{{.Job.City}}" autocomplete="off" placeholder="Paris" id="city" name="city" class="input " oninput="onCityInput(this)" onchange="onSaveSnapshot(this)"/>
                    </div>
                    <div class="field" style="flex: .3;">
                        <label for="adresse" class="label">Postal</label>
                        <input type="number" value="{{.Job.Postal}}"  autocomplete="off" placeholder="75001" id="postal" name="postal" class="input " onchange="onSaveSnapshot(this)" />
                    </div>
                    <div id="addr" >
                    </div>
                </div>
                <div>
                    <h2 class="subheader">Description</h2>
                    <div class="description" contenteditable="true" onblur="onSnapshotDescription(this)">{{.Job.Description}}"</div>
                </div>
            </div>
            <div class="newjob-section">
                <h3 class="newjob-section-name">Details</h3>
                <div style="display: flex;flex-wrap: wrap;">
                    <div class="newjob-select-category">
                        <h2 class="subheader">Categorie</h2>
                        <select id="category" name="category" class="select-element input" hx-patch="/job/creer" hx-trigger="change" 
                            hx-target="#subcategory" hx-swap="innerHTML" hx-params="category" hx-ext="ignore:json-enc">
                            {{range .Category}}
                                <option value="{{.Id}}" id="category-{{.Id}}" >{{.Name}}</option>
                            {{end}}
                        </select>
                    </div>
                    <div id="subcategory" class="newjob-select-category">
                        
                    </div>
                </div>
                <div class="newjob-flex-form">
                    <div class="field">
                        <label for="minSalary" class="label">Salaire Min</label>
                        <input type="number" value="{{index .Job.Salary 0}}" placeholder="1500" step=".01"  id="minSalary" class="input salary " onchange="onSaveSnapshot(this)" />
                    </div>
                    <div class="field">
                        <label for="maxSalary" class="label">Salaire Max</label>
                        <input type="number" value="{{index .Job.Salary 1}}" placeholder="2000" step=".01" id="maxSalary" class="input salary " onchange="onSaveSnapshot(this)" />
                    </div>
                </div>
                <div class="newjob-flex-form" style="flex-wrap: wrap;justify-content: space-between;">
                    <div>
                        <h2 class="subheader">Contrat</h2>
                        <select name="contract" class="input" onchange="onContractChange(this)">
                            {{range .Contract}}
                                <option value="{{.Id}}">{{.Name}}</option>
                            {{end}}
                        </select>
                    </div>
                    <div class="field" style="width: 75px;">
                        <label for="weeklyWorkTime" class="label">Heures</label>
                        <input type="number" placeholder="35" value="{{.Job.WeeklyWorkTime}}"  id="weeklyWorkTime" name="weeklyWorkTime" class="input" onchange="onSaveSnapshot(this)" />
                    </div>
                    <div class="field" style="width: 75px;">
                        <label for="exp" class="label">Experience</label>
                        <input type="number" placeholder="2" id="exp" name="exp" value="{{.Job.Experience}}" class="input " onchange="onSaveSnapshot(this)" />
                    </div>
                    <div>
                        <h3 class="subheader">Duré</h3>
                        <div class="newjob-time-date">
                            <div class="newjob-time-date-field">
                                <p>Du</p>
                                <input type="date" name="startDate" class="input" value="{{.Job.StartDate}}"/>
                            </div>
                            <div class="newjob-time-date-field hidden" id="endDate">
                                <p>AU</p>
                                <input type="date" name="endDate"class="input" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="newjob-section" id="advantage">
                <h2 class="newjob-section-name">Avantages</h2>
                {{range .Job.Advantage}}
                    <div class="advantage">
                        <input type="text" value="{{.}}" autocomplete="off" placeholder="Titre de transport" name="advantage" class="input advantage" onchange="onSaveAdvantageSnapshot()" onkeyup="onNewAdvantageByKey(event)" />
                        <button type="button" class="deleteBtn" onclick="onDeleteAdvantage(this)">
                            <svg fill="#000000" width="800px" height="800px" viewBox="0 0 41.336 41.336" class="deleteIcon">
                                <g>
                                    <path d="M36.335,5.668h-8.167V1.5c0-0.828-0.672-1.5-1.5-1.5h-12c-0.828,0-1.5,0.672-1.5,1.5v4.168H5.001c-1.104,0-2,0.896-2,2
                                            s0.896,2,2,2h2.001v29.168c0,1.381,1.119,2.5,2.5,2.5h22.332c1.381,0,2.5-1.119,2.5-2.5V9.668h2.001c1.104,0,2-0.896,2-2
                                            S37.438,5.668,36.335,5.668z M14.168,35.67c0,0.828-0.672,1.5-1.5,1.5s-1.5-0.672-1.5-1.5v-21c0-0.828,0.672-1.5,1.5-1.5
                                            s1.5,0.672,1.5,1.5V35.67z M22.168,35.67c0,0.828-0.672,1.5-1.5,1.5s-1.5-0.672-1.5-1.5v-21c0-0.828,0.672-1.5,1.5-1.5
                                            s1.5,0.672,1.5,1.5V35.67z M25.168,5.668h-9V3h9V5.668z M30.168,35.67c0,0.828-0.672,1.5-1.5,1.5s-1.5-0.672-1.5-1.5v-21
                                            c0-0.828,0.672-1.5,1.5-1.5s1.5,0.672,1.5,1.5V35.67z"/>
                                </g>
                            </svg>
                        </button>
                    </div>
                {{end}}
                <button type="button" class="addBtn" onclick="onNewAdvantage(this)"></button>
            </div>
            <div class="newjob-section" id="skill"> 
                <h2 class="newjob-section-name">Compétances</h2>
                <div class="profile">
                    <h3 class="header">Titre</h3>
                    <h3 class="header">Nécessaire</h3>
                    <h3></h3>
                </div>
                    {{range .Job.SkillNeeded}}
                        <div class="profile skill">
                            <input type="text" autocomplete="off"  placeholder="Experience 2 ans"  class="input" value="{{.Label}}" />
                            <toggle-btn on="{{.Required}}"></toggle-btn>
                            <button type="button" onclick="onDeleteProfil(this)"  class="deleteBtn">
                                <svg fill="#000000" width="800px" height="800px" viewBox="0 0 41.336 41.336" class="deleteIcon">
                                    <g>
                                        <path d="M36.335,5.668h-8.167V1.5c0-0.828-0.672-1.5-1.5-1.5h-12c-0.828,0-1.5,0.672-1.5,1.5v4.168H5.001c-1.104,0-2,0.896-2,2
                                                s0.896,2,2,2h2.001v29.168c0,1.381,1.119,2.5,2.5,2.5h22.332c1.381,0,2.5-1.119,2.5-2.5V9.668h2.001c1.104,0,2-0.896,2-2
                                                S37.438,5.668,36.335,5.668z M14.168,35.67c0,0.828-0.672,1.5-1.5,1.5s-1.5-0.672-1.5-1.5v-21c0-0.828,0.672-1.5,1.5-1.5
                                                s1.5,0.672,1.5,1.5V35.67z M22.168,35.67c0,0.828-0.672,1.5-1.5,1.5s-1.5-0.672-1.5-1.5v-21c0-0.828,0.672-1.5,1.5-1.5
                                                s1.5,0.672,1.5,1.5V35.67z M25.168,5.668h-9V3h9V5.668z M30.168,35.67c0,0.828-0.672,1.5-1.5,1.5s-1.5-0.672-1.5-1.5v-21
                                                c0-0.828,0.672-1.5,1.5-1.5s1.5,0.672,1.5,1.5V35.67z"/>
                                    </g>
                                </svg>
                            </button>
                        </div>
                    {{end}}
                <button type="button" class="addBtn" onclick="onNewProfil(this)"></button>
            </div>
            <button type="submit" class="newjob_submitBtn" >Creer</button>
        `)
        if err := templ.Execute(route.Response, data); err != nil{
            log.Println(err)
        }
    })

    route.Post(&job, func() {
        if job.Salary[0] > job.Salary[1]{
            route.Notification("error", "erreur dans le salaire")
            return
        }
        conn := store.GetDBPoolConn()
        defer conn.Close()
        job.EntrepriseId = route.User.Id
        job, err := job.CreateJob(conn)
        if  err != nil{
            log.Println(err)
            route.Notification("error", "Error dans la creationn de l'annonce")
            return
        }
        route.Response.Header().Add("HX-Redirect", fmt.Sprintf("/job/%v", job))
        route.Response.WriteHeader(http.StatusTemporaryRedirect)
    })
}

