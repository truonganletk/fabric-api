document.addEventListener('DOMContentLoaded', function() {
    // Get the form element
    const form = document.getElementById('upload-form');

    // Add a submit event listener
    form.addEventListener('submit', function(event) {
        event.preventDefault(); // Prevent default form submission

        // Create a FormData object from the form
        const formData = new FormData(form);

        // Send a POST request to the upload endpoint
        fetch('http://localhost:8080/erc20/getBalance', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (response.ok) {
                return response.text(); // Process the response text
            } else {
                throw new Error('File upload failed');
            }
        })
        // .then(response => response.json())
        .then(data => {
            // Select the container where data will be displayed
            const container = document.getElementById('data-container');

            // Clear previous content
            container.innerHTML = '';
            const div = document.createElement('div');
            div.textContent = `${data.toString()}`;
            container.appendChild(div);
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
});
