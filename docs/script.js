document.addEventListener('DOMContentLoaded', () => {
    initTabs();
    initCopy();
    initTheme();
    initScroll();
    initObservers();
    fetchStars();
    
    // Dynamic Year
    const yearEl = document.getElementById('year');
    if (yearEl) yearEl.innerText = new Date().getFullYear();
});

// Tab Switching & Keyboard Navigation
function initTabs() {
    const tabs = document.querySelectorAll('.tab-link');
    const panels = document.querySelectorAll('.tab-content');

    tabs.forEach(tab => {
        tab.addEventListener('click', () => switchTab(tab));
        
        tab.addEventListener('keydown', (e) => {
            const index = Array.from(tabs).indexOf(tab);
            if (e.key === 'ArrowRight') {
                tabs[(index + 1) % tabs.length].focus();
                switchTab(tabs[(index + 1) % tabs.length]);
            } else if (e.key === 'ArrowLeft') {
                tabs[(index - 1 + tabs.length) % tabs.length].focus();
                switchTab(tabs[(index - 1 + tabs.length) % tabs.length]);
            }
        });
    });

    function switchTab(selectedTab) {
        tabs.forEach(t => {
            t.classList.remove('active');
            t.setAttribute('aria-selected', 'false');
        });
        panels.forEach(p => {
            p.classList.remove('show');
            p.setAttribute('hidden', '');
        });

        selectedTab.classList.add('active');
        selectedTab.setAttribute('aria-selected', 'true');
        const panelId = selectedTab.getAttribute('aria-controls');
        const activePanel = document.getElementById(panelId);
        activePanel.classList.add('show');
        activePanel.removeAttribute('hidden');
    }
}

// Modern Clipboard API
function initCopy() {
    document.querySelectorAll('.copy-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const text = btn.getAttribute('data-copy');
            navigator.clipboard.writeText(text).then(() => {
                const originalText = btn.innerText;
                btn.innerText = "Copied!";
                btn.style.backgroundColor = "#10b981";
                setTimeout(() => {
                    btn.innerText = originalText;
                    btn.style.backgroundColor = "";
                }, 2000);
            });
        });
    });
}

// Theme Management
function initTheme() {
    const themeToggle = document.getElementById('theme-toggle');
    const getPreferredTheme = () => {
        const saved = localStorage.getItem('theme');
        if (saved) return saved;
        return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    };

    const setTheme = (theme) => {
        if (theme === 'dark') {
            document.documentElement.setAttribute('data-theme', 'dark');
        } else {
            document.documentElement.removeAttribute('data-theme');
        }
        localStorage.setItem('theme', theme);
    };

    setTheme(getPreferredTheme());

    themeToggle.addEventListener('click', () => {
        const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
        setTheme(isDark ? 'light' : 'dark');
    });
}

// Scroll Handling
function initScroll() {
    const backToTop = document.getElementById('back-to-top');
    window.addEventListener('scroll', () => {
        if (window.scrollY > 400) {
            backToTop.classList.add('visible');
        } else {
            backToTop.classList.remove('visible');
        }
    }, { passive: true });

    backToTop.addEventListener('click', () => {
        window.scrollTo({ top: 0, behavior: 'smooth' });
    });
}

// Consolidated Intersection Observer
function initObservers() {
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                if (entry.target.id === 'terminal-body') {
                    if (!entry.target.dataset.animated) {
                        entry.target.dataset.animated = "true";
                        typeTerminal();
                    }
                } else {
                    entry.target.style.opacity = "1";
                    entry.target.style.transform = "translateY(0)";
                    observer.unobserve(entry.target);
                }
            }
        });
    }, { threshold: 0.1 });

    document.querySelectorAll('.feature-card, .installation, #terminal-body').forEach(el => {
        if (el.id !== 'terminal-body') {
            el.style.opacity = "0";
            el.style.transform = prefersReducedMotion ? "none" : "translateY(30px)";
            el.style.transition = "all 0.8s cubic-bezier(0.22, 1, 0.36, 1)";
        }
        observer.observe(el);
    });
}

// Terminal Animation with Syntax Highlighting
const terminalLines = [
    { type: 'input', content: 'errpipe' },
    { type: 'system', content: 'errpipe session activated. Monitoring for errors...' },
    { type: 'input', content: 'python3 app.py' },
    { type: 'error', content: 'Traceback (most recent call last):\n  File "app.py", line 5, in <module>\n    result = data["user"]["id"]\nKeyError: \'user\'' },
    { type: 'ai', content: '<strong>AI Analysis:</strong> <span class="hl-type">KeyError</span> detected. You\'re accessing <span class="hl-key">"user"</span> but the dictionary is empty. <br>Add a check: <code>if "user" in data:</code>' }
];

async function typeTerminal() {
    const terminalBody = document.getElementById('terminal-body');
    if (!terminalBody) return;
    terminalBody.innerHTML = '';
    
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    const speedMultiplier = prefersReducedMotion ? 0 : 1;

    for (const line of terminalLines) {
        const div = document.createElement('div');
        div.style.marginBottom = '0.75rem';
        
        if (line.type === 'input') {
            div.innerHTML = `<span class="prompt">$</span> <span class="command"></span><span class="cursor"></span>`;
            terminalBody.appendChild(div);
            const cmdSpan = div.querySelector('.command');
            
            for (let i = 0; i < line.content.length; i++) {
                cmdSpan.textContent += line.content[i];
                if (speedMultiplier > 0) await new Promise(r => setTimeout(r, 45 + Math.random() * 45));
            }
            div.querySelector('.cursor').remove();
        } else {
            if (line.type === 'system') div.className = 'system-msg';
            if (line.type === 'error') {
                div.className = 'error';
                // Basic trace highlighting
                div.innerHTML = line.content
                    .replace(/"(.*\.py)"/g, '<span class="hl-file">"$1"</span>')
                    .replace(/line (\d+)/g, 'line <span class="hl-line">$1</span>');
            }
            if (line.type === 'ai') {
                div.className = 'ai-response';
                div.innerHTML = line.content;
                if (speedMultiplier > 0) {
                    div.style.opacity = '0';
                    div.style.transform = 'translateX(-10px)';
                    div.style.transition = 'all 0.5s ease-out';
                }
            } else if (line.type !== 'error') {
                div.textContent = line.content;
            }
            
            terminalBody.appendChild(div);
            if (speedMultiplier > 0) {
                await new Promise(r => setTimeout(r, line.type === 'ai' ? 200 : 600));
                if (line.type === 'ai') {
                    div.style.opacity = '1';
                    div.style.transform = 'translateX(0)';
                }
            }
        }
        if (speedMultiplier > 0) await new Promise(r => setTimeout(r, 800));
    }
}

// GitHub Stars
async function fetchStars() {
    const starCount = document.getElementById('star-count');
    try {
        const response = await fetch('https://api.github.com/repos/diamondosas/errpipe');
        const data = await response.json();
        const count = data.stargazers_count;
        if (count !== undefined) {
            starCount.classList.remove('skeleton-text');
            starCount.innerText = count >= 1000 ? (count / 1000).toFixed(1) + 'k' : count;
        }
    } catch (error) {
        starCount.classList.remove('skeleton-text');
        starCount.innerText = 'Star';
    }
}
