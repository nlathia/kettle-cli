from flask import jsonify

from src.model.model import predict, load_model

model = load_model()


def {{.FunctionName}}(request):
    """Responds to any HTTP request.
    Args:
        request (flask.Request): HTTP request object.
    Returns:
        The response text or any set of values that can be turned into a
        Response object using
        `make_response <http://flask.pocoo.org/docs/1.0/api/#flask.Flask.make_response>`.
    """
    request_json = request.get_json()
    # @TODO add any request validations that you need
    result = predict(model=model, input=request_json)
    return jsonify(result)
