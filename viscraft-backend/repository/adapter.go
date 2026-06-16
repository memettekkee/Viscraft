package repository

// This file contains adapters for satisfying service-layer interfaces.
// The adapters bridge repository types to the interfaces defined in the service package
// without creating circular imports.

// ProjectImageFinderAdapter wraps ProjectRepository to provide image ID lookups
// for project deletion cleanup. It satisfies interfaces expecting
// FindImagesByProjectId(projectId string) ([]string, error).
// Since ProjectRepository.FindImagesByProjectId already returns ([]string, error),
// ProjectRepository can be used directly without an adapter.

// Note: The ImageFinder adapter (for service.ImageFinder interface) must be wired
// at the composition root (main package or wire package) because it requires importing
// both the repository and service packages. The ImageRepository.FindImagesByUserId
// method returns []repository.ImageRecord which must be converted to []service.ImageRecord.
