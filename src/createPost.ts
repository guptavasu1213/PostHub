function redirectToAdminPage(json: any) {
    let link = '/posts/' + json.editLink;
    console.log("redirect to:", link, json)
    window.location.href = link;
}

// Create post and send to server
function sendPostToServer(): void {
    let title = (<HTMLInputElement>document.querySelector('#title')).value;
    let body = (<HTMLInputElement>document.querySelector('#body')).value;
    let scope = (<HTMLInputElement>document.querySelector('input[name="scope"]:checked')).value;

    if (title.length === 0 || body.length === 0 || scope === null) {
        return
    }

    // Create JSON from the elements
    let postBody = JSON.stringify({ title: title, body: body, scope: scope });

    console.log(postBody);

    // POST the form data to the server
    fetch('/api/v1/posts', {
        method: 'POST',
        body: postBody,
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(resp => {
        if (resp.ok) {
            return resp.json()
        } else {
            alert("Error: The post could not be created!")
            console.log("post creation error:", resp.status, resp.statusText);
        }
    }).then(json =>{
        if (json !== undefined){
            alert("Post Creation Successful!")
            redirectToAdminPage(json);
        }
    }).catch(error => {
        alert("Error: The post could not be created!")
        console.log(error);
    });
}

// Attach button listeners on creation page
function attachCreationListeners() {
    let submitBtn = document.querySelector("#submit");

    submitBtn?.addEventListener("click", function (e) {
        e.preventDefault();
        sendPostToServer();
    });
}

attachCreationListeners();