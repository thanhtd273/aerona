import requests
import json

data = []
for i in range(1, 30):
    resp = requests.get('https://api.aviationstack.com/v1/cities?access_key=a8adeb12e4c186739b81eb72c343f62e&offset=' + str(i))
    print("Called successfully" + str(i) + "/30 API")
    temp = resp.json()
    data.extend(temp["data"])
json_data = json.dumps(data, default=lambda o: o.__dict__, indent=4)
print(len(data))
with open("airline.json", "w") as outfile:
    outfile.write(json_data)

