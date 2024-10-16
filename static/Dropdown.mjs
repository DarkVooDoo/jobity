customElements.define("dropdown-ele", class Dropdown extends HTMLElement{
    constructor(){
        super()
        this.style.display = "inline-block"
        this.shadow = this.attachShadow({mode: "open"})
        this.items = JSON.parse(this.getAttribute("items"))
        if (!this.items || !this.id) throw new Error("needs require fields {items || id}")
        this.itemsElement = ""
        for(let i = 0; i < this.items.length; i++){
            this.itemsElement += `<p class="dropdown-option-item" data-id="${this.items[i].id}"  onclick="(function(ele){
                const dropdown = ele.closest('.dropdown')
                const input = dropdown.querySelector('.dropdown-input')
                input.value = ele.textContent
                input.dataset.id = ele.dataset.id
                input.dispatchEvent(new CustomEvent('change', {detail: ele.textContent}))
            })(this)">${this.items[i].value}</p>`
        }
        this.shadow.innerHTML = `
            <style>
            *{
                margin: 0;
                padding: 0;
                box-sizing: border-box;
            }
            .dropdown{
                position: relative;
                width: 100%;
                height: 100%;
                .dropdown-input{
                    height: 100%;
                    width: 100%;
                    &:focus + .dropdown-option{
                        display: block;
                    }
                }
                    .dropdown-option{
                        position: absolute;
                        top: calc(100% + 1px);
                        width: 100%;
                        background: white;
                        border: 1px solid darkgray;
                        display: none;
                        &:has(> .dropdown-option-item:hover){
                            display: block;
                        }
                            .dropdown-option-item{
                                padding: 3px 5px;
                                &:hover{
                                    background: lightgray;
                                }
                            }
                    }
            }
            </style>
            <div class="dropdown">
                <input type="text" class="dropdown-input" readonly />
                <div class="dropdown-option">
                    ${this.itemsElement}
                </div>
            </div>
            `
    }
})
