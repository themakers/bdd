# BDD

[![codecov](https://codecov.io/gh/themakers/bdd/branch/master/graph/badge.svg)](https://codecov.io/gh/themakers/bdd)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/themakers/bdd)

Package bdd is a very simple but expressive BDD-flavoured testing framework

Documentation is incomplete
Refer to [bdd_test.go](bdd_test.go) for usage

## Rationale

The most of present Golang BDD-style testing frameworks are:
1) too big => hard to maintain
2) overcomplicated => same here
3) bloated => goes against golang philosophy
4) unsupported
5) not enough expressive
6) provides less control of execution order

This one is:
1) dead-simple
2) easy to maintain
3) more expressive
4) better control of test blocks execution order

## The Goal

1) Be minimalistic, yet expressive and fully functional
2) Stay within 500LOC
3) Be easily maintainable
4) Be compatible with all the popular ecosystem tools (without sacrificing the above)

No matcher library provided; There are plenty of them.
i.e.: [github.com/stretchr/testify/assert]([https://github.com/stretchr/testify/tree/master/assert](https://github.com/stretchr/testify/tree/master/assert))

## Functions

### func [Act](/bdd.go#L308)

`func Act(t *testing.T, name string, fn BlockFunc)`

Act is just a structural block; You use it to make your test structure
more readable and explanatory

### func [Fork](/bdd.go#L321)

`func Fork(t *testing.T, branches ...BranchConstructorFunc)`

Fork comes into a business when you want a separate sub-scenarios in your tests.
For each Branch in a Fork the framework will reexecute all the blocks
defined prior to the fork

TODO: Rewrite this shame

### func [Scenario](/bdd.go#L281)

`func Scenario(t *testing.T, title string, fn ScenarioFunc)`

Scenario is a root of your scenario.
Scenario blocks can not be nested.

### func [Test](/bdd.go#L340)

`func Test(t *testing.T, what string, fn BlockFunc)`

Test is a block where you write your actual test code and run your assertions.
Remember: you write ALL the logic in Test blocks,
except for some error-free-by-default initializations
Test blocks can be nested.
Test blocks can NOT contain other block types.

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
