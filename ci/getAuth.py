import requests
import os
import argparse
import sys

DEFAULT_HOST = "X.X.X.X:8089"

def login(host, user, pwd):
    url = "http://%s/path/getAuthCode" %host
    data = {"username": (None, user), "password": (None, pwd)}
    resp = requests.post(url, files=data)

    if DEBUG:
        print(resp.text)

    if resp.ok:
        headers = resp.headers
        return headers
    else:
        return False

def Run(project, host, user ,pwd):
    headers = login(host, user, pwd)
    if not headers:
        return False
    url = "http://%s/admin/smoke" %host
    cookie = headers["Set-Cookie"].split(";")[0]
    headers = {"Cookie": cookie}
    data = {"project": (None, project)}
    resp = requests.post(url, headers=headers, files=data)
    if DEBUG:
        print(resp.text)
    if not resp.ok:
        return False

    return True

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description="Start Smoke Test in Data4Test by project")
    parser.add_argument('--host', dest="host", action="store", default=DEFAULT_HOST, help="Data4Test hostip or domain, default: %s" %DEFAULT_HOST)
    parser.add_argument('--username', dest="username", action="store", default='admin', help="user for login TestMgmt, default: admin")
    parser.add_argument('--password', dest="pwd", action="store", default='admin', help="password for username, default: admin")
    parser.add_argument('--project', dest="project", action="store", default='', help="Input project name for smoke test")
    parser.add_argument('--debug', dest="debug", action="store", default='N', help="[Y/N] Y as Yes, N as Not, default: N")
    args = parser.parse_args()

    if len(args.__dict__["project"]) == 0:
        parser.print_help()
        sys.exit(1)

    global DEBUG
    if args.debug.upper() == "Y":
        DEBUG = True
    else:
        DEBUG = False

    ret = Run(args.project, args.host, args.username, args.pwd)
    if ret:
        sys.exit(0)
    else:
        sys.exit(1)



