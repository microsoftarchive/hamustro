import uuid
import json
import time
import base64
import random
import hashlib
import datetime
from payload import *
from protobuf_to_dict import protobuf_to_dict as to_dict

class Message(object):
    def __init__(self, random_payload=True):
        self.random_payload = random_payload
        self.collection = self.set_collection()
        self.time = 1454514088

    def get_payload(self):
        p = Payload()
        p.at = int(time.time())
        p.event = 'Event.{}'.format(random.randint(10000,99999))
        p.nr = random.randint(1,1000)
        p.user_id = 'Event.{}'.format(random.randint(1,99000000))
        p.ip = '{}.{}.{}.{}'.format(random.randint(1,255), random.randint(1,255), random.randint(1,255), random.randint(1,255))
        return p

    def set_collection(self):
        c = Collection()
        c.env = 4
        c.device_id = hashlib.sha256(str(uuid.uuid4())).hexdigest()
        c.client_id = hashlib.md5(str(random.randint(1,1000000))).hexdigest()[:20]
        c.system_version = '{}.{}'.format(random.randint(1,5), random.randint(1,50))
        c.product_version = '{}.{}'.format(random.randint(1,5), random.randint(1,50))
        c.session = hashlib.md5("{device_id}:{client_id}:{system_version}:{product_version}".format(
            device_id=c.device_id, client_id=c.client_id, system_version=c.system_version,
            product_version=c.product_version)).hexdigest()
        c.system = ['OSX','Windows','iOS','Android'][random.randint(0,3)]
        number = random.randint(1,25) \
            if self.random_payload \
            else 1

        for _ in range(number):
            c.payloads.extend([self.get_payload()])

        return c

    @property
    def pb(self):
        return self.collection.SerializeToString()
    
    @property
    def json(self):
        return json.dumps(to_dict(self.collection)).decode('utf-8')
    
    def get_body(self, output_format = 'protobuf'):
        return {'protobuf': self.pb, 'json': self.json}[output_format]

    def signature(self, shared_secret, output_format = 'protobuf'):
        md5hex = hashlib.md5(self.get_body(output_format)).hexdigest()
        bytehash = hashlib.sha256("{time}|{md5hex}|{shared_secret}" \
            .format(time=self.time, md5hex=md5hex, shared_secret=shared_secret)) \
            .digest()
        return base64.b64encode(bytehash).decode('utf-8')
