import logging
import os
import sys

logger = logging.getLogger("ai-assistant-api")

def setup_logger():
    """
    Sets up a simple, text-based logger that writes to standard output.

    Log level is configurable via the `LOG_LEVEL` environment variable.
    (e.g., DEBUG, INFO, WARNING). Defaults to INFO.

    Includes a global exception hook to ensure all uncaught exceptions are logged.
    """
    log_level = os.getenv("LOG_LEVEL", "INFO").upper()

    logger.setLevel(logging.DEBUG)

    if logger.hasHandlers():
        logger.handlers.clear()

    log_handler = logging.StreamHandler(sys.stdout)
    log_handler.setLevel(log_level)

    formatter = logging.Formatter(
        "%(asctime)s - %(levelname)-8s - [%(module)s:%(lineno)d] - %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S"
    )

    log_handler.setFormatter(formatter)
    logger.addHandler(log_handler)

    def handle_exception(exc_type, exc_value, exc_traceback):
        """
        Log any uncaught exceptions using the configured logger.
        """
        if issubclass(exc_type, KeyboardInterrupt):
            sys.__excepthook__(exc_type, exc_value, exc_traceback)
            return

        logger.critical("Uncaught exception:", exc_info=(exc_type, exc_value, exc_traceback))

    sys.excepthook = handle_exception

setup_logger()