import requests

url = "http://localhost:2000/new_img"
body = "the Kate's email is read"

response = requests.post(url, data=body)

print("Status:", response.status_code)
print("Response:", response.text)
