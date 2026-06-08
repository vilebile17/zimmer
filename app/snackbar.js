function snackbar(content) {
        let snackbar = document.getElementById("snackbar");
        snackbar.textContent = content;
        snackbar.className = "show";

        setTimeout(function () {
                snackbar.className = snackbar.className.replace("show", "");
        }, 3000);
        return;
}

export { snackbar };
