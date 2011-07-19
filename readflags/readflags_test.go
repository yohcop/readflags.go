package readflags

import (
  "flag"
  "testing"
)

var a *string = flag.String("a", "a", "")
var b *int = flag.Int("b", 0, "")

func TestSimple(t *testing.T) {
  file := "a = TestSimple"
  if err := ReadFlagsFromString(file); err != nil {
    t.Error("Error: " + err.String())
  }
  if *a != "TestSimple" {
    t.Error("Flag not set: " + *a)
  }
}

func TestOtherTypes(t *testing.T) {
  file := "a = TestOtherTypes\nb=2"
  if err := ReadFlagsFromString(file); err != nil {
    t.Error("Error: " + err.String())
  }
  if *a != "TestOtherTypes" {
    t.Error("Flag not set: " + *a)
  }
  if *b != 2 {
    t.Error("Flag not set: " + string(*b))
  }
}


func TestIgnoreWhitespaces(t *testing.T) {
  file := "  a \t = \t TestIgnoreWhitespaces\n\t  \tb=3\t "
  if err := ReadFlagsFromString(file); err != nil {
    t.Error("Error: " + err.String())
  }
  if *a != "TestIgnoreWhitespaces" {
    t.Error("Flag not set: " + *a)
  }
  if *b != 3 {
    t.Error("Flag not set: " + string(*b))
  }
}

func TestIgnoreComments(t *testing.T) {
  // Also ignore indented comments.
  file := "#comment\na=TestIgnoreComments\n #foo\nb=4\n\t#cc"
  if err := ReadFlagsFromString(file); err != nil {
    t.Error("Error: " + err.String())
  }
  if *a != "TestIgnoreComments" {
    t.Error("Flag not set: " + *a)
  }
  if *b != 4 {
    t.Error("Flag not set: " + string(*b))
  }
}


func TestError(t *testing.T) {
  file := "a = TestError\nnonexisting=asd"
  if err := ReadFlagsFromString(file); err == nil {
    t.Error("Error: noexisting flag should not exist")
  }
}

func TestFile(t *testing.T) {
  err := ReadFlagsFromFile("testdata/f1.txt")
  if err != nil {
    t.Error("Unexpected error: " + err.String())
  }
  if *a != "f3" {
    t.Error("Flag not set: " + *a)
  }
  if *b != 111 {
    t.Error("Flag not set: " + string(*b))
  }
}

func TestFileError(t *testing.T) {
  err := ReadFlagsFromFile("testdata/err.txt")
  if err == nil {
    t.Error("Error should have triggered")
  }
}
