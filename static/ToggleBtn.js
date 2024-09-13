
export default class ToggleBtn extends HTMLElement{
    constructor(){
        super()

        this.onToggle = ()=>{
            this.state.classList.toggle("toggleBtnAnimation")
            this.container.classList.toggle("toggleBtnBG")
            const newState = this.getAttribute("on").trim() === "true" ? "false" : "true" 
            this.setAttribute("on", newState)
        }
        this.id = this.getAttribute("id")
        this.container = document.createElement("div")
        this.container.classList.add("toggleBtn")
        this.container.onclick = this.onToggle

        this.state = document.createElement("div")
        this.state.classList.add("toggleBtn_state")

        this.container.appendChild(this.state)

        this.appendChild(this.container)

        if(this.getAttribute("on").trim() === "true"){ 
            this.state.classList.toggle("toggleBtnAnimation")
            this.container.classList.toggle("toggleBtnBG")
        }
    }
}

customElements.define("toggle-btn", ToggleBtn)

