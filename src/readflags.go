package readflags

import (
  "flag"
  "io"
  "io/ioutil"
  "os"
  "strings"
)


func ReadFalgs(reader io.Reader) os.Error {
  if content, err := ioutil.ReadAll(reader); err != nil {
    return err
  } else {
    return ReadFlags(string(content))
  }
  return nil
}

func ReadFlags(content string) os.Error {
  lines := strings.Split(content, "\n", -1)

  for n, line := range lines {
    clean := strings.TrimSpace(line)
    if len(clean) <= 0 {
      continue
    }
    if clean[0] == '#' {
      continue
    }
    pieces := strings.Split(line, "=", 2)
    if len(pieces) != 2 {
      continue
    }
    key := strings.TrimSpace(pieces[0])
    val := strings.TrimSpace(pieces[1])
    if !flag.Set(key, val) {
      return os.NewError("readflags: no such flag: " + key +
          " (l." + string(n) + ")")
    }
  }
  return nil
}

