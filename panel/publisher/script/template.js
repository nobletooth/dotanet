document.addEventListener('DOMContentLoaded', function() {
    var adSeen = false;
    const publisherId = "__PUBLISHER_ID__";
    const adserverurl = "__ADSERVER_URL__";

    // Create the image div
    var imageDiv = document.createElement('div');
    imageDiv.className = 'image';

    // Create the image element
    var adImage = document.createElement('img');
    adImage.id = 'ad-image';
    adImage.alt = 'Ad Image';
    adImage.src = "";

    // Append the image to the image div
    imageDiv.appendChild(adImage);

    // Append the image div to the body
    document.body.appendChild(imageDiv);

    async function getAdInfo() {
        try {
            const response = await fetch(`${adserverurl}/getadinfo/${publisherId}`, {
                credentials: 'include'
            });
            const data = await response.json();
            adImage.src = data.ImageData;
            window.ImpressionsURL = `${data.ImpressionsURL}`;
            window.ClicksURL = `${data.ClicksURL}`;
        } catch (error) {
            console.error('Error:', error);
        }
    }

    function callAdSeenApi() {
        fetch(`${window.ImpressionsURL}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include'
        }).then(response => response.json())
        .catch(error => console.error('Error:', error));
    }

    adImage.addEventListener('click', function(event) {
        event.preventDefault();
        fetch(`${window.ClicksURL}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include'
        })
        .then(response => response.json())
        .then(data => {
            window.open(data.AdURL, '_blank');
        })
        .catch(error => console.error('Error:', error));
    });

    var observer = new IntersectionObserver(function(entries) {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                callAdSeenApi();
            }
        });
    }, { threshold: 0.5 });

    getAdInfo().then(() => {
        observer.observe(adImage);
    });
});
