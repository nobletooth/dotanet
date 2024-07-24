document.addEventListener('DOMContentLoaded', function() {
    var adSeen = false;

    // Create the image div
    var imageDiv = document.createElement('div');
    imageDiv.className = 'image';

    // Create the image element
    var adImage = document.createElement('img');
    adImage.id = 'ad-image';
    adImage.alt = 'Ad Image';

    // Append the image to the image div
    imageDiv.appendChild(adImage);

    // Append the image div to the body
    document.body.appendChild(imageDiv);

    function getAdInfo() {
        fetch('http://localhost:8080/getAd/3')
            .then(response => response.json())
            .then(data => {
                adImage.src = data.image;
                window.ImpressionsURL = data.ImpressionsURL;
                window.ClicksURL = data.ClicksURL;
            })
            .catch(error => console.error('Error:', error));
    }

    function callAdSeenApi() {
        fetch(window.ImpressionsURL, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(response => response.json())
            .then(data => console.log('Ad seen API response:', data))
            .catch(error => console.error('Error:', error));
    }

    adImage.addEventListener('click', function() {
        fetch(window.ClicksURL, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(response => response.json())
            .then(data => console.log('Ad click API response:', data))
            .catch(error => console.error('Error:', error));
    });

    var observer = new IntersectionObserver(function(entries) {
        entries.forEach(entry => {
            if (entry.isIntersecting && !adSeen) {
                adSeen = true;
                callAdSeenApi();
            }
        });
    }, { threshold: 0.5 });

    observer.observe(adImage);

    getAdInfo();
});
