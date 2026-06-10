function getClassID() {
        const text = document.getElementById("class-id").textContent;
        return text.substring(4);
}

async function fetchAssignments() {
        const classID = getClassID();
        const response = await fetch(`/api/classes/${classID}/assignments`, {
                credentials: "include",
        });
        return response;
}

function addCard(assignment) {
        const assignmentDiv = document.createElement("div");
        assignmentDiv.classList.add("card");

        const title = document.createElement("a");
        title.textContent = assignment.Title;
        title.classList.add("assignment-title");
        title.href = `/a/${assignment.ID}`;
        assignmentDiv.appendChild(title);

        const dueAt = document.createElement("p");
        dueAt.textContent = assignment.DueAt.Valid
                ? new Date(assignment.DueAt.Time).toUTCString()
                : "No due date";
        dueAt.classList.add("mini-text");
        assignmentDiv.appendChild(dueAt);

        document.body.insertBefore(assignmentDiv, null);
}

async function main() {
        const resp = await fetchAssignments();
        const assignments = await resp.json();
        if (!assignments) {
                return;
        }

        for (const a of assignments) {
                addCard(a);
        }
}

main();
