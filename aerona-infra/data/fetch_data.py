import requests
import json

data = []
for i in range(50, 90):
    resp = requests.get('https://api.aviationstack.com/v1/flights?access_key=a8adeb12e4c186739b81eb72c343f62e&offset=' + str(i))
    temp = resp.json()
    data.extend(temp["data"])
json_data = json.dumps(data, default=lambda o: o.__dict__, indent=4)
print(json_data)
with open("data2.json", "w") as outfile:
    outfile.write(json_data)

