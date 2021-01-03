from src.model.model import predict, load_model

model = load_model()


def {{.FunctionName}}(event, context):
    """Responds to any HTTP request.
    Args:
        event (usually a dict): An event is a JSON-formatted document that contains 
        data for a Lambda function to process.
        context: This object provides methods and properties that provide information 
        about the invocation, function, and runtime environment. 
    Returns:
        If the handler returns objects that can't be serialized by json.dumps, 
        the runtime returns an error. 
    """
    # @TODO add any request validations that you need
    return predict(model=model, input=event)
