FROM python:3.9-alpine

RUN mkdir /app
WORKDIR /app
COPY . .

RUN pip install -r requirements.txt

ENV {{.ApiEndpointEnv}}=:80
ENV {{.MaintenanceApiEndpointEnv}}=:8000
EXPOSE 80/tcp
EXPOSE 8000/tcp
ENTRYPOINT ["python", "-m", "{{.ServiceName}}"]
