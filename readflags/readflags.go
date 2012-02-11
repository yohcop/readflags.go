package readflags

import (
  "flag"
  "fmt"
  "io"
  "io/ioutil"
  "path"
  "path/filepath"
  "strings"
)

// Reads the entire Reader and parse flags from it.
func ReadFlags(reader io.Reader) error {
  if content, err := ioutil.ReadAll(reader); err != nil {
    return err
  } else {
    return ReadFlagsFromString(string(content))
  }
  return nil
}

// Reads the file and parses flags from it. This also
// understands the %include lines.
func ReadFlagsFromFile(file string) error {
  absPath, err := filepath.Abs(file)
  if err != nil {
    return err
  }
  absPath, _ = path.Split(absPath)

  if bytes, err := ioutil.ReadFile(file); err != nil {
    return err
  } else {
    content := string(bytes)
    lines := strings.Split(content, "\n")
    for n, line := range lines {
      err := parseLine(line)
      if err != nil {
        return fmt.Errorf("%s (line %d)", err.Error(), n)
      }
      err = parseCommand(line, absPath)
      if err != nil {
        return fmt.Errorf("%s (line %d)", err.Error(), n)
      }
    }
  }
  return nil
}

// Reads and sets flag from the given string.
func ReadFlagsFromString(content string) error {
  lines := strings.Split(content, "\n")

  for n, line := range lines {
    err := parseLine(line)
    if err != nil {
      return fmt.Errorf("%s (line %d)", err.Error(), n)
    }
  }
  return nil
}

func parseLine(line string) error {
  clean := strings.TrimSpace(line)
  if len(clean) <= 0 {
    return nil
  }
  if clean[0] == '#' {
    return nil
  }
  if clean[0] == '%' {
    return nil
  }
  pieces := strings.SplitN(line, "=", 2)
  if len(pieces) != 2 {
    return fmt.Errorf("readflags: misformatted line: %s", line)
  }
  key := strings.TrimSpace(pieces[0])
  val := strings.TrimSpace(pieces[1])
  if flag.Set(key, val) != nil {
    return fmt.Errorf("readflags: no such flag: %s", key)
  }
  return nil
}

func parseCommand(line string, path string) error {
  clean := strings.TrimSpace(line)
  if len(clean) == 0 {
    return nil
  }
  if clean[0] != '%' {
    return nil
  }
  if len(clean) < 10 {
    return fmt.Errorf("readflags: %%include needs a file name")
  }
  if clean[1:8] == "include" {
    return ReadFlagsFromFile(path + "/" + clean[9:])
  }
  return nil
}
