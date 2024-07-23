<script>
    document.addEventListener('DOMContentLoaded', function() {
        var adLink = document.getElementById('ad-link');
        var adImage = document.getElementById('ad-image');
        var adSeen = false;

        function getAdInfo() {
            fetch('https://your-api-url.com/get-ad-info')
                .then(response => response.json())
                .then(data => {
                    adLink.href = data.url;
                    adImage.src = data.image;
                })
                .catch(error => console.error('Error:', error));
        }

        function callAdSeenApi() {
            fetch('https://your-api-url.com/ad-seen', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ ad: 'your-ad-info' })
            })
            .then(response => response.json())
            .then(data => console.log('Ad seen API response:', data))
            .catch(error => console.error('Error:', error));
        }

        adLink.addEventListener('click', function() {
            fetch('https://your-api-url.com/ad-click', {
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
