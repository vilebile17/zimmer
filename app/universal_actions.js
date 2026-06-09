import { goToDashboard, toggleTheme } from "/functions.js";

document.getElementById("navbar-logo").addEventListener("click", goToDashboard);
document.getElementById("theme-switcher").addEventListener(
        "click",
        toggleTheme,
);
