document.addEventListener('DOMContentLoaded', function() {
    var adSeen = false;

    // Create the image div
    var imageDiv = document.createElement('div');
    imageDiv.className = 'image';

    // Create the image element
    var adImage = document.createElement('img');
    adImage.id = 'ad-image';
    adImage.alt = 'Ad Image';
    adImage.src=""

    // Append the image to the image div
    imageDiv.appendChild(adImage);

    // Append the image div to the body
    document.body.appendChild(imageDiv);

    function getAdInfo() {
        fetch('http://localhost:8081/getadinfo/1')
            .then(response => response.json())
            .then(data => {
                adImage.src = data.ImageData;
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
    }

    adImage.addEventListener('click', function() {
        fetch(window.ClicksURL, {
            method: 'GET',
            mode:'no-cors',
                headers: {
                'Content-Type': 'application/json'
            }
        })
    });

    var observer = new IntersectionObserver(function(entries) {
        entries.forEach(entry => {
            if (entry.isIntersecting && !adSeen) {
                adSeen = true;
                callAdSeenApi();
            }
        });
    }, { threshold: 0.5 });
    getAdInfo();
    observer.observe(adImage);

});
