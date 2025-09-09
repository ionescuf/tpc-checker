# Use the official Python image as the base image
FROM python:3.9-slim

# Set work directory inside the container
WORKDIR ./app

# Copy the current directory contents to /app in the container
COPY . .

# Install any application dependencies (if needed)
# Here we use requirements.txt in case you add external libraries
# If there are no dependencies, you can skip this step
RUN pip install --no-cache-dir -r requirements.txt || true

# Set the command to run the application
CMD ["sh", "-c", "python main.py && tail -f /dev/null"]