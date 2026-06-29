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

        const button = document.createElement("button");
        button.id = "edit-button";
        button.textContent = "Edit Profile";
        button.onclick = () => {
                const modal = document.getElementById("modal");
                modal.style.display = "block";
        };

        outerDiv.appendChild(button);
        document.body.insertBefore(outerDiv, null);
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

main();
document.getElementById("send-the-update-button").onclick = updateProfile;
