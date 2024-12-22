local boundary = "--TEST_BOUNDARY"  -- Разделитель multipart данных
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

