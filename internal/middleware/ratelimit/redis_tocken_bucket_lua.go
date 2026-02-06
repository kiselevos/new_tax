package ratelimit

const redisTokenBucketLua = `
-- KEYS[1] = bucket key (hash)
-- ARGV[1] = rate tokens/sec
-- ARGV[2] = burst (capacity)
-- ARGV[3] = ttl_ms

local key = KEYS[1]
local rate = tonumber(ARGV[1])
local burst = tonumber(ARGV[2])
local ttl_ms = tonumber(ARGV[3])

if rate == nil or burst == nil or ttl_ms == nil then
  return {0, 0}
end
if rate <= 0 or burst <= 0 or ttl_ms <= 0 then
  return {0, 0}
end

local t = redis.call("TIME")
local now_ms = (tonumber(t[1]) * 1000) + math.floor(tonumber(t[2]) / 1000)

local data = redis.call("HMGET", key, "tokens", "ts")
local tokens = tonumber(data[1])
local last_ms = tonumber(data[2])

if tokens == nil then tokens = burst end
if last_ms == nil then last_ms = now_ms end

local elapsed_ms = now_ms - last_ms
if elapsed_ms < 0 then elapsed_ms = 0 end

local refill = (elapsed_ms / 1000.0) * rate
tokens = math.min(burst, tokens + refill)

local allowed = 0
if tokens >= 1.0 then
  allowed = 1
  tokens = tokens - 1.0
end

redis.call("HMSET", key, "tokens", tokens, "ts", now_ms)
redis.call("PEXPIRE", key, ttl_ms)

return {allowed, tokens}
`
