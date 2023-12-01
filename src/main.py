import os
import logging
from functools import lru_cache

import uvicorn
from fastapi import FastAPI
from google.cloud.logging_v2 import Client


from starlette.middleware.cors import CORSMiddleware
from starlette.requests import Request

from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse

from src import models, routes
from src.config import local, env_settings, project_settings, routing




project: project_settings.ProjectSettings = project_settings.get_project_settings()

fastapi_app = FastAPI(
    title=project.name,
    description=project.description,
    version=project.version,
    contact=project.authors[0]
)

fastapi_app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)



if env_settings.get_env_settings().debug:
    @fastapi_app.exception_handler(RequestValidationError)
    async def validation_exception_handler(request: Request, exc: RequestValidationError):
        exc_str = f'{exc}'.replace('\n', ' ').replace('   ', ' ')
        logging.warning(f"{request}: {exc_str}")
        content = {'status_code': 10422, 'message': exc_str, 'data': None}
        return JSONResponse(content=content, status_code=422)


@fastapi_app.on_event("startup")
async def app_init():
    """
    Setups logging, connects to database, initializes ODM and include routers in the app.
    """
    settings = env_settings.get_env_settings()
    project = project_settings.get_project_settings()
    if settings.deploy_environment.is_local:
        local.setup_logging()
    else:
        Client().setup_logging(log_level=logging.DEBUG)

    
    for router, prefix in routing.get_routers(local=settings.deploy_environment.is_local):
        fastapi_app.include_router(router, prefix=prefix)
    
    
    logging.info(f"{project.name}: Finished initializing")


app = fastapi_app


if __name__ == "__main__":  # for debugging locally
    uvicorn.run("src.main:app",
                host="127.0.0.1",
                port=8000,
                
    )
    # run with `python -m src.main`
