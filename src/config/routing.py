from fastapi import APIRouter

from src import routes
from src.graphql import graphql_app

# TODO-CONFIG new routers have to be registered here
ROUTERS: list[tuple[APIRouter, str]] = [
    (routes.main_router, ""),  # (router, url prefix path)
    
]


def get_routers(local: bool = False):
    return ROUTERS + [(routes.local_router, "/local")] if local else ROUTERS
