import random
from locust import HttpUser, task, between

class QuickstartUser(HttpUser):
    wait_time = between(5, 9)

    @task
    def index_page(self):
        self.client.get("/air/version")
        self.client.get("/air/aqi")

    @task(3)
    def view_item(self):
        cities = ["beijing", 
                    "chengdu", 
                    "auckland", 
                    "london", 
                    "shanghai", 
                    "tianjing", 
                    "xian", 
                    "dalian", 
                    "shenzheng",
                    "guangzhou",
                    "wuhan",
                    "xiamen"]
        id = random.randint(0, len(cities)-1)
        self.client.get("/air/city/"+cities[id], name="/air/city")

