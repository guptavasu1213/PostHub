"use strict";
var _a;
const limit = 10;
let offset = 0;
function fillTableData(jsonArray) {
    const template = document.querySelector("#posts-template").innerHTML;
    const renderFn = doT.template(template);
    const result = renderFn(jsonArray);
    document.querySelector("#table-content").innerHTML = result;
    let currentPageNumber = Math.floor(offset / limit) + 1;
    document.querySelector("#table-page-info").innerHTML = `Displaying Page ${currentPageNumber}`;
    attachButtonHandlers();
}
function convertUnixTimestampToDate(jsonArray) {
    console.log(jsonArray);
    for (let i in jsonArray) {
        jsonArray[i].epoch = new Date(jsonArray[i].epoch * 1000).toLocaleTimeString("en-US", { month: "long", day: "numeric", year: "numeric" });
    }
    return jsonArray;
}
function getDataAndFillTable() {
    fetch('/api/v1/posts?offset=' + offset + '&limit=' + limit, {
        method: 'GET'
    }).then(resp => {
        if (resp.ok) {
            return resp.json();
        }
        else {
            console.log("posts retrieval error:", resp.status, resp.statusText);
            alert("Error: The posts cannot be retrieved");
        }
    }).then(jsonArray => {
        if (jsonArray !== undefined) {
            console.log(jsonArray);
            if (jsonArray.length === 0) {
                offset -= limit;
                alert("This is the last page. Cannot go forward.");
            }
            else {
                jsonArray = convertUnixTimestampToDate(jsonArray);
                fillTableData(jsonArray);
            }
        }
    }).catch(error => {
        console.log(error);
    });
}
function nextTablePage() {
    offset += limit;
    getDataAndFillTable();
}
function previousTablePage() {
    offset -= limit;
    if (offset < 0) {
        offset = 0;
        alert("This is the first page. Cannot go back.");
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
(_a = document.querySelector("#create-new")) === null || _a === void 0 ? void 0 : _a.addEventListener("click", redirectToPostCreationPage);
getDataAndFillTable();
