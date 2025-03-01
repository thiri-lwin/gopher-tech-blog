document.addEventListener("DOMContentLoaded", function () {
    fetchPostData();

    // Handle Like Button
    document.getElementById("like-button").addEventListener("click", function () {
        if (!isAuthenticated) {
            window.location.href = "/signin";
            console.log("user not logged in");
        } else {
            fetch(`/posts/${postID}/like-toggle`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
            })
            .then(response => response.json())
            .then(data => {
                document.getElementById("like-count").textContent = data.likes;
                document.getElementById("like-toggle").textContent = data.liked ? "Unlike" : "Like";
            });
        }
        
    });

    // Handle Comment Form Submission
    document.getElementById("comment-form").addEventListener("submit", function (event) {
        if (!isAuthenticated) {
            window.location.href = "/signin";
            console.log("user not logged in");
        } else {
            event.preventDefault();
            const content = document.getElementById("comment-content").value;
    
            fetch(`/posts/${postID}/comment`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ content })
            })
            .then(response => response.json())
            .then(data => {
                document.getElementById("comments").innerHTML += `
                    <div class="comment">
                        <p><strong>${data.user_name}</strong>: ${data.content}</p>
                    </div>
                `;
                document.getElementById("comment-form").reset();
            });
        }
    });
});

// Fetch likes and comments on page load
function fetchPostData() {
    fetch(`/posts/${postID}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById("like-count").textContent = data.likes;
            const commentsContainer = document.getElementById("comments");
            commentsContainer.innerHTML = '';
            data.comments.forEach(comment => {
                commentsContainer.innerHTML += `
                    <div class="comment">
                        <p><strong>${comment.author}</strong>: ${comment.content}</p>
                    </div>
                `;
            });
        });
}
