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

async function main() {
        const classes = await fetchClasses();
        const user = await fetchUserData();
        const classData = await classes.json();
        const userData = await user.json();

        console.log(classData);
        console.log(userData);
        document.getElementById("username").textContent = userData.name;
        document.getElementById("numClasses").textContent =
                classData.classesAsStudent.length +
                classData.classesAsTeacher.length;

        if (classData.classesAsStudent.length != 0) {
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

        if (classData.classesAsTeacher.length != 0) {
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

main();
