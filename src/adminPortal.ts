// Delete the post using the ID
function deletePost(): void {
    fetch('/api/v1' + window.location.pathname, {
        method: 'DELETE'
    }).then(resp => {
        if (resp.ok) {
            alert("Post successfully deleted");
            window.location.href = "/posts";
        } else {
            alert("Error: The post did not get deleted");
            console.log("post deletion error:", resp.status, resp.statusText);
        }
    }).catch(error => {
        console.log(error);
    });
}

// Update the post based on the form values
function updatePost(): void {
    let title = (<HTMLInputElement>document.querySelector('#title')).value;
    let body = (<HTMLInputElement>document.querySelector('#body')).value;
    let scope = (<HTMLInputElement>document.querySelector('input[name="scope"]:checked')).value;

    if (title.length === 0 || body.length === 0 || scope === null) {
        return;
    }

    // Create JSON from the elements
    let postBody = JSON.stringify({ title: title, body: body, scope: scope });

    console.log(postBody);

    // POST the form data to the server
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
        } else {
            alert("Error: The post did not get updated");
            console.log("post updation error:", resp.status, resp.statusText);
        }
    }).catch(error => {
        console.log(error);
    });
}

// Attach button listeners on administrator page
function attachUpdationListeners(): void {
    let postForm = document.querySelector("#post-submission-form");
    let deleteBtn = document.querySelector("#delete");
    console.log("s")
    postForm?.addEventListener("submit", function (e) {
        e.preventDefault();
        updatePost();
    });

    deleteBtn?.addEventListener("click", deletePost);
}

attachUpdationListeners();