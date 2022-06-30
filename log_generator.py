import sys
import os
import random 
import uuid 
import json 
import time 

from typing import Dict, List, IO 


FILE_NAME = "log.log" 
COUNT = 10000
BAD_REQUEST_FREQUENCY = 10  # So, one time in BAD_REQUEST_FREQUENCY count requests is going to be a bad request
REQUESTS_DELAY = 0.01

BAD_REQUESTS = [
        {"ip": "88.216.205.14", "method": "POST", "path": "/api/posts/add_like/", "request_uuid": "25635759-32a4-4504-a21a-f3c6e75bce5f"},
        {"ip": "23.83.244.182", "method": "POST", "path": "/api/authors/get_authors/", "request_uuid": "867be6a4-a76b-4161-9c85-99592ccacb1a"}
]

def _generate_random_ip():
    nums = []
    for i in range(4):
        nums.append(random.randint(0, 255))
    return ".".join(map(lambda a: str(a), nums))

def _generate_random_request_uuid():
    return str(uuid.uuid4())

def _generate_random_path():
    return random.choice([
        "/api/auth/login/", "/api/auth/signup/", "/api/auth/logout/", "/api/author/get_authors/", "/api/posts/edit_post/", "/api/posts/add_like/", "/api/posts/remove_like/", "/api/posts/create/", "/api/posts/delete/"  
    ])

def _generate_random_method():
    return random.choice(["POST", "GET"])


def _render_unit(request_uuid: uuid.UUID, ip: str, path: str, method: str) -> Dict[str, str]:
    return {
            'request_uuid': request_uuid,
            'ip': ip,
            'path': path,
            'method': method
    }


def _append_unit_to_file(file: IO, unit: Dict[str, str]):
    file.write(json.dumps(unit) + '\r\n')





def generate_logs(file_name: str, count: int):
    global REQUESTS_DELAY
    if not os.path.exists(file_name):
        raise FileNotFoundError(f"File {file_name} does not exist")

    file = open(file_name, 'w')
    
    now = time.time()

    for i in range(count):
        now += REQUESTS_DELAY 
        if i % BAD_REQUEST_FREQUENCY == 0:
            unit = random.choice(BAD_REQUESTS)
            unit['time'] = now
            _append_unit_to_file(file, unit)
        else:
            ip = _generate_random_ip()
            path = _generate_random_path()
            request_uuid = _generate_random_request_uuid()
            method = _generate_random_method()
        
            unit = _render_unit(ip=ip, path=path, request_uuid=request_uuid, method=method)
            unit['time'] = now 
            _append_unit_to_file(file, unit)


def main():
    global COUNT, FILE_NAME

    args = sys.argv
    if len(args) > 1:
        FILE_NAME = args[1]
    if len(args) > 2:
        COUNT = int(args[2])

    generate_logs(FILE_NAME, COUNT)


if __name__ == "__main__":
    main()


