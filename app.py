from flask import Flask, render_template, request, jsonify
import json
import openai
import os
from dotenv import load_dotenv
import bjoern

# Load environment variables from .env file
load_dotenv()

app = Flask(__name__)

# Load or initialize answers
try:
    with open("answers.json", "r") as f:
        answers = json.load(f)
except FileNotFoundError:
    answers = []

# Set up OpenAI API key from .env
openai.api_key = os.getenv("OPENAI_API_KEY")

def generate_answer(question):
    try:
        response = openai.ChatCompletion.create(
            model="gpt-3.5-turbo",  # Or another compatible model
            messages=[
                {"role": "system", "content": "You are a mysterious Ouija board answering questions with brief, mystical responses."},
                {"role": "user", "content": question}
            ]
        )
        answer = response['choices'][0]['message']['content'].strip()
    except Exception as e:
        print(f"Error fetching answer from ChatGPT: {e}")
        answer = "I am unable to answer at the moment."
    
    return answer


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
    bjoern.run(app, "0.0.0.0", 8000)
