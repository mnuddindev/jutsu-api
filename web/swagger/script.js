// JsonSchemaViewer Component - Fixed
class JsonSchemaViewer {
    constructor(containerElement) {
        this.container = containerElement;
        this.endpointId = containerElement.getAttribute('data-endpoint');
        this.init();
    }

    init() {
        this.extractJsonFromExample();
    }

    extractJsonFromExample() {
        try {
            console.log('Looking for endpoint:', this.endpointId);

            // Find the endpoint card - FIXED: Use the actual ID format
            const endpointCard = document.getElementById(`endpoint-${this.endpointId}`);
            console.log('Found endpoint card:', endpointCard);

            if (!endpointCard) {
                console.error('Endpoint card not found for:', `endpoint-${this.endpointId}`);
                this.renderError('Endpoint not found');
                return;
            }

            // Find the code block with JSON
            const codeBlock = endpointCard.querySelector('.code-block code.language-json');
            console.log('Found code block:', codeBlock);

            if (!codeBlock) {
                this.renderError('No example JSON found');
                return;
            }

            const jsonText = codeBlock.textContent.trim();
            console.log('Raw JSON text:', jsonText);

            // Use the JSON directly since it's already valid
            const jsonData = JSON.parse(jsonText);
            console.log('Parsed JSON successfully:', jsonData);

            this.render(jsonData);

        } catch (error) {
            console.error('Error processing JSON:', error);
            // Use fallback data
            const fallbackData = this.createValidExampleJson();
            this.render(fallbackData);
        }
    }

    createValidExampleJson() {
        return {
            "success": true,
            "results": {
                "spotlights": [
                    {
                        "id": "frieren-beyond-journeys-end-18542",
                        "data_id": "string",
                        "title": "Frieren: Beyond Journey's End",
                        "japanese_title": "string",
                        "poster": "http://example.com",
                        "description": "string",
                        "tvInfo": {
                            "showType": "string",
                            "duration": "string",
                            "releaseDate": "string",
                            "quality": "string",
                            "episodeInfo": []
                        }
                    }
                ],
                "trending": [
                    {
                        "id": "string",
                        "data_id": "string",
                        "title": "string",
                        "japanese_title": "string",
                        "poster": "http://example.com",
                        "duration": "string",
                        "type": "string",
                        "rating": "string",
                        "episodes": {
                            "sub": 0,
                            "dub": 0
                        }
                    }
                ],
                "topTen": {
                    "today": [
                        {
                            "id": "string",
                            "data_id": "string",
                            "number": 0,
                            "name": "string",
                            "poster": "http://example.com",
                            "tvInfo": {}
                        }
                    ],
                    "week": [
                        {
                            "id": "string",
                            "data_id": "string",
                            "number": 0,
                            "name": "string",
                            "poster": "http://example.com",
                            "tvInfo": {}
                        }
                    ],
                    "month": [
                        {
                            "id": "string",
                            "data_id": "string",
                            "number": 0,
                            "name": "string",
                            "poster": "http://example.com",
                            "tvInfo": {}
                        }
                    ]
                },
                "today": {
                    "schedule": [
                        {
                            "id": "string",
                            "data_id": "string",
                            "title": "string",
                            "japanese_title": "string",
                            "releaseDate": "string",
                            "time": "string",
                            "episode_no": 0
                        }
                    ]
                },
                "genres": ["string"]
            }
        };
    }

    render(jsonData) {
        console.log('Rendering schema with data:', jsonData);

        const schemaHTML = this.generateSchemaHTML(jsonData);
        this.container.innerHTML = `
            <div class="json-schema-viewer">
                ${schemaHTML}
            </div>
        `;

        console.log('Schema rendered successfully');
        this.attachEventListeners();
    }

    renderError(message) {
        console.error('Rendering error:', message);
        this.container.innerHTML = `
            <div class="json-schema-viewer error">
                <div class="schema-error">${message}</div>
            </div>
        `;
    }

    generateSchemaHTML(data, level = 0) {
        if (!data || typeof data !== 'object') {
            return '<div class="schema-item">Invalid data</div>';
        }

        let html = '';
        const entries = Object.entries(data);

        for (const [key, value] of entries) {
            const type = this.getType(value);
            const example = this.getExample(value, key);
            const isExpandable = (type === 'object' || type === 'array') && value !== null;

            html += `
                <div class="schema-item ${isExpandable ? 'expandable' : ''}" data-level="${level}">
                    <div class="schema-field" data-key="${key}">
                        ${isExpandable ? '<div class="expand-icon">▶</div>' : '<div class="expand-spacer"></div>'}
                        <div class="field-info">
                            <div class="field-name-type">
                                <span class="field-name">${key}</span>
                                <span class="field-type">${this.formatType(type, value)}</span>
                            </div>
                            ${example ? `<div class="field-example">Example: ${example}</div>` : ''}
                        </div>
                    </div>
                    ${isExpandable ? this.generateChildrenHTML(value, level + 1) : ''}
                </div>
            `;
        }

        return html;
    }

    generateChildrenHTML(value, level) {
        const type = this.getType(value);

        if (type === 'array') {
            if (value.length > 0 && typeof value[0] === 'object' && value[0] !== null) {
                return `
                    <div class="schema-children hidden">
                        <div class="schema-item">
                            <div class="schema-field">
                                <div class="expand-spacer"></div>
                                <div class="field-info">
                                    <div class="field-name-type">
                                        <span class="field-name">[array item]</span>
                                        <span class="field-type">object</span>
                                    </div>
                                </div>
                            </div>
                            <div class="schema-children">
                                ${this.generateSchemaHTML(value[0], level)}
                            </div>
                        </div>
                    </div>
                `;
            } else {
                const itemType = value.length > 0 ? this.getType(value[0]) : 'any';
                return `
                    <div class="schema-children hidden">
                        <div class="schema-item">
                            <div class="schema-field">
                                <div class="expand-spacer"></div>
                                <div class="field-info">
                                    <div class="field-name-type">
                                        <span class="field-name">[array items]</span>
                                        <span class="field-type">${itemType}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                `;
            }
        } else if (type === 'object') {
            return `
                <div class="schema-children hidden">
                    ${this.generateSchemaHTML(value, level)}
                </div>
            `;
        }

        return '';
    }

    getType(value) {
        if (value === null) return 'null';
        if (Array.isArray(value)) return 'array';
        return typeof value;
    }

    formatType(type, value) {
        if (type === 'array') {
            const itemType = value.length > 0 ? this.getType(value[0]) : 'any';
            return `array[${itemType}]`;
        }
        return type;
    }

    getExample(value, key) {
        const type = this.getType(value);

        const specialExamples = {
            'id': 'frieren-beyond-journeys-end-18542',
            'data_id': '18542',
            'title': 'Frieren: Beyond Journey\'s End',
            'japanese_title': '葬送のフリーレン',
            'poster': 'https://example.com/image.jpg',
            'description': 'The story of an elf mage...',
            'showType': 'TV',
            'duration': '24m',
            'releaseDate': '2024-01-01',
            'quality': 'HD',
            'type': 'TV',
            'rating': '4.8',
            'name': 'Episode 1',
            'time': '12:00',
            'episode_no': 1,
            'number': 1,
            'sub': 12,
            'dub': 12
        };

        if (specialExamples[key]) {
            return specialExamples[key];
        }

        switch (type) {
            case 'string':
                return 'text';
            case 'number':
                return '123';
            case 'boolean':
                return value.toString();
            case 'array':
                return `[${value.length} items]`;
            case 'object':
                return '{...}';
            default:
                return String(value);
        }
    }

    attachEventListeners() {
        this.container.addEventListener('click', (e) => {
            const field = e.target.closest('.schema-field');
            if (field) {
                const schemaItem = field.closest('.expandable');
                if (schemaItem) {
                    const children = schemaItem.querySelector('.schema-children');
                    if (children) {
                        schemaItem.classList.toggle('expanded');
                        children.classList.toggle('hidden');
                    }
                }
            }
        });
    }
}

document.addEventListener('DOMContentLoaded', () => {
    hljs.highlightAll();
    const searchInput = document.getElementById('sidebar-search');
    if (searchInput) {
        searchInput.addEventListener('input', (event) => {
            const term = event.target.value.toLowerCase();
            document.querySelectorAll('.docs-sidebar a').forEach((link) => {
                const text = link.textContent.toLowerCase();
                const match = text.includes(term);
                link.parentElement.style.display = match ? 'block' : 'none';
            });
        });
    }

    document.querySelectorAll('.copy-btn').forEach((btn) => {
        btn.addEventListener('click', () => {
            const target = document.getElementById(btn.dataset.target);
            if (!target) return;
            navigator.clipboard.writeText(target.textContent.trim()).then(() => {
                btn.textContent = 'Copied';
                setTimeout(() => (btn.textContent = 'Copy'), 1500);
            });
        });
    });

    // Json Schema Viewer Initialization
    console.log('DOM loaded, initializing JsonSchemaViewer...');

    const containers = document.querySelectorAll('.json-schema-viewer-container');
    console.log('Found containers:', containers.length);

    containers.forEach(container => {
        console.log('Initializing container:', container);
        new JsonSchemaViewer(container);
    });

    // Send Request functionality
    document.querySelectorAll('.send-request-btn').forEach(button => {
        button.addEventListener('click', async function () {
            const endpoint = this.getAttribute('data-endpoint');
            const endpointCard = this.closest('.endpoint-card');
            const responseCard = endpointCard.querySelector('.response-card');
            const responseOutput = endpointCard.querySelector('.response-output');
            const responseStatus = endpointCard.querySelector('.response-status');
            const responseContent = endpointCard.querySelector('.response-content');

            this.classList.add('loading');
            this.textContent = 'Loading...';

            try {
                const response = await fetch(`http://localhost:8080${endpoint}`);
                const data = await response.json();

                responseStatus.textContent = `${response.status} ${response.statusText}`;
                responseStatus.className = 'response-status';
                responseOutput.textContent = JSON.stringify(data, null, 2);

                // Make sure the card is visible
                responseCard.classList.remove('hidden');
                responseStatus.classList.remove('hidden');
                responseContent.classList.remove('hidden');

                // Update collapse button to minus since content is visible
                const collapseBtn = responseCard.querySelector('.collapse-btn');
                collapseBtn.textContent = '−';

                // Remove the highlighted flag and re-highlight
                delete responseOutput.dataset.highlighted;
                setTimeout(() => {
                    hljs.highlightElement(responseOutput);
                }, 50);

            } catch (error) {
                responseStatus.textContent = `Error: ${error.message}`;
                responseStatus.className = 'response-status error';
                responseOutput.textContent = JSON.stringify({
                    error: 'Failed to fetch data',
                    message: error.message
                }, null, 2);

                // Make sure the card is visible
                responseCard.classList.remove('hidden');
                responseStatus.classList.remove('hidden');
                responseContent.classList.remove('hidden');

                // Update collapse button to minus since content is visible
                const collapseBtn = responseCard.querySelector('.collapse-btn');
                collapseBtn.textContent = '−';

                // Remove the highlighted flag and re-highlight
                delete responseOutput.dataset.highlighted;
                setTimeout(() => {
                    hljs.highlightElement(responseOutput);
                }, 50);
            } finally {
                this.classList.remove('loading');
                this.textContent = 'Send Request';
            }
        });
    });

    // Collapse/Expand response cards
    document.querySelectorAll('.collapse-btn').forEach(button => {
        button.addEventListener('click', function () {
            const responseCard = this.closest('.response-card');
            const responseStatus = responseCard.querySelector('.response-status');
            const responseContent = responseCard.querySelector('.response-content');

            // Check if content is currently visible
            const isContentVisible = !responseContent.classList.contains('hidden');

            if (isContentVisible) {
                // Hide status and content
                responseStatus.classList.add('hidden');
                responseContent.classList.add('hidden');
                this.textContent = '+';
            } else {
                // Show status and content
                responseStatus.classList.remove('hidden');
                responseContent.classList.remove('hidden');
                this.textContent = '−';
            }
        });
    });
});