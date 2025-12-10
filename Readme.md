# WikiHunt-Bot

![Go](https://img.shields.io/badge/Go-1.24+-blue?logo=go)
![Python](https://img.shields.io/badge/Python-3.8+-blue?logo=python)
![License](https://img.shields.io/badge/License-MIT-green)

Find the shortest path between any two Wikipedia articles using intelligent semantic similarity matching.

## 📋 Table of Contents

- [What is WikiHunt-Bot?](#what-is-wikihunt-bot)
- [Key Features](#key-features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
- [Usage](#usage)
- [Architecture](#architecture)
- [How It Works](#how-it-works)
- [Contributing](#contributing)
- [Support](#support)

## What is WikiHunt-Bot?

WikiHunt-Bot is a tool that finds the shortest path between two Wikipedia articles by intelligently traversing Wikipedia links. It uses semantic similarity matching to determine which links are most likely to lead toward your target article, making the search process fast and efficient.

Think of it like a game of "Wikipedia Degrees of Separation" - starting from any article, the bot will navigate through Wikipedia links to reach your target article in the fewest steps possible.

## Key Features

- **Semantic Intelligence**: Uses Sentence Transformers to understand semantic similarity between article titles and find the most promising path
- **Fast Traversal**: Implements smart link selection to minimize the number of hops needed
- **Interactive Search**: User-friendly interface with autocomplete for Wikipedia article selection
- **Bidirectional Search**: Can search forward (via links) or backward (via backlinks)
- **Color-Coded Output**: Visual feedback with colored terminal output for better readability
- **GPU Acceleration**: Optional CUDA support for faster semantic similarity calculations

## Getting Started

### Prerequisites

- **Go** 1.24.2 or higher
- **Python** 3.8 or higher
- **pip** (Python package manager)
- **Internet connection** (for Wikipedia API access)

Optional:
- **CUDA** (for GPU acceleration of semantic similarity)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/indenigrate/WikiHunt-Bot.git
   cd WikiHunt-Bot
   ```

2. **Install Go dependencies**:
   ```bash
   go mod download
   ```

3. **Set up Python environment**:
   ```bash
   python3 -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   # If you're using fish shell: source venv/bin/activate.fish
   pip install -r requirements.txt
   ```

   Or install dependencies manually:
   ```bash
   pip install fastapi uvicorn sentence-transformers torch
   ```

### Running the Application

The application consists of two components that need to run together:

1. **Terminal 1 - Start the semantic similarity server**:
   ```bash
   source venv/bin/activate
   # If using fish shell: source venv/bin/activate.fish
   uvicorn sbert-server:app --reload
   ```

   The server will start on `http://localhost:8000`

2. **Terminal 2 - Run the main application**:
   ```bash
   go run .
   ```

   Follow the prompts to enter your starting and ending Wikipedia articles.

## Usage

Once the application is running, simply:

1. Enter the title of your starting Wikipedia article (with autocomplete suggestions)
2. Enter the title of your target Wikipedia article
3. Watch as WikiHunt-Bot navigates through Wikipedia to find a path between them

**Example**:
```
Title of starting page: Cat
Title of ending page: Eiffel Tower
```

The bot will display the path it takes with colored output showing each step of the journey.

## Architecture

The application is built with a hybrid architecture:

### Go Backend
- **main.go**: Entry point and user interaction
- **takeInput.go**: Handles user input with Wikipedia article autocomplete suggestions
- **execute.go**: Main search algorithm implementation (DFS/BFS traversal)
- **compare.go**: Wikipedia API integration and link extraction

### Python Semantic Server
- **sbert-server.py**: FastAPI server providing semantic similarity calculations using Sentence Transformers
  - Uses the `all-MiniLM-L6-v2` model for efficient embeddings
  - Provides `/similarity` endpoint for bulk similarity computation
  - Supports GPU acceleration via CUDA

### Key Dependencies

**Go:**
- `github.com/manifoldco/promptui`: Interactive CLI prompts with autocomplete
- `github.com/fatih/color`: Colored terminal output

**Python:**
- `fastapi`: Web framework for similarity API
- `sentence-transformers`: Pre-trained models for semantic similarity
- `torch`: Deep learning framework with GPU support

## How It Works

1. **Article Selection**: User provides starting and target Wikipedia articles with autocomplete assistance
2. **Link Fetching**: WikiHunt fetches all links/backlinks from the current Wikipedia article via the Wikipedia API
3. **Similarity Ranking**: Candidate links are sent to the semantic server which ranks them by similarity to the target article
4. **Path Selection**: The most semantically similar link is selected as the next step
5. **Iteration**: Steps 2-4 repeat until the target article is reached
6. **Path Completion**: The complete path is displayed to the user

The semantic similarity approach is much smarter than random clicking - it understands that "Feline" is closer to "Cat" than random articles, dramatically reducing the search space.

## Contributing

Contributions are welcome! Whether it's bug fixes, feature improvements, or optimizations, please feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Areas for improvement:
- Optimizing search algorithms (A* search, bidirectional search)
- Adding path visualization
- Implementing caching for frequently accessed articles
- Improving performance with advanced similarity models
- Adding configuration options

## Support

- **Issues**: Found a bug? [Open an issue](https://github.com/indenigrate/WikiHunt-Bot/issues) on GitHub
- **Questions**: Have a question? Check existing issues or create a new one with the `question` label
- **Wikipedia API**: For API-related questions, see [Wikipedia API Documentation](https://www.mediawiki.org/wiki/API)
- **Sentence Transformers**: For model information, visit [Sentence Transformers](https://www.sbert.net/)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Happy Wikipedia hunting! 🔍**
