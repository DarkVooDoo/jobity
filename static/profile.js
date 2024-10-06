const months = ["Jan", "Fév", "Mar", "Avr", "Mai", "Jun", "Jul", "Aoû", "Sep", "Oct", "Nov", "Déc"]
let fetchAdresseHtml
/**@type {number} */
let fetchTimer 
/**@type {HTMLImageElement} */
const picture = document.getElementById("picture_ele")
/**@type {HTMLDialogElement} */
const cvDialog = document.getElementById("profile_cvPopover")
/**@type {HTMLInputElement} */
const cvName = document.getElementById("profile_cvPopover_filename")
const suggestAdr = document.getElementById("suggest_adr")
const postalField = document.getElementById("postal")
const cityField = document.getElementById("city")
/**@param {HTMLInputElement} ele */
const onPictureChange = (ele)=>{
    const fileReader = new FileReader()
    fileReader.addEventListener("load", (e)=>{
        picture.src = e.target.result

    })
    if(ele.files.length > 0 && ele.files[0].size < 50000){
        fileReader.readAsDataURL(ele.files[0])
    }
}

/**@param {HTMLInputElement} ele */
const onFetchAddr = (ele)=>{
    fetchAdresseHtml = ""
    if (fetchTimer) clearTimeout(fetchTimer)
    fetchTimer = setTimeout(async ()=>{
        try{
            const fetchAddr = await fetch(`https://api-adresse.data.gouv.fr/search/?q=${ele.value}&type=municipality&limit=5`)
            /**@type {{features: {geometry:{coordinates: string[]},properties: {name: string, postcode: string, city: string}}[]}} */
                const result = await fetchAddr.json()
            for(const adr of result.features){
                fetchAdresseHtml += `<p class="item" data-lat="${adr.geometry.coordinates[1]}" data-long="${adr.geometry.coordinates[0]}"  onclick="onSelectAdresse(this)">${adr.properties.name}, ${adr.properties.postcode}</p>` 
            }
            suggestAdr.innerHTML = fetchAdresseHtml
            suggestAdr.style.display = "block"
            if(result.features.length > 0){
                const firstELement = result.features[0]
                propagateProfileAddr(firstELement.properties.city, firstELement.properties.postcode, firstELement.geometry.coordinates[1], firstELement.geometry.coordinates[0])
            }
        }catch(err){}
    }, 500)
}

/**
* @param {string} city    
* @param {string} postal    
* @param {string} lat    
* @param {string} long    
* */
const propagateProfileAddr = (city, postal, lat, long)=>{
    postalField.value = parseInt(postal)
    cityField.value = city
    cityField.dataset.lat = lat
    cityField.dataset.long = long

}

/**@param {HTMLParagraphElement} ele */
const onSelectAdresse = (ele)=>{
    const [city, postal] = ele.textContent.split(",")
    propagateProfileAddr(city, postal.trim(), ele.dataset.lat, ele.dataset.long)
    suggestAdr.style.display = "none"
}

/**@param {Event} ev */
const onCVChange = (ev)=>{
    cvDialog.showModal()
    if (ev.target.files.length < 1) return
    cvName.value = ev.target.files[0].name
}

const onCloseCVDialog = ()=>{
    cvDialog.close()
}

const profile = document.querySelector(".profile_editable")
const scenes = document.getElementsByClassName("scene")
const tabs = document.getElementsByClassName("tab")
const description = document.querySelector("#description textarea")
const descriptionLength = document.getElementById("length")
descriptionLength.textContent = 250 - description.textContent.length

const calcPopoverPosition = ()=>{
    const curriculumActionBtn = document.querySelector("#cv .submitBtn")
    const popover = document.querySelector("#cv #curriculum_action")

    const {left, top} = curriculumActionBtn.getBoundingClientRect()
    popover.style.left = `calc(${left}px - 100px  + 2rem)`
    popover.style.top = `calc(${top}px + 2.3rem)`
}

const onHideCurriculumAction = (ele)=>{
    ele.parentElement.hidePopover()
}

const onDescriptionChange = (ele)=>{
    const length = ele.value.length
    if (length > 249){
        ele.value = ele.value.substring(0, 249)
    }
    descriptionLength.textContent = 250 - length
}
const onChangeScene = (ev, scene)=>{
    const currentTab = ev.target.closest(".tab")
    for(const tab of tabs){
        const selection = tab.querySelector(".bg")
        selection.classList.remove("selected")
        tab.style.outline = "1px solid lightgray"
    }
    const selectedTab = currentTab.querySelector(".bg")
    selectedTab.classList.add("selected")
    currentTab.style.outline = "none"
    for(const scene of scenes){
        scene.classList.add("hidden")
    }
    if (scene === "CV"){
        cv.classList.toggle("hidden")
    }else{
        profile.classList.toggle("hidden")
    }
}

const onDateChange = (ele)=>{
    const sibling = ele.nextElementSibling ? ele.nextElementSibling : ele.previousElementSibling
    if (ele.nextElementSibling){
        sibling.setAttribute("min", ele.value)
    }else{
        sibling.setAttribute("max", ele.value)
    }

}

const isFullCapacity = (type)=>{
    const cards = document.querySelectorAll(`#${type} .card`)
    const newBtn = document.querySelector(`#${type} .templateBtn`)
    if (cards.length === 3){
        newBtn.classList.add("hidden")
    }else{
        newBtn.classList.remove("hidden")
    }
}

const onNewWork = (ele)=>{
    const cards = document.querySelectorAll("#work .card")
    if (cards.length > 2) return
    const card = document.createElement("div")
    card.classList.add("card")
    card.dataset.name = "work"
    card.setAttribute("ondrop", "onDrop(event)")
    card.setAttribute("ondragover", "onDragOver(event)")
    card.setAttribute("draggable", "true")
    card.setAttribute("ondragstart", "onDragStart(event)")
    card.innerHTML = `
        <div>
        <input type="text" class="input position" placeholder="Web Dev Full Stack" required />
        <input type="text" class="input entreprise" placeholder="Netflix" required />
        <div class="dates">
            <div class="picker start_date">
                <div class="picker-display" onclick="toggleModal(this)">Date</div>
                    <div class="modal" style="display: none;">
                    <div class="modal-content">
                        <span class="close" onclick="toggleModal(this)">&times;</span>
                        <div class="picker-header">
                            <span class="picker-title">Select Month</span>
                        </div>
                        <div class="picker-grid"></div>
                    </div>
                </div>
            </div>
            <div class="picker end_date">
                <div class="picker-display" onclick="toggleModal(this)">Date</div>
                <div class="modal" style="display: none;">
                    <div class="modal-content">
                        <span class="close" onclick="toggleModal()">&times;</span>
                        <div class="picker-header">
                            <span class="picker-title">Select Month</span>
                        </div>
                        <div class="picker-grid"></div>
                    </div>
                </div>
            </div>
        </div>
        <div class="task">
            <h2 >Description du poste</h2>
            <div class="description">
                <button type="button" class="closeBtn" onclick="onDeleteWorkDescription(this)">X</button>
                <div contenteditable="true" class="field" required placeholder="Backend testing"></div>
            </div>
            <button type="button" class="newDescriptionPlaceholder" onclick="onAddDescription(this)"></button>
        </div>
        <button type="button" class="deleteBtn" onclick="onDeleteCard(this)">Supprimer</button>
        </div>`
    ele.insertAdjacentElement('beforebegin', card)
    isFullCapacity('work')
}

const onDeleteCard = (ele)=>{
    const card = ele.closest(".card")
    const parentName = card.parentElement.id
    card.remove()
    isFullCapacity(parentName)
}
const onDeleteWorkDescription = (ele)=>{
    const description = ele.closest(".description")
    description.remove()
}

const onNewSchool = (ele)=>{
    const cards = document.querySelectorAll("#diploma .card")
    if (cards.length > 2) return
    const schoolCard = document.createElement("div")
    schoolCard.classList.add("card")
    schoolCard.dataset.name = "school"
    schoolCard.setAttribute("ondrop", "onDrop(event)")
    schoolCard.setAttribute("ondragover", "onDragOver(event)")
    schoolCard.setAttribute("draggable", "true")
    schoolCard.setAttribute("ondragstart", "onDragStart(event)")
    schoolCard.innerHTML = `
        <div>
        <input type="text" class="input name" placeholder="Baccalauréat" required /> 
        <input type="text" class="input establishment" placeholder="Lycée Jean Villar" required /> 
            <div class="dates">
                <div class="picker start_date">
                    <div class="picker-display" onclick="toggleModal(this)">Date</div>
                    <div class="modal" style="display: none;">
                        <div class="modal-content">
                            <span class="close" onclick="toggleModal(this)">&times;</span>
                            <div class="picker-header">
                                <span class="picker-title">Select Month</span>
                            </div>
                            <div class="picker-grid"></div>
                        </div>
                    </div>
                </div>
                <div class="picker end_date">
                    <div class="picker-display" onclick="toggleModal(this)">Date</div>
                    <div class="modal" style="display: none;">
                        <div class="modal-content">
                            <span class="close" onclick="toggleModal()">&times;</span>
                            <div class="picker-header">
                                <span class="picker-title">Select Month</span>
                            </div>
                            <div class="picker-grid"></div>
                        </div>
                    </div>
                </div>
            </div>
        <button type="button" class="deleteBtn" onclick="onDeleteCard(this)">Supprimer</button>
        </div>`
    ele.insertAdjacentElement('beforebegin', schoolCard)
    isFullCapacity('diploma')
}
const onAddDescription = (ele)=>{
    const newDescription = document.createElement("div")
    newDescription.classList.add("description", "descriptionFadeIn")
    newDescription.innerHTML = `
        <button type="button" class="closeBtn" onclick="onDeleteWorkDescription(this)">X</button>
        <div contenteditable="true" class="field" required placeholder="Backend testing"></div>
        `
    newDescription.addEventListener("animationend", ()=>{
        newDescription.classList.remove("descriptionFadeIn")
    })
    ele.insertAdjacentElement('beforebegin', newDescription)
}
const onNewSkill = (ele, name)=>{
    const newSkill = document.createElement("div")
    newSkill.classList.add("field")
    newSkill.innerHTML = `
        <input type="text" class="input ${name}" autocomplete="off" />
        <button type="button" class="closeBtn" onclick="onRemoveLangue(this)">X</button>
        `
    ele.before(newSkill)
}
const onRemoveLangue = (ele)=>{
    ele.parentElement.remove()
}

const getProfileValues = ()=>{
    return {lat: parseFloat(cityField.dataset.lat), long: parseFloat(cityField.dataset.long)}
}

// Curriculum values for save
const getValues = ()=>{
    const workCards = document.querySelectorAll("#work .card")
    const schoolCards = document.querySelectorAll("#diploma .card")
    const skills = document.querySelectorAll("#skill .skill")
    const interests = document.querySelectorAll("#interest .interest")

    let skillList = []
    let interestList = []
    let schoolList = []
    let workList = []
    let taskList = []
    for(const work of workCards){
        taskList.length = 0
        const position = work.querySelector(".position")
        const entreprise = work.querySelector(".entreprise")
        const start = work.querySelector(".start_date .picker-display")
        const end = work.querySelector(".end_date .picker-display")
        const tasks = work.querySelectorAll(".task .field")
        for(const task of tasks){
            taskList.push(task.innerText)
        }
        workList.push({position: position.value, entreprise: entreprise.value, start_date: start.textContent, end_date:end.textContent, description: [...taskList]})
    }
    for(const school of schoolCards){
        const name = school.querySelector(".name") 
        const establishment = school.querySelector(".establishment") 
        const start = school.querySelector(".start_date .picker-display") 
        const end = school.querySelector(".end_date .picker-display") 
        schoolList.push({name: name.value, establishment: establishment.value, start_date: start.textContent, end_date: end.textContent})
    }

    for(const skill of skills){
        skillList.push(skill.value)
    }
    for(const interest of interests){
        interestList.push(interest.value)
    }

    console.log(workList)
    return {school: schoolList, work: workList, interest: interestList, skill: skillList}
}


// Drag and drop handler for cards
let lastDrag
const onDragStart = (ev)=>{
    lastDrag = ev.target
    ev.dataTransfer.setData("text/plain", ev.target.dataset.name)
}
const onDragOver = (ev)=>{
    ev.preventDefault()
}
const onDrop = (ev)=>{
    ev.preventDefault()
    ev.stopPropagation()
    const type = ev.dataTransfer.getData("text/plain")
    const dropZone = ev.target.closest(".card") || ev.target
    const lastData = lastDrag.firstElementChild
    const dropData = dropZone.firstElementChild
    if (dropZone.dataset.name === type){
        lastDrag.replaceChild(dropData, lastData)
        dropZone.appendChild(lastData)
    }
}

// Custom Month, Year picker
let currentView = 'month';
let selectedMonth = null;
let selectedYear = null;

const currentYear = new Date().getFullYear();
const years = Array.from({ length: 60 }, (_, i) => currentYear - i);

function toggleModal(ele) {

    const modal = ele.closest(".picker").querySelector(".modal")
    modal.style.display = modal.style.display === 'none' ? 'block' : 'none';
    if (modal.style.display === 'block') {
        currentView = 'month';
        renderGrid(ele);
    }
}

function renderGrid(ele) {
    const pickerGrid = ele.parentElement.querySelector(".picker-grid")
    const pickerTitle = ele.parentElement.querySelector(".picker-title")
    pickerGrid.innerHTML = '';
    const items = currentView === 'month' ? months : years;
    pickerTitle.textContent = currentView === 'month' ? 'Mois' : 'Année';

    items.forEach((item, index) => {
        const div = document.createElement('div');
        div.classList.add('picker-item');
        div.textContent = item;
        div.addEventListener('click', () => selectItem(index, ele));
        pickerGrid.appendChild(div);
    });
}

function selectItem(index, ele) {
    if (currentView === 'month') {
        selectedMonth = index;
        currentView = 'year';
        renderGrid(ele);
    } else {
        selectedYear = years[index];
        updateDisplay(ele);
        toggleModal(ele);
        currentView = 'month';
    }
}

function updateDisplay(ele) {
    if (selectedMonth !== null && selectedYear !== null) {
        ele.textContent = `${months[selectedMonth]} ${selectedYear}`;
    }
}

const onPickerTitle = ()=>{
    if (currentView === 'year') {
        currentView = 'month';
        renderGrid();
    }
}












