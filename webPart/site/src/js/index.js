let tableBodyElem = document.getElementById("table_body")
let tableElem = document.getElementById("main_table")
let addTaskForm = document.getElementById("addTaskForm");
let todoList = []
async function main() {
    await promiseOfGetUsername
    loadList()
    if (userGet) {
        addTaskForm.hidden = false
    }
}

main()



function createTableElemHTML(num, listElem) {
    let res = document.createElement("tr")
    let th = document.createElement("th")
    th.setAttribute("scope", "row")
    th.innerText = num

    let tdInput = document.createElement("td")
    let input = document.createElement("input")
    input.setAttribute("type", "checkbox")
    input.checked = listElem.done
    input.addEventListener('change', () => {
        updateListItem(listElem.id, input.checked)
    })
    tdInput.appendChild(input)


    let tdName = document.createElement("td")
    tdName.innerText = listElem.taskname
    let tdDescription = document.createElement("td")
    tdDescription.innerText = listElem.description
    let tdDeleteButton = document.createElement("td")
    let deleteButton = document.createElement("button")
    deleteButton.setAttribute("type", "button")
    deleteButton.setAttribute("class", "btn btn-danger")
    deleteButton.setAttribute("style", "float:right;")
    deleteButton.innerText = "Delete"
    deleteButton.addEventListener("click", async () => {
        let done = await deleteTaskByID(listElem.id)
        console.log(done)
        if (done)
            tableBodyElem.removeChild(res)
    })
    tdDeleteButton.appendChild(deleteButton)
    res.appendChild(th)
    res.appendChild(tdInput)
    res.appendChild(tdName)
    res.appendChild(tdDescription)
    res.appendChild(tdDeleteButton)
    return res
}


function loadList() {
    let session_token = getCookie("session_token")
    if (session_token === undefined) {
        return
    }
    fetch('http://localhost:8080/get_tasks_of_user', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': 'bearer ' + session_token
        }
    }).then(async r => {
        if (r.status === 200) {
            let todolistText = await r.text()
            if (todolistText === "null") {
                return
            }
            todoList = JSON.parse(todolistText)
            todoList.forEach((listElem, i) => {
                let newDiv = createTableElemHTML(i + 1, listElem)
                tableBodyElem.appendChild(newDiv)
                tableElem.hidden = false
            })

        }
    })
}


async function deleteTaskByID(id) {
    let res = false
    let session_token = getCookie("session_token")
    if (session_token === undefined) {
        return
    }
    await fetch('http://localhost:8080/delete_task_by_id', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': 'bearer ' + session_token,
            'TaskID': id
        }
    }).then(async r => {
        if (r.status === 200) {
            res = true
            return
        }
        if (await handleRequestAndUpdateTokens(r)) {
            res = await deleteTaskByID(id)
        }

    })
    return res
}

function updateListItem(id, done) {
    console.log(done)
    let curElem = todoList.find(element => element.id === id)
    console.log(curElem)
    let session_token = getCookie("session_token")
    if (session_token === undefined) {
        return
    }
    fetch('http://localhost:8080/update_todo_list_item', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            taskID: id,
            token: session_token,
            taskname: curElem.taskname,
            description: curElem.description,
            done: done
        })
    }).then(async r => {
        if (r.status === 200) {
            return
        }
        if (await handleRequestAndUpdateTokens(r)) {
            await updateListItem(id, done)
            return
        }
        console.error("ERROR DURING UPDATE")
    })

}


addTaskForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    let taskname = document.getElementById("validationDefault01");
    let description = document.getElementById("validationDefault02");
    let results = await addTaskByID(taskname.value, description.value)
    if (results.done) {
        let newCurElem = {}
        newCurElem.id = results.lastID
        newCurElem.taskname = taskname.value
        newCurElem.description = description.value
        newCurElem.done = false
        todoList.push(newCurElem)
        let a = createTableElemHTML(todoList.length, newCurElem)
        tableBodyElem.appendChild(a)
        tableElem.hidden = false
    }
});


async function addTaskByID(taskname, description) {
    let res = false
    let lastID = 0
    let session_token = getCookie("session_token")
    if (session_token === undefined) {
        return res
    }
    await fetch('http://localhost:8080/add_task_for_user', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            token: session_token,
            taskname: taskname,
            description: description,
            done: false
        })
    }).then(async r => {
        if (r.status === 200) {
            res = true
            let lastIDSTR = await r.text()
            lastID = parseInt(lastIDSTR, 10)
            return
        }
        if (await handleRequestAndUpdateTokens(r)) {
            res = await addTaskByID(taskname, description)
        }

    })
    return {done: res, lastID: lastID}
}