function redirectToAdminPage(json: any) {
    let link = '/posts/' + json.editLink;
    console.log("redirect to:", link, json)
    window.location.href = link;
}

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
            console.log("post creation error:", resp.status, resp.statusText);
        }
    }).then(json =>{
        if (json !== undefined){
            redirectToAdminPage(json);
        }
    }).catch(error => {
        console.log(error);
    });
}

function attachFormListeners() {
    let submitBtn = document.querySelector("#submit");

    submitBtn?.addEventListener("click", function (e) {
        e.preventDefault();
        sendPostToServer();
    });
}

attachFormListeners();