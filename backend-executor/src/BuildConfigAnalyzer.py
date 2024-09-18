import json
import logging
from structures.BuildConfig import BuildConfig

logger = logging.getLogger(__name__)


class BuildConfigAnalyzer:
    @staticmethod
    def save_script(commands: list):
        script_content = "#!/bin/bash\n" + "\n".join(commands)
        with open("run.sh", "w", encoding="utf-8") as f:
            f.write(script_content)

    @staticmethod
    def analyze(json_data: str) -> BuildConfig:
        try:
            # Parse the JSON data
            data = json.loads(json_data)

            # Extract the image and commands
            image = data.get("image")
            commands = data.get("commands")
            config = BuildConfig(image, commands)

            # Perform analysis on the extracted data
            # For example, you can print them or perform any other operations
            BuildConfigAnalyzer.save_script(commands)
            logger.info("Got: %s", str(config))
            return config

        except json.JSONDecodeError as e:
            logger.error("Invalid JSON: %s", str(e))
            return None
