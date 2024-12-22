local counter = 0
local max_requests = 100000

request = function()
    if counter >= max_requests then
        wrk.thread:stop()  -- Остановка при достижении 100,000 запросов
    end
    counter = counter + 1
    local id = math.random(1, 100000)
    local path = "/events/" .. id
    return wrk.format("GET", path)
end
