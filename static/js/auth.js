function handleCredentialResponse(response) {
    const idToken = response.credential;

    // Send the ID token to the backend for verification
    fetch("/auth/google-signin", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token: idToken })
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
}
