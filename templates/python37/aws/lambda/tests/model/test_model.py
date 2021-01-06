import pytest
from src.model import model


@pytest.fixture
def test_model():
    return model.load_model()


def test_load_model(test_model):
    # @TODO update this test to check
    # For example, assert isinstance(clf, RandomForestRegressor)
    assert test_model is None


@pytest.mark.parametrize(
    "input_dict,expected",
    [
        ({}, {"prediction": "hello world!"}),
    ],
)
def test_predict(input_dict, expected, test_model):
    result = model.predict(test_model, input_dict)
    assert result == expected
