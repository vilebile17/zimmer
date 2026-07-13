const theme = localStorage.getItem("theme");
if (theme) {
        console.log(`Setting theme to ${theme}`);
        document.documentElement.dataset.theme = theme;
} else {
        const darkModeMql = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)');
        if (darkModeMql && darkModeMql.matches) {
                console.log("dark theme detected");
                document.documentElement.dataset.theme = "mocha";
                localStorage.theme = "mocha";
        } else {
                console.log("light theme detected");
                document.documentElement.dataset.theme = "latte";
                localStorage.theme = "latte";
        }
}
