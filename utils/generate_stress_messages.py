import io
import os
import json
import argparse
from message import Message

def generate(N, dir, shared_secret, random_payload):
    for i in range(N):
        msg = Message(random_payload)
        with io.open(os.path.join(dir, '{}.pb'.format(i+1)), 'wb') as fd:
            fd.write(msg.body)
        with io.open(os.path.join(dir, '{}.signature'.format(i+1)), 'w') as fd:
            fd.write(msg.signature(shared_secret))

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-N', type=int, help='number of message to generate', default=100)
    parser.add_argument('-r', action='store_true', help='generate messages with multiple payload')
    parser.add_argument('CONFIG', type=argparse.FileType('r'), help="configuration file")
    parser.add_argument('DIR', help="output directory for the messages")
    args = parser.parse_args()

    config = json.load(args.CONFIG)
    shared_secret = config.get('shared_secret', 'ultrasafesecret')

    if not os.path.exists(args.DIR):
        os.mkdir(args.DIR)
    generate(args.N, args.DIR, shared_secret, random_payload=args.r)