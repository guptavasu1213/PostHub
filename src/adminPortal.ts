// Delete the post using the ID
function deletePost(): void {
	fetch(`/api/v1/posts/${window.location.pathname.split("/").pop()}`, {
		method: 'DELETE'
	}).then(resp => {
		if (resp.ok) {
			alert("Post successfully deleted");
			redirectToViewAllPostsPage();
		} else {
			console.log("post dWeletion error:", resp.status, resp.statusText);
			alert("Error: The post did not get deleted");
		}
	}).catch(error => {
		console.log(error);
		alert("Error: The post did not get deleted");
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
		} else {
			console.log("post updation error:", resp.status, resp.statusText);
			alert("Error: The post did not get updated");
		}
	}).catch(error => {
		console.log(error);
		alert("Error: The post did not get updated");
	});
}

// Attach button listeners on administrator page
function attachUpdationListeners(): void {
	let postForm = document.querySelector("#post-submission-form");
	let deleteBtn = document.querySelector("#delete");

	postForm?.addEventListener("submit", function (e) {
		e.preventDefault();
		updatePost();
	});

	deleteBtn?.addEventListener("click", deletePost);
}

attachUpdationListeners();