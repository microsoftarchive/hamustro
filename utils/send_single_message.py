from __future__ import print_function
import os
import json
import datetime
import argparse
import requests
from message import Message

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('CONFIG', type=argparse.FileType('r'), help="configuration file")
    parser.add_argument('URL', help="tavis url")
    args = parser.parse_args()

    config = json.load(args.CONFIG)
    shared_secret = config.get('shared_secret', 'ultrasafesecret')

    msg = Message(random_payload=False)
    msg.time = datetime.datetime.utcnow().isoformat()
    resp = requests.post(args.URL, headers={
        'X-Hamustro-Time': msg.time,
        'X-Hamustro-Signature': msg.signature(shared_secret)
    }, data=msg.body)
    print('Response code: {}'.format(resp.status_code))