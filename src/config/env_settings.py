"""
Configuration of .env variables and logging, as well as other settings constants
"""
from functools import lru_cache

from pydantic_settings import BaseSettings, SettingsConfigDict

from src.config import DEPLOYMENT_DIR, DeployEnvironments


class EnvSettings(BaseSettings):
    """
    read env variables, cast and validate them.
    variable names are case insensitive.
    """

    model_config = SettingsConfigDict(env_file=DEPLOYMENT_DIR / ".env", validate_default=True)

    deploy_environment: DeployEnvironments = DeployEnvironments.LOCAL

    @property
    def debug(self) -> bool:
        return self.deploy_environment != DeployEnvironments.PRODUCTION


@lru_cache(maxsize=1)
def get_env_settings() -> EnvSettings:
    return EnvSettings()
