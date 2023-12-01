"""
Additionnal routes that are only added to the app when testing it automatically.
Useful for language / framework / tools testing.
"""
from typing import Annotated

from fastapi import APIRouter, Security

test_router = APIRouter(tags=["Automated Testing Routes"])


@test_router.get("/openroute/")
async def route_unrestricted():
    return {"success": True}


