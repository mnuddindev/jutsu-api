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
            const endpointCard = this.container.closest('.endpoint-card')[0];
            if (!endpointCard) {
                this.renderError('Endpoint not found');
                return;
            }

            const codeBlock = $(endpointCard).find('.examples-panel .code-block code.language-json')[0];
            if (!codeBlock) {
                this.renderError('No example JSON found');
                return;
            }

            const jsonText = $(codeBlock).text().trim();
            const jsonData = JSON.parse(jsonText);
            this.render(jsonData);

        } catch (error) {
            this.renderError('Failed to parse JSON example: ' + error.message);
        }
    }

    render(jsonData) {
        const schemaHTML = this.generateSchemaHTML(jsonData);
        this.container.html(`
            <div class="json-schema-viewer">
                ${schemaHTML}
            </div>
        `);
        this.attachEventListeners();
    }

    renderError(message) {
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
        switch (type) {
            case 'string':
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

// Parameter Input Generator
class ParameterInputs {
    static generateInputs($endpointCard) {
        const $curlCode = $endpointCard.find('.examples-panel .code-block code.language-bash');
        if ($curlCode.length === 0) {
            return '';
        }

        const curlExample = $curlCode.text().trim();
        const params = this.extractParametersFromCurl(curlExample);

        if (params.length === 0) {
            return `
                
            `;
        }

        let inputsHTML = `
            <div class="parameter-inputs">
                <div class="parameter-header">
                    <h4>Request Parameters</h4>
                </div>
        `;

        params.forEach(param => {
            inputsHTML += `
                <div class="parameter-field ${param.required ? 'required' : 'optional'}">
                    <label for="param-${param.name}">
                        ${param.name}
                        <span class="parameter-badge ${param.required ? 'required' : 'optional'}">${param.required ? 'Required' : 'Optional'}</span>
                    </label>
                    <input type="text" 
                           id="param-${param.name}" 
                           class="parameter-input" 
                           data-param="${param.name}"
                           placeholder="${param.example || 'Enter value'}"
                           value="${param.example || ''}">
                    ${param.example ? `<div class="parameter-example">Example: ${param.example}</div>` : ''}
                </div>
            `;
        });

        inputsHTML += '</div>';
        return inputsHTML;
    }

    static extractParametersFromCurl(curlExample) {
        const params = [];

        let url = '';
        const urlMatch1 = curlExample.match(/--url '([^']+)'/);
        if (urlMatch1) {
            url = urlMatch1[1];
        } else {
            const urlMatch2 = curlExample.match(/curl '([^']+)'/);
            if (urlMatch2) {
                url = urlMatch2[1];
            } else {
                const urlMatch3 = curlExample.match(/curl "([^"]+)"/);
                if (urlMatch3) {
                    url = urlMatch3[1];
                }
            }
        }

        if (!url) {
            return params;
        }

        try {
            const urlObj = new URL(url);
            const urlParams = new URLSearchParams(urlObj.search);

            for (const [name, value] of urlParams) {
                params.push({
                    name: name,
                    example: value,
                    required: true,
                    description: `Parameter for ${name}`
                });
            }
        } catch (error) {
            console.log('Error parsing URL:', error);
        }

        return params;
    }
}

// Main jQuery code
$(document).ready(function () {
    // Initialize syntax highlighting
    hljs.highlightAll();

    // Generate parameter inputs for each endpoint
    $('.endpoint-card').each(function () {
        const $endpointCard = $(this);
        const $sendButton = $endpointCard.find('.send-request-btn');
        const parameterInputsHTML = ParameterInputs.generateInputs($endpointCard);

        if (parameterInputsHTML) {
            $sendButton.before(parameterInputsHTML);
        }
    });

    // Initialize JsonSchemaViewer
    $('.json-schema-viewer-container').each(function () {
        new JsonSchemaViewer(this);
    });

    // Search functionality
    $('#sidebar-search').on('input', function () {
        const term = $(this).val().toLowerCase();
        $('.docs-sidebar a').each(function () {
            const text = $(this).text().toLowerCase();
            const match = text.includes(term);
            $(this).parent().toggle(match);
        });
    });

    // Copy button functionality
    $('.copy-btn').on('click', function () {
        const $target = $('#' + $(this).data('target'));
        if (!$target.length) return;
        navigator.clipboard.writeText($target.text().trim()).then(() => {
            const $btn = $(this);
            $btn.text('Copied');
            setTimeout(() => $btn.text('Copy'), 1500);
        });
    });

    // Send Request functionality with parameters
    $(document).on('click', '.send-request-btn', async function () {
        const $button = $(this);
        const $endpointCard = $button.closest('.endpoint-card');
        const endpoint = $button.data('endpoint');
        const $responseCard = $endpointCard.find('.response-card');
        const $responseOutput = $endpointCard.find('.response-output');
        const $responseStatus = $endpointCard.find('.response-status');
        const $responseContent = $endpointCard.find('.response-content');

        // Validate required parameters
        let hasErrors = false;
        $endpointCard.find('.parameter-field.required .parameter-input').each(function () {
            const $input = $(this);
            if (!$input.val().trim()) {
                $input.css('border-color', '#ef4444');
                hasErrors = true;
            } else {
                $input.css('border-color', '');
            }
        });

        if (hasErrors) {
            $responseStatus.text('Error: Missing required parameters');
            $responseStatus.removeClass().addClass('response-status error');
            $responseOutput.text(JSON.stringify({
                error: 'Validation failed',
                message: 'Please fill in all required parameters'
            }, null, 2));

            $responseCard.removeClass('hidden');
            $responseStatus.removeClass('hidden');
            $responseContent.removeClass('hidden');
            return;
        }

        $button.addClass('loading');
        $button.text('Loading...');

        try {
            let url = `http://localhost:8080${endpoint}`;
            const params = new URLSearchParams();

            $endpointCard.find('.parameter-input').each(function () {
                const $input = $(this);
                const paramName = $input.data('param');
                const paramValue = $input.val().trim();

                if (paramValue) {
                    params.append(paramName, paramValue);
                }
            });

            const queryString = params.toString();
            if (queryString) {
                url += '?' + queryString;
            }

            const response = await fetch(url, {
                headers: {
                    'Accept': 'application/json'
                }
            });

            const data = await response.json();

            $responseStatus.text(`${response.status} ${response.statusText}`);
            $responseStatus.removeClass().addClass('response-status');
            $responseOutput.text(JSON.stringify(data, null, 2));

            $responseCard.removeClass('hidden');
            $responseStatus.removeClass('hidden');
            $responseContent.removeClass('hidden');

            const $collapseBtn = $responseCard.find('.collapse-btn');
            $collapseBtn.text('−');

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

            $responseCard.removeClass('hidden');
            $responseStatus.removeClass('hidden');
            $responseContent.removeClass('hidden');

            const $collapseBtn = $responseCard.find('.collapse-btn');
            $collapseBtn.text('−');

            delete $responseOutput[0].dataset.highlighted;
            setTimeout(() => {
                hljs.highlightElement($responseOutput[0]);
            }, 50);
        } finally {
            $button.removeClass('loading');
            $button.text('Send Request');
        }
    });

    // Clear validation errors when typing
    $(document).on('input', '.parameter-input', function () {
        $(this).css('border-color', '');
    });

    // Collapse/Expand response cards
    $('.collapse-btn').on('click', function () {
        const $responseCard = $(this).closest('.response-card');
        const $responseStatus = $responseCard.find('.response-status');
        const $responseContent = $responseCard.find('.response-content');
        const isContentVisible = !$responseContent.hasClass('hidden');

        if (isContentVisible) {
            $responseStatus.addClass('hidden');
            $responseContent.addClass('hidden');
            $(this).text('+');
        } else {
            $responseStatus.removeClass('hidden');
            $responseContent.removeClass('hidden');
            $(this).text('−');
        }
    });
});