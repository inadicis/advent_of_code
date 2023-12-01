from functools import lru_cache

import toml
from pydantic import BaseModel, Extra, ConfigDict
from src.config import BASE_DIR


class ProjectSettings(BaseModel):
    model_config = ConfigDict(extra=Extra.ignore)
    name: str
    version: str
    description: str
    authors: list[dict[str, str]]

    def __init__(self, **kwargs):
        with open(BASE_DIR / "pyproject.toml", "r") as f:
            project_settings: dict = toml.loads(f.read())["project"]
        project_settings.update(kwargs)
        super().__init__(**project_settings)

    @property
    def slug(self) -> str:
        return self.name.lower().replace(" ", "_")


@lru_cache(maxsize=1)
def get_project_settings() -> ProjectSettings:
    return ProjectSettings()
