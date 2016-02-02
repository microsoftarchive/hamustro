-- Note: This file is _not_ using any LUA dependency on purpose. (like: md5, lfs)
--   Please keep in mind when you modify this script.

-- Initialize the pseudo random number generator
math.randomseed(os.time())

MESSAGE_PATH = ""
METHOD = "POST"
NR_FILES = 100

-- Load body informations from files
bodies = {}
for i=1,NR_FILES do
	f = io.open(MESSAGE_PATH .. i .. ".pb", "rb")
	bodies[i] = f:read("*all")
	f:close()
end

-- Load signature informations from files (shared_secret: "ultrasafesecret")
signatures = {}
for i=1,NR_FILES do
	f = io.open(MESSAGE_PATH .. i .. ".signature", "r")
	signatures[i] = f:read("*all")
	f:close()
end

-- Send a random messages
request = function()
	idx = math.random(1,NR_FILES)
	body = bodies[idx]
	headers = {}
	headers["Content-Type"] = "application/x-google-protobuf"
	headers["X-Hamustro-Time"] = "2016-01-01T00:00:00"
	headers["X-Hamustro-Signature"] = signatures[idx]
	return wrk.format(METHOD, wrk.url, headers, bodies[idx])
end