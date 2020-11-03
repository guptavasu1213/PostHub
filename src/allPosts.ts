// CMPT 315 (Fall 2020)
// Assignment 2
// Author: Vasu Gupta

const limit = 10;
let offset = 0;

// Fill the table with the json data
function fillTableData(jsonArray: Array<any>): void {
	// Get the template from the DOM.
	const template = (<HTMLOutputElement>document.querySelector("#posts-template")).innerHTML;

	// // Create a render function for the template with doT.template.
	const renderFn = doT.template(template);

	// // Use the render function to render the data.
	const result = renderFn(jsonArray);

	// // Insert the result into the DOM (inside the <div> with the ID log.
	(<HTMLOutputElement>document.querySelector("#table-content")).innerHTML = result;

	let currentPageNumber = Math.floor(offset / limit) + 1;
	(<HTMLOutputElement>document.querySelector("#table-page-info")).innerHTML = `Displaying Page ${currentPageNumber}`;

	attachButtonHandlers();
}

// Convert Unix time stamp from the database to human legible date and time
function convertUnixTimestampToDate(jsonArray: Array<any>): Array<any> {
	console.log(jsonArray);
	for (let i in jsonArray) {
		jsonArray[i].epoch = new Date(jsonArray[i].epoch * 1000).toLocaleTimeString("en-US", { month: "long", day: "numeric", year: "numeric" });
	}
	return jsonArray;
}

// Get the data from the server and fill it in the table
function getDataAndFillTable(): void {
	fetch('/api/v1/posts?offset=' + offset + '&limit=' + limit, { // Make request
		method: 'GET'
	}).then(resp => { // Check response
		if (resp.ok) {
			return resp.json();
		} else {
			console.log("posts retrieval error:", resp.status, resp.statusText);
			alert("Error: The posts cannot be retrieved");
		}
	}).then(jsonArray => { // Check JSON
		if (jsonArray !== undefined) {
			console.log(jsonArray);

			if (jsonArray.length === 0) {
				offset -= limit;
				alert("This is the last page. Cannot go forward.");
			} else {
				jsonArray = convertUnixTimestampToDate(jsonArray);
				fillTableData(jsonArray);
			}
		}
	}).catch(error => {
		console.log(error);
	});
}

function nextTablePage(): void {
	offset += limit;
	getDataAndFillTable();
}

function previousTablePage(): void {
	offset -= limit;

	if (offset < 0) {
		offset = 0;
		alert("This is the first page. Cannot go back.");
	} else {
		getDataAndFillTable();
	}
}

// Attach handlers for buttons underneath the table
function attachButtonHandlers(): void {
	let nextBtn = document.querySelector("#next-page");
	let previousBtn = document.querySelector("#previous-page");

	nextBtn?.addEventListener("click", nextTablePage);
	previousBtn?.addEventListener("click", previousTablePage);
}

document.querySelector("#create-new")?.addEventListener("click", redirectToPostCreationPage);

getDataAndFillTable();