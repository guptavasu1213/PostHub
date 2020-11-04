"use strict";
function deletePost() {
    fetch(`/api/v1/posts/${window.location.pathname.split("/").pop()}`, {
        method: 'DELETE'
    }).then(resp => {
        if (resp.ok) {
            alert("Post successfully deleted");
            redirectToViewAllPostsPage();
        }
        else {
            console.log("post dWeletion error:", resp.status, resp.statusText);
            alert("Error: The post did not get deleted");
        }
    }).catch(error => {
        console.log(error);
        alert("Error: The post did not get deleted");
    });
}
function updatePost() {
    let title = document.querySelector('#title').value;
    let body = document.querySelector('#body').value;
    let scope = document.querySelector('input[name="scope"]:checked').value;
    if (title.length === 0 || body.length === 0 || scope === null) {
        return;
    }
    let postBody = JSON.stringify({ title: title, body: body, scope: scope });
    console.log(postBody);
    fetch(`/api/v1/posts/${window.location.pathname.split("/").pop()}`, {
        method: 'PUT',
        body: postBody,
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(resp => {
        if (resp.ok) {
            alert("Post successfully updated");
            location.reload();
        }
        else {
            console.log("post updation error:", resp.status, resp.statusText);
            alert("Error: The post did not get updated");
        }
    }).catch(error => {
        console.log(error);
        alert("Error: The post did not get updated");
    });
}
function attachUpdationListeners() {
    let postForm = document.querySelector("#post-submission-form");
    let deleteBtn = document.querySelector("#delete");
    postForm === null || postForm === void 0 ? void 0 : postForm.addEventListener("submit", function (e) {
        e.preventDefault();
        updatePost();
    });
    deleteBtn === null || deleteBtn === void 0 ? void 0 : deleteBtn.addEventListener("click", deletePost);
}
attachUpdationListeners();
