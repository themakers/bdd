// Package bdd is a very simple but expressive BDD-flavoured testing framework
//
// Documentation is incomplete
// Refer to [bdd_test.go](bdd_test.go) for usage
//
// ## Rationale
//
// The most of present Golang BDD-style testing frameworks are:
// 1) too big => hard to maintain
// 2) overcomplicated => same here
// 3) bloated => goes against golang philosophy
// 4) unsupported
// 5) not enough expressive
// 6) provides less control of execution order
//
// This one is:
// 1) dead-simple
// 2) easy to maintain
// 3) more expressive
// 4) better control of test blocks execution order
//
// ## The Goal
//
// 1) Be minimalistic, yet expressive and fully functional
// 2) Stay within 500LOC
// 3) Be easily maintainable
// 4) Be compatible with all the popular ecosystem tools (without sacrificing the above)
//
//
// No matcher library provided; There are plenty of them.
// i.e.: [github.com/stretchr/testify/assert](https://github.com/stretchr/testify/tree/master/assert)
//
package bdd

import (
	"encoding/json"
	"fmt"
	"hash/crc64"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/fatih/color"
)

/***************************************************************
** Printing Helpers
**/

// For debugging purposes; To print objects...
func jsonf(v interface{}) string {
	if data, err := json.MarshalIndent(v, "", "  "); err != nil {
		panic(err)
	} else {
		return string(data)
	}
}

func printService(s *scenario, msg ...interface{}) {
	s.t.Helper()
	s.t.Log(color.New(color.FgCyan).Sprint("<BDD> ", fmt.Sprint(msg...)))
}

func printScenario(s *scenario, msg ...interface{}) {
	s.t.Helper()
	s.t.Log(color.New(color.FgBlue).Sprint("<", s.runID, "/S> ", fmt.Sprint(msg...)))
}

func printAct(s *scenario, msg ...interface{}) {
	s.t.Helper()
	s.t.Log(color.New(color.FgYellow).Sprint("<", s.runID, "/A> ", fmt.Sprint(msg...)))
}

func printBranch(s *scenario, msg ...interface{}) {
	s.t.Helper()
	s.t.Log(color.New(color.FgHiYellow).Sprint("<", s.runID, "/B> ", fmt.Sprint(msg...)))
}

func printTest(s *scenario, msg ...interface{}) {
	s.t.Helper()
	var c = color.FgGreen
	if s.t.Failed() {
		c = color.FgRed
	}
	s.t.Log(color.New(c).Sprint("<", s.runID, "/T> ", fmt.Sprint(msg...)))
}

/***************************************************************
** Core Helpers
**/

var runCounter int64 = 0

func runID() string { return fmt.Sprint(atomic.AddInt64(&runCounter, 1)) }

var (
	scenarios     = map[uint64]*scenario{}
	scenariosLock sync.RWMutex
)

func testID(t *testing.T) uint64 { return uint64(uintptr(unsafe.Pointer(t))) }
func getScenario(t *testing.T) *scenario {
	scenariosLock.RLock()
	defer scenariosLock.RUnlock()

	return scenarios[testID(t)]
}

func sum(name []byte) string {
	return fmt.Sprintf("%x", crc64.Checksum(name, crc64.MakeTable(crc64.ECMA)))
}

func branchID(forkFile string, forkLine uint64, branchName string) string {
	return sum(append(
		append(
			[]byte(sum([]byte(forkFile))),
			[]byte{
				byte(0xff & forkLine),
				byte(0xff & (forkLine >> 8)),
				byte(0xff & (forkLine >> 16)),
				byte(0xff & (forkLine >> 24)),
				byte(0xff & (forkLine >> 32)),
				byte(0xff & (forkLine >> 40)),
				byte(0xff & (forkLine >> 48)),
				byte(0xff & (forkLine >> 56)),
			}...,
		),
		branchName...,
	))
}

func addPath(path, node string) string { return fmt.Sprintf("%s/%s", path, node) }

func whereIsTheFork() (file string, line int) {
	pc := [1]uintptr{}
	runtime.Callers(4, pc[:])
	f := runtime.FuncForPC(pc[0])
	return f.FileLine(pc[0])
}

/***************************************************************
** Core implementation
**/

type tree struct {
	id       string
	children map[string]*tree
}

type scenario struct {
	t *testing.T

	runID string

	path string
	tree map[string]bool
}

func (s *scenario) run(title string, fn ScenarioFunc) {
	s.t.Helper()

L:
	for {
		printService(s, strings.Repeat("•", 64))

		s.runID = runID()
		s.path = ""

		printScenario(s, title)

		fn(s.t, s.runID)

		{ //> Are there any unvisited paths?
			thereIs := false
			for _, v := range s.tree {
				if !v {
					thereIs = true
				}
			}
			if !thereIs {
				break L
			}
		}
	}

}

func (s *scenario) visit(id string) bool {

	var (
		oldPath  = s.path
		nextPath = addPath(s.path, id)
	)

	res := map[string]bool{}
	for p, v := range s.tree {
		if strings.HasPrefix(p, nextPath) {
			res[p] = v
		}
	}

	if len(res) == 0 {
		s.tree[nextPath] = true
		s.path = nextPath
		delete(s.tree, oldPath)
		return true
	}

	for _, v := range res {
		if !v {
			s.tree[nextPath] = true
			s.path = nextPath
			delete(s.tree, oldPath)
			return true
		}
	}

	return false
}

func (s *scenario) blockAct(name string, fn BlockFunc) {
	s.t.Helper()

	printAct(s, name)
	fn()
}

func (s *scenario) blockFork(branches []BranchConstructorFunc) {
	s.t.Helper()

	var (
		oldPath            = s.path
		visited            = false
		forkFile, forkLine = whereIsTheFork()
	)

	for _, branchFn := range branches {
		name, fn := branchFn()

		// TODO: Check for duplicated names

		id := branchID(forkFile, uint64(forkLine), name)

		if visited {
			s.tree[addPath(oldPath, id)] = false
		} else if s.visit(id) {
			visited = true

			printBranch(s, name)

			(func() {
				defer func() {
					s.path = strings.TrimSuffix(s.path, fmt.Sprintf("/%s", id))
				}()
				fn()
			})()
		}
	}
}

func (s *scenario) blockTest(what string, fn BlockFunc) {
	s.t.Helper()

	printTest(s, what)
	fn()
}

/***************************************************************
** Public API
**/

type BlockFunc func()

type ScenarioFunc func(t *testing.T, runID string)

// Scenario is a root of your scenario.
// Scenario blocks can not be nested.
func Scenario(t *testing.T, title string, fn ScenarioFunc) {
	t.Helper()

	scenario := &scenario{
		t:    t,
		tree: map[string]bool{},
	}

	(func() {
		scenariosLock.Lock()
		defer scenariosLock.Unlock()

		scenarios[testID(t)] = scenario
	})()

	defer (func() {
		scenariosLock.Lock()
		defer scenariosLock.Unlock()

		delete(scenarios, testID(t))
	})()

	scenario.run(title, fn)
}

// Act is just a structural block; You use it to make your test structure
// more readable and explanatory
func Act(t *testing.T, name string, fn BlockFunc) {
	t.Helper()

	getScenario(t).blockAct(name, fn)
}

type BranchConstructorFunc func() (string, BlockFunc)

// Fork comes into a business when you want a separate sub-scenarios in your tests.
// For each Branch in a Fork the framework will reexecute all the blocks
// defined prior to the fork
//
// TODO: Rewrite this shame
func Fork(t *testing.T, branches ...BranchConstructorFunc) {
	t.Helper()

	getScenario(t).blockFork(branches)
}

// Branch is a branch of a Fork; That is.
// Branch block names should NOT have duplicates within same Fork
func Branch(name string, fn BlockFunc) BranchConstructorFunc {
	return func() (string, BlockFunc) {
		return name, fn
	}
}

// Test is a block where you write your actual test code and run your assertions.
// Remember: you write ALL the logic in Test blocks,
// except for some error-free-by-default initializations
// Test blocks can be nested.
// Test blocks can NOT contain other block types.
func Test(t *testing.T, what string, fn BlockFunc) {
	t.Helper()

	getScenario(t).blockTest(what, fn)
}
