import { snackbar } from "/functions.js";

async function signup() {
        let email = document.getElementById("email").value;
        let password = document.getElementById("password").value;
        let password2 = document.getElementById("confirm-password").value;
        let name = document.getElementById("name").value;

        if (email == "" || password == "" || name == "") {
                snackbar(
                        "Name, email and password parameters cannot be empty!",
                );
                return;
        }

        if (password != password2) {
                snackbar("Passwords don't match!");
                return;
        }

        if (password.length < 8) {
                snackbar("Password must be at least 8 characters long");
                return;
        }

        let response = await fetch("/api/users", {
                method: "POST",
                body: JSON.stringify({
                        email,
                        password,
                        name,
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

        console.log(response.status);

        response = await fetch("/api/login", {
                method: "POST",
                body: JSON.stringify({
                        email,
                        password,
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
}

document.getElementById("signup-button").addEventListener("click", signup);
