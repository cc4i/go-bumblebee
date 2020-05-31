FROM locustio/locust:1.0.2

WORKDIR /workspace
COPY ./locustfile.py /workspace/locustfile.py

CMD ["-f", "/workspace/locustfile.py"]

ENV PYTHONUNBUFFERED=1