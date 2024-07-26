<script>
    document.addEventListener('DOMContentLoaded', function() {
    //var adLink = document.getElementById('ad-link');
    var adImage = document.getElementById('ad-image');
    var adSeen = false;

    function getAdInfo() {
    fetch('http://localhost:8080/getadinfo/2')
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
    method: 'POST',
    headers: {
    'Content-Type': 'application/json'
},
    body: JSON.stringify({ ad: 'your-ad-info' })
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
</script>