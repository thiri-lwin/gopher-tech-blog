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
                    <div class="comment d-flex align-items-center">
                        <span class="mb-0"><strong>${data.user_name}</strong>: ${data.content}</span>
                        <button class="btn btn-link text-danger ms-2 p-0 delete-comment" data-comment-id="${data.id}">
                            <i class="fas fa-trash-alt"></i>
                        </button>
                    </div>
                `;
                
                document.getElementById("comment-form").reset();
            });
        }
    });

    document.addEventListener("click", function (event) {
        const deleteButton = event.target.closest(".delete-comment"); // Find closest delete button
        if (deleteButton) {
            const commentID = deleteButton.dataset.commentId; 
    
            if (!commentID) return;
    
            // Confirm before deletion
            if (!confirm("Are you sure you want to delete this comment?")) {
                return;
            }
    
            fetch(`/comments/${commentID}`, {
                method: "DELETE",
                headers: { "Content-Type": "application/json" },
            })
            .then(response => {
                if (response.ok) {
                    deleteButton.closest(".comment").remove(); 
                } else {
                    return response.json(); 
                }
            })
            .catch(error => console.error("Error deleting comment:", error));
        }
    });
    
});

//Fetch likes and comments on page load
function fetchPostData() {
    fetch(`/posts/${postID}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById("like-count").textContent = data.likes;
            const commentsContainer = document.getElementById("comments");
            commentsContainer.innerHTML = '';
            data.comments.forEach(comment => {
                commentsContainer.innerHTML += `
                    <div class="comment d-flex align-items-center">
                        <span class="mb-0"><strong>${comment.author}</strong>: ${comment.content}</span>
                        {{ if eq $.AuthUserID .UserID }}
                            <button class="btn btn-link text-danger ms-2 p-0 delete-comment" data-comment-id="{{ .ID }}">
                                <i class="fas fa-trash-alt"></i>
                            </button>
                        {{ end }}
                    </div>
                `;
            });
        });
}

