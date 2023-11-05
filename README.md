# 🔛 arrow

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](http://opensource.org/licenses/MIT)

![image](./images/arrow.jpg)

arrow is a CLI tool specialized for moving directories that can be used as a `cd` replacement.

To move directories, simply use the up, down, left, and right arrow keys to select the directory and press the Enter key.

## Install

```bash
$ go install github.com/harehare/arrow@latest
```

### Zsh

```bash
function ao() {
  cd $(arrow --icons)
}
```

## Usage

| Key binding  | Description                                  |
| ------------ | -------------------------------------------- |
| `Up`, `Down` | Move cursor                                  |
| `Right`      | Move directory                               |
| `Enter`      | Select directory                             |
| `Shift+Down` | Change order (directory name, modified time) |
| `Ctrl+c`     | Exit                                         |

```
USAGE:
   arrow [options]

OPTIONS:
   --all, -a                Show hidden files. (default: false)
   --icons, -i              Display icons. (default: false)
   --query value, -q value  Specifies a query to search the directory.
   --help, -h               show help
   --version, -V            print only the version (default: false)
```

## Customization

ANSI 256 Colors or HEX

```bash
export ARROW_BORDER_COLOR="80"
export ARROW_CURRENT_DIRECTORY_COLOR="57"
export ARROW_CURSOR_COLOR="57"
export ARROW_DISABLED_COLOR="240"
export ARROW_FOREGROUND_COLOR="#009CD1"
export ARROW_HIGHLIGHT_COLOR="80"
export ARROW_PROMPT_COLOR="36"
export ARROW_SYMLINK_COLOR="36"
```

## Run

```sh
just run
```

## License

MIT
