function openTab(event, tabID) {
        let tabContent = document.getElementsByClassName("tab-content");
        for (let i = 0; i < tabContent.length; i++) {
                tabContent[i].style.display = "none";
        }

        let tabLinks = document.getElementsByClassName("tab-links");
        for (i = 0; i < tabLinks.length; i++) {
                tabLinks[i].className = tabLinks[i].className.replace(
                        " active",
                        "",
                );
        }

        document.getElementById(tabID).style.display = "block";
        event.currentTarget.className += " active";
}

function getClassID() {
        let text = document.getElementById("class-id").textContent;
        text = text.trim();
        return text.substring(4);
}

async function fetchAssignments() {
        const classID = getClassID();
        const response = await fetch(`/api/classes/${classID}/assignments`, {
                credentials: "include",
        });
        return response;
}

async function fetchMembers() {
        const classID = getClassID();
        const response = await fetch(`/api/classes/${classID}/members`, {
                credentials: "include",
        });
        return response;
}

function addCard(assignment, grandadDiv) {
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

        grandadDiv.appendChild(assignmentDiv);
}

async function main() {
        const resp = await fetchAssignments();
        const assignments = await resp.json();
        if (!assignments) {
                return;
        }

        const assignmentsDiv = document.createElement("div");
        assignmentsDiv.id = "assignments";
        assignmentsDiv.classList.add("tab-content");
        for (const a of assignments) {
                addCard(a, assignmentsDiv);
        }
        document.body.insertBefore(assignmentsDiv, null);

        const students = (await (await fetchMembers()).json()).students;
        const studentsDiv = document.createElement("div");
        studentsDiv.id = "students";
        studentsDiv.classList.add("tab-content");
        studentsDiv.classList.add("card");

        for (const s of students) {
                let studentPoint = document.createElement("ul");
                let studentName = document.createElement("a");
                studentName.textContent = s.Name;
                studentName.classList.add("student-name");
                studentName.href = `/u/${s.ID}`;
                studentPoint.appendChild(studentName);
                studentsDiv.appendChild(studentPoint);
        }
        document.body.insertBefore(studentsDiv, null);

        document.getElementById("default").click();
}

main();
