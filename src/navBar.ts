// Redirect to Post Creation
function redirectToPostCreationPage(): void {
    window.location.href = "/";
}

// Redirect to public view displaying all posts
function redirectToViewAllPostsPage(): void {
    window.location.href = "/posts";
}

// Attach listeners for the navigation bar
function attachNavBarListeners(): void {
    let createNewPostsBtn = document.querySelector("#create-new");
    let viewAllPostsBtn = document.querySelector("#view-all");

    createNewPostsBtn?.addEventListener("click", redirectToPostCreationPage);
    viewAllPostsBtn?.addEventListener("click", redirectToViewAllPostsPage);
}

attachNavBarListeners();