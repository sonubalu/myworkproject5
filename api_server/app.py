from flask import Flask, request, jsonify
import json

app = Flask(__name__)

@app.route('/summary', methods=['POST'])
def write_summary():
    data = request.json
    with open('parking_summary.txt', 'a') as f:
        f.write(json.dumps(data) + '\n')
    return jsonify({"status": "success"}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
