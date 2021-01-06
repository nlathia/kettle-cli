import os

from src.model.artifacts import files


def test_get_path():
    # get_path() returns the path of a file in the artifacts directory
    # we test this by loading the path to the .py file itself
    path = files.get_path("files.py")
    assert os.path.exists(path)
