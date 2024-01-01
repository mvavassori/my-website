document
  .getElementById("subscribe-form")
  .addEventListener("submit", function (event) {
    event.preventDefault();

    var emailInput = document.getElementById("email");

    fetch("/subscribe", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email: emailInput.value }),
    })
      .then((response) => response.text())
      .then((data) => {
        console.log(data);
        alert(data);
        // Clear the input field after successful submission
        emailInput.value = "";
      })
      .catch((error) => {
        console.error("Error:", error);
        alert("Error: Subscription failed.");
      });
  });
