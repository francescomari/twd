# Temporary Working Directory

Temporary Workind Directory (twd) creates a temporary directory and executes a
command into it. Since the implementation uses `execve(2)` to replace the
current process with the provided command, the temporary directory is not
removed afterward. See [chain
loading](https://en.wikipedia.org/wiki/Chain_loading) for further information
about the design choice.

## Install

    go get github.com/francescomari/twd

## Usage

    twd [-print] [-root root] [-prefix prefix] command...

If `-print` is specified, the path to the temporary directory is printed on the
standard output before executing the command. `root` and `prefix` are used when
creating the temporary directory. If specified, the temporary directory will be
created under `root`. If `-root` is not specified, the default system directory
for temporary files is used. `prefix` is a prefix for the name of the temporary
directory. If `-prefix` is not specified, the default prefix of `twd-` is used.

## Example

The following command shows how to run `touch test.txt` in a temporary
directory. The directory will be crated in the current directory with a prefix
of `whatever`.

    $ twd -root . -prefix whatever touch test.txt
    $ ls whatever*
    test.txt

The following command shows how to run `touch test.txt` in a temporary directory
created in the default system directory for temporary files and with the default
prefix. The path to the temporary folder is printed before the command is
executed.

    $ twd -print touch test.txt
    /var/folders/s5/337j21rj3m7cws2bql5m8_rc0000gp/T/twd-520522576
    $ ls /var/folders/s5/337j21rj3m7cws2bql5m8_rc0000gp/T/twd-520522576
    test.txt

## License

This software is released under the MIT license.