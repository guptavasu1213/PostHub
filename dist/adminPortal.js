"use strict";
function deletePost() {
    fetch('/api/v1' + window.location.pathname, {
        method: 'DELETE'
    }).then(resp => {
        if (resp.ok) {
            alert("Post successfully deleted");
            window.location.href = "/";
        }
        else {
            alert("Error: The post did not get deleted");
            console.log("post deletion error:", resp.status, resp.statusText);
        }
    }).catch(error => {
        console.log(error);
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
    fetch('/api/v1' + window.location.pathname, {
        method: 'POST',
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
            alert("Error: The post did not get updated");
            console.log("post updation error:", resp.status, resp.statusText);
        }
    }).catch(error => {
        console.log(error);
    });
}
function attachUpdationListeners() {
    let submitBtn = document.querySelector("#update");
    let deleteBtn = document.querySelector("#delete");
    submitBtn === null || submitBtn === void 0 ? void 0 : submitBtn.addEventListener("click", function (e) {
        e.preventDefault();
        updatePost();
    });
    deleteBtn === null || deleteBtn === void 0 ? void 0 : deleteBtn.addEventListener("click", deletePost);
}
attachUpdationListeners();
