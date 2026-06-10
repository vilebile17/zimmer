import { snackbar } from "/functions.js";

async function login() {
        let email = document.getElementById("email").value;
        let password = document.getElementById("password").value;

        if (email == "" || password == "") {
                snackbar("Email and password parameters cannot be empty!");
                return;
        }

        const response = await fetch("/api/login", {
                method: "POST",
                body: JSON.stringify({
                        email: email,
                        password: password,
                }),
                headers: {
                        "Content-Type": "application/json",
                },
        });

        if (response.status >= 400) {
                let error = await response.json();
                snackbar(error.error);
                return;
        }
        window.location.replace("/dashboard");
        window.location.href = "/dashboard";
        console.log(response.status);
}

document.getElementById("login-button").addEventListener("click", login);
