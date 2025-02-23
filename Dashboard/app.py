from flask import Flask, render_template, jsonify, request
import requests

app = Flask(__name__)

# In-memory store for experiments (for demonstration purposes)
experiments = []


@app.route("/")
def index():
    return render_template("index.html", experiments=experiments)


@app.route("/api/experiments", methods=["GET"])
def get_experiments():
    return jsonify({"experiments": experiments})


@app.route("/api/experiments", methods=["POST"])
def add_experiment():
    data = request.json
    experiments.append(data)
    return jsonify({"status": "success", "experiment": data}), 201


@app.route("/refresh")
def refresh_experiments():
    try:
        resp = requests.get("http://localhost:8080/experiments")
        experiments = resp.json()
    except Exception as e:
        experiments = []
        print("Error fetching experiments:", e)
    return render_template("index.html", experiments=experiments)


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5500, debug=True)
