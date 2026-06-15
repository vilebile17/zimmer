async function handIn(work) {
        const classID = getClassID();
        console.log(classID);
        const assignmentID = getAssignmentID();
        console.log(assignmentID);
        const response = await fetch(
                `/api/classes/${classID}/assignments/${assignmentID}/submissions`,
                {
                        method: "POST",
                        body: JSON.stringify({
                                answers: work,
                        }),
                        headers: {
                                "Content-Type": "application/json",
                        },
                        credentials: "include",
                },
        );
        return await response.json();
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
        const work = document.getElementById("student-work");
        console.log(await handIn(work.value));
}

document.getElementById("hand-in-button").addEventListener("click", main);
