var alreadyHandedIn = false;

async function handIn() {
        const work = document.getElementById("student-work");
        const classID = getClassID();
        const assignmentID = getAssignmentID();
        var response;

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
        }

        console.log(await response.json());
        alreadyHandedIn = true;
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

async function main() {
        const response = await fetch(
                `/api/classes/${getClassID()}/assignments/${getAssignmentID()}/submissions`,
                {
                        method: "GET",
                        credentials: "include",
                },
        );

        const userWork = await response.json();
        if (userWork?.Answers) {
                const workBox = document.getElementById("student-work");
                workBox.value = userWork.Answers.String;
                alreadyHandedIn = true;
        } else {
                console.log("No old submission found");
        }
}

document.getElementById("hand-in-button").addEventListener("click", handIn);
main();
