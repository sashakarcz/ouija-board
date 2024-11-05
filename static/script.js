// Positions for each letter, number, and special words on the board
const letterPositions = {
    // Letters A-Z
    'A': { x: 27, y: 70 },
    'B': { x: 34, y: 65 },
    'C': { x: 40, y: 60 },
    'D': { x: 46, y: 58 },
    'E': { x: 52, y: 56 },
    'F': { x: 58, y: 54 },
    'G': { x: 65, y: 55 },
    'H': { x: 71, y: 56 },
    'I': { x: 77, y: 58 },
    'J': { x: 82, y: 58 },
    'K': { x: 88, y: 60 },
    'L': { x: 94, y: 63 },
    'M': { x: 101, y: 69 },
    'N': { x: 27, y: 88 },
    'O': { x: 32, y: 82 },
    'P': { x: 37, y: 77 },
    'Q': { x: 43, y: 74 },
    'R': { x: 50, y: 71 },
    'S': { x: 56, y: 69 },
    'T': { x: 62, y: 68 },
    'U': { x: 68, y: 68 },
    'V': { x: 76, y: 69 },
    'W': { x: 84, y: 72 },
    'X': { x: 92, y: 76 },
    'Y': { x: 98, y: 80 },
    'Z': { x: 102, y: 88 },
    
    // Numbers 1-0
    '1': { x: 39, y: 96 },
    '2': { x: 43, y: 96 },
    '3': { x: 49, y: 96 },
    '4': { x: 55, y: 96 },
    '5': { x: 61, y: 96 },
    '6': { x: 67, y: 96 },
    '7': { x: 73, y: 96 },
    '8': { x: 78, y: 96 },
    '9': { x: 84, y: 96 },
    '0': { x: 90, y: 96 },

    // Special positions
    'YES': { x: 38, y: 38 },
    'NO': { x: 92, y: 39 },
    'GOOD BYE': { x: 65, y: 112 }
};


// Handle form submission
document.getElementById("questionForm").addEventListener("submit", function(event) {
    event.preventDefault();
    const question = document.getElementById("questionInput").value;

    // Send the question to the server and get the answer
    fetch("/ask", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ question: question })
    })
    .then(response => response.json())
    .then(data => {
        // Animate the planchette to spell out the answer
        animatePlanchette(data.answer);
        // Display the final answer below the board
        displayAnswer(data.answer);
    });
});

// Function to display the final answer below the Ouija board
function displayAnswer(answer) {
    const answerElement = document.getElementById("answer");
    answerElement.innerText = "Answer: " + answer;
    answerElement.classList.add("show-answer"); // Trigger fade-in animation
}


const planchetteSize = 30; // Adjust this based on the actual size of the planchette in vw or vh
const planchetteWindowOffset = planchetteSize / 2; // Offset to center the window

// Function to animate the planchette to each letter in the answer
function animatePlanchette(answer) {
    const planchette = document.getElementById("planchette");

    // Convert the answer to uppercase for consistent matching
    const upperAnswer = answer.toUpperCase().trim();

    // Check for special cases if the answer contains "YES", "NO", or "BYE"
    let specialPosition;
    if (upperAnswer.includes("YES")) {
        specialPosition = letterPositions['YES'];
    } else if (upperAnswer.includes("NO")) {
        specialPosition = letterPositions['NO'];
    } else if (upperAnswer.includes("BYE", "GOOD BYE", "FAREWELL", "NEXT TIME")) {
        specialPosition = letterPositions['GOOD BYE'];
    }

    if (specialPosition) {
        // Move planchette directly to the special position
        planchette.style.left = (specialPosition.x - planchetteWindowOffset) + "%";
        planchette.style.top = (specialPosition.y - planchetteWindowOffset) + "%";
        return; // Exit function after moving to special position
    }

    // If it's not a special case, proceed with spelling out each letter
    const letters = upperAnswer.split(""); // Split answer into individual letters
    let index = 0;

    const interval = setInterval(() => {
        if (index < letters.length) {
            const letter = letters[index];
            const position = letterPositions[letter] || letterPositions[' ']; // Default to a central position for unknown characters

            if (position) {
                // Center the planchette over the letter's position
                planchette.style.left = (position.x - planchetteWindowOffset) + "%";
                planchette.style.top = (position.y - planchetteWindowOffset) + "%";
            }

            index++; // Move to the next letter
        } else {
            clearInterval(interval); // Stop the animation when done
        }
    }, 700); // Adjust the timing for each movement as needed
}