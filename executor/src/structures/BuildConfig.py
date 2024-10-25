class BuildConfig:
    def __init__(self, image: str, commands: list, timeout: int):
        self.image = image
        self.commands = commands
        self.timeout = timeout

    def __str__(self):
        return (
            f"Image: {self.image}, Commands: {self.commands}, Timeout: {self.timeout}"
        )
