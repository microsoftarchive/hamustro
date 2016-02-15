# Pseudo Client for Tracking

These are the functions you have to implement to write a compatible client.

## Initialization

Create a `ClientTracker` object on application open.

```cpp
t = ClientTracker(
  collector_url string, // required, set from config
  shared_secret_key string, // required, set from config
  device_id string, // required, please sha256hex it.
  client_id string, // required
  system_version string, // required
  product_version string, // required
  system string,
  product_git_hash string,
  queue_size int, // set from config, default: 20
  queue_retention int // set from config, default: 1440 (minutes = 24 hours)
)
```

It will generate pre-populated information for new events so it should not be calculated on adding each event.

```cpp
// Generated as md5hex(sha256hex(device_id) + ":" + client_id + ":" + system_version + ":" + product_version)
t.GenerateSession()
```

Loading information from persistent storage.

```cpp
tc.LoadCollections() // not sent events
tc.LoadNumberPerSession() // events per session
tc.LoadLastSyncTime() // timestamp for events sent last time
```

## Track events

To track events you should call the following:

```cpp
t.TrackEvent(
  event string, // required
  user_id int,
  params string,
  is_testing bool // default: false
)
```

Please set automatically 
- the `at` uint64 attribute for the events - it must contain an EPOCH UTC timestamp, 
- the `nr` integer attribute - it must contain the serial number of this event within the session all time. So, it starts counting from the first event in the session and it never defaults for that session, not even new application open,
- the `ip` string attribute for actual IP address.

It'll queue up the events within the `ClientTracker`. You should store this in a persistent storage. Please make sure to save the `ClientTracker`'s attributes with the event because you will need to send to the `collector_url` by session.

## Send events to the collector

For sending messages we're using [Protobuf](https://developers.google.com/protocol-buffers/?hl=en). This is quicker to handle and smaller than JSON.

This is the message [format](../proto/payload.proto) we're using:

```protobuf
message Payload {
  required uint64 at = 1;
  required string event = 2;
  required uint32 nr = 3;
  optional uint32 user_id = 4;
  optional string ip = 5;
  optional string parameters = 6;
  optional bool is_testing = 7;
}

message Collection {
  required string device_id = 1;
  required string client_id = 2;
  required string session = 3;
  required string system_version = 4;
  required string product_version = 5;
  optional string system = 6;
  optional string product_git_hash = 7;
  repeated Payload payloads = 8;
}
```

You can send multiple `Payloads`, but as you can see you have to send these by session. Normally you won't have multiple sessions in your code but it can happen with bad connection and updates happening. Please sign them as 'Sent', but don't delete them before getting a response of '200'.

The collector only accepts `POST` with a body of a valid Protobuf bytestream.

Please wait for `200` response code before you delete the already sent payloads. After that update the last sync time. Do not remove the events receiving a different error code.

Please define the following headers for sending:

```
X-Hamustro-Time: EPOCH UTC timestamp
X-Hamustro-Signature: base64(sha256(X-Time + "|" + md5hex(request.body) + "|" + t.shared_secret_key))
```

You can check out the proper signature generation in [Python](https://github.com/sub-ninja/hamustro/blob/master/utils/message.py#L46-L51).

Send this information to `/api/v1/track`.

## Trigger the sending

It will triggered by the `t.TrackEvent()` function. It is going to send the message to the Collector if
- the unsent payloads number is bigger or equal to the `ClientTracker`'s `queue_size` attribute,
- or the last sync time has happened before current time minus the `ClientTracker`'s `queue_retention`.
