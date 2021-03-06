// Report the post
function reportPost(): void {
	let reportReason = (<HTMLInputElement>document.querySelector('input[name="report"]:checked')).value;

	if (reportReason === null) {
		return;
	}

	// When other is selected, extract the user input
	if (reportReason === "Other") {
		reportReason = (<HTMLInputElement>document.querySelector('#other-report')).value;

		if (reportReason.trim().length === 0) {
			alert("Reason is required if 'Other' option is selected");
			return;
		}
	}

	// Create JSON from the elements
	let reportBody = JSON.stringify({ reason: reportReason });

	console.log(reportBody);

	// POST the form data to the server
	fetch(`/api/v1/posts/${window.location.pathname.split("/").pop()}/report`, {
		method: 'POST',
		body: reportBody,
		headers: {
			'Content-Type': 'application/json'
		}
	}).then(resp => {
		if (resp.ok) {
			alert("Post Reporting Successful!");
			redirectToViewAllPostsPage();
		} else {
			console.log("post reporting error:", resp.status, resp.statusText);
			alert("Error: The post could not be reported");
		}
	}).catch(error => {
		console.log(error);
		alert("Error: The post could not be reported");
	});
}

// Show Post Reporting Popup
function showPopup(): void {
	(<HTMLInputElement>document.querySelector(".popup")).style.display = 'flex';
}

// Hide Post Reporting Popup
function hidePopup(): void {
	(<HTMLInputElement>document.querySelector(".popup")).style.display = 'none';
}

// Attach button listeners on public view page
function attachPublicPortalListeners(): void {
	let reportBtn = document.querySelector("#report");
	let submitBtn = document.querySelector("#submit");
	let cancelBtn = document.querySelector("#cancel");

	reportBtn?.addEventListener("click", showPopup);
	submitBtn?.addEventListener("click", reportPost);
	cancelBtn?.addEventListener("click", hidePopup);
}

attachPublicPortalListeners();