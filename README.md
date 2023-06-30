# lsmake

[![License](https://img.shields.io/badge/license-LGPL--3.0-blue.svg)](https://github.com/servel333/lsmake/blob/main/LICENSE)

`lsmake` is a command-line utility written in Go that lists the targets in a Makefile. It provides a convenient way to view the available targets without having to manually inspect the Makefile.

## Features

- Retrieve a list of targets defined in a Makefile
- Simplify navigation and understanding of Makefile targets
- Streamline build process by quickly identifying available options

<!--
## Precompiled Binary

You can download precompiled binaries for `lsmake` from the [Releases](https://github.com/servel333/lsmake/releases) page.
Z-->

## Build from Source

To build `lsmake` from source, ensure you have Go installed and then run:

```bash
go get github.com/servel333/lsmake
```

## Usage

To use lsmake, simply navigate to the directory containing your Makefile and run:

```bash
lsmake
```

This will display a list of targets defined in the Makefile.

```bash
lsmake --help
```

This will display the help message.

```bash
lsmake ~/workspace/my-project/Makefile
```

This will display the targets of `~/workspace/my-project/Makefile`

## Contributing

Contributions are welcome! If you find a bug or have a suggestion, please open an issue or submit a pull request.

## License

This project is licensed under the GNU Lesser General Public License v3.0 (LGPL-3.0).
