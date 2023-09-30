from flask import (
        Flask,
        jsonify,
        render_template,
        request
)
from werkzeug.datastructures import Headers
from redis import Redis
import logging
import os

app = Flask(__name__)

gunicorn_error_logger = logging.getLogger('gunicorn.error')
app.logger.handlers.extend(gunicorn_error_logger.handlers)
app.logger.setLevel(logging.INFO)

# get environment variables
redis_url = os.getenv("REDIS_URL")
websocket_url = os.getenv("WEBSOCKET_URL")

redisDb = Redis(host="redis", db=0, socket_timeout=5)

@app.route('/')
def home():
    ip_addr = request.environ.get('HTTP_X_FORWARDED_FOR', request.remote_addr)
    app.logger.info('Received visit from %s', ip_addr)
    redisDb.incr('visitorCount')
    return render_template('index.html', websocket_url=websocket_url)

if __name__ == '__main__':
    app.run(host="0.0.0.0", debug=True)
