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

// Extended Parameter Configuration
class ParameterConfig {
    static getEndpointParameters(endpointPath) {
        // Define additional parameters for each endpoint that aren't in curl examples
        const endpointConfigs = {
            '/api/info': {
                additionalParams: [
                    {
                        name: 'fields',
                        example: 'title,poster,description',
                        required: false,
                        description: 'Comma-separated list of fields to include in response (optional)'
                    }
                ]
            },
            '/api/search': {
                additionalParams: [
                    {
                        name: 'type',
                        example: 'tv',
                        required: false,
                        description: 'Filter by anime type: tv, movie, ova, ona, special (optional)'
                    },
                    {
                        name: 'status',
                        example: 'ongoing',
                        required: false,
                        description: 'Filter by status: ongoing, completed (optional)'
                    },
                    {
                        name: 'genre',
                        example: 'action,adventure',
                        required: false,
                        description: 'Comma-separated genre IDs (optional)'
                    },
                    {
                        name: 'year',
                        example: '2024',
                        required: false,
                        description: 'Filter by release year (optional)'
                    },
                    {
                        name: 'season',
                        example: 'winter',
                        required: false,
                        description: 'Filter by season: winter, spring, summer, fall (optional)'
                    },
                    {
                        name: 'sort',
                        example: 'popularity',
                        required: false,
                        description: 'Sort by: popularity, rating, latest, title (optional)'
                    },
                    {
                        name: 'page',
                        example: '1',
                        required: false,
                        description: 'Page number for pagination (optional)'
                    }
                ]
            },
            '/api/filter': {
                additionalParams: [
                    {
                        name: 'type',
                        example: 'tv',
                        required: false,
                        description: 'Filter by anime type (optional)'
                    },
                    {
                        name: 'status',
                        example: 'ongoing',
                        required: false,
                        description: 'Filter by status (optional)'
                    },
                    {
                        name: 'genre',
                        example: '1,2,3',
                        required: false,
                        description: 'Comma-separated genre IDs (optional)'
                    },
                    {
                        name: 'rating',
                        example: '7.5',
                        required: false,
                        description: 'Minimum rating score (optional)'
                    },
                    {
                        name: 'year',
                        example: '2024',
                        required: false,
                        description: 'Release year (optional)'
                    },
                    {
                        name: 'language',
                        example: 'sub',
                        required: false,
                        description: 'Language: sub or dub (optional)'
                    },
                    {
                        name: 'page',
                        example: '1',
                        required: false,
                        description: 'Page number (optional)'
                    }
                ]
            },
            '/api/stream/:id': {
                additionalParams: [
                    {
                        name: 'server',
                        example: 'hd-1',
                        required: false,
                        description: 'Preferred streaming server (optional)'
                    },
                    {
                        name: 'type',
                        example: 'sub',
                        required: false,
                        description: 'Content type: sub or dub (optional)'
                    },
                    {
                        name: 'quality',
                        example: '1080p',
                        required: false,
                        description: 'Video quality preference (optional)'
                    }
                ]
            },
            '/api/episodes/:id': {
                additionalParams: [
                    {
                        name: 'page',
                        example: '1',
                        required: false,
                        description: 'Page number for episode list (optional)'
                    }
                ]
            },
            '/api/schedule': {
                additionalParams: [
                    {
                        name: 'date',
                        example: '2024-01-18',
                        required: false,
                        description: 'Date in YYYY-MM-DD format (optional, defaults to today)'
                    },
                    {
                        name: 'tzOffset',
                        example: '-330',
                        required: false,
                        description: 'Timezone offset in minutes (optional)'
                    }
                ]
            },
            '/api/top-airing': {
                additionalParams: [
                    {
                        name: 'page',
                        example: '1',
                        required: false,
                        description: 'Page number (optional)'
                    }
                ]
            },
            '/api/most-popular': {
                additionalParams: [
                    {
                        name: 'page',
                        example: '1',
                        required: false,
                        description: 'Page number (optional)'
                    }
                ]
            },
            '/api/genre/:slug': {
                additionalParams: [
                    {
                        name: 'page',
                        example: '1',
                        required: false,
                        description: 'Page number (optional)'
                    }
                ]
            },
            '/api/producer/:id': {
                additionalParams: [
                    {
                        name: 'page',
                        example: '1',
                        required: false,
                        description: 'Page number (optional)'
                    }
                ]
            },
            '/api/studio/:id': {
                additionalParams: [
                    {
                        name: 'page',
                        example: '1',
                        required: false,
                        description: 'Page number (optional)'
                    }
                ]
            },
            '/api/character/list/:id': {
                additionalParams: [
                    {
                        name: 'page',
                        example: '1',
                        required: false,
                        description: 'Page number for character list (optional)'
                    }
                ]
            },
            '/api/watchlist/:userId': {
                additionalParams: [
                    {
                        name: 'page',
                        example: '1',
                        required: false,
                        description: 'Page number for watchlist (optional)'
                    }
                ]
            }
        };

        // Find matching endpoint configuration
        for (const [configPath, config] of Object.entries(endpointConfigs)) {
            if (this.pathsMatch(configPath, endpointPath)) {
                return config.additionalParams || [];
            }
        }

        return [];
    }

    static pathsMatch(configPath, endpointPath) {
        // Convert both paths to regex patterns to handle dynamic segments
        const configPattern = configPath.replace(/:[^/]+/g, '([^/]+)').replace(/\{([^}]+)\}/g, '([^/]+)');
        const regex = new RegExp(`^${configPattern}$`);
        return regex.test(endpointPath);
    }
}

// Updated Parameter Input Generator
class ParameterInputs {
    static generateInputs($endpointCard) {
        const $curlCode = $endpointCard.find('.examples-panel .code-block code.language-bash');
        const endpointPath = $endpointCard.find('.path').text().trim();

        let params = [];

        // Extract parameters from curl example
        if ($curlCode.length > 0) {
            const curlExample = $curlCode.text().trim();
            params = this.extractParametersFromCurl(curlExample, $endpointCard);
        }

        // Add additional configured parameters
        const additionalParams = ParameterConfig.getEndpointParameters(endpointPath);
        params = [...params, ...additionalParams];

        // Remove duplicates (in case curl already has some of the additional params)
        const uniqueParams = this.removeDuplicateParams(params);

        if (uniqueParams.length === 0) {
            return '';
        }

        let inputsHTML = `
            <div class="parameter-inputs">
                <div class="parameter-header">
                    <h4>Request Parameters</h4>
                    <div class="parameter-subtitle">
                        ${this.getParameterSummary(uniqueParams)}
                    </div>
                </div>
        `;

        uniqueParams.forEach(param => {
            const isAdditional = additionalParams.some(p => p.name === param.name && !params.some(ep => ep.name === p.name));

            inputsHTML += `
                <div class="parameter-field ${param.required ? 'required' : 'optional'} ${isAdditional ? 'additional' : ''}">
                    <label for="param-${param.name}">
                        ${param.name}
                        <span class="parameter-badge ${param.required ? 'required' : 'optional'}">${param.required ? 'Required' : 'Optional'}</span>
                        ${isAdditional ? '<span class="parameter-badge info">Documentation</span>' : ''}
                    </label>
                    <input type="text" 
                           id="param-${param.name}" 
                           class="parameter-input" 
                           data-param="${param.name}"
                           placeholder="${param.example || 'Enter value'}"
                           value="${param.example || ''}"
                           ${param.required ? 'required' : ''}>
                    ${param.example ? `<div class="parameter-example">Example: ${param.example}</div>` : ''}
                    ${param.description ? `<div class="parameter-description">${param.description}</div>` : ''}
                </div>
            `;
        });

        inputsHTML += '</div>';
        return inputsHTML;
    }

    static removeDuplicateParams(params) {
        const seen = new Set();
        return params.filter(param => {
            if (seen.has(param.name)) {
                return false;
            }
            seen.add(param.name);
            return true;
        });
    }

    static getParameterSummary(params) {
        const requiredCount = params.filter(p => p.required).length;
        const optionalCount = params.filter(p => !p.required).length;

        const parts = [];
        if (requiredCount > 0) parts.push(`${requiredCount} required`);
        if (optionalCount > 0) parts.push(`${optionalCount} optional`);

        return parts.join(', ') + ` parameter${params.length > 1 ? 's' : ''}`;
    }

    static extractParametersFromCurl(curlExample, $endpointCard) {
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

            // Extract path parameters from endpoint path
            const endpointPath = $endpointCard.find('.path').text().trim();
            const pathParams = this.extractPathParameters(endpointPath, urlObj.pathname);

            // Add path parameters first (they're always required)
            pathParams.forEach(param => {
                params.push({
                    name: param.name,
                    example: param.value,
                    required: true,
                    description: `Path parameter - ${param.name}`,
                    type: 'path'
                });
            });

            // Add query parameters with smart required detection
            for (const [name, value] of urlParams) {
                const isRequired = this.isParameterRequired(name, value, endpointPath);
                params.push({
                    name: name,
                    example: value,
                    required: isRequired,
                    description: this.getParameterDescription(name),
                    type: 'query'
                });
            }
        } catch (error) {
            console.log('Error parsing URL:', error);
        }

        return params;
    }

    static extractPathParameters(endpointPath, actualPath) {
        const params = [];

        // Check if endpoint path has placeholders like :id, {id}, etc.
        const pathSegments = endpointPath.split('/');
        const actualSegments = actualPath.split('/');

        for (let i = 0; i < pathSegments.length; i++) {
            const segment = pathSegments[i];
            if ((segment.startsWith(':') ||
                (segment.startsWith('{') && segment.endsWith('}'))) &&
                i < actualSegments.length) {

                const paramName = segment.replace(/[:{}]/g, '');
                const paramValue = actualSegments[i];

                params.push({
                    name: paramName,
                    value: paramValue,
                    required: true
                });
            }
        }

        return params;
    }

    static isParameterRequired(name, exampleValue, endpointPath) {
        // Core parameters that are always required
        const alwaysRequired = ['id', 'ep', 'episode', 'slug', 'key'];

        // Pagination/sorting parameters that are usually optional
        const usuallyOptional = ['page', 'limit', 'offset', 'sort', 'order', 'per_page'];

        // Streaming/media parameters that are usually optional
        const streamingOptional = ['server', 'type', 'quality', 'sub', 'dub', 'language'];

        // Filter parameters that are usually optional
        const filterOptional = ['genre', 'status', 'year', 'season', 'rating'];

        if (alwaysRequired.includes(name.toLowerCase())) {
            return true;
        }

        if (usuallyOptional.includes(name.toLowerCase()) ||
            streamingOptional.includes(name.toLowerCase()) ||
            filterOptional.includes(name.toLowerCase())) {
            return false;
        }

        // If the parameter appears in multiple endpoints with different values, it's likely optional
        if (exampleValue && exampleValue.includes('...')) {
            return false;
        }

        // Default to required for safety
        return true;
    }

    static getParameterDescription(name) {
        const descriptions = {
            // Path parameters
            'id': 'Unique identifier for the resource',
            'slug': 'URL-friendly identifier',

            // Required query parameters
            'ep': 'Episode number or identifier',
            'episode': 'Episode number or identifier',

            // Optional query parameters
            'page': 'Page number for pagination (optional)',
            'limit': 'Number of items per page (optional)',
            'offset': 'Number of items to skip (optional)',
            'sort': 'Field to sort by (optional)',
            'order': 'Sort order: asc or desc (optional)',
            'server': 'Streaming server name',
            'type': 'Content type: sub or dub (optional)',
            'quality': 'Video quality preference (optional)',
            'language': 'Language preference (optional)',
            'genre': 'Filter by genre (optional)',
            'status': 'Filter by status: ongoing, completed (optional)',
            'year': 'Filter by release year (optional)',
            'season': 'Filter by season (optional)',
            'rating': 'Filter by minimum rating (optional)',
            'keyword': 'Search keyword (optional)'
        };

        return descriptions[name.toLowerCase()] || `Parameter for ${name}`;
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
        let endpoint = $button.data('endpoint');
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
            // Replace path parameters in the endpoint
            let finalEndpoint = endpoint;
            const params = new URLSearchParams();

            $endpointCard.find('.parameter-input').each(function () {
                const $input = $(this);
                const paramName = $input.data('param');
                const paramValue = $input.val().trim();

                if (paramValue) {
                    // Check if this is a path parameter (replace in endpoint)
                    if (finalEndpoint.includes(`:${paramName}`) || finalEndpoint.includes(`{${paramName}}`)) {
                        finalEndpoint = finalEndpoint.replace(`:${paramName}`, paramValue)
                            .replace(`{${paramName}}`, paramValue);
                    } else {
                        // It's a query parameter
                        params.append(paramName, paramValue);
                    }
                }
            });

            // Build final URL
            let url = `http://localhost:8080${finalEndpoint}`;
            const queryString = params.toString();
            if (queryString) {
                url += '?' + queryString;
            }

            console.log('Making request to:', url);

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