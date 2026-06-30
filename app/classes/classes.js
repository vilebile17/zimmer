import { snackbarSuccess, snackbarWarning } from "/functions.js";

function createDefaultTab(text, id) {
        const outerDiv = document.createElement("div");
        outerDiv.classList.add("tab-content");
        outerDiv.id = id;
        defaultTabHelper(text, outerDiv);

        if (id === "students") {
                console.log("adding button");
                createDangerButton(outerDiv);
        }

        return outerDiv;
}

function defaultTabHelper(text, outerDiv) {
        const innerDiv = document.createElement("div");
        innerDiv.classList.add("card");
        const unorderedList = document.createElement("ul");
        const textDOM = document.createElement("p");

        textDOM.classList.add("card-heading");
        textDOM.textContent = text;
        unorderedList.appendChild(textDOM);

        innerDiv.appendChild(unorderedList);
        outerDiv.appendChild(innerDiv);
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

async function fetchResources() {
        const classID = getClassID();
        const response = await fetch(`/api/classes/${classID}/resources`, {
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

        const assignmentsDiv = document.createElement("div");
        assignmentsDiv.id = "assignments";
        assignmentsDiv.classList.add("tab-content");

        if (await isUserTeacher()) {
                const createAssignmentButton = document.createElement("button");
                createAssignmentButton.className = "centered-buttons";
                createAssignmentButton.textContent = "Create Assignment";
                createAssignmentButton.onclick = showCreateAssignmentModal;

                const createButton = document.getElementById("create");
                createButton.onclick = createAssignment;

                assignmentsDiv.appendChild(createAssignmentButton);
        }

        if (!assignments) {
                defaultTabHelper("No assignments yet...", assignmentsDiv);
                document.body.insertBefore(assignmentsDiv, null);
                return;
        }

        for (const a of assignments) {
                addCard(a, assignmentsDiv);
        }
        document.body.insertBefore(assignmentsDiv, null);
}

function showCreateAssignmentModal() {
        const modal = document.getElementById("create-assignment-modal");
        modal.style.display = "block";
}

async function isUserTeacher() {
        const classID = getClassID();
        let response = await fetch(`/api/classes/${classID}`, {
                credentials: "include",
        });
        const classObj = await response.json();

        response = await fetch(`/api/users`, {
                credentials: "include",
        });
        const user = await response.json();

        return classObj.TeacherID === user.id;
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

        const grandadDiv = document.createElement("div");
        grandadDiv.appendChild(studentsDiv);
        createDangerButton(grandadDiv);
        document.body.insertBefore(
                grandadDiv,
                document.getElementsByClassName("modal")[0],
        );
}

async function createDangerButton(outerDiv) {
        const leaveButton = document.createElement("button");
        leaveButton.id = "leave-button";
        leaveButton.className = "centered-buttons";

        if (await isUserTeacher()) {
                leaveButton.textContent = "Delete Class";
                leaveButton.onclick = deleteClass;
        } else {
                leaveButton.textContent = "Leave Class";
                leaveButton.onclick = leaveClass;
        }

        outerDiv.id = "students";
        outerDiv.className = "tab-content";

        outerDiv.appendChild(leaveButton);
}

async function leaveClass() {
        if (!confirm("Are you sure you want to leave this class?")) {
                return;
        }

        let response = await fetch("/api/users", {
                credentials: "include",
        });
        const user = await response.json();
        const classID = getClassID();

        response = await fetch(`/api/classes/${classID}/members/${user.id}`, {
                credentials: "include",
                method: "DELETE",
        });
        if (!response.ok) {
                const error = await response.json();
                snackbarWarning(error.error);
                return;
        }

        snackbarSuccess("Successfully left class");
        setTimeout(() => {
                window.location.replace("/dashboard");
                window.location.href = "/dashboard";
        }, 2000);
}

async function deleteClass() {
        if (
                !confirm(
                        "Are you sure that you want to delete this class?\nThis action CANNOT be undone.",
                )
        ) {
                return;
        }

        const classID = getClassID();
        const response = await fetch(`/api/classes/${classID}`, {
                credentials: "include",
                method: "DELETE",
        });
        if (!response.ok) {
                const error = await response.json();
                snackbarWarning(error.error);
                return;
        }

        snackbarSuccess("Successfully deleted class");
        setTimeout(() => {
                window.location.replace("/dashboard");
                window.location.href = "/dashboard";
        }, 2000);
}

async function createAllResources() {
        const resources = await (await fetchResources()).json();
        const classID = getClassID();

        if (!resources) {
                document.body.insertBefore(
                        createDefaultTab("No resources yet...", "resources"),
                        null,
                );
                return;
        }

        const grandadDiv = document.createElement("div");
        grandadDiv.id = "resources";
        grandadDiv.classList.add("tab-content");

        for (const r of resources) {
                const resourceDiv = document.createElement("div");
                resourceDiv.classList.add("card");

                const title = document.createElement("a");
                title.textContent = r.title;
                title.classList.add("card-heading");
                title.onclick = showResource(classID, r.id);
                resourceDiv.appendChild(title);

                const createdAt = document.createElement("p");
                createdAt.textContent = `Posted on ${new Date(r.created_at).toUTCString()}`;
                createdAt.classList.add("mini-text");
                resourceDiv.appendChild(createdAt);

                grandadDiv.appendChild(resourceDiv);
        }

        document.body.insertBefore(grandadDiv, null);
}

function showResource(classID, resourceID) {
        return async function () {
                const modal = document.getElementById("resources-modal");
                modal.style.display = "block";

                const response = await fetch(
                        `/api/classes/${classID}/resources/${resourceID}`,
                        {
                                credentials: "include",
                        },
                );
                const resource = await response.json();

                const modalHeader = document.getElementById("resource-title");
                modalHeader.textContent = resource.title;
                const modalContent =
                        document.getElementById("resource-content");
                modalContent.textContent = resource.content;
        };
}

async function createAssignment() {
        const title = document.getElementById("title-input").value;
        const instructions =
                document.getElementById("instructions-input").value;
        const due_at = document.getElementById("due-at-input").value;
        const allow_late = document.getElementById("allow-late-input").checked;
        const classID = getClassID();

        console.log(title);
        console.log(instructions);
        console.log(due_at);
        console.log(allow_late);

        const resp = await fetch(`/api/classes/${classID}/assignments`, {
                credentials: "include",
                method: "POST",
                body: JSON.stringify({
                        title,
                        instructions,
                        due_at,
                        allow_late,
                }),
        });

        const respObj = await resp.json();
        if (!resp.ok) {
                snackbarWarning(respObj?.error);
                return;
        }
        snackbarSuccess("Successfully made assignment");

        setTimeout(() => {
                location.reload();
        }, 2000);
}

function setUpResourcesModal() {
        var modal = document.getElementById("resources-modal");
        var span = document.getElementById("resources-close");

        span.onclick = function () {
                modal.style.display = "none";
        };
        window.onclick = function (event) {
                if (event.target == modal) {
                        modal.style.display = "none";
                }
        };
}

function setUpCreateAssignmentModal() {
        var modal = document.getElementById("create-assignment-modal");
        var span = document.getElementById("create-close");

        span.onclick = function () {
                modal.style.display = "none";
        };
        window.onclick = function (event) {
                if (event.target == modal) {
                        modal.style.display = "none";
                }
        };
}

async function main() {
        await createAllAssignments();
        await createAllStudents();
        await createAllResources();
        setUpResourcesModal();
        setUpCreateAssignmentModal();

        document.getElementById("default-tab").click();
}

main();
