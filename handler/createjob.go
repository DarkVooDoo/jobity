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
    Templates string
    Contract string
}

var CreateJobHandler = func(res http.ResponseWriter, req *http.Request){
    route, err := NewRoute(res, req)
    job := store.Job{EntrepriseId: route.User.Id}
    route.Get(func() {
        if err != nil{
            route.Response.Header().Add("Location", "/")
            route.Response.WriteHeader(http.StatusTemporaryRedirect)
        }
        templ, err := job.GetTemplates()
        if err != nil{
            log.Println(err)
            //route.Notification("warning", "impossible dans les modéles")
        }
        stringTemplate := store.TemplatesIntoString(templ)
        contract := store.GetContracts()
        data := NewJobPage{
            Contract: contract,
            Templates: stringTemplate,
            RequireData: RequireData{Search: SearchQuery{Query: ""}},
        }
        route.Render(data, "route/protemplate.html", "route/create_job.html")
    })

    route.Put(nil, func() {
        templateId := route.UrlEncoded["template"]
        job.GetJobByTemplateId(templateId)
        job.ContractArray = store.GetContracts()
        templ, _ := template.New("form").Parse(`
        <form class="newjob" hx-post="/job/creer" hx-ext="json-enc" hx-swap="none" id="offert" hx-vals='js:{...getValues()}' onclick="onSubmitForm()">
            <div class="newjob_field">
                <input type="text" autocomplete="off" value="{{.Title}}"  required id="title"  name="title" class="newjob_field_input newjob_field_input_withLabel" onchange="onSaveSnapshot(this)" />
                <label for="title" class="newjob_field_label">Titre</label>
            </div>
            <div class="newjob_flex" style="position: relative;">
                <div class="newjob_field" style="flex: 1;">
                    <input type="text" autocomplete="off" value="{{.City}}" placeholder="Paris" id="city" name="city" class="newjob_field_input newjob_field_input_withLabel" oninput="onCityInput(this)" onchange="onSaveSnapshot(this)"/>
                    <label for="adresse" class="newjob_field_label">Departement</label>
                </div>
                <div class="newjob_field" style="flex: .3;">
                    <input type="number" autocomplete="off" value="{{.Postal}}"  placeholder="75001" id="postal" name="postal" class="newjob_field_input newjob_field_input_withLabel" onchange="onSaveSnapshot(this)" />
                    <label for="adresse" class="newjob_field_label">Postal</label>
                </div>
                <div id="newjob_addr"></div>
            </div>
            <h2 class="subheader">Salaire</h2>
            <div class="newjob_flex">
                <div class="newjob_field">
                    <input type="number" placeholder="1500" value="{{index .Salary 0}}"  step=".01" id="minSalary" class="newjob_field_input salary newjob_field_input_withLabel" />
                    <label for="minSalary" class="newjob_field_label">Min</label>
                </div>
                <div class="newjob_field">
                    <input type="number" placeholder="2000" value="{{index .Salary 1}}"  step=".01" id="maxSalary" class="newjob_field_input salary newjob_field_input_withLabel" />
                    <label for="maxSalary" class="newjob_field_label">Max</label>
                </div>
            </div>
            <h2 class="subheader">Contrat</h2>
            <div class="newjob_flex">
                <dropdown-ele id="contract" value="{{.Contract}}"  array="{{.ContractArray}}"></dropdown-ele>
                <div class="newjob_field">
                    <input type="number" placeholder="35" value="{{.WeeklyWorkTime}}" id="weeklyWorkTime" name="weeklyWorkTime" class="newjob_field_input newjob_field_input_withLabel" />
                    <label for="weeklyWorkTime" class="newjob_field_label">Heures</label>
                </div>
                <div class="newjob_field">
                    <input type="date" placeholder="2000" id="startDate" name="startDate" class="newjob_field_input newjob_field_input_withLabel" />
                    <label for="startDate" class="newjob_field_label">Debut</label>
                </div>
            </div>
            <div class="newjob_section" id="advantage">
                <div class="newjob_headerWithButton">
                    <h2 class="subheader">Avantages</h2>
                    <button type="button" class="newjob_addBtn" onclick="onNewAdvantage(this)">
                        <svg
                            width="63.999996"
                            height="63.999996"
                            style="width: 60%;height: 60%;"
                            viewbox="0 0 16.933332 16.933332">
                            <g
                                transform="translate(-49.871622,-101.59117)">
                            <rect
                                style="fill:#000000;stroke-width:0.0701647"
                                id="rect113"
                                width="2.6458337"
                                height="15.875"
                                x="108.73493"
                                y="-66.275787"
                                ry="1.3229169"
                                transform="rotate(90)" />
                            <rect
                                style="fill:#000000;stroke-width:0.0701647"
                                id="rect273"
                                width="2.6458337"
                                height="15.875"
                                x="57.015373"
                                y="102.12034"
                                ry="1.3229169" />
                            </g>
                        </svg>
                    </button>
                </div>
                {{range .Advantage}}
                <div class="newjob_advantage">
                    <input type="text" autocomplete="off" value="{{.}}" name="advantage" placeholder="Titre de transport" class="newjob_field_input advantage" />
                    <button type="button" class="newjob_profile_deleteBtn" onclick="onDeleteAdvantage(this)">
                        <svg
                            width="63.999996"
                            height="63.999996"
                            style="width: 60%;height: 60%;"
                            viewBox="0 0 16.933332 16.933332">
                            <g
                                transform="translate(-51.903267,-103.62282)">
                                <rect
                                    style="fill:#000000;stroke-width:0.0928013"
                                    width="3.4994376"
                                    height="20.996624"
                                    x="120.19751"
                                    y="26.072937"
                                    ry="1.7497188"
                                    transform="rotate(45)" />
                                <rect
                                    style="fill:#000000;stroke-width:0.0928013"
                                    id="rect273"
                                    width="3.4994376"
                                    height="20.996624"
                                    x="-38.320965"
                                    y="111.44891"
                                    ry="1.7497188"
                                    transform="rotate(-45)" />
                                </g>
                        </svg>
                    </button>
                </div>
                {{end}}
            </div>
            <h2 class="subheader">Description</h2>
            <pre class="newjob_description" contenteditable="true">{{.Description}}</pre>
            <div class="newjob_section" id="skill"> 
                <div class="newjob_headerWithButton">
                    <h2 class="subheader">Compétances</h2>
                    <button type="button" class="newjob_addBtn" onclick="onNewProfil(this)">
                        <svg
                            width="63.999996"
                            height="63.999996"
                            style="width: 60%;height: 60%;"
                            viewbox="0 0 16.933332 16.933332">
                            <g
                                transform="translate(-49.871622,-101.59117)">
                            <rect
                                style="fill:#000000;stroke-width:0.0701647"
                                id="rect113"
                                width="2.6458337"
                                height="15.875"
                                x="108.73493"
                                y="-66.275787"
                                ry="1.3229169"
                                transform="rotate(90)" />
                            <rect
                                style="fill:#000000;stroke-width:0.0701647"
                                id="rect273"
                                width="2.6458337"
                                height="15.875"
                                x="57.015373"
                                y="102.12034"
                                ry="1.3229169" />
                            </g>
                        </svg>
                    </button>
                </div>
                <div class="newjob_profile">
                    <h3 class="newjob_profile_header">Titre</h3>
                    <h3 class="newjob_profile_header">Nécessaire</h3>
                    <h3></h3>
                </div>

                {{range .SkillNeeded}}
                <div class="newjob_profile skill">
                    <input type="text" autocomplete="off" value="{{.Label}}"  placeholder="Experience 2 ans"  class="newjob_field_input" />
                    <toggle-btn on="{{if .Required}}true{{else}}false{{end}}"></toggle-btn>
                    <button type="button" onclick="onDeleteProfil(this)"  class="newjob_profile_deleteBtn">
                        <svg
                            width="63.999996"
                            height="63.999996"
                            style="width: 60%;height: 60%;"
                            viewBox="0 0 16.933332 16.933332">
                            <g transform="translate(-51.903267,-103.62282)">
                                <rect
                                    style="fill:#000000;stroke-width:0.0928013"
                                    width="3.4994376"
                                    height="20.996624"
                                    x="120.19751"
                                    y="26.072937"
                                    ry="1.7497188"
                                    transform="rotate(45)" />
                                <rect
                                    style="fill:#000000;stroke-width:0.0928013"
                                    id="rect273"
                                    width="3.4994376"
                                    height="20.996624"
                                    x="-38.320965"
                                    y="111.44891"
                                    ry="1.7497188"
                                    transform="rotate(-45)" />
                            </g>
                        </svg>
                    </button>
                </div>
                {{end}}
            </div>
            <button type="submit" class="newjob_submitBtn" >Creer</button>
        </form>
        `)
        if err := templ.Execute(route.Response, job); err != nil{
            log.Println(err)
        }
    })

    route.Post(&job, func() {
        if job.Salary[0] > job.Salary[1]{
            route.Notification("error", "erreur dans le salaire")
            return
        }
        job.EntrepriseId = route.User.Id
        job, err := job.CreateJob()
        if  err != nil{
            log.Println(err)
            route.Notification("error", "Error dans la creationn de l'annonce")
            return
        }
        route.Response.Header().Add("HX-Redirect", fmt.Sprintf("/job/%v", job))
        route.Response.WriteHeader(http.StatusTemporaryRedirect)
    })
}

