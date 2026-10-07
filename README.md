Who needs [neofetch](https://en.wikipedia.org/wiki/Neofetch) when you have rain?

![Rain](https://raw.githubusercontent.com/Mortimyrrh/rain/main/exsample.gif)

✯ Butter smooth 120fps
✯ Coherent coloring, using the terminal ansi colors.
✯ Customizable via cli args
✯ Built on ![bubbletea](https://github.com/charmbracelet/bubbletea)

Live webpage https://mortimyrrh.github.io/rain/ (sometimes it needs a refresh)
NB: Renders much faster in a real terminal


### Customization

use `--density (float 0-1)` for more or less lines, 0.05 is default.
use `--speed (float)` to set speed. 1 for default speed, 2 = double speed, .5 half speed etc.
use `--multi` flag for multicolour mode
![Multi](https://raw.githubusercontent.com/Mortimyrrh/rain/main/multi.png)
As colours are rendered by the terminal, they can be set by changing your terminal themes ansi colors.


### Dependencies

- Go 1.25 needs to be installed and in your PATH see [Go install](https://go.dev/doc/install) for instructions.


### Install

From source: Clone the repo and follow [this guide](https://go.dev/doc/tutorial/compile-install).

Direct install: 'go install gitlab.com/mortimyrrh/rain@latest'

Add `rain` to the end of your ~/.<shell>rc of choice for cool splash screen when opening a new terminal.


### Troubleshooting

A GPU accelerated terminal can help with rendering on older machines
