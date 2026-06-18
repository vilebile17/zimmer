import { snackbarSuccess, snackbarWarning } from "/functions.js";

function getMark() {
        return Number(document.getElementById("marks").value);
}

function getAssignmentID() {
        const linkToAss = document.getElementById("back-button").href;
        for (var i = linkToAss.length - 1; i > 0; i--) {
                if (linkToAss[i] == "/") {
                        return linkToAss.substring(i + 1);
                }
        }
}

async function grade() {
        const mark = getMark();
        if (mark < 0 || mark > 100) {
                snackbarWarning("marks must be in the range 0-100");
                return;
        }

        const classID = document.getElementById("class-id").textContent;
        const assignmentID = getAssignmentID();
        const submissionID = document
                .getElementById("submission-id")
                .textContent.trimStart()
                .substring(4);

        const response = await fetch(
                `/api/classes/${classID}/assignments/${assignmentID}/submissions/${submissionID}`,
                {
                        method: "PUT",
                        credentials: "include",
                        body: JSON.stringify({
                                score: mark,
                        }),
                },
        );

        if (response.ok) {
                snackbarSuccess(`Successfully graded the work (${mark}/100)`);
        } else {
                const error = await response.json();
                snackbarWarning(error?.error);
        }
}

document.getElementById("grade-button").addEventListener("click", grade);
