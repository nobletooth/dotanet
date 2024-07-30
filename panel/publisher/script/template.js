document.addEventListener('DOMContentLoaded', function() {
    var adSeen = false;
    const publisherId = "__PUBLISHER_ID__";
    const adserverurl = "__ADSERVER_URL__";
    let impressionTime;

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
            const response = await fetch(`${adserverurl}/getadinfo/${publisherId}`);
            const data = await response.json();
            adImage.src = data.ImageData;
            window.ImpressionURL = `${data.ImpressionURL}/${data.ImpressionId}/
            ${encodeURIComponent(new Date().toISOString())}`;
            window.ClickURL = `${data.ClickURL}/${data.ClickId}/
            ${data.ImpressionId}/${encodeURIComponent(new Date().toISOString())}`;
        } catch (error) {
            console.error('Error:', error);
        }
    }

    function callAdSeenApi() {
        impressionTime = new Date()
        fetch(window.ImpressionURL, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => response.json())
        .catch(error => console.error('Error:', error));
    }

    adImage.addEventListener('click', function(event) {
        event.preventDefault();
    
        const clickTime = new Date();
        const timeDiff = (clickTime - impressionTime) / 1000; // difference in seconds
    
        if (timeDiff > 10 && timeDiff < 30) {
            fetch(window.ClickURL, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => response.json())
            .then(data => {
                window.open(data.AdURL, '_blank');
            })
            .catch(error => console.error('Error:', error));
        } else {
            console.log('Click time is not in the valid range');
        }
    });

    var observer = new IntersectionObserver(function(entries) {
        entries.forEach(entry => {
            if (entry.isIntersecting && !adSeen) {
                adSeen = true;
                callAdSeenApi();
            }
        });
    }, { threshold: 0.5 });

    getAdInfo().then(() => {
        observer.observe(adImage);
    });
});
