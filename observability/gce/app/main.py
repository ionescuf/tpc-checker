from flask import Flask
import logging
import os

app = Flask(__name__)

def setup_logging():
    """
    Configures the logging settings.
    Logs are written to both the console and a file.
    """
    # Create logs directory if it doesn't exist
    log_dir = "logs"
    os.makedirs(log_dir, exist_ok=True)
    
    # Log file path
    log_file = os.path.join(log_dir, "application.log")
    
    # Define the logging format
    log_format = "%(asctime)s [%(levelname)s] %(message)s"
    
    # Set up logging configuration
    logging.basicConfig(
        level=logging.DEBUG,  # Log level: DEBUG, INFO, WARNING, ERROR, or CRITICAL
        format=log_format,
        handlers=[
            logging.StreamHandler(),    # Output logs to console
            logging.FileHandler(log_file, mode='a')  # Append logs to the file
        ]
    )


def perform_operations():
    """
    Simulates operations that generate logs, even when errors occur.
    """
    logging.debug("This is a debug message, useful for troubleshooting.")
    logging.info("Application started successfully.")
    
    try:
        logging.info("Performing a division operation.")
        numerator = 10
        denominator = 0  # Intentional divide by zero
        result = numerator / denominator
        logging.info(f"Division result: {result}")
    except ZeroDivisionError:
        logging.error("Division by zero encountered!", exc_info=True)
        
        # Handle the failure gracefully by setting a default or skipping further operations
        logging.info("Default result will be set to 0 to continue operations.")
        result = 0

    # Continue with the rest of the logic, regardless of errors
    logging.info(f"Continuing with result = {result}.")
    
    # Simulate other operations
    logging.warning("This is a warning message to signal potential issues.")
    logging.critical("Critical issue detected, application needs attention!")

    logging.info("All operations completed without interruption.")

@app.route("/")
def home():
    return "The application is running!"

def main():
    setup_logging()
    app.run(host="0.0.0.0", port=5000)

if __name__ == "__main__":
    main()