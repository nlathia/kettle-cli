import os

from flask import Flask, jsonify, request

from src.model.model import predict, load_model

app = Flask(__name__)
model = load_model()


@app.route('/', methods=['POST'])
def {{.FunctionName}}():
    request_json = request.get_json()
    # @TODO add any request validations that you need
    result = predict(model=model, input=request_json)
    return jsonify(result)


if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=int(os.environ.get('PORT', 8080)))
