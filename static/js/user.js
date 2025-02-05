document.addEventListener("DOMContentLoaded", function () {
    // Handle SignIn Button
    document.getElementById("signin-form").addEventListener("submit", function (event) {
            event.preventDefault();
            const email = document.getElementById("email").value;
            const password = document.getElementById("password").value;
    
            fetch(`/signin`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: email, password: password })
            })
            .then(response => {
                if (response.status === 200) {
                    window.location.href = "/"; // Redirect to the index page
                } else {
                    return response.json().then(data => {
                        alert("Sign-in failed: " + data.error);
                    });
                }
            })
            .catch(error => {
                console.error("Error:", error);
                alert("An error occurred. Please try again.");
            });
        });
});

