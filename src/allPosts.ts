const limit = 10;
let totalEntries = 10000000000;/////////////////////////////////////
let offset = 0;

// Fill the table with the json data
function fillTableData(jsonArray: Array<any>): void {
    // jsonArray[0].link = "sndldsfd";
    // console.log("adsnsakdnj",jsonArray);

    // Get the template from the DOM.
    const template = (<HTMLOutputElement>document.querySelector("#posts-template")).innerHTML;

    // // Create a render function for the template with doT.template.
    const renderFn = doT.template(template);

    // // Use the render function to render the data.
    const result = renderFn(jsonArray);

    // // Insert the result into the DOM (inside the <div> with the ID log.
    (<HTMLOutputElement>document.querySelector("#table-content")).innerHTML = result;

    let currentPageNumber = Math.floor(offset / limit) + 1;
    (<HTMLOutputElement>document.querySelector("#table-page-info")).innerHTML = `Displaying ${currentPageNumber} of ${totalEntries} pages`;

    attachButtonHandlers();
}

// Get the data from the server and fill it in the table
function getDataAndFillTable(): void {
    fetch('/api/v1/posts?offset=' + offset + '&limit=' + limit, {
        method: 'GET'
    }).then(resp => {
        if (resp.ok) {
            return resp.json();
        } else {
            alert("Error: The posts cannot be retrieved");
            console.log("posts retrieval error:", resp.status, resp.statusText);
        }
    }).then(jsonArray => {
        if (jsonArray !== undefined) {
            console.log("-----------++", jsonArray);
            fillTableData(jsonArray);
        }
    }).catch(error => {
        console.log(error);
    });
}

function nextTablePage(): void {
    offset += limit;
    if (offset > totalEntries) {
        offset = totalEntries;
    } else {
        getDataAndFillTable();
    }
}

function previousTablePage(): void {
    offset -= limit;
    if (offset < 0) {
        offset = 0;
    } else {
        getDataAndFillTable();
    }
}

function attachButtonHandlers(): void {
    let nextBtn = document.querySelector("#next-page");
    let previousBtn = document.querySelector("#previous-page");

    nextBtn?.addEventListener("click", nextTablePage);
    previousBtn?.addEventListener("click", previousTablePage);
}

function redirectToPostCreationPage(): void{
    window.location.href = "/";
}

document.querySelector("#create-post")?.addEventListener("click", redirectToPostCreationPage)

getDataAndFillTable();