import os.path
from logging.config import dictConfig

from src.config import BASE_DIR
from src.config.project_settings import get_project_settings

LOG_FORMAT: str = "%(levelprefix)s | %(asctime)s | %(message)s"
LOGS_DIR = BASE_DIR / "logs"
LOGGER_NAME = f"{get_project_settings().slug}log"
LOG_FILE_NAME = f"{get_project_settings().slug}.log"


def setup_logging():
    dictConfig(
        {
            "LOG_FORMAT": "%(levelprefix)s | %(asctime)s | %(message)s",
            "version": 1,
            "disable_existing_loggers": False,
            "formatters": {
                "default": {
                    "()": "uvicorn.logging.DefaultFormatter",
                    "fmt": LOG_FORMAT,
                    "datefmt": "%Y-%m-%d %H:%M:%S",
                },
                "file_formatter": {
                    "class": "logging.Formatter",
                    "format": "%(asctime)s\t%(levelname)s\t%(filename)s\t%(message)s",
                    "datefmt": "%d %b %y %H:%M:%S",
                },
            },
            "handlers": {
                "console": {
                    "class": "logging.StreamHandler",
                    "formatter": "default",
                    "level": "INFO",
                    "stream": "ext://sys.stdout",
                },
                "file": {
                    "class": "logging.handlers.RotatingFileHandler",
                    "formatter": "file_formatter",
                    "level": "DEBUG",
                    "filename": LOGS_DIR / LOG_FILE_NAME,
                    "mode": "a",
                    "encoding": "utf-8",
                    "maxBytes": 500000,
                    "backupCount": 4,
                },
            },
            "loggers": {LOGGER_NAME: {"handlers": ["console", "file"], "level": "DEBUG"}},
            "root": {"level": "DEBUG", "handlers": ["file"]},
        }
    )


if not os.path.exists(LOGS_DIR):
    os.mkdir(LOGS_DIR)
