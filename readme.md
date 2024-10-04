Sure, here's an improved version of your `readme.md` documentation:

# AI Code Reviewer

## Prerequisite

Prepare a `.env` file with the following content:

```ini
AZURE_OPENAI_ENDPOINT=https://xxx.openai.azure.com
AZURE_OPENAI_KEY=
AZURE_OPENAI_MODEL=
AZURE_OPENAI_API_VERSION=
```

## How to Use

### 1. Compile the Code

#### 1.1 Download Dependencies

Run the following command to download all the dependencies specified in your `go.mod` file:

```sh
go mod tidy
```

#### 1.2 Compile the Code

To compile your Go program, run:

- **For Windows:** 
  ```sh
  go build -o code-reviewer.exe
  ```
- **For Mac/Linux:** 
  ```sh
  go build -o code-reviewer
  ```

#### 1.3 Optional: Add Binary to PATH

You can add the path of the binary to your `PATH` environment variable for easier access.

## Features

### Version 0.1

- Supported AI Generator: Azure OpenAI

## Usage

Run the compiled binary with the following options:

```sh
Usage: code-reviewer [options]

Options:
  -scandir string
        Directory to scan (default ".")
  -exclude string
        Comma-separated list of directories, files, extensions, or regex patterns to exclude
  -extensions string
        Comma-separated list of extensions used
  -keywords string
        Comma-separated list of keywords to filter files
```

## Example

To scan the current directory and exclude certain files or directories, you can use:

```sh
./code-reviewer -scandir . -exclude "dir1,file1.go,*.tmp" -extensions ".go,.py" -keywords "TODO,FIXME"
```

This will scan the current directory, excluding `dir1`, `file1.go`, and any `.tmp` files, looking for files with `.go` or `.py` extensions that contain the keywords `TODO` or `FIXME`.

---

Feel free to reach out if you have any questions or need further assistance!