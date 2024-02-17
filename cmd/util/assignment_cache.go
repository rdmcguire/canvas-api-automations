package util

import (
	"slices"
	"sync"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
)

// Thread-safe storage of assignments with modules mapped to
// ridiculous netacad names (or any other name)
type AssignmentCache struct {
	modules     map[string]*canvasauto.Module     // Found modules by netacad name
	assignments map[string]*canvasauto.Assignment // Found assignments by netacad name
	lostCauses  []string                          // Netacad grade items that can't be located
	lock        sync.Mutex
}

// unnecessary API requests attempting to locate
// assignments / modules by dumbass netacad export names
func NewAssignmentCache() *AssignmentCache {
	return &AssignmentCache{
		assignments: make(map[string]*canvasauto.Assignment, 0),
		modules:     make(map[string]*canvasauto.Module, 0),
	}
}

// Safely retrieves by netacad assignment name
func (c *AssignmentCache) Get(name string) (*canvasauto.Assignment, *canvasauto.Module) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if a, set := c.assignments[name]; set {
		return a, c.modules[name]
	}
	return nil, nil
}

// Sets an assignment with netacad assignment name
func (c *AssignmentCache) Set(name string, assignment *canvasauto.Assignment, module *canvasauto.Module) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.assignments[name] = assignment
	c.modules[name] = module
}

// Returns true if an item is a lost cause
func (c *AssignmentCache) LostCause(name string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if slices.Contains(c.lostCauses, name) {
		return true
	}

	return false
}

// Registers a lost cause
func (c *AssignmentCache) IsLostCause(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !slices.Contains(c.lostCauses, name) {
		c.lostCauses = append(c.lostCauses, name)
	}
}
