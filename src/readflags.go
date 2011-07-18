package readflags

import (
  "fmt"
  "flag"
  "io"
  "io/ioutil"
  "os"
  "strings"
  "path"
  "path/filepath"
)


func ReadFlags(reader io.Reader) os.Error {
  if content, err := ioutil.ReadAll(reader); err != nil {
    return err
  } else {
    return ReadFlagsFromString(string(content))
  }
  return nil
}

func ReadFlagsFromFile(file string) os.Error {
  absPath, err := filepath.Abs(file)
  if err != nil {
    return err
  }
  absPath, _ = path.Split(absPath)

  if bytes, err := ioutil.ReadFile(file); err != nil {
    return err
  } else {
    content := string(bytes)
    lines := strings.Split(content, "\n", -1)
    for n, line := range lines {
      err := parseLine(line)
      if err != nil {
        return fmt.Errorf("%s (line %d)", err.String(), n)
      }
      err = parseCommand(line, absPath)
      if err != nil {
        return fmt.Errorf("%s (line %d)", err.String(), n)
      }
    }
  }
  return nil
}

func ReadFlagsFromString(content string) os.Error {
  lines := strings.Split(content, "\n", -1)

  for n, line := range lines {
    err := parseLine(line)
    if err != nil {
      return fmt.Errorf("%s (line %d)", err.String(), n)
    }
  }
  return nil
}

func parseLine(line string) os.Error {
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
  pieces := strings.Split(line, "=", 2)
  if len(pieces) != 2 {
    return fmt.Errorf("readflags: misformatted line: %s", line)
  }
  key := strings.TrimSpace(pieces[0])
  val := strings.TrimSpace(pieces[1])
  if !flag.Set(key, val) {
    return fmt.Errorf("readflags: no such flag: %s", key)
  }
  return nil
}

func parseCommand(line string, path string) os.Error {
  clean := strings.TrimSpace(line)
  if len(clean) == 0 {
    return nil
  }
  if clean[0] != '%' {
    return nil
  }
  if len(clean) < 10 {
    return nil
  }
  if clean[1:8] == "include" {
    return ReadFlagsFromFile(path + "/" + clean[9:])
  }
  return nil
}

