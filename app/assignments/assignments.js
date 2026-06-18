import { snackbarDanger, snackbarSuccess } from "/functions.js";

var alreadyHandedIn = false;

function removeItem(array, itemToRemove) {
        const index = array.indexOf(itemToRemove);

        if (index !== -1) {
                array.splice(index, 1);
        }

        return array;
}

async function handIn() {
        const classID = getClassID();
        const assignmentID = getAssignmentID();
        const work = document.getElementById("student-work");
        var response;
        var result;

        if (!alreadyHandedIn) {
                response = await fetch(
                        `/api/classes/${classID}/assignments/${assignmentID}/submissions`,
                        {
                                method: "POST",
                                body: JSON.stringify({
                                        answers: work.value,
                                }),
                                headers: {
                                        "Content-Type": "application/json",
                                },
                                credentials: "include",
                        },
                );
                if (response.ok) {
                        snackbarSuccess("successfully handed in!");
                } else {
                        result = await response.json();
                        snackbarDanger(result?.error);
                }
        } else {
                response = await fetch(
                        `/api/classes/${classID}/assignments/${assignmentID}/submissions`,
                        {
                                method: "PUT",
                                body: JSON.stringify({
                                        answers: work.value,
                                }),
                                headers: {
                                        "Content-Type": "application/json",
                                },
                                credentials: "include",
                        },
                );
                if (response.ok) {
                        snackbarSuccess("successfully updated submission!");
                } else {
                        result = await response.json();
                        snackbarDanger(result?.error);
                }
        }

        if (result) {
                console.log(result);
        }
        alreadyHandedIn = true;
        const gradeSpan = document.getElementById("status");
        gradeSpan.textContent =
                gradeSpan.textContent == "assigned"
                        ? "handed in"
                        : gradeSpan.textContent;
}

function getClassID() {
        const linkToClass = document.getElementById("back-button").href;
        for (var i = linkToClass.length - 1; i > 0; i--) {
                if (linkToClass[i] == "/") {
                        return linkToClass.substring(i + 1);
                }
        }
}

function getAssignmentID() {
        const idText = document.getElementById("id");
        return idText.textContent.substring(4).trim();
}

async function loadStudentStuff() {
        console.log("this guy is a student");
        const gradeSpan = document.getElementById("status");
        const response = await fetch(
                `/api/classes/${getClassID()}/assignments/${getAssignmentID()}/submissions`,
                {
                        method: "GET",
                        credentials: "include",
                },
        );

        const userWork = await response.json();
        console.log(userWork);
        if (userWork?.Answers) {
                const workBox = document.getElementById("student-work");
                workBox.value = userWork.Answers.String;
                alreadyHandedIn = true;

                gradeSpan.textContent = userWork?.Score.Valid
                        ? `${userWork.Score.Int32}/100`
                        : "handed in";
        } else {
                console.log("No old submission found");
        }
}

async function createYetToHandIn(students) {
        if (students.length === 0) {
                return;
        }

        const bigDiv = document.createElement("div");
        bigDiv.classList.add("card");
        const header = document.createElement("h3");
        header.textContent = "Yet to Hand in";
        bigDiv.appendChild(header);

        for (const studentID of students) {
                const user = await (
                        await fetch(`/api/users/${studentID}`)
                ).json();

                let studentPoint = document.createElement("ul");
                studentPoint.textContent = user.name;

                bigDiv.appendChild(studentPoint);
        }
        document.body.insertBefore(
                bigDiv,
                document.getElementById("snackbar-success"),
        );
}

async function loadTeacherStuff(classID, assignmentID) {
        console.log("This guy is the teacher");
        document.getElementById("submission-card").remove();
        const submissions = await (
                await fetch(
                        `/api/classes/${classID}/assignments/${assignmentID}/submissions`,
                        {
                                credentials: "include",
                        },
                )
        ).json();

        const studentsJSON = (
                await (
                        await fetch(`/api/classes/${classID}/members`, {
                                credentials: "include",
                        })
                ).json()
        ).students;

        let students = [];
        for (const student of studentsJSON) {
                students.push(student.ID);
        }

        if (!submissions) {
                createYetToHandIn(students);
                return;
        }

        const submissionsDiv = document.createElement("div");
        submissionsDiv.classList.add("card");
        const header = document.createElement("h3");
        header.textContent = "Submissions";
        submissionsDiv.appendChild(header);

        for (const submission of submissions) {
                const user = await (
                        await fetch(`/api/users/${submission.UserID}`)
                ).json();
                students = removeItem(students, user.id);

                let studentPoint = document.createElement("ul");
                let studentName = document.createElement("a");

                studentName.textContent = user.name;
                studentName.href = `/s/${submission.ID}`;
                studentName.classList.add("submission-link");

                studentPoint.classList.add("mini-text");
                studentPoint.appendChild(studentName);
                studentPoint.appendChild(
                        document.createTextNode(
                                ` - on ${new Date(submission.UpdatedAt).toUTCString()}`,
                        ),
                );

                submissionsDiv.appendChild(studentPoint);
        }
        document.body.insertBefore(
                submissionsDiv,
                document.getElementById("snackbar-success"),
        );
        createYetToHandIn(students);
}

async function main() {
        const classID = getClassID();
        const assignmentID = getAssignmentID();
        const user = await (
                await fetch("/api/users", {
                        credentials: "include",
                })
        ).json();
        const classs = await (
                await fetch(`/api/classes/${classID}`, {
                        credentials: "include",
                })
        ).json();

        console.log(user.id);
        console.log(classs);
        if (user.id === classs.TeacherID) {
                loadTeacherStuff(classID, assignmentID);
        } else {
                loadStudentStuff(classID, assignmentID);
        }
}

document.getElementById("hand-in-button").addEventListener("click", handIn);
main();
