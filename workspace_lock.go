package otf

import "errors"

var (
	ErrWorkspaceAlreadyLocked         = errors.New("workspace already locked")
	ErrWorkspaceAlreadyUnlocked       = errors.New("workspace already unlocked")
	ErrWorkspaceLockedByDifferentUser = errors.New("workspace locked by different user")
)

// WorkspaceLock is the lock for the workspace.
type WorkspaceLock struct {
	Locker WorkspaceLocker
}

// WorkspaceLocker is the entity that has locked a workspace.
type WorkspaceLocker interface {
	GetID() string
	String() string
}

// Lock locks the workspace with the specified locker.
func (l *WorkspaceLock) Lock(locker WorkspaceLocker) error {
	if l.IsLocked() {
		return ErrWorkspaceAlreadyLocked
	}

	l.Locker = locker

	return nil
}

// Unlock unlocks the workspace with the specified locker. Only the original
// locker can unlock the workspace, unless force is true.
func (l *WorkspaceLock) Unlock(unlocker WorkspaceLocker, force bool) error {
	if !l.IsLocked() {
		return ErrWorkspaceAlreadyUnlocked
	}

	if force {
		l.Locker = nil
		return nil
	}

	if l.Locker.GetID() != unlocker.GetID() {
		return ErrWorkspaceLockedByDifferentUser
	}

	l.Locker = nil
	return nil
}

// IsLocked queries whether the workspace lock is locked or unlocked.
func (l WorkspaceLock) IsLocked() bool {
	return l.Locker != nil
}
