import pytest

from src.config.env_settings import EnvSettings, get_env_settings
from src.config import DeployEnvironments


@pytest.mark.asyncio
async def test_settings(override_settings):
    e = EnvSettings()
    assert e.deploy_environment == DeployEnvironments.TESTING
    assert get_env_settings().deploy_environment == DeployEnvironments.TESTING
