FROM python:3.12-alpine

RUN pip install uv

WORKDIR /backend-executor

COPY . /backend-executor

RUN uv sync

CMD ["sh", "-c", "source .venv/bin/activate && python3 src/main.py"]
