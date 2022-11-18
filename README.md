# Nap

<img width="1200" alt="Nap" src="https://user-images.githubusercontent.com/42545625/202545409-eb53f92a-233a-4f78-b598-a59c65248ad3.png">

<sub><sub>z</sub></sub><sub>z</sub>z

Nap is a code snippet manager for your terminal. Create and access new snippets
quickly with the command-line interface or browse, manage, and organize them with the
text-user interface.

Keep your code snippets safe, sound, and well-rested, right from your terminal.

<br />

<p align="center">
<img width="1000" src="https://user-images.githubusercontent.com/42545625/202577549-f2e0887a-b740-41f4-9408-c2f53673503f.gif" />
</p>

<br />

## Text-based User Interface

Launch the interactive interface:

```bash
nap
```

<img width="1000" src="https://user-images.githubusercontent.com/42545625/202766527-06b01d6a-9ad9-4728-8183-4efb9e05b335.gif" />

<details>

<summary>Key Bindings</summary>

<br />

| Action | Key |
| :--- | :--- |
| Create a new snippet | <kbd>n</kbd> |
| Edit selected snippet (in `$EDITOR`) | <kbd>e</kbd> |
| Copy selected snippet to clipboard | <kbd>c</kbd> |
| Paste clipboard to selected snippet | <kbd>p</kbd> |
| Delete selected snippet | <kbd>x</kbd> |
| Rename selected snippet | <kbd>r</kbd> |
| Set folder of selected snippet | <kbd>f</kbd> |
| Set language of selected snippet | <kbd>L</kbd> |
| Move to next pane | <kbd>tab</kbd> |
| Move to previous pane | <kbd>shift+tab</kbd> |
| Search for snippets | <kbd>/</kbd> |
| Toggle help | <kbd>?</kbd> |
| Quit application | <kbd>q</kbd> <kbd>ctrl+c</kbd> |

</details>

## Command Line Interface

Create new snippets:

```bash
# Quick save an untitled snippet.

nap < main.go

# From a file, specify Notes/ folder and Go language.
nap Notes/FizzBuzz.go < main.go

# Save some code from the internet for later.
curl https://example.com/main.go | nap Notes/FizzBuzz.go

# Works great with GitHub gists
gh gist view 4ff8a6472247e6dd2315fd4038926522 | nap
```

<img width="600" src="https://user-images.githubusercontent.com/42545625/202767159-134d679f-490f-4ad2-8875-cda604aa7b13.gif" />

Output saved snippets:

```bash
# Fuzzy find snippet.
nap fuzzy

# Write snippet to a file.
nap go/boilerplate > main.go

# Copy snippet to clipboard.
nap foobar | pbcopy
nap foobar | xclip
```

<img width="600" src="https://user-images.githubusercontent.com/42545625/202240249-d724fd73-2f90-4036-b9fc-6d2ccef982b3.gif" />

List snippets:

```bash
nap list
```
<img width="600" src="https://user-images.githubusercontent.com/42545625/202242653-1696dda6-2527-4c38-b673-74d67ad1517f.gif" />

Fuzzy find a snippet (with [Gum](https://github.com/charmbracelet/gum)).

```bash
nap $(nap list | gum filter)
```

<img width="600" src="https://user-images.githubusercontent.com/42545625/202240268-3a71fde6-73c3-4b0a-b129-f87ec1bb1b88.gif" />

## Installation

<!--

Use a package manager:

```bash
# macOS
brew install nap

# Arch
yay -S nap

# Nix
nix-env -iA nixpkgs.nap
```

-->

Install with Go:

```sh
go install github.com/maaslalani/nap@latest
```

Or download a binary from the [releases](https://github.com/maaslalani/nap/releases).


## Customization

Nap is customized through environment variables:
* `NAP_HOME`, the folder where your snippets will rest. Defaults to `$XDG_DATA_HOME/snooze`.
* `NAP_DEFAULT_LANGUAGE`, the language your snippets will use by default.
* `NAP_THEME`, the theme to highlight code. Defaults to `dracula`.

* `NAP_PRIMARY_COLOR` / `NAP_PRIMARY_COLOR_SUBDUED`, the color to use for the active pane title bars.
* `NAP_RED` / `NAP_BRIGHT_RED`, the colors to use for the selected item being deleted.
* `NAP_GREEN` / `NAP_BRIGHT_GREEN`, the colors to use for the selected item being copied.
* `NAP_FOREGROUND` / `NAP_BACKGROUND`, the colors to use for the foreground and background colors.
* `NAP_BLACK` / `NAP_WHITE` / `NAP_GRAY`, the colors to use for the unselected items.

<br />

<p align="center">
  <img
    width="1000"
    alt="image"
    src="https://user-images.githubusercontent.com/42545625/202580527-54c0364f-9f85-45ca-bd63-ad71462040a4.png"
  />
</p>

## License

[MIT](https://github.com/maaslalani/nap/blob/master/LICENSE)

## Feedback

I'd love to hear your feedback on improving `nap`.

Feel free to reach out via:
* [Email](mailto:maas@lalani.dev) 
* [Twitter](https://twitter.com/maaslalani)
* [GitHub issues](https://github.com/maaslalani/nap/issues/new)

---

<sub><sub>z</sub></sub><sub>z</sub>z
