# Testing Instructions for GitHub PR Checkout Helper

## How to Load the Extension in Chrome

1. Open Chrome and go to `chrome://extensions/`
2. Enable "Developer mode" (toggle in top right)
3. Click "Load unpacked"
4. Select this directory: `/Users/sakaeshinya/src/ai-kata/chrome-ext`

## How to Test

1. Navigate to any GitHub Pull Request page (e.g., `https://github.com/owner/repo/pull/123`)
2. Look for the "Copy checkout" button in the PR header area
3. Click the button to copy `gh pr checkout 123` to clipboard
4. The button should show "Copied!" feedback for 2 seconds
5. Test navigation between different PRs to ensure the button updates correctly

## Expected Behavior

- Button appears on all GitHub PR pages matching `github.com/*/pull/*`
- Button copies the correct PR number for each page
- Visual feedback shows success/failure states
- Button persists through GitHub's AJAX navigation

## Troubleshooting

- If button doesn't appear, check browser console for errors
- Ensure the extension has clipboard permissions
- Try refreshing the page if button doesn't update after navigation