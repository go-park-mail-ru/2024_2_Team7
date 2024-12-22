math.randomseed(os.time())  

request = function()
    local id = math.random(190260, 308230)  
    local path = "/events/" .. id      
    return wrk.format("GET", path)
end
