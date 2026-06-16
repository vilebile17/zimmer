import { goToDashboard, changeTheme } from "/functions.js";

document.getElementById("navbar-logo").addEventListener("click", goToDashboard);

document.getElementById("theme-switcher").addEventListener(
        "click",
        changeTheme,
);

const theme = localStorage.getItem("theme");
if (theme) {
        document.getElementById("theme-switcher").value = theme;
}

//const prefersLightScheme = window.matchMedia("(prefers-color-scheme: light)");
//if (prefersLightScheme.matches) {
//j       document.getElementById("theme-switcher").va;
//}
