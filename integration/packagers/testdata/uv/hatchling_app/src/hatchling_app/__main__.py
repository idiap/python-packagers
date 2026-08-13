from http.server import BaseHTTPRequestHandler, HTTPServer

from . import MESSAGE


class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(MESSAGE.encode())


HTTPServer(("0.0.0.0", 8080), Handler).serve_forever()
