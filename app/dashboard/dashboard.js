async function fetchClasses() {
        const response = await fetch("/api/classes", {
                credentials: "include",
        });
        return response;
}

async function main() {
        const response = await fetchClasses();
        const data = await response.json();
        console.log(data);
}

main();
