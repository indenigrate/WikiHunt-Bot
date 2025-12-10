document.addEventListener('DOMContentLoaded', () => {
    const startInput = document.getElementById('start-article');
    const endInput = document.getElementById('end-article');
    const findPathButton = document.getElementById('find-path-btn');
    const startServerButton = document.getElementById('start-server-btn');
    const serverStatus = document.getElementById('server-status');
    const resultsContainer = document.getElementById('results-container');
    const pathList = document.getElementById('path-list');
    const loadingSpinner = document.getElementById('loading-spinner');
    const errorMessage = document.getElementById('error-message');

    const startSuggestions = document.getElementById('start-suggestions');
    const endSuggestions = document.getElementById('end-suggestions');

    const GO_SERVER_BASE_URL = 'https://wikihunt-bot-1.onrender.com';
    const PYTHON_SERVER_BASE_URL = 'https://devanshsoni899-newenv--bge-similarity-api-fastapi-app.modal.run';

    // Health check function
    async function healthCheck() {
        try {
            const goServerResponse = await fetch(`${GO_SERVER_BASE_URL}/healthz`);
            const pythonServerResponse = await fetch(`${PYTHON_SERVER_BASE_URL}/healthz`);

            if (goServerResponse.ok && pythonServerResponse.ok) {
                serverStatus.textContent = 'Both servers are running!';
                serverStatus.className = 'text-green-500';
            } else {
                let statusMessage = 'Server status: ';
                statusMessage += `Go server ${goServerResponse.ok ? 'OK' : 'starting...'}. `;
                statusMessage += `Python server ${pythonServerResponse.ok ? 'OK' : 'starting...'}.`;
                serverStatus.textContent = statusMessage;
                serverStatus.className = 'text-orange-500';
            }
        } catch (error) {
            serverStatus.textContent = 'One or both servers are starting, please wait a couple of minutes...';
            serverStatus.className = 'text-orange-500';
        }
    }

    // Call health check on page load
    healthCheck();

    startServerButton.addEventListener('click', healthCheck);

    // 1. Debounce function to limit API calls
    function debounce(func, wait) {
        let timeout;
        return function (...args) {
            clearTimeout(timeout);
            timeout = setTimeout(() => func.apply(this, args), wait);
        };
    }
    // 2. Function to fetch Wikipedia suggestions
    async function fetchWikipediaSuggestions(query) {
        if (!query) return [];
        const url = `https://en.wikipedia.org/w/api.php?action=opensearch&format=json&origin=*&search=${encodeURIComponent(query)}&limit=5`;
        try {
            const response = await fetch(url);
            const data = await response.json();
            return data[1]; // Index 1 contains the list of titles
        } catch (error) {
            console.error("Wiki fetch error:", error);
            return [];
        }
    }

    // 3. Function to render suggestions
    function renderSuggestions(suggestions, container, inputElement) {
        const ul = container.querySelector('ul');
        ul.innerHTML = '';

        if (suggestions.length === 0) {
            container.classList.add('hidden');
            return;
        }

        suggestions.forEach(title => {
            const li = document.createElement('li');
            li.textContent = title;
            li.className = 'cursor-pointer select-none px-4 py-2 hover:bg-orange-100 text-text-main transition-colors';

            li.addEventListener('click', () => {
                inputElement.value = title;
                container.classList.add('hidden');
            });

            ul.appendChild(li);
        });

        container.classList.remove('hidden');
    }

    // 4. Setup Input Listeners
    function setupAutocomplete(inputElement, suggestionContainer) {
        let currentFocus = -1; // Track the currently selected item

        const handleInput = debounce(async (e) => {
            const query = e.target.value.trim();
            currentFocus = -1; // Reset focus on new typing

            if (query.length < 2) {
                suggestionContainer.classList.add('hidden');
                return;
            }

            const suggestions = await fetchWikipediaSuggestions(query);
            renderSuggestions(suggestions, suggestionContainer, inputElement);
        }, 300);

        inputElement.addEventListener('input', handleInput);

        // KEYBOARD NAVIGATION LOGIC
        inputElement.addEventListener('keydown', function (e) {
            let items = suggestionContainer.querySelectorAll('li');

            if (e.key === 'ArrowDown') { // Down Arrow
                currentFocus++;
                addActive(items);
            } else if (e.key === 'ArrowUp') { // Up Arrow
                currentFocus--;
                addActive(items);
            } else if (e.key === 'Enter') { // Enter
                e.preventDefault(); // Prevent form submission
                if (currentFocus > -1 && items) {
                    items[currentFocus].click(); // Simulate click on selected item
                }
            }
        });
        function addActive(items) {
            if (!items || items.length === 0) return false;

            // Wrap around logic
            if (currentFocus >= items.length) currentFocus = 0;
            if (currentFocus < 0) currentFocus = (items.length - 1);

            removeActive(items);

            // Add active class (matching your hover style)
            items[currentFocus].classList.add('bg-orange-100');

            // Ensure the selected item is visible in the scrollable area
            items[currentFocus].scrollIntoView({ block: 'nearest' });
        }

        function removeActive(items) {
            for (let i = 0; i < items.length; i++) {
                items[i].classList.remove('bg-orange-100');
            }
        }

        // Hide suggestions when clicking outside
        document.addEventListener('click', (e) => {
            if (!inputElement.contains(e.target) && !suggestionContainer.contains(e.target)) {
                suggestionContainer.classList.add('hidden');
            }
        });

        // Show suggestions again if focusing back on an input with text
        inputElement.addEventListener('focus', async () => {
            if (inputElement.value.trim().length >= 2) {
                const suggestions = await fetchWikipediaSuggestions(inputElement.value.trim());
                renderSuggestions(suggestions, suggestionContainer, inputElement);
            }
        });
    }
    // Initialize autocomplete for both inputs
    setupAutocomplete(startInput, startSuggestions);
    setupAutocomplete(endInput, endSuggestions);
    let pollingInterval;

    findPathButton.addEventListener('click', async () => {
        const startArticle = startInput.value.trim();
        const endArticle = endInput.value.trim();

        if (!startArticle || !endArticle) {
            showError('Please enter both a starting and a target article.');
            return;
        }

        hideError();
        showLoading();
        clearResults();
        if (pollingInterval) {
            clearInterval(pollingInterval);
        }

        try {
            const response = await fetch(`${GO_SERVER_BASE_URL}/path`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    start: startArticle,
                    end: endArticle,
                }),
            });

            if (response.status !== 202) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to start pathfinding task.');
            }

            const data = await response.json();
            startPolling(data.task_id);

        } catch (error) {
            showError(error.message);
            hideLoading();
        }
    });

    function startPolling(taskId) {
        pollingInterval = setInterval(async () => {
            try {
                const response = await fetch(`${GO_SERVER_BASE_URL}/path/${taskId}`);
                if (!response.ok) {
                    // Stop polling on error
                    clearInterval(pollingInterval);
                    const errorData = await response.json();
                    showError(errorData.error || 'An error occurred while fetching the path.');
                    hideLoading();
                    return;
                }

                const task = await response.json();
                displayPath(task.path); // Update the path display

                if (task.done) {
                    clearInterval(pollingInterval);
                    hideLoading();
                    if (task.error) {
                        showError(task.error);
                    }
                }
            } catch (error) {
                clearInterval(pollingInterval);
                showError('Failed to poll for task status.');
                hideLoading();
            }
        }, 2000); // Poll every 2 seconds
    }

    function showLoading() {
        loadingSpinner.classList.remove('hidden');
        findPathButton.disabled = true;
        findPathButton.querySelector('span').textContent = 'Searching...';
        resultsContainer.classList.remove('hidden'); // Show container for loading spinner
    }

    function hideLoading() {
        loadingSpinner.classList.add('hidden');
        findPathButton.disabled = false;
        findPathButton.querySelector('span').textContent = 'Find Path';
    }

    function showError(message) {
        errorMessage.textContent = message;
        errorMessage.classList.remove('hidden');
    }

    function hideError() {
        errorMessage.classList.add('hidden');
    }

    function clearResults() {
        pathList.innerHTML = '';
        resultsContainer.classList.add('hidden');
    }

    function displayPath(path) {
        if (!path || path.length === 0) {
            return; // Nothing to display yet
        }

        pathList.innerHTML = ''; // Clear and redraw path

        path.forEach((article, index) => {
            const li = document.createElement('li');
            li.textContent = article;
            li.className = 'path-item';

            // Highlight start and end articles
            const endArticle = endInput.value.trim();
            if (index === 0 || article === endArticle) {
                li.classList.add('start-end');
            }
            pathList.appendChild(li);

            if (index < path.length - 1) {
                const arrow = document.createElement('div');
                arrow.className = 'path-arrow material-symbols-outlined';
                arrow.textContent = 'arrow_forward';
                pathList.appendChild(arrow);
            }
        });

        resultsContainer.classList.remove('hidden');
    }
});
