window.addEventListener('load', function () {
    getVisitorCount();
})
function getVisitorCount() {
    var visitData = document.getElementById("visitor-count");
    var conn = new WebSocket(websocketURL);
    conn.onclose = function(evt) {
        visitData.textContent = 'Connection closed';
    }
    conn.onmessage = function(evt) {
        var jData = JSON.parse(evt.data);
        visitData.textContent = jData.visitors;
    }
}