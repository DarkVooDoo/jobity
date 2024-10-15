customElements.define("job-card", class JobCard extends HTMLElement{
    constructor(){
        super()
        this.shadow = this.attachShadow({mode: "open"})
        this.picture = !this.getAttribute("picture") ? "hidden" : ""
        this.shadow.innerHTML = `
            <style>
            .card {
                width: 100%;
                background-color: #fff;
                border-radius: 10px;
                box-shadow: var(--card_shadow);
                margin-bottom: 20px;
                overflow: hidden;
                .card-content {
                    padding: 15px;
                    .card-header {
                        display: flex;
                        gap: 16px;
                        .company-photo {
                            height: 3rem;
                            aspect-ratio: 1/1;
                            object-fit: cover;
                            border-radius: 10px;
                        }
                            .card-info {
                                flex-grow: 1;
                                .card-top{
                                    display:flex;
                                    gap:1rem;
                                    margin-bottom: .5rem;
                                    .card-title {
                                        view-transition-name: title;
                                        font-size: 1rem;
                                        font-weight: 600;
                                        text-wrap: balance;
                                        margin: 0;

                                    }
                                }
                                    .card-company-info {
                                        display: flex;
                                        flex-direction: column;
                                        height: 6rem;
                                        justify-content: space-evenly;
                                        .company-name, .card-location, .card-salary {
                                            display: flex;
                                            align-items: center;
                                            gap: 8px;
                                            color: #666;
                                            margin-bottom: .2rem;
                                            .salary, .location, .name{
                                                font-size: .85rem;
                                            }
                                        }
                                            .card-salary, .company-name{
                                                text-wrap: balance;
                                            }
                                    }
                                    .card-meta {
                                        display: flex;
                                        align-items: center;
                                        gap: 12px;
                                        margin-top: 8px;
                                    }
                            }
                    }
                }
            }
            .badge {
                background-color: var(--primary_light_color);
                color: black;
                font-size: 0.75rem;
                padding: 2px 8px;
                border-radius: 9999px;
            }
            .card-footer {
                background-color: var(--primary_light_color);
                padding: 16px 20px;
                display: flex;
                justify-content: space-between;
                align-items: center;
                .card-posted {
                    display: flex;
                    align-items: center;
                    gap: 8px;
                    font-size: 0.875rem;
                    color: black;
                }
                    .button {
                        background-color: transparent;
                        color: #204767;
                        border: none;
                        height: 32px;
                        width: 32px;
                        border-radius: 50%;
                        font-size: 0.875rem;
                        font-weight: 500;
                        cursor: pointer;
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        gap: 4px;
                        font-weight: bold;
                        background-color: rgb(230,230,230);
                        &:hover{
                            background-color: #b2c4d3;
                        }
                    }
            }
            .icon {
                width: 16px;
                height: 16px;
            }
            </style>
            <div class="card">
                <div class="card-content">
                    <div class="card-header">
                        <div class="card-info">
                            <div class="card-top" style=""> 
                                <img src=${this.getAttribute("picture")}  alt="DesignWave logo" class="company-photo m-photo ${this.picture}"/>
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
                    <a href="/job/${this.getAttribute("id")}?tp=${this.getAttribute("isThirdParty") === "true" ? "true" : "false"}" class="button">
                        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="9 18 15 12 9 6"></polyline>
                        </svg>
                    </a>
                </div>
            </div>
        `
    }
})

