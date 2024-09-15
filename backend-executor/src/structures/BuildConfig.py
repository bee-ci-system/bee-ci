class BuildConfig:
    def __init__(self, image, commands):
        self.image = image
        self.commands = commands
    def __str__(self):
        return f"Image: {self.image}, Commands: {self.commands}"