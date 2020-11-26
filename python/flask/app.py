#app.py

from flask import Flask, request, jsonify
from datetime import datetime

app = Flask(__name__)
app.config["DEBUG"] = True
app.config["PORT"] = 5000

@app.route('/', methods=['GET'])
def root():
    response = {"msg": "Aplicação de Exemplo para Integração com API Pública KeyAPP","timestamp": datetime.now()}
    return jsonify(response)

@app.route('/signed-actions', methods=['POST']) #GET requests will be blocked
def parse_payload():
    req_data = request.get_json()
    
    # TODO: get token and verify with public api
    
    return req_data

app.run()