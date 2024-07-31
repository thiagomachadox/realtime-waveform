let socket;

function startWaveform() {
    const youtubeUrl = document.getElementById('youtubeUrl').value;
    if (!youtubeUrl) {
        alert('Please enter a YouTube URL');
        return;
    }

    if (socket) {
        socket.close();
    }

    socket = new WebSocket('ws://localhost:8080/ws');

    socket.onopen = function(e) {
        console.log('WebSocket connection established');
        socket.send(youtubeUrl);
    };

    socket.onmessage = function(event) {
        document.getElementById('waveform').textContent = event.data;
    };

    socket.onclose = function(event) {
        if (event.wasClean) {
            console.log(`WebSocket connection closed cleanly, code=${event.code}, reason=${event.reason}`);
        } else {
            console.log('WebSocket connection died');
        }
    };

    socket.onerror = function(error) {
        console.log(`WebSocket error: ${error.message}`);
    };
}