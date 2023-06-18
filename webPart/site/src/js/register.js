let registerForm = document.getElementById("registerForm");

registerForm.addEventListener("submit", (e) => {
    e.preventDefault();

    let username = document.getElementById("usernameForm");
    let name = document.getElementById("nameForm");
    let password = document.getElementById("passwordForm");


    fetch('http://localhost:8080/sign_up', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({username: username.value, password: password.value, name: name.value})
    }).then(async r => {
        if (r.status !== 200) {
            alert(await r.text())
            return
        }
        alert("user successfully created!")
        location.href = 'login.html';
    })
});