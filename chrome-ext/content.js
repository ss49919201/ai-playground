// GitHub PR Checkout Helper - Content Script

function getPRNumber() {
  const urlMatch = window.location.pathname.match(/\/pull\/(\d+)/);
  return urlMatch ? urlMatch[1] : null;
}

function createCopyButton(prNumber) {
  const button = document.createElement('button');
  button.className = 'btn btn-sm';
  button.style.marginLeft = '8px';
  button.style.backgroundColor = '#238636';
  button.style.color = '#ffffff';
  button.style.border = '1px solid #30363d';
  button.style.borderRadius = '6px';
  button.style.padding = '8px 16px';
  button.style.fontSize = '14px';
  button.style.fontWeight = 'bold';
  button.style.cursor = 'pointer';
  button.style.boxShadow = '0 2px 4px rgba(0,0,0,0.3)';
  button.style.transition = 'all 0.2s ease';
  button.textContent = '📋 Copy checkout';
  button.title = `Copy 'gh pr checkout ${prNumber}' to clipboard`;
  
  button.addEventListener('click', async () => {
    const command = `gh pr checkout ${prNumber}`;
    
    try {
      await navigator.clipboard.writeText(command);
      
      // Visual feedback
      const originalText = button.textContent;
      button.textContent = 'Copied!';
      button.style.backgroundColor = '#238636';
      
      setTimeout(() => {
        button.textContent = originalText;
        button.style.backgroundColor = '#238636';
      }, 2000);
      
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
      
      // Fallback feedback for error
      const originalText = button.textContent;
      button.textContent = 'Failed';
      button.style.backgroundColor = '#da3633';
      
      setTimeout(() => {
        button.textContent = originalText;
        button.style.backgroundColor = '#238636';
      }, 2000);
    }
  });
  
  return button;
}

function insertCopyButton() {
  const prNumber = getPRNumber();
  if (!prNumber) return;
  
  // Look for the PR header area
  const prHeader = document.querySelector('.gh-header-meta');
  if (!prHeader) return;
  
  // Check if button already exists
  if (document.querySelector('.pr-checkout-copy-btn')) return;
  
  const copyButton = createCopyButton(prNumber);
  copyButton.classList.add('pr-checkout-copy-btn');
  
  prHeader.appendChild(copyButton);
}

// Initialize when page loads
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', insertCopyButton);
} else {
  insertCopyButton();
}

// Handle GitHub's AJAX navigation
let currentUrl = window.location.href;
const observer = new MutationObserver(() => {
  if (window.location.href !== currentUrl) {
    currentUrl = window.location.href;
    setTimeout(insertCopyButton, 500); // Small delay for DOM to update
  }
});

observer.observe(document.body, {
  childList: true,
  subtree: true
});