import pytest

from src.config import DeployEnvironments
from src.config.env_settings import EnvSettings, get_env_settings


@pytest.mark.asyncio
async def test_settings(override_settings):
    e = EnvSettings()
    assert e.deploy_environment == DeployEnvironments.TESTING
    assert get_env_settings().deploy_environment == DeployEnvironments.TESTING
