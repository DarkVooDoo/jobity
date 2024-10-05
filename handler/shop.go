package handler

import (
	"html/template"
	"job/store"
	"net/http"
)

type ShopPayload struct{
    User store.ConnectedUser
    Shop []store.Shop
    NotAssignEmploye []store.Employe
}

type ShopEmploye struct{
    Name string
    EmployeCount int
    Employe []store.Employe
}

var ShopHandler = func(res http.ResponseWriter, req *http.Request){
    route, _ := NewRoute(res, req)

    var shop store.Shop = store.Shop{EntrepriseId: route.User.Id}
    route.Get(func() {
        var page = ShopPayload{
            User:  route.User,
            Shop: shop.GetAll(),
            NotAssignEmploye: shop.UnassignEmploye(),
        }
        route.Render(page, "route/protemplate.html", "route/shop.html")       
    })

    route.Post(&shop, func() {
        employes := ShopEmploye{
            Name: shop.Name,
            EmployeCount: shop.EmployeCount,
            Employe: shop.Employes(),
        }
        temp, _ := template.New("Shop Employes").Parse(`
            <h1 class="employe-shop">{{.Name}}<span class="employe-count">({{.EmployeCount}})</span></h1>
                {{range .Employe}}
                    <div class="employe">
                        <img src="/static/profile-picture.jpeg" class="employe-photo"  />
                        <div class="employe-content">
                            <h1 class="employe-name">{{.Name}}</h1>
                            <p class="employe-date">Ancienneté: {{.Age}}</p>
                        </div>
                        <div class="employe-action">
                            <button type="button" class="actionBtn deleteBtn" hx-delete="/employe" hx-vals='{"id": "{{.Id}}"}' hx-target="closest .employe" hx-swap="delete">Supprimer</button>
                            <dropdown-ele array='[{"id": "", "value": "User"}, {"id": "", "value": "Admin"}]'></dropdown-ele>
                            <button type="button" class="actionBtn">Velizy 2</button>
                            <button type="button" class="actionBtn">Planning</button>
                        </div>
                        <button type="button" class="employe-moreBtn" popovertarget="employe-actionModal">
                            <div class="dot" ></div>
                            <div class="dot" ></div>
                            <div class="dot" ></div>
                        </button>
                        <div id="employe-actionModal" popover>
                            <form type="" class="actionBtn" >
                                <h2 class="actionModal-header">Role</h2>
                                <label for="user">User</label>
                                <input type="radio" name="role" id="user" value="User" />
                                <label for="admin">Admin</label>
                                <input type="radio" name="role" id="admin" value="Admin" />
                            </form>
                            <div>
                                <h2 class="actionModal-header">Boutique</h2>
                                <dropdown-ele array='[{"id": "", "value": "Boutique 1"}, {"id": "", "value": "Boutique 2"}]'></dropdown-ele>
                            </div>
                            <div>
                                <h2>Planning</h2>
                                <button type="button" class="actionBtn">Planning</button>
                            </div>
                            <button type="button" class="actionBtn deleteBtn" hx-delete="/employe" hx-vals='{"id": "{{.Id}}"}' hx-target="closest .employe" hx-swap="delete">Supprimer</button>
                        </div>
                    </div>
                {{end}}
        `)
        temp.Execute(route.Response, employes)
    })

    route.Put(&shop, func() {
        if err := shop.Create(); err != nil{
            route.Notification("error", "error dans la creating")
            return
        }
        temp, _ := template.New("Shop card").Parse(`
            <div class="card">
                <div class="card-data">
                    <div class="shop-icon">
                        <svg width="24" height="24" style="width: 100%;height: 100%;" xmlns="http://www.w3.org/2000/svg" fill-rule="evenodd" clip-rule="evenodd"><path d="M21 14h.004l1.996 4h-2v4h2v2h-22v-2h2v-4h-2l1.996-4h.004v-14h18v14zm-12 5h-4v4h4v-4zm10 0h-4v4h4v-4zm-5 0h-4v4h4v-4zm6.386-4h-16.772l-1 2h18.772l-1-2zm-1.386-13h-14v11h14v-11zm-12 7h2v2h-2v-2zm4 0h2v2h-2v-2zm6 0v2h-2v-2h2zm-10-3h2v2h-2v-2zm4 0h2v2h-2v-2zm4 0h2v2h-2v-2zm-8-3h2v2h-2v-2zm4 0h2v2h-2v-2zm4 0h2v2h-2v-2z"/></svg>
                    </div>
                    <div class="card-content">
                        <h1>{{.Name}}</h1>
                        <p>{{.Adresse}}, {{.Postal}}</p>
                    </div>
                </div>
                <div class="card-footer">
                    <button type="button" class="card-footer-deleteBtn" hx-delete="/boutiques" hx-vals='{"id": "{{.Id}}"}' hx-target="closest .card" hx-swap="delete">Supprimer</button>
                    <button type="button" class="card-footer-btn" hx-post="/boutiques" hx-ext="json-enc" hx-vals='{"id": "{{.Id}}"}' hx-target="#shop-employes" onclick="onSelectShop(this)">Ouvrir</button>
                </div>
            </div>
        `)
        route.Response.Header().Add("HX-Trigger", "closeModal")
        route.Notification("success", "établissement crée")
        temp.Execute(route.Response, shop)

    })

    route.Delete(nil, func() {
        shop.Id = route.UrlEncoded["id"]
        if err := shop.Delete(); err != nil{
            route.Notification("error", "erorr dans la suppresion")
        }
        route.Notification("success", "établissement supprimé")
    })
}
