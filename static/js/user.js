document.addEventListener("DOMContentLoaded", function () {
    // Handle SignIn Button
    const signinForm = document.getElementById("signin-form");
    if (signinForm) {
        signinForm.addEventListener("submit", function (event) {
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
    }

    // Handle SignUp Button
    const signupForm = document.getElementById("signup-form");
    if (signupForm) {
        signupForm.addEventListener("submit", function (event) {
            event.preventDefault();
            const firstName = document.getElementById("first_name").value;
            const lastName = document.getElementById("last_name").value;
            const email = document.getElementById("email").value;
            const password = document.getElementById("password").value;

            fetch(`/signup`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ first_name: firstName, last_name: lastName, email: email, password: password })
            })
            .then(response => {
                if (response.status === 200) {
                    window.location.href = "/signin"; // Redirect to the signin page
                } else {
                    return response.json().then(data => {
                        alert("Sign-up failed: " + data.error);
                    });
                }
            })
            .catch(error => {
                console.error("Error:", error);
                alert("An error occurred. Please try again.");
            });
        });
    }
});
