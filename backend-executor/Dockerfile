# Use the official Python image from the Docker Hub
FROM python:3.12-alpine

# Install any needed packages specified in requirements.txt
RUN pip install uv

# Set the working directory in the container
WORKDIR /backend-executor

# Copy the current directory contents into the container at /backend-executor
COPY . /backend-executor

# Sync the uv (assuming it's a Python package)
RUN uv sync

# Activate the virtual environment and run the application
CMD ["sh", "-c", "source .venv/bin/activate && python3 src/main.py"]