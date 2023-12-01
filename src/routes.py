import uuid
import base64
import json
import logging
from typing import Annotated
from functools import lru_cache

import pydantic
from pydantic import BaseModel, Extra, ConfigDict
from fastapi import APIRouter, Security, HTTPException, Depends
from starlette import status
from starlette.requests import Request


from starlette.responses import RedirectResponse

from src import config
from src.config.templates import templates
from src.config.env_settings import EnvSettings, get_env_settings
from src.config.project_settings import ProjectSettings, get_project_settings




main_router = APIRouter(tags=["Main Router"])


@main_router.get("/")
async def home_page():
    return RedirectResponse("/docs")








local_router = APIRouter(tags=["Routes only accessible when deployed locally"])


@local_router.get("/")
async def local_home_page():
    return "Hello localhost"


