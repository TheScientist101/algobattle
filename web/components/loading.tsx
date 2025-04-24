export default function LoadingScreen() {
    return (
      <div className="flex items-center justify-center h-screen w-screen">
        <div
          className="h-8 w-8 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"
          role="status"
          aria-label="Loading"
        >
          <span className="sr-only">Loading</span>
        </div>
      </div>
    );
  }
  