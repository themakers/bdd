package bdd

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/fatih/color"
)

/***************************************************************
** Test Helpers
**/

type shorthand [][2]string

func (s *shorthand) add(name, rid string) {
	*s = append(*s, [...]string{name, rid})
}

func (s shorthand) string1() string {
	var out strings.Builder

	rid := ""

	out.WriteString("{")
	for _, r := range s {
		if r[1] != rid {
			rid = r[1]
			out.WriteString(fmt.Sprintf("\n %s:", rid))
		}
		out.WriteString(fmt.Sprintf(" %s", r[0]))
	}
	out.WriteString("\n}")

	return out.String()
}

func (s shorthand) string2() string {
	var out strings.Builder

	rid := ""

	out.WriteString("{")
	for _, r := range s {
		if r[1] != rid {
			rid = r[1]
			out.WriteString("\n")
		}
		out.WriteString(fmt.Sprintf(" {%s:%s}", r[0], r[1]))
	}
	out.WriteString("\n}")

	return out.String()
}

func (s shorthand) String() string {
	return s.string1()
}

func (s shorthand) assert(t *testing.T, reference shorthand) {
	if !reflect.DeepEqual(s, reference) {

		t.Log(color.New(color.FgRed).Sprintf(
			"\n\nFAILURE:\n\nexpected: %s\n\nactual: %s", reference, s,
		))

		t.FailNow()
	}
}

/***************************************************************
** The Tests
**/

func TestBasic(t *testing.T) {
	runCounter = 0
	var s = &shorthand{}

	Scenario(t, "A", func(t *testing.T, runID string) {
		s.add("A", runID)

		Act(t, "B", func() {
			s.add("B", runID)

			Test(t, "C", func() {
				s.add("C", runID)
			})

			Fork(t,
				Branch("D", func() {
					s.add("D", runID)

					Test(t, "E", func() {
						s.add("E", runID)
					})

				}),
				Branch("F", func() {
					s.add("F", runID)

					Test(t, "G", func() {
						s.add("G", runID)
					})

				}),
			)

			Test(t, "H", func() {
				s.add("H", runID)

				Test(t, "I", func() {
					s.add("I", runID)
				})

			})

		})

	})

	s.assert(t, shorthand{
		{"A", "1"},
		{"B", "1"},
		{"C", "1"},
		{"D", "1"},
		{"E", "1"},
		{"H", "1"},
		{"I", "1"},
		{"A", "2"},
		{"B", "2"},
		{"C", "2"},
		{"F", "2"},
		{"G", "2"},
		{"H", "2"},
		{"I", "2"},
	})

}

func TestRealWorld(t *testing.T) {
	// TODO: Improve to reflect some real-world example

	runCounter = 0
	var s = &shorthand{}

	Scenario(t, "A", func(t *testing.T, runID string) {
		s.add("A", runID)
	})

	Scenario(t, "B", func(t *testing.T, runID string) {
		s.add("B", runID)

		Fork(t,
			Branch("C", func() {
				s.add("C", runID)
			}),
			RealWorldExternalBranch(t, runID, s),
		)
	})

	s.assert(t, shorthand{
		{"A", "1"},
		{"B", "2"},
		{"C", "2"},
		{"B", "3"},
		{"D", "3"},
		{"E", "3"},
	})
}

func RealWorldExternalBranch(t *testing.T, runID string, s *shorthand) BranchConstructorFunc {
	return Branch("D: External branch", func() {
		s.add("D", runID)

		Act(t, "E", func() {
			s.add("E", runID)
		})
	})
}

func TestBranching(t *testing.T) {
	runCounter = 0
	var s = &shorthand{}

	Scenario(t, "A", func(t *testing.T, runID string) {
		s.add("A", runID)

		Fork(t,
			Branch("B", func() {
				s.add("B", runID)
			}),
			Branch("C", func() {
				s.add("C", runID)
			}),
		)

		Fork(t,
			Branch("D", func() {
				s.add("D", runID)

				Fork(t,
					Branch("E", func() {
						s.add("E", runID)
					}),
					Branch("F", func() {
						s.add("F", runID)
					}),
				)
			}),
			Branch("G", func() {
				s.add("G", runID)
			}),
		)

	})

	s.assert(t, shorthand{
		{"A", "1"},
		{"B", "1"},
		{"D", "1"},
		{"E", "1"},

		{"A", "2"},
		{"C", "2"},
		{"D", "2"},
		{"F", "2"},

		{"A", "3"},
		{"G", "3"},
	})
}

func TestFull(t *testing.T) {
	runCounter = 0
	var s = &shorthand{}

	Scenario(t, "A", func(t *testing.T, runID string) {
		s.add("A", runID)

		Act(t, "B", func() {
			s.add("B", runID)

			Test(t, "C", func() {
				s.add("C", runID)
			})

			Fork(t,
				Branch("D", func() {
					s.add("D", runID)

					Test(t, "E", func() {
						s.add("E", runID)
					})

					Test(t, "F", func() {
						s.add("F", runID)

						Test(t, "G", func() {
							s.add("G", runID)
						})

					})

				}),
				Branch("H", func() {
					s.add("H", runID)

					Act(t, "I", func() {
						s.add("I", runID)

						Test(t, "J", func() {
							s.add("J", runID)
						})

						Act(t, "K", func() {
							s.add("K", runID)

							Test(t, "L", func() {
								s.add("L", runID)
							})

							Test(t, "M", func() {
								s.add("M", runID)
							})

						})

					})

				}),
				Branch("N", func() {
					s.add("N", runID)

					Test(t, "O", func() {
						s.add("O", runID)
					})

				}),
			)

			Test(t, "P", func() {
				s.add("P", runID)
			})

			Fork(t,
				Branch("Q", func() {
					s.add("Q", runID)

					Test(t, "R", func() {
						s.add("R", runID)
					})

				}),
				Branch("S", func() {
					s.add("S", runID)

					Act(t, "T", func() {
						s.add("T", runID)

						Fork(t,
							Branch("U", func() {
								s.add("U", runID)

								Test(t, "V", func() {
									s.add("V", runID)
								})

							}),
							Branch("W", func() {
								s.add("W", runID)

								Test(t, "X", func() {
									s.add("X", runID)
								})

							}),
						)

					})

				}),
			)

			Act(t, "Y", func() {
				s.add("Y", runID)

				Test(t, "Z", func() {
					s.add("Z", runID)
				})

			})

		})

	})

	s.assert(t, shorthand{
		{"A", "1"},
		{"B", "1"},
		{"C", "1"},
		{"D", "1"},
		{"E", "1"},
		{"F", "1"},
		{"G", "1"},
		{"P", "1"},
		{"Q", "1"},
		{"R", "1"},
		{"Y", "1"},
		{"Z", "1"},

		{"A", "2"},
		{"B", "2"},
		{"C", "2"},
		{"H", "2"},
		{"I", "2"},
		{"J", "2"},
		{"K", "2"},
		{"L", "2"},
		{"M", "2"},
		{"P", "2"},
		{"S", "2"},
		{"T", "2"},
		{"U", "2"},
		{"V", "2"},
		{"Y", "2"},
		{"Z", "2"},

		{"A", "3"},
		{"B", "3"},
		{"C", "3"},
		{"N", "3"},
		{"O", "3"},
		{"P", "3"},
		{"S", "3"},
		{"T", "3"},
		{"W", "3"},
		{"X", "3"},
		{"Y", "3"},
		{"Z", "3"},
	})
}

func TestMalformed(t *testing.T) {
	// TODO
}
