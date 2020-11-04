"use strict";
function redirectToAdminPage(json) {
    window.location.href = `/pastes/${json.editLink}`;
}
function sendPostToServer() {
    let title = document.querySelector('#title').value;
    let body = document.querySelector('#body').value;
    let scope = document.querySelector('input[name="scope"]:checked').value;
    if (title.length === 0 || body.length === 0 || scope === null) {
        return;
    }
    let postBody = JSON.stringify({ title: title, body: body, scope: scope });
    console.log(postBody);
    fetch('/api/v1/posts', {
        method: 'POST',
        body: postBody,
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(resp => {
        if (resp.ok) {
            return resp.json();
        }
        else {
            console.log("post creation error:", resp.status, resp.statusText);
            alert("Error: The post could not be created!");
        }
    }).then(json => {
        if (json !== undefined) {
            alert("Post Creation Successful!");
            redirectToAdminPage(json);
        }
    }).catch(error => {
        console.log(error);
        alert("Error: The post could not be created!");
    });
}
function attachCreationListeners() {
    let postForm = document.querySelector("#post-submission-form");
    postForm === null || postForm === void 0 ? void 0 : postForm.addEventListener("submit", function (e) {
        e.preventDefault();
        sendPostToServer();
    });
}
attachCreationListeners();
