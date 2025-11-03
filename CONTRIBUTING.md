# Contributing

Thank you for your interest in contributing to Owl! 🦉

## Prerequisites

1. [Install Go][go-install] (1.22 or later)
2. Fork the repository on GitHub
3. Clone your fork:

   ```bash
   git clone https://github.com/YOUR_USERNAME/owl.git
   cd owl
   git remote add upstream https://github.com/go-owl/owl.git
   ```

## Development Setup

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run specific example
cd _example/helloworld
go run main.go
```

## Submitting a Pull Request

1. [Fork the repository.][fork]
2. [Create a topic branch.][branch]
3. Add tests for your change.
4. Run `go test ./...` to ensure tests pass.
5. Implement the change and ensure tests still pass.
6. Run `goimports -w .` to format code.
7. [Commit and push your changes.][git-help]
8. [Submit a pull request.][pull-req]

## Guidelines

- Follow Go best practices and idioms
- Add tests for new features
- Update documentation as needed
- Keep commits focused and atomic
- Write clear commit messages

[go-install]: https://golang.org/doc/install
[fork]: https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo
[branch]: https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-branches
[git-help]: https://docs.github.com/en
[pull-req]: https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-pull-requests
