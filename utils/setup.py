from __future__ import print_function

import io
import os
import sys
import json
import multiprocessing

class Setup(object):
    OPT_BOOL = {"Y": True, "n": False}
    OPT_SIGNATURE = {"Y": "optional", "n": "required"}
    OPT_FILE_FORMAT = ["csv", "json"]
    OPT_DIALECTS = {
        "abs": "Azure Blob Storage",
        "aqs": "Azure Queue Storage",
        "s3": "Amazon Simple Storage Service (S3)",
        "sns": "Amazon SNS",
        "file": "Local file",
    }

    # void
    def __init__(self):
        self.config = {}
        self.dialect = None

    @property
    def buffered(self):
        return self.dialect in ('s3', 'abs', 'file')

    # type
    def q(self, text, default="", choices=None, required=False):
        if choices:
            text += " [{}]".format("/".join(choices))
        if default:
            text += " (default: {})".format(default)
        text += ": "

        while True:
            value = raw_input(text) or default
            if not value and required:
                continue
            if not choices:
                return value
            if value in choices:
                return value

            print("Supported options: {}".format(", ".join(choices)))

    # void
    def _logfile(self):
        self.config['logfile'] = self.q("Log filepath", default="output.log", required=True) \
            if self.OPT_BOOL[self.q("Do you want to log into a file", default="n", choices=self.OPT_BOOL.keys())] \
            else ""

    # void
    def _security(self):
        self.config['signature'] = self.OPT_SIGNATURE[self.q("Hamustro will be served behind HTTPS?", default="n", choices=self.OPT_SIGNATURE.keys())]
        self.config['shared_secret'] = self.q("Please set a shared secret key", required=True)

    # void
    def _dialect(self):
        print("\nHamustro can send your events into ...")
        for k,v in self.OPT_DIALECTS.items():
            print("  [{}] {}".format(k,v))

        self.dialect = self.q("Please choose a dialect", choices=self.OPT_DIALECTS.keys())
        self.config['dialect'] = self.dialect
        self._storage = {True: self._buffered, False: self._simple}[self.buffered]
        self._dialect_options = getattr(self, self.dialect)

    # void
    def _workers(self):
        rec_worker_size = multiprocessing.cpu_count()
        if not self.buffered:
            rec_worker_size *= 4
        rec_worker_size += 1

        print("\nHamustro is using multiple workers in parallel to process multiple requests at the same time.")
        self.config['max_worker_size'] = int(self.q("How many worker do you need?", default=rec_worker_size, required=True))
        self.config['max_queue_size'] = int(self.q("Queue size", default=rec_worker_size*20, required=True))

    # void
    def _flush(self):
        if not self.buffered:
            return

        print("\nHamustro can flush events with the flush API if maintenance key is configured.")
        self.config['maintenance_key'] = self.q("Maintenance key", default="mk", required=True) \
            if self.OPT_BOOL[self.q("Do you want to use the flush API?", default="n", choices=self.OPT_BOOL.keys())] \
            else ""
        print("\nHamustro can flush periodically if automatic flush interval is configured.")
        self.config['auto_flush_interval'] = int(self.q("Automatic flush interval in minutes", default=60, required=True) \
            if self.OPT_BOOL[self.q("Do you want to setup automatic flush?", default="n", choices=self.OPT_BOOL.keys())] \
            else 0)

    # void
    def _buffered(self):
        print("\nHamustro's workers collect events in the memory to increase the performance.")
        self.config['buffer_size'] = int(self.q("Define the buffer size/worker", required=True))
        self.config['spread_buffer_size'] = self.OPT_BOOL[self.q("Do you want to randomize the buffer size near your setting to avoid flush conflicts?", choices=self.OPT_BOOL.keys(), default='Y')]

    # void
    def _simple(self):
        print("\nHamustro send the incoming messages immediately to your selected target.")
        self.config['retry_attempt'] = int(self.q("When the saving has failed, how many times do you want to try again before we remove the event?", default=3, required=True))

    # void
    def s3(self):
        return {
            'access_key_id': self.q("Access Key ID", required=True),
            'secret_access_key': self.q("Secret Access Key", required=True),
            'region': self.q("Region", required=True),
            'bucket': self.q("Bucket", required=True),
            'blob_path': self.q("Blob path", required=True, default="{date}/"),
            'file_format': self.q("File output format", choices=self.OPT_FILE_FORMAT, default="json"),
            'endpoint': self.q("Endpoint", required=True),
        }

    # voifd
    def abs(self):
        return {
            'account': self.q("Account", required=True),
            'access_key': self.q("Access Key", required=True),
            'container': self.q("Container", required=True),
            'blob_path': self.q("Blob path", required=True, default="{date}/"),
            'file_format': self.q("File output format", choices=self.OPT_FILE_FORMAT, default="csv"),
        }

    # void
    def sns(self):
        return {
            'access_key_id': self.q("Access Key ID", required=True),
            'secret_access_key': self.q("Secret Access Key", required=True),
            'region': self.q("Region", required=True),
            'topic_arn': self.q("Topic ARN", required=True),
        }

    # void
    def aqs(self):
        return {
            'account': self.q("Account", required=True),
            'access_key': self.q("Access Key", required=True),
            'queue_name': self.q("Queue Name", required=True),
        }

    # void
    def file(self):
        return {
            'file_path': self.q("File path", required=True, default="{date}/"),
            'file_format': self.q("File output format", choices=self.OPT_FILE_FORMAT, default="csv"),
            'compress': self.q("Do you want to compress the output files?", default="n", choices=self.OPT_BOOL.keys())
        }

    # void
    def _options(self):
        print("\nPlease define the expected behavior of the collector")
        self.config['masked_ip'] = self.OPT_BOOL[self.q("Do you want to remove the last octet of incoming IP addresses?", choices=self.OPT_BOOL.keys(), default="n")]

    # void
    def run(self):
        self.path = self.q("Configuration path", default="config/config.json", required=True)
        if os.path.exists(self.path) and not self.OPT_BOOL[self.q("Configuration exists, do you want to overwrite?", choices=self.OPT_BOOL.keys(), required=True)]:
            return False

        self._logfile()
        self._security()
        self._dialect()
        self._workers()
        self._storage()
        self._flush()

        print("\nPlease set the credentials for the selected ({}) dialect:".format(self.dialect))
        self.config[self.dialect] = self._dialect_options()

        self._options()

        print("\nYour configuration file was created successfully!")
        return True

    # void
    def save(self):
        with io.open(self.path, "w", encoding="utf-8") as fd:
            fd.write(json.dumps(self.config, indent=2, sort_keys=True).decode('utf-8'))

if __name__ == '__main__':
    s = Setup()
    if not s.run():
        sys.exit(1)
    s.save()
