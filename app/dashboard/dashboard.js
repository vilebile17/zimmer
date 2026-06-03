async function fetchClasses() {
        const response = await fetch("/api/classes", {
                credentials: "include",
        });
        return response;
}

async function fetchUserData() {
        const response = await fetch("/api/users", {
                credentials: "include",
        });
        return response;
}

async function main() {
        const classes = await fetchClasses();
        const user = await fetchUserData();
        const classData = await classes.json();
        const userData = await user.json();
        console.log(classData);
        console.log(userData);
}

main();
