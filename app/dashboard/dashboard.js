import {
        snackbarSuccess,
        snackbarDanger,
        snackbarWarning,
} from "/functions.js";

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

async function fetchNumAssignmentsDue() {
        const response = await fetch("/api/numAssignmentsDue", {
                credentials: "include",
        });
        return response;
}

function setUpJoinModal() {
        var modal = document.getElementById("join-class-modal");
        var btn = document.getElementById("join-class-button");
        var join = document.getElementById("join");
        var span = document.getElementById("join-close");
        var classID = document.getElementById("class-id");

        btn.onclick = function () {
                modal.style.display = "block";
        };
        span.onclick = function () {
                modal.style.display = "none";
        };
        join.onclick = async function () {
                if (!classID.value) {
                        snackbarWarning("classID parameter must be filled");
                        return;
                }

                const response = await fetch(
                        `/api/classes/${classID.value}/members`,
                        {
                                method: "POST",
                                credentials: "include",
                        },
                );

                if (response.ok) {
                        modal.style.display = "none";
                        snackbarSuccess("successfully joined class!");
                        setTimeout(() => {
                                window.location.replace(`/c/${classID.value}`);
                                window.location.href = `/c/${classID.value}`;
                        }, 1000);
                } else {
                        const data = await response.json();
                        snackbarDanger(data?.error);
                }
        };
        window.onclick = function (event) {
                if (event.target == modal) {
                        modal.style.display = "none";
                }
        };
}

function setUpCreateModal() {
        var modal = document.getElementById("create-class-modal");
        var btn = document.getElementById("create-class-button");
        var create = document.getElementById("create");
        var span = document.getElementById("create-close");
        var className = document.getElementById("class-name");

        btn.onclick = function () {
                modal.style.display = "block";
        };
        span.onclick = function () {
                modal.style.display = "none";
        };
        create.onclick = async function () {
                if (!className.value) {
                        snackbarWarning("className parameter must be filled");
                        return;
                }

                const response = await fetch(`/api/classes`, {
                        method: "POST",
                        body: JSON.stringify({
                                name: className.value,
                        }),
                        headers: {
                                "Content-Type": "application/json",
                        },
                        credentials: "include",
                });

                if (response.ok) {
                        modal.style.display = "none";
                        snackbarSuccess("successfully created the class!");
                        setTimeout(() => {
                                location.reload();
                        }, 1500);
                } else {
                        const data = await response.json();
                        snackbarDanger(data?.error);
                }
        };
        window.onclick = function (event) {
                if (event.target == modal) {
                        modal.style.display = "none";
                }
        };
}

async function logout() {
        await fetch("/api/logout", {
                method: "POST",
                credentials: "include",
        });
        window.location.replace("/login");
        window.location.href = "/login";
}

function writeNumClasses(classData) {
        let total = 0;
        const studentNull = classData.classesAsStudent == null;
        const teacherNull = classData.classesAsTeacher == null;
        total += studentNull ? 0 : classData.classesAsStudent.length;
        total += teacherNull ? 0 : classData.classesAsTeacher.length;
        document.getElementById("numClasses").textContent = total;
}

function writeClasses(classData) {
        const holderDiv = document.createElement("div");
        holderDiv.id = "classes-holder";
        if (classData.classesAsStudent != null) {
                const studentDiv = document.createElement("div");
                studentDiv.classList.add("card");
                let title = document.createElement("h3");
                title.textContent = "Classes as a Student:";
                studentDiv.appendChild(title);

                for (let i = 0; i < classData.classesAsStudent.length; i++) {
                        let bulletPoint = document.createElement("ul");
                        let className = document.createElement("a");
                        className.href = `/c/${classData.classesAsStudent[i].ID}`;
                        className.textContent =
                                classData.classesAsStudent[i].Name;
                        bulletPoint.appendChild(className);
                        studentDiv.appendChild(bulletPoint);
                }

                holderDiv.appendChild(studentDiv);
        }

        if (classData.classesAsTeacher != null) {
                const teacherDiv = document.createElement("div");
                teacherDiv.classList.add("card");
                let title = document.createElement("h3");
                title.textContent = "Classes as a Teacher:";
                teacherDiv.appendChild(title);

                for (let i = 0; i < classData.classesAsTeacher.length; i++) {
                        let bulletPoint = document.createElement("ul");
                        let className = document.createElement("a");
                        className.href = `/c/${classData.classesAsTeacher[i].ID}`;
                        className.textContent =
                                classData.classesAsTeacher[i].Name;
                        bulletPoint.appendChild(className);
                        teacherDiv.appendChild(bulletPoint);
                }

                holderDiv.appendChild(teacherDiv);
        }

        document.body.insertBefore(
                holderDiv,
                document.getElementsByClassName("snackbar")[0],
        );
}

async function main() {
        let response = await fetchUserData();
        if (response.status == 401) {
                window.location.replace("/login");
                window.location.href = "/login";
                return;
        }
        const userData = await response.json();
        response = await fetchClasses();
        const classData = await response.json();
        response = await fetchNumAssignmentsDue();
        const numAss = await response.json();

        document.getElementById("username").textContent = userData.name;
        document.getElementById("username").href = `/u/${userData.id}`;
        document.getElementById("numAssignments").textContent = numAss.num;

        writeNumClasses(classData);

        if (
                classData.classesAsTeacher != null ||
                classData.classesAsStudent != null
        ) {
                writeClasses(classData);
        }

        setUpJoinModal();
        setUpCreateModal();
}

main();
document.getElementById("logout-button").addEventListener("click", logout);
