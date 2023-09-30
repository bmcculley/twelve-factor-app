# The Visitor WebSocket Image

This houses the WebSocket server that will poll the Redis table and push out the updated visitor count to the HTML page.

# Running

Set an environment varaible to pass configuration data to the image.

 - addr: this is where the WebSocket server should listen
 - origins: this is where the connections will originate from. If they are different than it's own host it will allow the server to respond to the requests.

```bash
export WS_CONFIG='{"origins": ["http://127.0.0.1:3000", "http://localhost:3000", "http://example.com:3000"],"addr":":8080"}' 
```