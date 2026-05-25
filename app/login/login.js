function printEmailAndPassword() {
        let email = document.getElementById("email").value;
        let password = document.getElementById("password").value;
        login(email, password);
}

async function login(email, password) {
        const response = fetch("http://localhost:8080/api/login", {
                method: "POST",
                body: JSON.stringify({
                        email: email,
                        password: password,
                }),
                headers: {
                        "Content-Type": "application/json",
                },
        })
        
        const status = await response;
        console.log(status);
}

document.getElementById("login-button").addEventListener("click", printEmailAndPassword);
