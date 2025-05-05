/**
 * Displays a spinning loader animation to indicate a loading state.
 * 
 * Renders:
 * - A circular spinner using Tailwind classes.
 * - Accessible attributes for screen readers.
 */
export default function LoadingScreen() {
    return (
        <div
          className="h-8 w-8 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"
          role="status"
          aria-label="Loading"
        >
          <span className="sr-only">Loading</span>
        </div>
    );
  }
  