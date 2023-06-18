
let userGet = false
async function getUsername() {
    let token = getCookie("session_token")
    let refreshToken = getCookie("refresh_token")
    if (token === undefined && refreshToken === undefined) {
        console.log("no session tokens!")
        return
    }
    if (token === undefined) {
        let done = await updateTokensByRefresh(refreshToken)
        if (done)
            await getUsername()
        return
    }

    await fetch('http://localhost:8080/get_username_by_token', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': 'bearer ' + token
        }
    }).then(async r => {
        if (r.status !== 200) {
            if (await handleRequestAndUpdateTokens(r)) {
                await getUsername()
            }
            return

        }
        if (!location.href.endsWith("index.html")) {
            window.location.href = "index.html"
            return
        }
        let username = await r.text()
        let loginElement = document.getElementById("login")
        loginElement.hidden = true
        let usernameElement = document.getElementById("username")
        usernameElement.innerHTML = username
        usernameElement.hidden = false
        userGet = true
        let navbar = document.getElementById("navbar")
        navbar.appendChild(signOutButton())
    })
}

function signOutButton() {
    let btn = document.createElement("a")
    btn.setAttribute("class", "getstarted scrollto")
    btn.setAttribute("href", "index.html")
    btn.innerText = "Sign Out"
    btn.addEventListener("click", ()=>{
        deleteCookie("session_token")
        deleteCookie("refresh_token")
    })
    return btn
}

async function handleRequestAndUpdateTokens(request) {
    let refreshToken = getCookie("refresh_token")
    if (refreshToken === undefined)
        return false
    let errMessage = await request.text()
    if (!errMessage.startsWith("token is expired by")) {
        return false
    }
    return await updateTokensByRefresh(refreshToken)
}

async function updateTokensByRefresh(refreshToken) {
    let res = false
    await fetch('http://localhost:8080/sign_in_by_refresh', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'RefreshToken': refreshToken
        }
    }).then(async r => {
        if (r.status !== 200) {
            return
        }
        res = true
        let tokensJSON = await r.text()
        let tokens = JSON.parse(tokensJSON)
        setCookie("session_token", tokens.session_token, {"max-age": 900})
        setCookie("refresh_token", tokens.refresh_token, {"max-age": 7200})
    })
    return res
}

promiseOfGetUsername = getUsername()