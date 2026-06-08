async function login() {
        let email = document.getElementById("email").value;
        let password = document.getElementById("password").value;

        if (email == "" || password == "") {
                let snackbar = document.getElementById("snackbar");
                snackbar.textContent =
                        "Email and password parameters cannot be empty!";
                snackbar.className = "show";

                setTimeout(function () {
                        snackbar.className = snackbar.className.replace(
                                "show",
                                "",
                        );
                }, 3000);

                return;
        }

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
        if (r.status >= 400) {
                let error = await r.json();
                let snackbar = document.getElementById("snackbar");
                snackbar.textContent = error.error;
                snackbar.className = "show";

                setTimeout(function () {
                        snackbar.className = snackbar.className.replace(
                                "show",
                                "",
                        );
                }, 3000);
        } else {
                window.location.replace("/dashboard");
                window.location.href = "/dashboard";
        }
        console.log(r.status);
}

document.getElementById("login-button").addEventListener("click", login);
