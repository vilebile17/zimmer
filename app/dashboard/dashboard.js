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

function writeNumClasses(classData) {
        var total = 0;
        const studentNull = classData.classesAsStudent == null;
        const teacherNull = classData.classesAsTeacher == null;
        total += studentNull ? 0 : classData.classesAsStudent.length;
        total += teacherNull ? 0 : classData.classesAsTeacher.length;
        document.getElementById("numClasses").textContent = total;
}

function writeClasses(classData) {
        if (classData.classesAsStudent != null) {
                const studentDiv = document.createElement("div");
                var title = document.createElement("h3");
                title.textContent = "Classes as a Student:";
                studentDiv.appendChild(title);

                for (var i = 0; i < classData.classesAsStudent.length; i++) {
                        var className = document.createElement("ul");
                        className.textContent =
                                classData.classesAsStudent[i].Name;
                        studentDiv.appendChild(className);
                }

                document.body.insertBefore(studentDiv, null);
        }

        if (classData.classesAsTeacher != null) {
                const teacherDiv = document.createElement("div");
                var title = document.createElement("h3");
                title.textContent = "Classes as a Teacher:";
                teacherDiv.appendChild(title);

                for (var i = 0; i < classData.classesAsTeacher.length; i++) {
                        var className = document.createElement("ul");
                        className.textContent =
                                classData.classesAsTeacher[i].Name;
                        teacherDiv.appendChild(className);
                }

                document.body.insertBefore(teacherDiv, null);
        }
}

async function main() {
        var response = await fetchNumAssignmentsDue();
        if (response.status == 401) {
                window.location.replace("/login");
                window.location.href = "/login";
        }
        const numAss = await response.json();
        response = await fetchClasses();
        const classData = await response.json();
        response = await fetchUserData();
        const userData = await response.json();

        console.log(classData);
        console.log(userData);
        console.log(numAss);
        document.getElementById("username").textContent = userData.name;
        document.getElementById("numAssignments").textContent = numAss.num;

        writeNumClasses(classData);
        writeClasses(classData);
}

main();
