package handler

import (
	"html/template"
	"job/store"
	"log"
	"net/http"
)

type ApplicationPage struct{
    RequireData
    Apps []store.JobApplication
    Fav []store.Bookmark
}

var ApplicationHandler = func(res http.ResponseWriter, req *http.Request){
    route, _ := NewRoute(res, req) 

    route.Get(func() {
        jobApplication := store.JobApplication{UserId: route.User.Id}
        applications := jobApplication.GetUserApplications()
        bookmarkStruct :=  store.Bookmark{UserId: route.User.Id}
        bookmark := bookmarkStruct.Get()
        page := ApplicationPage{
            RequireData{
                Search: SearchQuery{Query: ""},
                User: route.User,
            },
            applications,
            bookmark,
        }
        route.Render(page, "route/template.html", "route/my-applications.html")
    })
    route.Put(nil, func() {
        jobApplication := store.JobApplication{UserId: route.User.Id, JobId: route.Multipart.Body["jobId"]}
        if err := jobApplication.CreateJobApplication(); err != nil{
            route.Notification("error", "postulation error")
            return
        }
        temp, _ := template.New("checkmark").Parse(`
        <button type="button" id="submitBtn" hx-swap-oob="true" disabled  popovertarget="application" class="job_submitBtn">Postuler</button>
        <div class="successResponse">
            <div class="successResponse_bg"></div>
            <div class="successRsponse_circle">
            <svg
                width="64.000008"
                height="57.549763"
                viewBox="0 0 16.933335 15.226707">
                <g
                    transform="translate(-46.262468,-78.805815)">
                <path
                    style="fill:transparent;stroke:#ebebeb;stroke-width:3.361;stroke-opacity:1"
                    d="m 47.390491,85.138996 c 1.894167,1.715301 3.479146,3.648796 4.407963,5.488994 3.713518,-5.874505 7.072233,-8.107992 10.440608,-10.440608"
                    id="checkmark"/>
                </g>
            </svg>
            </div>
        </div>`)
        if err := temp.Execute(route.Response, nil); err != nil{
            log.Println(err)
        }
    })

    route.Delete(nil, func() {
        jobApplication := store.JobApplication{Id: route.UrlEncoded["id"]}        
        if err := jobApplication.Delete(); err != nil{
            route.Notification("error", "suppresion impossible")
            return
        }
        route.Response.Header().Add("HX-Trigger-After-Swap", "onDeleteCard")
        route.Notification("success", "suppresion")
    })
}
