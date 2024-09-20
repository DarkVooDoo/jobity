package handler

import (
	"html/template"
	"job/store"
	"net/http"
)

var BookmarkHandler = func(res http.ResponseWriter, req *http.Request){
    route, err := NewRoute(res, req)
    if err != nil{
        return
    }
    bookmark := store.Bookmark{UserId: route.User.Id}
    route.Put(nil, func() {
        bookmark.JobId = route.UrlEncoded["jobId"]
        if route.UrlEncoded["isThirdParty"] == "true"{
            bookmark.IsThirdParty = true
        }
        if err := bookmark.Create(); err != nil{
            route.Notification("warning", "il y eu un probleme")
            return
        }
        temp, _ := template.New("Bookmark").Parse(`
            <button class="bookmark" onclick="onBookmarkClick(this)" hx-delete="/bookmark" hx-vals='{"bookmarkId": "{{.Id}}", "jobId": "{{.JobId}}"}'>
                <svg
                    class="bookmark-icon"
                    width="12.837221mm"
                    height="16.933336mm"
                    viewBox="0 0 12.837221 16.933336"
                    version="1.1">
                    <g transform="translate(-74.365777,-89.4474)">
                        <path class="bookmark-path fill" style="stroke:#000000;stroke-width:0.998655;stroke-dasharray:none"
                        d="m 76.15008,90.744212 h 9.268616 a 1,1 45 0 1 1,1 l 0,12.805038 a 0.54317256,0.54317256 151.4904 0 1 -0.838853,0.45564 l -4.795455,-3.11192 -4.795455,3.11192 A 0.54317256,0.54317256 28.509595 0 1 75.15008,104.54925 V 91.744212 a 1,1 135 0 1 1,-1 z"/>
                    </g>
                </svg>
            </button>
        `)
        temp.Execute(route.Response, bookmark)
        route.Notification("success", "Job ajouté a vos favorites")
    })

    route.Delete(nil, func() {
        bookmark.Id = route.UrlEncoded["bookmarkId"]
        bookmark.JobId = route.UrlEncoded["jobId"]
        if err := bookmark.Delete(); err != nil{
            route.Notification("error", "il y a eu un probleme")
            return 
        }
        temp, _ := template.New("Bookmark").Parse(`
            <button class="bookmark" onclick="onBookmarkClick(this)" hx-put="/bookmark" hx-vals='{"jobId": "{{.}}"}' hx-swap="outerHTML">
                <svg
                    class="bookmark-icon"
                    width="12.837221mm"
                    height="16.933336mm"
                    viewBox="0 0 12.837221 16.933336"
                    version="1.1">
                    <g transform="translate(-74.365777,-89.4474)">
                        <path class="bookmark-path" style="stroke:#000000;stroke-width:0.998655;stroke-dasharray:none"
                        d="m 76.15008,90.744212 h 9.268616 a 1,1 45 0 1 1,1 l 0,12.805038 a 0.54317256,0.54317256 151.4904 0 1 -0.838853,0.45564 l -4.795455,-3.11192 -4.795455,3.11192 A 0.54317256,0.54317256 28.509595 0 1 75.15008,104.54925 V 91.744212 a 1,1 135 0 1 1,-1 z"/>
                    </g>
                </svg>
            </button>
        `)
        temp.Execute(route.Response, bookmark.JobId)
        route.Notification("success", "Job supprimer des vos favorites")
    })
}
