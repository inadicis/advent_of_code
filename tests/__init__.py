import logging
from importlib import metadata

from src.config import BASE_DIR

# Check that all required dependencies are installed

with open(BASE_DIR / "requirements.txt", "r") as f:
    requirements = f.readlines()
separators = [">=", "<=", "==", "#", ">", "<", "[", " ", ",", ";", "@"]
separators.sort(key=lambda x: -len(x))

missing_requirements = []
for requirement in requirements:
    requirement = requirement.strip()
    if not requirement or requirement.startswith("#"):
        continue

    if requirement.startswith('-') or requirement.startswith('--requirement'):
        logging.warning(f"Did not analyze line `{requirement}` (- not supported)")
        # TODO
        continue

    indexes_of_separators = [(requirement.index(s), s) for s in separators if s in requirement]

    if indexes_of_separators:
        first_separator_i = min(indexes_of_separators)[0]
        dependency = requirement[:first_separator_i]
        modifiers = requirement[first_separator_i:]
    else:
        dependency = requirement

    try:
        while "${" in dependency:
            dependency = dependency[dependency.index("}") + 1:]
        v = metadata.version(dependency)
    except metadata.PackageNotFoundError:
        missing_requirements.append(requirement)
        # TODO validate version

if missing_requirements:
    raise metadata.PackageNotFoundError(
        f"Required dependencies are not installed: {missing_requirements}"
    )
