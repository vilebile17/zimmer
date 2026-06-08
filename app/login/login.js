async function login() {
        let email = document.getElementById("email").value;
        let password = document.getElementById("password").value;

        if (email == "" || password == "") {
                window.alert("Email and password parameters cannot be empty!");
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
                error = await r.json();
                window.alert(error.error);
        } else {
                window.location.replace("/dashboard");
                window.location.href = "/dashboard";
        }
        console.log(r.status);
}

document.getElementById("login-button").addEventListener("click", login);
