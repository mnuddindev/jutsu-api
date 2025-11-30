// JsonSchemaViewer Component
class JsonSchemaViewer {
    constructor(containerElement) {
        this.container = $(containerElement);
        this.endpointId = containerElement.getAttribute('data-endpoint');
        this.init();
    }

    init() {
        this.extractJsonFromExample();
    }

    extractJsonFromExample() {
        try {
            console.log('Looking for endpoint with data-endpoint:', this.endpointId);

            // FIXED: Find the parent endpoint card directly
            const endpointCard = this.container.closest('.endpoint-card')[0];
            console.log('Found endpoint card:', endpointCard);

            if (!endpointCard) {
                console.error('Endpoint card not found for container');
                this.renderError('Endpoint not found');
                return;
            }

            // Find the code block with JSON in the examples-panel
            const codeBlock = $(endpointCard).find('.examples-panel .code-block code.language-json')[0];
            console.log('Found code block:', codeBlock);

            if (!codeBlock) {
                this.renderError('No example JSON found');
                return;
            }

            const jsonText = $(codeBlock).text().trim();
            console.log('Raw JSON text:', jsonText);

            // Use the JSON directly since it's already valid
            const jsonData = JSON.parse(jsonText);
            console.log('Parsed JSON successfully:', jsonData);

            this.render(jsonData);

        } catch (error) {
            console.error('Error processing JSON:', error);
            // REMOVED: Don't use fallback data
            this.renderError('Failed to parse JSON example: ' + error.message);
        }
    }

    // REMOVED: createValidExampleJson() method

    render(jsonData) {
        console.log('Rendering schema with data:', jsonData);

        const schemaHTML = this.generateSchemaHTML(jsonData);
        this.container.html(`
            <div class="json-schema-viewer">
                ${schemaHTML}
            </div>
        `);

        console.log('Schema rendered successfully');
        this.attachEventListeners();
    }

    renderError(message) {
        console.error('Rendering error:', message);
        this.container.html(`
            <div class="json-schema-viewer error">
                <div class="schema-error">${message}</div>
            </div>
        `);
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

        // Use ACTUAL values from the parsed JSON
        switch (type) {
            case 'string':
                // Show actual string value, truncate if too long
                return value.length > 50 ? value.substring(0, 50) + '...' : value;
            case 'number':
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
        this.container.on('click', '.schema-field', (e) => {
            e.stopPropagation();
            const $field = $(e.currentTarget);
            const $schemaItem = $field.closest('.schema-item.expandable');

            if ($schemaItem.length) {
                const $children = $schemaItem.find('> .schema-children');
                if ($children.length) {
                    $schemaItem.toggleClass('expanded');
                    $children.toggleClass('hidden');
                }
            }
        });
    }
}

$(document).ready(function () {
    hljs.highlightAll();

    const $searchInput = $('#sidebar-search');
    if ($searchInput.length) {
        $searchInput.on('input', function () {
            const term = $(this).val().toLowerCase();
            $('.docs-sidebar a').each(function () {
                const text = $(this).text().toLowerCase();
                const match = text.includes(term);
                $(this).parent().toggle(match);
            });
        });
    }

    $('.copy-btn').each(function () {
        $(this).on('click', function () {
            const $target = $('#' + $(this).data('target'));
            if (!$target.length) return;

            navigator.clipboard.writeText($target.text().trim()).then(() => {
                const $btn = $(this);
                $btn.text('Copied');
                setTimeout(() => $btn.text('Copy'), 1500);
            });
        });
    });

    // Json Schema Viewer Initialization
    console.log('DOM loaded, initializing JsonSchemaViewer...');

    const $containers = $('.json-schema-viewer-container');
    console.log('Found containers:', $containers.length);

    $containers.each(function () {
        console.log('Initializing container:', this);
        new JsonSchemaViewer(this);
    });

    // Send Request functionality
    $('.send-request-btn').each(function () {
        $(this).on('click', async function () {
            const endpoint = $(this).data('endpoint');
            const $endpointCard = $(this).closest('.endpoint-card');
            const $responseCard = $endpointCard.find('.response-card');
            const $responseOutput = $endpointCard.find('.response-output');
            const $responseStatus = $endpointCard.find('.response-status');
            const $responseContent = $endpointCard.find('.response-content');

            const $button = $(this);
            $button.addClass('loading');
            $button.text('Loading...');

            try {
                const response = await fetch(`http://localhost:8080${endpoint}`);
                console.log(`http://localhost:8080${endpoint}`)
                const data = await response.json();

                $responseStatus.text(`${response.status} ${response.statusText}`);
                $responseStatus.removeClass().addClass('response-status');
                $responseOutput.text(JSON.stringify(data, null, 2));

                // Make sure the card is visible
                $responseCard.removeClass('hidden');
                $responseStatus.removeClass('hidden');
                $responseContent.removeClass('hidden');

                // Update collapse button to minus since content is visible
                const $collapseBtn = $responseCard.find('.collapse-btn');
                $collapseBtn.text('−');

                // Remove the highlighted flag and re-highlight
                delete $responseOutput[0].dataset.highlighted;
                setTimeout(() => {
                    hljs.highlightElement($responseOutput[0]);
                }, 50);

            } catch (error) {
                $responseStatus.text(`Error: ${error.message}`);
                $responseStatus.removeClass().addClass('response-status error');
                $responseOutput.text(JSON.stringify({
                    error: 'Failed to fetch data',
                    message: error.message
                }, null, 2));

                // Make sure the card is visible
                $responseCard.removeClass('hidden');
                $responseStatus.removeClass('hidden');
                $responseContent.removeClass('hidden');

                // Update collapse button to minus since content is visible
                const $collapseBtn = $responseCard.find('.collapse-btn');
                $collapseBtn.text('−');

                // Remove the highlighted flag and re-highlight
                delete $responseOutput[0].dataset.highlighted;
                setTimeout(() => {
                    hljs.highlightElement($responseOutput[0]);
                }, 50);
            } finally {
                $button.removeClass('loading');
                $button.text('Send Request');
            }
        });
    });

    // Collapse/Expand response cards
    $('.collapse-btn').each(function () {
        $(this).on('click', function () {
            const $responseCard = $(this).closest('.response-card');
            const $responseStatus = $responseCard.find('.response-status');
            const $responseContent = $responseCard.find('.response-content');

            // Check if content is currently visible
            const isContentVisible = !$responseContent.hasClass('hidden');

            if (isContentVisible) {
                // Hide status and content
                $responseStatus.addClass('hidden');
                $responseContent.addClass('hidden');
                $(this).text('+');
            } else {
                // Show status and content
                $responseStatus.removeClass('hidden');
                $responseContent.removeClass('hidden');
                $(this).text('−');
            }
        });
    });
});