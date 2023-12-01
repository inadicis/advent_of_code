FROM python:3.11
# TODO-CONFIG might want to change python version

WORKDIR /app

COPY requirements.txt requirements.txt

RUN python3 -m pip install --no-cache-dir --upgrade -r requirements.txt

COPY pyproject.toml .
COPY src src

CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8080", "--workers", "1"]

# buildable locally from repo root as working directory with
# docker build -f deployment/Dockerfile -t ml_toolkit .
# then starting a container with
# docker run -p 127.0.0.1:80:8080/tcp ml_toolkit
# now service accessible at 127.0.0.1/docs
