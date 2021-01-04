import os


def get_path(file_name: str) -> str:
    """Returns a path to a file in this directory.
    Args:
        file_name (str): The name of a file in the same directory as
        this python file (e.g., "model.pkl").
    Returns:
        The full path to a file in this directory
        (e.g., "/path/to/model.pkl")
    """
    directory = os.path.dirname(os.path.realpath(__file__))
    return os.path.join(directory, file_name)