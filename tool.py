import time
import datetime
import json
import requests


def addRange(category, day, start, end):
    s = datetime.datetime(2021, 12, day, start)
    e = datetime.datetime(2021, 12, day, end)
    start_time = time.mktime(s.timetuple())
    end_time = time.mktime(e.timetuple())
    print(requests.post('http://localhost:8081/api/range/add',
                        json.dumps({'token': 'woshimima', 'range': {
                            'category': category,
                            'start': int(start_time)*1000,
                            'end': int(end_time)*1000
                        }})).text)

cat = 'T6'
for i in range(20, 25):
    addRange(cat, i, 8, 10)
    addRange(cat, i, 10, 12)
    addRange(cat, i, 12, 15)
    addRange(cat, i, 15, 18)
    addRange(cat, i, 18, 21)

for i in range(27, 32):
    addRange(cat, i, 8, 10)
    addRange(cat, i, 10, 12)
    addRange(cat, i, 12, 15)
    addRange(cat, i, 15, 18)
    addRange(cat, i, 18, 21)
cat = 'V26'
for i in range(20, 25):
    addRange(cat, i, 8, 10)
    addRange(cat, i, 10, 12)
    addRange(cat, i, 12, 15)
    addRange(cat, i, 15, 18)
    addRange(cat, i, 18, 21)

for i in range(27, 32):
    addRange(cat, i, 8, 10)
    addRange(cat, i, 10, 12)
    addRange(cat, i, 12, 15)
    addRange(cat, i, 15, 18)
    addRange(cat, i, 18, 21)
