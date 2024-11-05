from flask import Flask, render_template, request, jsonify
import json
import os
import requests  # Use requests to call Ollama API
import bjoern

# OpenTelemetry imports
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.flask import FlaskInstrumentor

app = Flask(__name__)

# Optional OpenTelemetry instrumentation
if os.environ.get("ENABLE_OTEL", "").lower() == "true":
    # Setup OpenTelemetry tracing
    resource = Resource(attributes={
        "service.name": "ouija-flask-app"
    })
    provider = TracerProvider(resource=resource)
    exporter = OTLPSpanExporter()
    span_processor = BatchSpanProcessor(exporter)
    provider.add_span_processor(span_processor)

    trace.set_tracer_provider(provider)

    # Instrument Flask
    FlaskInstrumentor().instrument_app(app)
# Load or initialize answers
try:
    with open("answers.json", "r") as f:
        answers = json.load(f)
except FileNotFoundError:
    answers = []

# Ollama instance URL
OLLAMA_INSTANCE_URL = "http://nn.starnix.net:11435/api/generate"

def generate_answer(question):
    # Define the mystical prompt for the Ouija board
    mystical_prompt = f"Pretend that you are a Ouija board. As a mystical Ouija board, answer the following question in a short answer. Respond without using any actions, such as *smiles*, *laughs*, or any text within asterisks. If the question is a yes or no question, answer with a yes or a no. Question: {question}"
    
    answer = ""  # Initialize an empty answer string to collect the streamed responses
    
    try:
        # Send the request to the Ollama instance and enable streaming
        with requests.post(
            OLLAMA_INSTANCE_URL,
            json={
                "model": "llama3.1",
                "prompt": mystical_prompt,
                "options": {
                    "num_predict": 10
                }
            },
            stream=True  # Enable streaming to handle line-by-line response
        ) as response:
            response.raise_for_status()  # Raise an error for bad responses (4xx or 5xx)

            # Process each line in the response stream
            for line in response.iter_lines():
                if line:
                    # Parse each line as a JSON object
                    line_data = json.loads(line)
                    # Append the "response" part to the answer
                    answer += line_data.get("response", "")

    except requests.exceptions.RequestException as e:
        print(f"Error contacting Ollama instance: {e}")
        answer = "The spirits cannot answer at this time. Try again later."

    return answer.strip()  # Return the fully concatenated answer

@app.route("/")
def index():
    return render_template("index.html")

@app.route("/ask", methods=["POST"])
def ask():
    question = request.json.get("question", "")
    answer = generate_answer(question)
    answers.append({"question": question, "answer": answer})
    with open("answers.json", "w") as f:
        json.dump(answers, f)
    return jsonify({"answer": answer})

@app.route("/history")
def history():
    return jsonify(answers)

if __name__ == "__main__":
    # Run the app using Bjoern
    bjoern.run(app, "0.0.0.0", 8080)