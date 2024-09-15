import json
from structures.BuildConfig import BuildConfig

class BuildConfigAnalyzer:
    @staticmethod
    def save_script(commands: list):
        # Save the script content to default file
        script_content = "#!/bin/bash\n" + "\n".join(commands)
        with open("run.sh", "w") as f:
            f.write(script_content)
    @staticmethod
    def analyze(json_data) -> BuildConfig:
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
            print(str(config))
            return config

        except json.JSONDecodeError as e:
            print("Invalid JSON:", str(e))
            return None