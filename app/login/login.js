async function login() {
        let email = document.getElementById("email").value;
        let password = document.getElementById("password").value;
        const response = fetch("/api/login", {
                method: "POST",
                body: JSON.stringify({
                        email: email,
                        password: password,
                }),
                headers: {
                        "Content-Type": "application/json",
                },
        });

        let r = await response;
        console.log(r.status);
}

document.getElementById("login-button").addEventListener("click", login);
