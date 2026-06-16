const theme = localStorage.getItem("theme");
if (theme) {
        console.log(`Setting theme to ${theme}`);
        document.documentElement.dataset.theme = theme;
} else {
        document.documentElement.dataset.theme = "mocha";
}
