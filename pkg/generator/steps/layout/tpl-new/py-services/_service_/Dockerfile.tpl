FROM python:3.9-alpine

RUN mkdir /app
WORKDIR /app
COPY . .

RUN pip install -r requirements.txt

ENV {{.ApiEndpointEnv}}=:80
EXPOSE 80/tcp
ENTRYPOINT ["python", "-m", "{{.ServiceName}}"]
