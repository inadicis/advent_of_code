from enum import Enum
from pathlib import Path
from typing import Self

BASE_DIR = Path(__file__).resolve().parent.parent.parent
APP_DIR = BASE_DIR / "src"
TEST_DIR = BASE_DIR / "tests"
DOCS_DIR = BASE_DIR / "docs"
DEPLOYMENT_DIR = BASE_DIR / "deployment"


class DeployEnvironments(str, Enum):
    DEV = "DEV"
    LOCAL = "LOCAL"
    STAGE = "STAGE"
    PRODUCTION = "PROD"
    TESTING = "TEST"

    @classmethod
    def is_env_remote(cls, environment: str | Self) -> bool:
        return environment in cls.remote_environments()

    @property
    def is_remote(self) -> bool:
        return self.is_env_remote(self)

    @property
    def is_local(self) -> bool:
        return not self.is_remote

    @classmethod
    def remote_environments(cls) -> set[Self]:
        return {cls.DEV, cls.STAGE, cls.PRODUCTION}
