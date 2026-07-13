import { snackbarSuccess, snackbarDanger } from "/functions.js";

async function fetchRequester() {
        const response = await fetch("/api/users", {
                credentials: "include",
        });
        return await response.json();
}

function getUserID() {
        const idDOM = document.getElementById("user-id");
        return idDOM.textContent.substring(4);
}

function setUpModal() {
        var modal = document.getElementById("modal");
        var span = document.getElementById("close-modal");

        span.onclick = function () {
                modal.style.display = "none";
        };
        window.onclick = function (event) {
                if (event.target == modal) {
                        modal.style.display = "none";
                }
        };
}

async function updateProfile() {
        const response = await fetch("/api/users/profile", {
                credentials: "include",
                method: "PUT",
                body: JSON.stringify({
                        name: document.getElementById("name-input").value,
                        bio: document.getElementById("bio-input").value,
                }),
        });

        const responseObj = await response.json();
        if (!response.ok) {
                snackbarDanger(responseObj.error);
                return;
        }

        snackbarSuccess("Successfully updated profile!");
        setTimeout(() => {
                location.reload();
        }, 2000);
}

async function main() {
        const userID = getUserID();
        const requester = await fetchRequester();

        if (userID != requester.id) {
                const modal = document.getElementById("modal");
                modal.remove();
                return;
        }

        setUpModal();
        const outerDiv = document.createElement("div");
        outerDiv.id = "holder-div";

        const editButton = document.createElement("button");
        editButton.id = "edit-button";
        editButton.textContent = "Edit Profile";
        editButton.onclick = () => {
                const modal = document.getElementById("modal");
                modal.style.display = "block";
        };

        const deleteButton = document.createElement("button");
        deleteButton.className = "danger-button";
        deleteButton.textContent = "Delete Account";
        deleteButton.onclick = async () => {
                if (
                        !window.confirm(
                                "Are you 100% sure you want to delete your account?\nThis action CANNOT be undone.",
                        )
                ) {
                        return;
                }

                const resp = await fetch("/api/users", {
                        credentials: "include",
                        method: "DELETE",
                });

                const respObj = await resp.json();
                if (resp.ok) {
                        snackbarSuccess("It's a shame to see you go :(");
                        setTimeout(() => {
                                window.location.replace("/");
                                window.location.href = "/";
                        }, 1500);
                        return;
                }
                snackbarDanger(respObj?.error);
        };

        outerDiv.appendChild(editButton);
        outerDiv.appendChild(deleteButton);
        document.body.insertBefore(outerDiv, null);
}

main();
document.getElementById("send-the-update-button").onclick = updateProfile;
