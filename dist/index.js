"use strict";
var _a;
const limit = 10;
let totalEntries = 10000000000;
let offset = 0;
function fillTableData(jsonArray) {
    jsonArray[0].link = "sndldsfd";
    console.log("adsnsakdnj", jsonArray);
    const template = document.querySelector("#posts-template").innerHTML;
    const renderFn = doT.template(template);
    const result = renderFn(jsonArray);
    document.querySelector("#table-content").innerHTML = result;
    let currentPageNumber = Math.floor(offset / limit) + 1;
    document.querySelector("#table-page-info").innerHTML = `Displaying ${currentPageNumber} of ${totalEntries} pages`;
    attachButtonHandlers();
}
function getDataAndFillTable() {
    fetch('/api/v1/posts?offset=' + offset + '&limit=' + limit, {
        method: 'GET'
    }).then(resp => {
        if (resp.ok) {
            return resp.json();
        }
        else {
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
function nextTablePage() {
    offset += limit;
    if (offset > totalEntries) {
        offset = totalEntries;
    }
    else {
        getDataAndFillTable();
    }
}
function previousTablePage() {
    offset -= limit;
    if (offset < 0) {
        offset = 0;
    }
    else {
        getDataAndFillTable();
    }
}
function attachButtonHandlers() {
    let nextBtn = document.querySelector("#next-page");
    let previousBtn = document.querySelector("#previous-page");
    nextBtn === null || nextBtn === void 0 ? void 0 : nextBtn.addEventListener("click", nextTablePage);
    previousBtn === null || previousBtn === void 0 ? void 0 : previousBtn.addEventListener("click", previousTablePage);
}
function redirectToPostCreationPage() {
    window.location.href = "/createPost.html";
}
(_a = document.querySelector("#create-post")) === null || _a === void 0 ? void 0 : _a.addEventListener("click", redirectToPostCreationPage);
getDataAndFillTable();
