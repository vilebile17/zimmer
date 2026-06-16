import { snackbar } from "/functions.js";

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

function setUpModal() {
        var modal = document.getElementById("join-class-modal");
        var btn = document.getElementById("join-class-button");
        var join = document.getElementById("join");
        var span = document.getElementsByClassName("close")[0];
        var classID = document.getElementById("class-id");

        btn.onclick = function () {
                modal.style.display = "block";
        };
        span.onclick = function () {
                modal.style.display = "none";
        };
        join.onclick = async function () {
                console.log(classID.value);
                const response = await fetch(
                        `/api/classes/${classID.value}/members`,
                        {
                                method: "POST",
                                credentials: "include",
                        },
                );

                if (response.ok) {
                        modal.style.display = "none";
                        snackbar("successfully joined class!");
                        setTimeout(() => {
                                location.reload();
                        }, 1000);
                } else {
                        const data = await response.json();
                        snackbar(data?.error);
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
                document.getElementById("snackbar"),
        );
}

async function main() {
        setUpModal();

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

        console.log(classData);
        console.log(userData);
        console.log(numAss);
        document.getElementById("username").textContent = userData.name;
        document.getElementById("numAssignments").textContent = numAss.num;

        writeNumClasses(classData);

        if (
                classData.classesAsTeacher != null ||
                classData.classesAsStudent != null
        ) {
                writeClasses(classData);
        }
}

main();
document.getElementById("logout-button").addEventListener("click", logout);
