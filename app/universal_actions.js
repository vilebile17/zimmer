import { goToDashboard, toggleTheme } from "/functions.js";

document.getElementById("navbar-logo").addEventListener("click", goToDashboard);
document.getElementById("theme-switcher").addEventListener(
        "click",
        toggleTheme,
);

const prefersLightScheme = window.matchMedia("(prefers-color-scheme: light)");
if (prefersLightScheme.matches) {
        document.getElementById("theme-switcher").click();
}
