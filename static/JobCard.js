customElements.define("job-card", class JobCard extends HTMLElement{
    constructor(){
        super()
        this.innerHTML = `
            <div class="card">
                <div class="card-content">
                    <div class="card-header">
                        <img src="/static/euro.svg" alt="DesignWave logo" class="company-photo w-photo">
                        <div class="card-info">
                            <div style="display:flex;gap:1rem;align-items:center"> 
                                <img src="/static/euro.svg" alt="DesignWave logo" class="company-photo m-photo">
                                <h3 class="card-title">${this.getAttribute("title")}</h3>
                            </div>
                            <div class="card-company-info">
                                <div class="company-name">
                                    <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                        <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
                                        <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
                                        <line x1="12" y1="22.08" x2="12" y2="12"></line>
                                    </svg>
                                    <span class="name">${this.getAttribute("entreprise") ? this.getAttribute("entreprise") : "N/A"}</span>
                                </div>
                                <div class="card-salary">
                                    <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                        <line x1="12" y1="1" x2="12" y2="23"></line>
                                        <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
                                    </svg>
                                    <span class="salary">${this.getAttribute("salary") ? this.getAttribute("salary") : "N/A"}</span>
                                </div>
                                <div class="card-location">
                                    <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                        <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                                        <circle cx="12" cy="10" r="3"></circle>
                                    </svg>
                                    <span class="location">${this.getAttribute("fulladdr")}</span>
                                </div>
                            </div>
                            <div class="card-meta">
                                <span class="badge">${this.getAttribute("contract")}</span>
                                <span class="badge">${this.getAttribute("fulltime")}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card-footer">
                    <div class="card-posted">
                        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="10"></circle>
                            <polyline points="12 6 12 12 16 14"></polyline>
                        </svg>
                        <span>${this.getAttribute("date")}</span>
                    </div>
                    <a href="/job/${this.getAttribute("id")}?tp=${this.getAttribute("isThirdParty") === "true" ? "true" : "false"}" class="button">Postuler
                        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="9 18 15 12 9 6"></polyline>
                        </svg>
                    </a>
                </div>
            </div>
        `
    }

})

