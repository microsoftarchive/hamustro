import uuid
import random
import hashlib
import datetime
from payload import *

class Message(object):
    def __init__(self, random_payload=True):
        self.random_payload = random_payload
        self.collection = self.set_collection()
        self.time = datetime.datetime(2016,1,1).isoformat()

    def get_payload(self):
        p = Payload()
        p.at = datetime.datetime.utcnow().isoformat()
        p.event = 'Event.{}'.format(random.randint(10000,99999))
        p.nr = random.randint(1,1000)
        p.user_id = random.randint(1,10000)
        p.is_testing = True
        return p

    def set_collection(self):
        c = Collection()
        c.device_id = str(uuid.uuid4())
        c.client_id = hashlib.md5(str(random.randint(1,1000000))).hexdigest()
        c.session = hashlib.md5(str(random.randint(1,1000000))).hexdigest()
        c.system_version = '{}.{}'.format(random.randint(1,5), random.randint(1,50))
        c.product_version = '{}.{}'.format(random.randint(1,5), random.randint(1,50))
        c.system = ['OSX','Windows','iOS','Android'][random.randint(0,3)]

        number = random.randint(1,25) \
            if self.random_payload \
            else 1

        for _ in range(number):
            c.payloads.extend([self.get_payload()])

        return c

    @property
    def body(self):
        return self.collection.SerializeToString()

    def signature(self, shared_secret):
        return hashlib.md5("{time}|{md5body}|{shared_secret}" \
            .format(time=self.time, md5body=hashlib.md5(self.body).hexdigest(), shared_secret=shared_secret)) \
            .hexdigest().decode('utf-8')
