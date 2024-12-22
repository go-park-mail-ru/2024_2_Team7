local boundary = "--TEST_BOUNDARY"  -- Разделитель multipart данных
local max_requests = 100000         -- Общее количество запросов
local counter = 0                   -- Счётчик выполненных запросов

wrk.method = "POST"
wrk.headers["Content-Type"] = "multipart/form-data; boundary=" .. boundary

wrk.body =
    "--" .. boundary .. "\r\n" ..
    "Content-Disposition: form-data; name=\"json\"\r\n\r\n" ..
    [[{
        "title": "Test Event",
        "description": "Sample description",
        "location": "Test Location",
        "event_start": "2024-12-25T12:00:00Z",
        "event_end": "2024-12-25T14:00:00Z",
        "category_id": 1,
        "capacity": 100,
        "tag": ["tag1", "tag2"]
    }]] .. "\r\n" ..
    "--" .. boundary .. "\r\n"

request = function()
    if counter >= max_requests then
        wrk.thread:stop()  -- Остановка после выполнения 100,000 запросов
    end
    counter = counter + 1
    return wrk.format(nil, "/events")  -- Отправка POST-запроса на /events
end
