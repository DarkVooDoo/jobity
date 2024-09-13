
export default class Dropdown extends HTMLElement{
    constructor(){
        super()
        this.maxWidth = 0
    }
    connectedCallback() {
        /**@param {PointerEvent} ele */
        this.onElementClick = (ele)=>{
            const value = ele.currentTarget.textContent
            this.display.value = value
            this.display.dataset.id = ele.currentTarget.dataset.id
            this.option.toggleAttribute("open")
            this.display.dispatchEvent(new Event("valueChange"))
        }
        this.array = JSON.parse(this.getAttribute("array"))
        this.value = this.getAttribute("value") || this.array[0].value

        this.container = document.createElement("div")
        this.container.classList.add("dropdown")
        this.container.style.height = this.getAttribute("height") || "100%"
        this.option = document.createElement("div")
        this.option.classList.add("dropdown_option")
        for(const arrayValue of this.array){
            if(arrayValue.id === this.value){
                this.value = arrayValue.value
            }
            const eleValue = document.createElement("p")
            eleValue.addEventListener("click", this.onElementClick)
            eleValue.classList.add("dropdown_option_value")
            eleValue.textContent = arrayValue.value
            eleValue.dataset.id = arrayValue.id
            if(this.maxWidth < arrayValue.value.length){
                this.maxWidth = arrayValue.value.length
                this.container.style.width = `${this.maxWidth+3}ch`
            }
            this.option.appendChild(eleValue)
        }
        this.display = document.createElement("input")
        this.display.style.backgroundColor = this.getAttribute("bgColor") || "white"
        this.display.readOnly = true
        this.display.name = this.getAttribute("id")
        this.display.addEventListener("click", ()=>{
            this.option.toggleAttribute("open")
        })
        this.display.classList.add("dropdown_display")
        this.display.value = this.value
        this.display.dataset.id = this.array[0].id
        this.container.appendChild(this.display)
        this.container.appendChild(this.option)
        this.appendChild(this.container)
    }
}

customElements.define("dropdown-ele", Dropdown)
