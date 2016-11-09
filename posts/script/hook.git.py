#!/usr/bin/env python
# coding=utf-8

# start:server.py --port=81 --log-file-prefix=/data/www/mgr/log/tornado.81.log

import os
import tornado.web
import tornado.ioloop
import tornado.options
import tornado.httpserver
from tornado.options import define, options

class DefaultHandler(tornado.web.RequestHandler) :
    def get(self):
        self.post()

    def post(self):
        try :
            project = self.get_argument("project")
            path = "/data/www/%s" % (project,)
            if not os.path.exists(path) :
                return self.write("Error");
            if project == 'asm_api' :
                os.system("cd %s && git checkout 4.0.0 && git reset --hard HEAD && git pull" % (path,))
            else :
                os.system("cd %s && git checkout dev && git reset --hard HEAD && git pull" % (path,))
        except :
            return self.write("Hello World")

router = [
    (r'/', DefaultHandler),
]

define("port", default = 8011, help = "run on the given port", type = int)

setting = {
    'debug' : True
}

def main():
    tornado.options.parse_command_line()
    http_server = tornado.httpserver.HTTPServer(tornado.web.Application(
        handlers = router,
        **setting
    ))
    http_server.listen(options.port)

    print("Development server is running at http://127.0.0.1:%s" % options.port)
    print("Quit the server with Control-C")

    tornado.ioloop.IOLoop.instance().start()

if __name__ == "__main__":
    main()
