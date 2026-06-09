function snackbar(content) {
        let snackbar = document.getElementById("snackbar");
        snackbar.textContent = content;
        snackbar.className = "show";

        setTimeout(function () {
                snackbar.className = snackbar.className.replace("show", "");
        }, 3000);
        return;
}

function toggleTheme() {
        document.body.classList.toggle("light");
}

function goToDashboard() {
        window.location.replace("/dashboard");
        window.location.href = "/dashboard";
}

export { snackbar, toggleTheme, goToDashboard };
