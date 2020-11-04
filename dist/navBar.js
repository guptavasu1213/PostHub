"use strict";
function redirectToPostCreationPage() {
    window.location.href = "/";
}
function redirectToViewAllPostsPage() {
    window.location.href = "/pastes";
}
function attachNavBarListeners() {
    let createNewPostsBtn = document.querySelector("#create-new");
    let viewAllPostsBtn = document.querySelector("#view-all");
    createNewPostsBtn === null || createNewPostsBtn === void 0 ? void 0 : createNewPostsBtn.addEventListener("click", redirectToPostCreationPage);
    viewAllPostsBtn === null || viewAllPostsBtn === void 0 ? void 0 : viewAllPostsBtn.addEventListener("click", redirectToViewAllPostsPage);
}
attachNavBarListeners();
