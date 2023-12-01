import os
from logging.config import dictConfig
from typing import Callable, IO
from functools import lru_cache
import json

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from strawberry import Schema



from socketio import ASGIApp, async_manager

from src import config, main, authentication, constants, local_logging, routes
from src.config import DeployEnvironments, TEST_DIR
from src.config.env_settings import EnvSettings, get_env_settings
from src.config.local import LOGS_DIR, LogConfig
from src.config.routing import ROUTERS
from src.main import fastapi_app





from tests.additional_routes import test_router



test_logs_file = "tests.log"


@pytest.fixture(autouse=True)
def override_settings(monkeypatch) -> None:
    
    monkeypatch.setenv("DEPLOY_ENVIRONMENT", "TEST")
    get_env_settings.cache_clear()


@pytest.fixture()
def override_logs_file(monkeypatch) -> None:
    monkeypatch.setattr(config.local, "LOG_FILE_NAME", test_logs_file)


@pytest.fixture()
def include_test_routers(monkeypatch) -> None:
    def get_test_routers(local: bool = True):
        return [*ROUTERS, (routes.local_router, "/local"), (test_router, "/tests")]
    monkeypatch.setattr(config.routing, "get_routers", get_test_routers)




@pytest.fixture()
def app(override_settings, override_logs_file) -> FastAPI:
    fastapi_app.include_router(test_router, prefix="/tests")
    return fastapi_app



@pytest.fixture()
def log_file() -> IO:
    with open(LOGS_DIR / test_logs_file, "r") as f:
        f.read()  # set pointer to end of file
        yield f









@pytest.fixture(autouse=True)  # auto use to init fastapi
async def client(app,) -> TestClient:
    # admin has all permissions
    with TestClient(app=app, base_url="http://test", ) as c:
        yield c




@pytest.fixture()
def random_str_id() -> Callable[..., str]:
    i = 0

    def _counter() -> str:
        nonlocal i
        i += 1
        return f"{i:08d}"

    yield _counter

