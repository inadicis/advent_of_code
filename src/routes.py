from fastapi import APIRouter
from starlette.responses import RedirectResponse

main_router = APIRouter(tags=["Main Router"])


@main_router.get("/")
async def home_page():
    return RedirectResponse("/docs")


local_router = APIRouter(tags=["Routes only accessible when deployed locally"])


@local_router.get("/")
async def local_home_page():
    return "Hello localhost"
