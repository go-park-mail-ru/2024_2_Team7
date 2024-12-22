math.randomseed(os.time())
request = function()
    local id = math.random(1, 100000)
    local path = "/events/" .. id
    return wrk.format("GET", path)
end