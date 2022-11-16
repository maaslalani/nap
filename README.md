# Nap

![](https://user-images.githubusercontent.com/42545625/202176622-a822dc10-f9fb-46bf-943e-a6dc20016618.png)

Nap is a code snippet manager for your terminal. Access and create new snippets
quickly with the command-line interface or browse and manage them with the
text-user interface.

Keep your code snippets safe, sound, and well-rested, right from your terminal.

<br />

<p align="center">
<img width="1000" src="https://user-images.githubusercontent.com/42545625/202191135-0eca8fbc-a216-4c00-a3f2-ef6ce41f011f.gif" />
</p>

<br />

## Text-based User Interface

Launch the interactive interface:

```bash
nap
```

![](https://user-images.githubusercontent.com/42545625/202177235-b34a8e0b-9d35-48aa-998b-a9b0583c46d5.png)

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
```

Output saved snippets:

```bash
# Write snippet to a file.
nap go/boilerplate > main.go

# Copy snippet to clipboard.
nap foobar | pbcopy
nap foobar | xclip
```

List snippets:

```bash
nap list
```

List all snippets, with `--filter`:

```bash
nap list -f search
```

Fuzzy find a snippet (with [Gum](https://github.com/charmbracelet/gum)).

```bash
nap $(nap list | gum filter)
```

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

## License

[MIT](https://github.com/maaslalani/nap/blob/master/LICENSE)


## Feedback

I'd love to hear your feedback on improving `nap`.
Feel free to reach out via [email](mailto:maas@lalani.dev) or [twitter](https://twitter.com/maaslalani), or [create an issue](https://github.com/maaslalani/nap/issues/new)!

---

<sub><sub>z</sub></sub><sub>z</sub>z
