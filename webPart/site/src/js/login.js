let loginForm = document.getElementById("loginForm");

loginForm.addEventListener("submit", (e) => {
    e.preventDefault();

    let username = document.getElementById("usernameForm");
    let password = document.getElementById("passwordForm");


    fetch('http://localhost:8080/sign_in', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({username: username.value, password: password.value})
    }).then(async r => {
        if (r.status === 200) {
            let tokensJSON = await r.text()
            let tokens = JSON.parse(tokensJSON)
            setCookie("session_token", tokens.session_token, {"max-age":900})
            setCookie("refresh_token", tokens.refresh_token, {"max-age":7200})
            location.href = 'index.html';
        }
    })
});

