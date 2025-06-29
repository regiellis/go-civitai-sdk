// Civitai API Tester - WebSocket-based Real-time Updates with Smart UI Updates

// Global state
let appState = {
    results: [],
    summary: { total: 0, passed: 0, failed: 0, running: 0 },
    testsStarted: false,
    isLoading: false,
    lastUpdated: 'Loading...',
    expandedItems: new Set(),
    ws: null,
    reconnectAttempts: 0,
    maxReconnectAttempts: 5,
    testElements: new Map() // Cache test elements to prevent re-rendering
};

// Initialize the application
function initApiTester() {
    console.log('🚀 Civitai API Tester initialized (WebSocket + Smart Updates)');
    connectWebSocket();
    loadInitialData();
}

// WebSocket connection
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    console.log('🔌 Connecting to WebSocket:', wsUrl);
    
    appState.ws = new WebSocket(wsUrl);
    
    appState.ws.onopen = function() {
        console.log('✅ WebSocket connected');
        appState.reconnectAttempts = 0;
        document.getElementById('last-updated').textContent = 'Connected - Waiting for updates...';
    };
    
    appState.ws.onmessage = function(event) {
        try {
            const message = JSON.parse(event.data);
            console.log('📡 WebSocket message received:', message);
            
            if (message.type === 'update' && message.data) {
                handleRealTimeUpdate(message.data);
            }
        } catch (error) {
            console.error('❌ Error parsing WebSocket message:', error);
        }
    };
    
    appState.ws.onclose = function(event) {
        console.log('❌ WebSocket connection closed:', event.code, event.reason);
        
        // Don't reconnect if it was a normal close (tests finished)
        if (event.code === 1000) {
            document.getElementById('last-updated').textContent = 'Connection closed normally';
            return;
        }
        
        document.getElementById('last-updated').textContent = 'Connection lost - Attempting to reconnect...';
        
        // Attempt to reconnect with exponential backoff
        if (appState.reconnectAttempts < appState.maxReconnectAttempts) {
            appState.reconnectAttempts++;
            const delay = Math.min(1000 * Math.pow(2, appState.reconnectAttempts), 10000);
            console.log(`🔄 Reconnecting... (attempt ${appState.reconnectAttempts}/${appState.maxReconnectAttempts}) in ${delay}ms`);
            setTimeout(connectWebSocket, delay);
        } else {
            document.getElementById('last-updated').textContent = 'Connection failed - Please refresh the page';
        }
    };
    
    appState.ws.onerror = function(error) {
        console.error('❌ WebSocket error:', error);
    };
}

// Handle real-time updates from WebSocket with smart updates
function handleRealTimeUpdate(data) {
    console.log('📊 Real-time update received:', data);
    
    const hasResults = data.results && data.results.length > 0;
    const wasEmpty = appState.results.length === 0;
    
    // Update state
    const oldResults = [...appState.results];
    appState.results = data.results || [];
    appState.summary = data.summary || { total: 0, passed: 0, failed: 0, running: 0 };
    appState.lastUpdated = `Last updated: ${new Date().toLocaleString()}`;
    
    if (hasResults) {
        appState.testsStarted = true;
        
        // Show progress bar if tests are running
        updateGlobalProgress();
        
        // Update UI components
        updateSummaryCards();
        updateButtons();
        
        // Smart update: only update changed tests
        if (wasEmpty) {
            // First load - render all tests
            updateTestResultsInitial();
        } else {
            // Smart update - only update changed tests
            updateTestResultsSmart(oldResults, appState.results);
        }
        
        // Start animations for new running tests
        setTimeout(() => {
            animateProgressBars();
            animateNewResults();
        }, 100);
    } else {
        // No results yet
        appState.results = [];
        appState.summary = { total: 0, passed: 0, failed: 0, running: 0 };
        appState.lastUpdated = 'No tests run yet - click "Start Tests" to begin';
        
        hideGlobalProgress();
        updateSummaryCards();
        updateButtons();
    }
    
    // Update last updated time
    document.getElementById('last-updated').textContent = appState.lastUpdated;
    
    console.log('📊 Dashboard updated via WebSocket (smart):', {
        hasResults,
        testsStarted: appState.testsStarted,
        summary: appState.summary,
        resultCount: appState.results.length
    });
}

// Update global progress bar
function updateGlobalProgress() {
    const progressContainer = document.getElementById('global-progress-container');
    const progressBar = document.getElementById('global-progress-bar');
    const progressText = document.getElementById('progress-text');
    
    if (appState.summary.total > 0) {
        const completed = appState.summary.passed + appState.summary.failed;
        const percentage = (completed / appState.summary.total) * 100;
        
        progressContainer.style.display = 'block';
        progressBar.style.width = `${percentage}%`;
        progressText.textContent = `${completed} / ${appState.summary.total} tests completed`;
        
        // Animate progress bar if running
        if (appState.summary.running > 0) {
            progressBar.classList.add('animate-pulse');
        } else {
            progressBar.classList.remove('animate-pulse');
        }
    }
}

// Hide global progress bar
function hideGlobalProgress() {
    document.getElementById('global-progress-container').style.display = 'none';
}

// Smart update: only update tests that have changed
function updateTestResultsSmart(oldResults, newResults) {
    const container = document.getElementById('test-results');
    
    // Compare old and new results
    for (let i = 0; i < newResults.length; i++) {
        const newTest = newResults[i];
        const oldTest = oldResults[i];
        
        // Check if test changed or is new
        const testChanged = !oldTest || 
            oldTest.status !== newTest.status ||
            oldTest.message !== newTest.message ||
            oldTest.error !== newTest.error;
        
        if (testChanged) {
            const existingElement = container.children[i];
            const newElement = createTestElement(newTest, i);
            
            if (existingElement) {
                // Replace existing element with smooth transition
                existingElement.style.opacity = '0.5';
                setTimeout(() => {
                    container.replaceChild(newElement, existingElement);
                    newElement.style.opacity = '0';
                    setTimeout(() => {
                        newElement.style.opacity = '1';
                    }, 50);
                }, 100);
            } else {
                // Add new element
                newElement.style.opacity = '0';
                container.appendChild(newElement);
                setTimeout(() => {
                    newElement.style.opacity = '1';
                }, 50);
            }
        }
    }
    
    // Remove any extra elements if list got shorter
    while (container.children.length > newResults.length) {
        container.removeChild(container.lastChild);
    }
}

// Initial render of all test results
function updateTestResultsInitial() {
    const container = document.getElementById('test-results');
    container.innerHTML = '';
    
    appState.results.forEach((test, index) => {
        const testElement = createTestElement(test, index);
        container.appendChild(testElement);
    });
}

// Load initial data (fallback for HTTP API)
async function loadInitialData() {
    try {
        const response = await fetch('/api/results', {
            headers: {
                'Accept': 'application/json',
                'Cache-Control': 'no-cache'
            }
        });
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        const data = await response.json();
        console.log('📡 Initial data loaded:', data);
        
        // Only use this if WebSocket hasn't provided data yet
        if (appState.results.length === 0) {
            handleRealTimeUpdate(data);
        }
        
    } catch (error) {
        console.error('❌ Error loading initial data:', error);
        if (appState.results.length === 0) {
            document.getElementById('last-updated').textContent = `Error loading data: ${error.message}`;
        }
    }
}

// Update summary cards
function updateSummaryCards() {
    document.getElementById('total-count').textContent = appState.summary.total;
    document.getElementById('passed-count').textContent = appState.summary.passed;
    document.getElementById('failed-count').textContent = appState.summary.failed;
    document.getElementById('running-count').textContent = appState.summary.running;
}

// Update button visibility
function updateButtons() {
    const startBtn = document.getElementById('start-btn');
    const refreshBtn = document.getElementById('refresh-btn');
    
    if (appState.testsStarted && appState.results.length > 0) {
        startBtn.style.display = 'none';
        refreshBtn.style.display = 'flex';
        document.getElementById('no-tests-message').style.display = 'none';
        document.getElementById('test-results').style.display = 'block';
    } else {
        startBtn.style.display = 'flex';
        refreshBtn.style.display = 'none';
        document.getElementById('no-tests-message').style.display = 'block';
        document.getElementById('test-results').style.display = 'none';
    }
}

// Create a test element
function createTestElement(test, index) {
    const testDiv = document.createElement('div');
    testDiv.className = 'border-b border-gray-700 last:border-b-0 hover:bg-gray-750 transition-all duration-200 service-item';
    testDiv.style.transition = 'opacity 0.3s ease-in-out';
    
    const isExpanded = appState.expandedItems.has(index);
    
    testDiv.innerHTML = `
        <div class="p-6 cursor-pointer" onclick="toggleDetails(${index})">
            <div class="flex items-center justify-between">
                <div class="flex-1">
                    <h4 class="text-xl font-semibold mb-2">${test.name}</h4>
                    <p class="text-gray-400 mb-2">${test.message}</p>
                    
                    ${test.status === 'running' ? `
                    <div class="w-full bg-gray-700 rounded-full h-2 mb-2">
                        <div class="bg-yellow-500 h-2 rounded-full progress-bar animate-pulse" style="width: 0%"></div>
                    </div>
                    ` : ''}
                    
                    ${test.error ? `
                    <div class="mt-3 p-3 bg-red-900 border border-red-700 rounded-lg text-red-200 text-sm">
                        ${test.error}
                    </div>
                    ` : ''}
                </div>
                
                <div class="flex items-center gap-4 ml-6">
                    <div class="flex items-center gap-3">
                        <div class="w-3 h-3 rounded-full transition-all duration-300 ${getStatusIndicatorClass(test.status)}"></div>
                        <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase tracking-wide transition-all duration-300 ${getStatusBadgeClass(test.status)}">
                            ${test.status}
                        </span>
                    </div>
                    
                    <div class="transform transition-transform duration-200 ${isExpanded ? 'rotate-180' : ''}">
                        <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                        </svg>
                    </div>
                </div>
            </div>
        </div>
        
        <div class="details-container ${isExpanded ? 'expanded' : ''}" style="max-height: ${isExpanded ? '24rem' : '0'}; opacity: ${isExpanded ? '1' : '0'}; overflow: hidden; transition: max-height 0.3s ease-in-out, opacity 0.3s ease-in-out;">
            <div class="bg-gray-750 px-6 pb-6">
                <ul class="space-y-2 pt-4">
                    ${test.details ? test.details.map(detail => `
                        <li class="text-sm text-gray-400 py-2 border-b border-gray-600">${detail}</li>
                    `).join('') : ''}
                </ul>
            </div>
        </div>
    `;
    
    return testDiv;
}

// Get status indicator CSS class
function getStatusIndicatorClass(status) {
    switch (status) {
        case 'passed': return 'bg-green-500 shadow-green-500/50';
        case 'failed': return 'bg-red-500 shadow-red-500/50';
        case 'running': return 'bg-yellow-500 animate-pulse shadow-yellow-500/50';
        default: return 'bg-gray-500';
    }
}

// Get status badge CSS class
function getStatusBadgeClass(status) {
    switch (status) {
        case 'passed': return 'bg-green-900 text-green-300 border border-green-700';
        case 'failed': return 'bg-red-900 text-red-300 border border-red-700';
        case 'running': return 'bg-yellow-900 text-yellow-300 border border-yellow-700 animate-pulse';
        default: return 'bg-gray-900 text-gray-300 border border-gray-700';
    }
}

// Toggle test details (preserve expanded state during updates)
function toggleDetails(index) {
    console.log(`🔄 Toggling details for test ${index}`);
    
    if (appState.expandedItems.has(index)) {
        appState.expandedItems.delete(index);
    } else {
        appState.expandedItems.add(index);
    }
    
    // Only update this specific test element
    const container = document.getElementById('test-results');
    const element = container.children[index];
    if (element) {
        const newElement = createTestElement(appState.results[index], index);
        container.replaceChild(newElement, element);
    }
}

// Start tests
async function startTests() {
    const btn = document.getElementById('start-btn');
    appState.isLoading = true;
    btn.disabled = true;
    btn.innerHTML = '🚀 Starting Tests...';
    
    console.log('▶️ Starting tests...');
    
    try {
        const response = await fetch('/api/refresh', { 
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        const result = await response.json();
        console.log('✅ Tests started:', result);
        
        appState.testsStarted = true;
        
        // WebSocket will handle all updates now
        // Animate the button
        animateButtonSuccess();
        
    } catch (error) {
        console.error('❌ Error starting tests:', error);
        alert(`Failed to start tests: ${error.message}`);
        btn.disabled = false;
        btn.innerHTML = '▶️ Start Tests';
    } finally {
        appState.isLoading = false;
    }
}

// Refresh tests
async function refreshTests() {
    const btn = document.getElementById('refresh-btn');
    appState.isLoading = true;
    btn.disabled = true;
    btn.innerHTML = '🔄 Refreshing...';
    
    console.log('🔄 Refreshing tests...');
    
    try {
        const response = await fetch('/api/refresh', { 
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        const result = await response.json();
        console.log('✅ Tests refreshed:', result);
        
        // Clear expanded items and reset progress
        appState.expandedItems.clear();
        
        // WebSocket will handle all updates now
        // Animate the button
        animateButtonSuccess();
        
    } catch (error) {
        console.error('❌ Error refreshing tests:', error);
        alert(`Failed to refresh tests: ${error.message}`);
    } finally {
        setTimeout(() => {
            appState.isLoading = false;
            btn.disabled = false;
            btn.innerHTML = '🔄 Refresh Tests';
        }, 2000);
    }
}

// Animations
function animateProgressBars() {
    setTimeout(() => {
        const progressBars = document.querySelectorAll('.progress-bar');
        progressBars.forEach(bar => {
            if (typeof anime !== 'undefined') {
                try {
                    anime({
                        targets: bar,
                        width: ['0%', '100%'],
                        duration: 3000,
                        easing: 'easeInOutQuad',
                        loop: true
                    });
                } catch (error) {
                    console.log('Progress animation skipped:', error.message);
                }
            }
        });
    }, 50);
}

function animateNewResults() {
    const serviceItems = document.querySelectorAll('.service-item');
    if (typeof anime !== 'undefined' && serviceItems.length > 0) {
        try {
            anime({
                targets: serviceItems,
                translateY: [20, 0],
                opacity: [0, 1],
                duration: 500,
                delay: anime.stagger(100),
                easing: 'easeOutQuad'
            });
        } catch (error) {
            console.log('Result animation skipped:', error.message);
        }
    }
}

function animateButtonSuccess() {
    const buttons = document.querySelectorAll('button[id$="-btn"]');
    if (typeof anime !== 'undefined' && buttons.length > 0) {
        try {
            anime({
                targets: buttons,
                scale: [1, 1.05, 1],
                duration: 300,
                easing: 'easeInOutQuad'
            });
        } catch (error) {
            console.log('Button animation skipped:', error.message);
        }
    }
}

// Cleanup function
function cleanup() {
    if (appState.ws && appState.ws.readyState === WebSocket.OPEN) {
        appState.ws.close();
    }
}

// Error handling
window.addEventListener('error', (event) => {
    // Only log serious errors
    const errorMessage = event.error?.message || event.message || '';
    const filename = event.filename || '';
    
    // Skip known harmless errors
    const harmlessErrors = ['Script error', 'anime', 'cdn', 'websocket'];
    const isHarmless = harmlessErrors.some(err => 
        errorMessage.toLowerCase().includes(err) || 
        filename.toLowerCase().includes(err)
    );
    
    if (!isHarmless) {
        console.error('❌ JavaScript error:', {
            message: errorMessage,
            filename: event.filename,
            lineno: event.lineno
        });
    }
});

// Cleanup on page unload
window.addEventListener('beforeunload', cleanup);

// Debug utilities
window.debugTester = {
    getState: () => appState,
    getResults: () => fetch('/api/results').then(r => r.json()),
    startTests: startTests,
    refreshTests: refreshTests,
    reconnectWS: connectWebSocket,
    clearConsole: () => {
        console.clear();
        console.log('🧹 Console cleared');
    }
};

console.log('🎯 WebSocket API Tester with Smart Updates loaded!');
console.log('🔧 Debug utilities available at window.debugTester');