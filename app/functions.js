function snackbarHelper(content, id) {
        let snackbar = document.getElementById(id);
        snackbar.textContent = content;
        snackbar.classList.add("show");

        setTimeout(function () {
                snackbar.className = "snackbar";
        }, 2900);
        return;
}

function snackbar(content) {
        snackbarHelper(content, "snackbar");
}
function snackbarSuccess(content) {
        snackbarHelper(content, "snackbar-success");
}
function snackbarWarning(content) {
        snackbarHelper(content, "snackbar-warning");
}
function snackbarDanger(content) {
        snackbarHelper(content, "snackbar-danger");
}
function snackbarInfo(content) {
        snackbarHelper(content, "snackbar-info");
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

export {
        snackbar,
        changeTheme,
        goToDashboard,
        snackbarDanger,
        snackbarInfo,
        snackbarSuccess,
        snackbarWarning,
};
