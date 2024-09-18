class BuildConfig:
    def __init__(self, image: str, commands: list):
        self.image = image
        self.commands = commands

    def __str__(self):
        return f"Image: {self.image}, Commands: {self.commands}"
