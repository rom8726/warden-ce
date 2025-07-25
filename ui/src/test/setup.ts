import '@testing-library/jest-dom';

// Mock matchMedia for components that use media queries
window.matchMedia = window.matchMedia || (function() {
  return {
    matches: false,
    media: '',
    onchange: null,
    addListener: function() {},
    removeListener: function() {},
    addEventListener: function() {},
    removeEventListener: function() {},
    dispatchEvent: function() {
      return false;
    },
  } as MediaQueryList;
}) as (query: string) => MediaQueryList;

// Suppress React 18 console errors/warnings
const originalConsoleError = console.error;
console.error = (...args) => {
  if (
    /Warning.*not wrapped in act/i.test(args[0]) ||
    /Warning.*ReactDOM.render is no longer supported/i.test(args[0])
  ) {
    return;
  }
  originalConsoleError(...args);
};