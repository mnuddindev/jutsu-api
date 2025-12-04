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
        this.container.html(`<div class="json-schema-viewer">${schemaHTML}</div>`);
        this.attachEventListeners();
    }

    renderError(message) {
        this.container.html(`<div class="json-schema-viewer error"><div class="schema-error">${message}</div></div>`);
    }

    generateSchemaHTML(data, level = 0) {
        if (!data || typeof data !== 'object') return '<div class="schema-item">Invalid data</div>';

        let html = '';
        for (const [key, value] of Object.entries(data)) {
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
            if (value.length > 0 && typeof value[0] === 'object') {
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
                            <div class="schema-children">${this.generateSchemaHTML(value[0], level)}</div>
                        </div>
                    </div>
                `;
            }
        } else if (type === 'object') {
            return `<div class="schema-children hidden">${this.generateSchemaHTML(value, level)}</div>`;
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
            case 'string': return value.length > 50 ? value.substring(0, 50) + '...' : value;
            case 'number': case 'boolean': return value.toString();
            case 'array': return `[${value.length} items]`;
            case 'object': return '{...}';
            default: return String(value);
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
        const endpointConfigs = {
            '/api/info': {
                additionalParams: [
                    { name: 'fields', example: 'title,poster,description', required: false, description: 'Comma-separated list of fields to include in response (optional)' }
                ]
            },
            '/api/filter': {
                additionalParams: [
                    { name: 'type', example: '1-6', required: false, description: 'Filter by anime type: tv, movie, ova, ona, special, music (optional)' },
                    { name: 'status', example: '1-3', required: false, description: 'Filter by status: Finished airing, currently airing, not yet aired (optional)' },
                    { name: 'rated', example: '1-6', required: false, description: 'Rating of anime (e.g., G, PG, PG-13, R, R+, Rx etc.) (optional)' },
                    { name: 'score', example: '1-10', required: false, description: 'Score rating (e.g., 1 to 10) (optional)' },
                    { name: 'season', example: '1-4', required: false, description: 'Filter by season: winter, spring, summer, fall (optional)' },
                    { name: 'language', example: '1-3', required: false, description: 'Language of anime (e.g., sub, dub, sub-dub) (optional)' },
                    { name: 'genres', example: '1,2,3-32', required: false, description: 'Comma-separated list of genres (e.g., action, comedy) (optional)' },
                    { name: 'sort', example: 'popularity', required: false, description: 'Sort by: popularity, rating, latest, title (optional)' },
                    { name: 'page', example: '1', required: false, description: 'Page number for pagination (optional)' },
                    { name: 'sy', example: '2024', required: false, description: 'Start year (optional)' },
                    { name: 'sm', example: '2024', required: false, description: 'Start month (optional)' },
                    { name: 'sd', example: '2024', required: false, description: 'Start day (optional)' },
                    { name: 'ey', example: '2024', required: false, description: 'End year (optional)' },
                    { name: 'em', example: '2024', required: false, description: 'End month (optional)' },
                    { name: 'ed', example: '2024', required: false, description: 'End day (optional)' },
                    { name: 'keyword', example: 'one piece', required: false, description: 'Search Keyword' },
                ]
            },
            '/api/stream/:id': {
                additionalParams: [
                    { name: 'server', example: 'hd-1', required: false, description: 'Preferred streaming server (optional)' },
                    { name: 'type', example: 'sub', required: false, description: 'Content type: sub or dub (optional)' },
                    { name: 'quality', example: '1080p', required: false, description: 'Video quality preference (optional)' }
                ]
            },
            '/api/episodes/:id': {
                additionalParams: [
                    { name: 'page', example: '1', required: false, description: 'Page number for episode list (optional)' }
                ]
            },
            '/api/top-airing': { additionalParams: [{ name: 'page', example: '1', required: false, description: 'Page number (optional)' }] },
            '/api/most-popular': { additionalParams: [{ name: 'page', example: '1', required: false, description: 'Page number (optional)' }] },
            '/api/genre/:slug': { additionalParams: [{ name: 'page', example: '1', required: false, description: 'Page number (optional)' }] },
            '/api/producer/:id': { additionalParams: [{ name: 'page', example: '1', required: false, description: 'Page number (optional)' }] },
            '/api/studio/:id': { additionalParams: [{ name: 'page', example: '1', required: false, description: 'Page number (optional)' }] },
            '/api/character/list/:id': { additionalParams: [{ name: 'page', example: '1', required: false, description: 'Page number for character list (optional)' }] },
            '/api/watchlist/:userId': { additionalParams: [{ name: 'page', example: '1', required: false, description: 'Page number for watchlist (optional)' }] }
        };

        for (const [configPath, config] of Object.entries(endpointConfigs)) {
            const configPattern = configPath.replace(/:[^/]+/g, '([^/]+)').replace(/\{([^}]+)\}/g, '([^/]+)');
            if (new RegExp(`^${configPattern}$`).test(endpointPath)) {
                return config.additionalParams || [];
            }
        }
        return [];
    }
}

// Parameter Input Generator
class ParameterInputs {
    static generateInputs($endpointCard) {
        const $curlCode = $endpointCard.find('.examples-panel .code-block code.language-bash');
        const endpointPath = $endpointCard.find('.path').text().trim();

        let params = [];
        if ($curlCode.length > 0) {
            const curlExample = $curlCode.text().trim();
            params = this.extractParametersFromCurl(curlExample, $endpointCard);
        }

        const additionalParams = ParameterConfig.getEndpointParameters(endpointPath);
        const uniqueParams = [...params, ...additionalParams].filter((param, index, self) =>
            index === self.findIndex(p => p.name === param.name)
        );

        if (uniqueParams.length === 0) return '';

        let inputsHTML = `
            <div class="parameter-inputs">
                <div class="parameter-header">
                    <h4>Request Parameters</h4>
                    <div class="parameter-subtitle">${this.getParameterSummary(uniqueParams)}</div>
                </div>
        `;

        uniqueParams.forEach(param => {
            const isAdditional = additionalParams.some(p => p.name === param.name);
            // For optional parameters, use placeholder instead of value
            const inputValue = param.required ? (param.example || '') : '';
            const placeholder = param.required ? (param.example || 'Enter value') : (param.example || 'Enter value');

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
                           placeholder="${placeholder}"
                           value="${inputValue}"
                           ${param.required ? 'required' : ''}>
                    ${param.example ? `<div class="parameter-example">Example: ${param.example}</div>` : ''}
                    ${param.description ? `<div class="parameter-description">${param.description}</div>` : ''}
                </div>
            `;
        });

        return inputsHTML + '</div>';
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

        const urlMatch = curlExample.match(/--url '([^']+)'/) || curlExample.match(/curl '([^']+)'/) || curlExample.match(/curl "([^"]+)"/);
        if (urlMatch) url = urlMatch[1];
        if (!url) return params;

        try {
            const urlObj = new URL(url);
            const endpointPath = $endpointCard.find('.path').text().trim();

            // Extract path parameters
            const pathSegments = endpointPath.split('/');
            const actualSegments = urlObj.pathname.split('/');
            for (let i = 0; i < pathSegments.length; i++) {
                const segment = pathSegments[i];
                if ((segment.startsWith(':') || (segment.startsWith('{') && segment.endsWith('}'))) && i < actualSegments.length) {
                    const paramName = segment.replace(/[:{}]/g, '');
                    params.push({
                        name: paramName,
                        example: actualSegments[i],
                        required: true,
                        description: `Path parameter - ${paramName}`,
                        type: 'path'
                    });
                }
            }

            // Extract query parameters
            for (const [name, value] of new URLSearchParams(urlObj.search)) {
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

    static isParameterRequired(name, exampleValue, endpointPath) {
        const alwaysRequired = ['id', 'ep', 'episode', 'slug', 'key'];
        const usuallyOptional = ['page', 'limit', 'offset', 'sort', 'order', 'per_page', 'server', 'type', 'quality', 'sub', 'dub', 'language', 'genre', 'status', 'year', 'season', 'rating'];

        if (alwaysRequired.includes(name.toLowerCase())) return true;
        if (usuallyOptional.includes(name.toLowerCase())) return false;
        if (exampleValue && exampleValue.includes('...')) return false;
        return true;
    }

    static getParameterDescription(name) {
        const descriptions = {
            'id': 'Unique identifier for the resource', 'slug': 'URL-friendly identifier',
            'ep': 'Episode number or identifier', 'episode': 'Episode number or identifier',
            'page': 'Page number for pagination (optional)', 'limit': 'Number of items per page (optional)',
            'offset': 'Number of items to skip (optional)', 'sort': 'Field to sort by (optional)',
            'order': 'Sort order: asc or desc (optional)', 'server': 'Streaming server name',
            'type': 'Content type: sub or dub (optional)', 'quality': 'Video quality preference (optional)',
            'language': 'Language preference (optional)', 'genre': 'Filter by genre (optional)',
            'status': 'Filter by status: ongoing, completed (optional)', 'year': 'Filter by release year (optional)',
            'season': 'Filter by season (optional)', 'rating': 'Filter by minimum rating (optional)',
            'keyword': 'Search keyword'
        };
        return descriptions[name.toLowerCase()] || `Parameter for ${name}`;
    }

    static updateCurlExample($endpointCard) {
        const $curlCode = $endpointCard.find('.examples-panel .code-block code.language-bash');
        if ($curlCode.length === 0) return;

        // Get original curl to extract the base structure
        const originalCurl = $curlCode.text().trim();

        let baseUrl = '';
        const urlMatch = originalCurl.match(/--url '([^']+)'/) || originalCurl.match(/curl '([^']+)'/) || originalCurl.match(/curl "([^"]+)"/);
        if (urlMatch) baseUrl = urlMatch[1];
        if (!baseUrl) return;

        const originalUrl = new URL(baseUrl);
        let finalPath = originalUrl.pathname;
        const params = new URLSearchParams(originalUrl.search);

        // Get path parameter names from the endpoint path
        const pathParamNames = [];
        const endpointPath = $endpointCard.find('.path').text().trim();
        const pathSegments = endpointPath.split('/');
        pathSegments.forEach(segment => {
            if (segment.startsWith(':') || (segment.startsWith('{') && segment.endsWith('}'))) {
                pathParamNames.push(segment.replace(/[:{}]/g, ''));
            }
        });

        // Remove query parameters that duplicate path parameters
        pathParamNames.forEach(paramName => {
            if (params.has(paramName)) {
                params.delete(paramName);
            }
        });

        // Update ONLY query parameters, leave path parameters as-is
        $endpointCard.find('.parameter-input').each(function () {
            const $input = $(this);
            const paramName = $input.data('param');
            const paramValue = $input.val().trim();

            // Only update query parameters, skip path parameters
            if (paramValue && !pathParamNames.includes(paramName)) {
                params.set(paramName, paramValue);
            } else if (params.has(paramName) && !$input.closest('.parameter-field').hasClass('required')) {
                params.delete(paramName);
            }
        });

        // Build the new URL (path remains unchanged)
        let url = `http://localhost:8080${finalPath}`;
        const queryString = params.toString();
        if (queryString) url += '?' + queryString;

        // Build the colored curl command
        let newCurl;
        if (originalCurl.includes("--url ")) {
            newCurl = `curl --url '<span class="hljs-string">${url}</span>'`;
        } else {
            newCurl = `curl '<span class="hljs-string">${url}</span>'`;
        }

        // Add other curl options if they exist in original
        if (originalCurl.includes('--header')) {
            newCurl += ` --header '<span class="hljs-string">Accept: application/json</span>'`;
        }

        $curlCode.html(newCurl);
    }
}

// Real-time curl updates
$(document).on('input', '.parameter-input', function () {
    const $endpointCard = $(this).closest('.endpoint-card');
    $(this).css('border-color', '');
    ParameterInputs.updateCurlExample($endpointCard);
});

function headerHideHandler() {
    let lastScrollTop = 0;
    const $header = $('.docs-header');

    $header.css({
        'transition': 'all 0.3s ease-in-out',
        'position': 'sticky',
        'top': '0'
    });

    $(window).scroll(function () {
        const scrollTop = $(this).scrollTop();

        if (scrollTop > lastScrollTop && scrollTop > 50) {
            $header.css({
                'opacity': '0',
                'transform': 'translateY(-20px)',
                'pointer-events': 'none'
            });
        } else {
            $header.css({
                'opacity': '1',
                'transform': 'translateY(0)',
                'pointer-events': 'auto'
            });
        }

        lastScrollTop = scrollTop;
    });
}

// Main jQuery code
$(document).ready(function () {
    hljs.highlightAll();
    headerHideHandler();

    $('.endpoint-meta .endpoint-path .path').each(function () {
        $('<span class="server">http://localhost:8080</span>').insertBefore($(this));
    });


    // Generate parameter inputs
    $('.endpoint-card').each(function () {
        const $endpointCard = $(this);
        const $sendButton = $endpointCard.find('.send-request-btn');
        const parameterInputsHTML = ParameterInputs.generateInputs($endpointCard);
        if (parameterInputsHTML) $sendButton.before(parameterInputsHTML);
    });

    // Initialize components
    $('.json-schema-viewer-container').each(function () { new JsonSchemaViewer(this); });

    // Update curl examples after parameters are generated
    setTimeout(() => {
        $('.endpoint-card').each(function () {
            if ($(this).find('.parameter-input').length > 0) {
                ParameterInputs.updateCurlExample($(this));
            }
        });
    }, 100);

    // Search functionality
    $('#sidebar-search').on('input', function () {
        const term = $(this).val().toLowerCase();
        $('.docs-sidebar a').each(function () {
            const text = $(this).text().toLowerCase();
            $(this).parent().toggle(text.includes(term));
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

    // Send Request functionality
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
            $responseStatus.text('Error: Missing required parameters').removeClass().addClass('response-status error');
            $responseOutput.text(JSON.stringify({ error: 'Validation failed', message: 'Please fill in all required parameters' }, null, 2));
            $responseCard.removeClass('hidden');
            $responseStatus.removeClass('hidden');
            $responseContent.removeClass('hidden');
            return;
        }

        $button.addClass('loading').text('Loading...');

        try {
            let finalEndpoint = endpoint;
            const params = new URLSearchParams();

            $endpointCard.find('.parameter-input').each(function () {
                const $input = $(this);
                const paramName = $input.data('param');
                const paramValue = $input.val().trim();
                if (paramValue) {
                    if (finalEndpoint.includes(`:${paramName}`) || finalEndpoint.includes(`{${paramName}}`)) {
                        finalEndpoint = finalEndpoint.replace(`:${paramName}`, paramValue).replace(`{${paramName}}`, paramValue);
                    } else {
                        params.append(paramName, paramValue);
                    }
                }
            });

            let url = `http://localhost:8080${finalEndpoint}`;
            const queryString = params.toString();
            if (queryString) url += '?' + queryString;

            const response = await fetch(url, { headers: { 'Accept': 'application/json' } });
            const data = await response.json();

            $responseStatus.text(`${response.status} ${response.statusText}`).removeClass().addClass('response-status');
            $responseOutput.text(JSON.stringify(data, null, 2));
            $responseCard.removeClass('hidden');
            $responseStatus.removeClass('hidden');
            $responseContent.removeClass('hidden');
            $responseCard.find('.collapse-btn').text('−');

            delete $responseOutput[0].dataset.highlighted;
            setTimeout(() => hljs.highlightElement($responseOutput[0]), 50);

        } catch (error) {
            $responseStatus.text(`Error: ${error.message}`).removeClass().addClass('response-status error');
            $responseOutput.text(JSON.stringify({ error: 'Failed to fetch data', message: error.message }, null, 2));
            $responseCard.removeClass('hidden');
            $responseStatus.removeClass('hidden');
            $responseContent.removeClass('hidden');
            $responseCard.find('.collapse-btn').text('−');
            delete $responseOutput[0].dataset.highlighted;
            setTimeout(() => hljs.highlightElement($responseOutput[0]), 50);
        } finally {
            $button.removeClass('loading').text('Send Request');
        }
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