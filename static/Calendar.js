class Calendar extends HTMLElement{
    constructor(){
        super()
        this.bgColor = this.getAttribute("bg") || "white"
        this.shadow = this.attachShadow({mode: "open"})
        this.value = this.getAttribute("value") || "Choissir une date"
        this.key = this.getAttribute("key")
        this.shadow.innerHTML = `
            <style>
            *{
                margin: 0;
                padding: 0;
                box-sizing: border-box;
            }
        input:focus{
            outline: none;
        }
            .date-time{
                .calendar{
                    position: relative;
                    width: 200px;
                    .date-input{
                        width: 100%;
                        height: 2rem;
                        border: 1px solid darkgray;
                        border-radius: 5px;
                        text-align: center;
                        &:focus + .date-option{
                            @starting-style{
                                opacity: 0;
                            }
                            opacity: 1;
                            display: block;
                        }
                    }
                        .date-option{
                            position: absolute;
                            top: calc(100% + 9px);
                            width: 100%;
                            display: none;
                            background-color: white;
                            border: 1px solid darkgray;
                            border-radius: 5px;
                            padding: 5px;
                            opacity: 0;
                            transition: opacity 150ms linear, display 150ms linear allow-discrete;
                            z-index: 1;
                            &:hover{
                                opacity: 1;
                                display: block;
                            }
                                &:has(.date-change-picker > * > .date-change-month:focus), &:has(.date-change-picker > * > .date-change-year:focus){
                                    display: block;
                                    opacity: 1;
                                }
                                &::after{
                                    content: "";
                                    position: absolute;
                                    top: -15px;
                                    left: calc(50% - 6px);
                                    border: 6px solid transparent;
                                    border-bottom: 6px solid darkgray;
                                    border-right: 6px solid darkgray;
                                    border-bottom-right-radius: 2px;
                                    rotate: 45deg;
                                    z-index: 1;
                                }
                                .date-change{
                                    display: grid;
                                    grid-template-columns: .3fr 1fr .3fr;
                                    height: 1.5lh;
                                    gap: .2rem;
                                    .date-change-picker{
                                        display: flex;
                                        gap: .3rem;
                                        .date-change-switch{
                                            position: relative;
                                            display: flex;
                                            align-items: center;
                                            justify-content: center;
                                            .date-change-month, .date-change-year{
                                                width: 100%;
                                                height: 1.5lh;
                                                border: none;
                                                text-align: center;
                                                &:focus + *{
                                                    display: block;
                                                    opacity: 1;
                                                    @starting-style{
                                                        opacity: 0;
                                                    }
                                                }
                                            }

                                                .date-change-picker-month, .date-change-picker-year{
                                                    position: absolute;
                                                    top: 100%;
                                                    width: 100%;
                                                    border: 1px solid darkgray;
                                                    display: none;
                                                    max-height: calc((1lh + 10px) * 10);
                                                    background-color: white;
                                                    overflow-y: scroll;
                                                    z-index: 1;
                                                    transition: opacity 150ms linear, display 150ms linear allow-discrete;
                                                    opacity: 0;
                                                    &:hover{
                                                        opacity: 1;
                                                        display: block;
                                                    }
                                                        .date-change-picker-cell{
                                                            padding: 5px;
                                                            text-align: center;
                                                            &:hover{
                                                                background-color: rgb(220, 220, 220);
                                                            }
                                                        }
                                                }
                                        }
                                    }
                                        .date-change-btn{
                                            width: 100%;
                                            height: 100%;
                                            display: flex;
                                            align-items: center;
                                            justify-content: center;
                                            background-color: transparent;
                                            border: 1px solid darkgray;
                                            border-radius: 5px;
                                            &::before{
                                                display: block;
                                                content: "";
                                                width: .5lh;
                                                aspect-ratio: 1/1;
                                                border-bottom: 3px solid black;
                                                border-right: 3px solid black;
                                            }
                                                &:first-child::before{
                                                    rotate: 135deg;
                                                }
                                                &:last-child::before{
                                                    rotate: -45deg;
                                                }
                                        }
                                }
                                .date-calendar{
                                    display: grid;
                                    justify-content: center;
                                    align-items: center;
                                    grid-template-columns: repeat(7, 1fr);
                                    border-radius: 10px;
                                    .date-calendar-cell{
                                        position: relative;
                                        width: 100%;
                                        aspect-ratio: 1/1;
                                        display: flex;
                                        background-color: transparent;
                                        border: none;
                                        align-items: center;
                                        justify-content: center;
                                    }
                                        .select{
                                            &::before{
                                                position: absolute;
                                                content: attr(data-day);
                                                color: white;
                                                display: flex;
                                                align-items: center;
                                                justify-content: center;
                                                width: 2lh;
                                                aspect-ratio: 1/1;
                                                border-radius: 50%;
                                                background-color: rgb(0, 0, 0);
                                            }
                                        }
                                }
                        }
                }
            }
            </style>
            <div class="date-time">
            <div class="calendar">
            <input type="text" readonly class="date-input" name="${this.key}" />
            <div class="date-option">
            <div class="date-change">
            <button type="button" class="date-change-btn" data-type="minus"></button>
            <div class="date-change-picker">
            <div class="date-change-switch">
            <input type="text" class="date-change-month" readonly />
            <div class="date-change-picker-month"></div>
            </div>
            <div class="date-change-switch">
            <input type="text" class="date-change-year" readonly />
            <div class="date-change-picker-year"></div>
            </div>
            </div>
            <button type="button" class="date-change-btn" data-type="plus" ></button>
            </div>
            <div class="date-calendar"></div>
            </div>
            </div>
            </div>
            `
        const shadowRoot = this.shadowRoot
        const format = this.getAttribute("format") || "yyyy-mm-dd"
        const monthTotalDay = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31]
        const week = ["L", "M", "M", "J", "V", "S", "D"]
        const month = ["Jan", "Fév", "Mars", "Avr", "Mai", "Juin", "Juil", "Août", "Sept", "Oct", "Nov", "Déc"]
        const now = new Date()
        const changeBtns = shadowRoot.querySelectorAll(".date-change-btn")
        for(const btn of changeBtns){
            btn.addEventListener("click", function(){
                const operation = btn.dataset.type
                const month = parseInt(dateMonth.dataset.month)
                const year = parseInt(dateYear.value)
                if (operation === "minus"){
                    const [newYear, newMonth] = month === 0 ? [year - 1, 11] : [year, month - 1]
                    renderCalendar(newYear, newMonth)
                }else{
                    const [newYear, newMonth] = month === 11 ? [year + 1, 0] : [year, month + 1]
                    renderCalendar(newYear, newMonth)
                }
            })
        }
        const dateDisplay = shadowRoot.querySelector(".date-input")
        const dateOption = shadowRoot.querySelector(".date-option")
        const dateCalendar = shadowRoot.querySelector(".date-calendar")
        const dateMonth = dateOption.querySelector(".date-change-month")
        const dateYear = dateOption.querySelector(".date-change-year")
        const monthPicker = dateOption.querySelector(".date-change-picker-month")
        const yearPicker = dateOption.querySelector(".date-change-picker-year")
        for(let m = 0; m < month.length; m++){
            const newCell = document.createElement("p")
            newCell.classList.add("date-change-picker-cell")
            newCell.textContent = month[m]
            newCell.addEventListener("click", function(){
                renderCalendar(parseInt(dateYear.value), parseInt(m))
            })
            monthPicker.appendChild(newCell)
        }
        for (let y = now.getUTCFullYear() - 60; y < now.getUTCFullYear() + 20; y++){
            const newCell = document.createElement("p")
            newCell.classList.add("date-change-picker-cell")
            newCell.textContent = y
            newCell.addEventListener("click", function(){
                renderCalendar(parseInt(y), parseInt(dateMonth.dataset.month))
            })
            yearPicker.appendChild(newCell)
        }

        const days = new Array(7*6).fill(0)
        
        const displayValue = (day, month, year)=>{
            const monthLabel = month[month]
            const monthNumber = month < 10 ? `0${month}` : month
            switch (format){
                case "dd-MM-yyyy":
                    dateDisplay.value = `${day} ${monthLabel}, ${year}`
                    break;
                case "dd-mm-yyyy":
                    dateDisplay.value = `${day}-${monthNumber}-${year}`
                    break;
                default:
                    dateDisplay.value = `${year}-${monthNumber}-${day}`

            }
        }

        dateDisplay.value = this.value
        const renderCalendar = (year, currentMonth)=>{
            dateCalendar.innerHTML = ""

            for(let i = 0; i < week.length; i++){
                dateCalendar.innerHTML += `<b class="date-calendar-cell">${week[i]}</b>`
            }
            const firstMonthDay = new Date(year, currentMonth, 1)
            const monthWeekStartDay = firstMonthDay.getUTCDay()
            dateMonth.value = month[currentMonth]
            dateMonth.dataset.month = currentMonth
            dateYear.value = year
            for(let day = 0; day < days.length; day++){
                const monthDay = day - monthWeekStartDay
                const newBtn = document.createElement("button")
                newBtn.setAttribute("type", "button")
                newBtn.classList.add("date-calendar-cell")
                if (day < monthWeekStartDay){
                    const lastMonthTotalDay = currentMonth === 0 ?  monthTotalDay[11] : monthTotalDay[currentMonth-1]
                    const lastMonthDay = lastMonthTotalDay  - monthWeekStartDay + day + 1
                    newBtn.dataset.day = lastMonthDay
                    newBtn.textContent = lastMonthDay
                    newBtn.setAttribute("disabled", "")
                    dateCalendar.appendChild(newBtn)
                    continue
                }else if(monthDay + 1 > monthTotalDay[currentMonth]){
                    const nextMonthDay = day - monthWeekStartDay - monthTotalDay[currentMonth] + 1
                    newBtn.dataset.day = nextMonthDay
                    newBtn.textContent = nextMonthDay
                    newBtn.setAttribute("disabled", "")
                    dateCalendar.appendChild(newBtn)
                    continue
                }
                newBtn.dataset.day = monthDay + 1
                newBtn.textContent = monthDay + 1
                if (monthDay + 1 === now.getUTCDate() && now.getUTCMonth() === currentMonth && year === now.getUTCFullYear()) newBtn.classList.add("select")
                dateCalendar.appendChild(newBtn)
                newBtn.addEventListener("click", function(){
                    displayValue(newBtn.textContent, currentMonth, year)
                    const selects = shadowRoot.querySelectorAll('.select')
                    for(const select of selects){
                        select.classList.remove('select')
                    }
                    newBtn.classList.add('select')
                })
            }
        }
        renderCalendar(now.getUTCFullYear(), now.getUTCMonth())
        htmx.process(this.shadow) 
    }
}
customElements.define("calendar-ele", Calendar)
