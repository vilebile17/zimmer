function snackbar(content) {
        let snackbar = document.getElementById("snackbar");
        snackbar.textContent = content;
        snackbar.className = "show";

        setTimeout(function () {
                snackbar.className = snackbar.className.replace("show", "");
        }, 3000);
        return;
}

function goToDashboard() {
        window.location.replace("/dashboard");
        window.location.href = "/dashboard";
}

function changeTheme() {
        const theme = document.getElementById("theme-switcher").value;
        document.documentElement.dataset.theme = theme;
        localStorage.setItem("theme", theme);
}

export { snackbar, changeTheme, goToDashboard };
