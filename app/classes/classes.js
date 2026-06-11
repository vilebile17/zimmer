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

function createDefaultTab(text, id) {
        const outerDiv = document.createElement("div");
        outerDiv.classList.add("card");
        outerDiv.classList.add("tab-content");
        outerDiv.id = id;

        const unorderedList = document.createElement("ul");
        const textDOM = document.createElement("p");
        textDOM.classList.add("card-heading");
        textDOM.textContent = text;
        unorderedList.appendChild(textDOM);
        outerDiv.appendChild(unorderedList);
        return outerDiv;
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
        title.classList.add("card-heading");
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

async function createAllAssignments() {
        const resp = await fetchAssignments();
        const assignments = await resp.json();
        if (!assignments) {
                document.body.insertBefore(
                        createDefaultTab(
                                "No assignments yet...",
                                "assignments",
                        ),
                        null,
                );
                return;
        }

        const assignmentsDiv = document.createElement("div");
        assignmentsDiv.id = "assignments";
        assignmentsDiv.classList.add("tab-content");
        for (const a of assignments) {
                addCard(a, assignmentsDiv);
        }
        document.body.insertBefore(assignmentsDiv, null);
}

async function createAllStudents() {
        const students = (await (await fetchMembers()).json()).students;

        if (!students) {
                document.body.insertBefore(
                        createDefaultTab("No students yet...", "students"),
                        null,
                );
                return;
        }

        const studentsDiv = document.createElement("div");
        studentsDiv.id = "students";
        studentsDiv.classList.add("tab-content");
        studentsDiv.classList.add("card");

        for (const s of students) {
                let studentPoint = document.createElement("ol");
                let studentName = document.createElement("a");
                studentName.textContent = s.Name;
                studentName.classList.add("card-heading");
                studentName.href = `/u/${s.ID}`;
                studentPoint.appendChild(studentName);
                studentsDiv.appendChild(studentPoint);
        }
        document.body.insertBefore(studentsDiv, null);
}

async function main() {
        await createAllAssignments();
        await createAllStudents();
        document.getElementById("default").click();
}

main();
