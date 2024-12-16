# Contributing to Walrus

Thank you for considering contributing to Walrus! Your support will help make this language better for everyone.

## Ways to Contribute

### 1. Report Issues
- Found a bug or an unexpected behavior? Please [open an issue](https://github.com/itsfuad/walrus/issues).
- When reporting an issue, include:
  - Clear steps to reproduce the issue.
  - The version of Walrus or commit hash.
  - Relevant code snippets, if applicable.

### 2. Suggest Features
- Have an idea to improve Walrus? Feel free to [submit a feature request](https://github.com/itsfuad/walrus/issues).
- Be as descriptive as possible and explain the motivation behind the feature.

### 3. Submit Pull Requests
- Contributions are welcome for bug fixes, new features, documentation updates, and more.
- Make sure your changes align with the project's [design principles](https://github.com/itsfuad/walrus#language-design-principles).

#### Steps to Submit a Pull Request:
1. Fork the repository and clone it locally.
2. Create a new branch for your changes:
   ```bash
   git checkout -b feature-name
   ```
3. Make your changes and commit them with a clear message:
   ```bash
   git commit -m "Add: Description of your change"
   ```
4. Push your changes to your forked repository:
   ```bash
   git push origin feature-name
   ```
5. Open a pull request to the main repository and describe your changes in detail.

### 4. Improve Documentation
- Help keep the documentation accurate and up-to-date.
- Typos, missing examples, and clarifications are welcome contributions.

## Guidelines for Code Contributions

### Code Style
- Follow existing coding conventions used in the project.
- Write clear and concise code.
- Avoid introducing unnecessary dependencies.

### Testing
- Ensure all existing tests pass before submitting changes.
- Add new tests to cover your changes if applicable:
  ```bash
  go test ./...
  ```

### Commit Messages
- Use meaningful and descriptive commit messages.
- Example: `Fix: Correct type casting for arrays`.

### Pull Request Checklist
- Ensure your code compiles and passes all tests.
- Avoid including unrelated changes.
- Provide a clear description of the change and its purpose.

## Getting Started with the Codebase

1. Install [Go](https://golang.org/dl/) and ensure it is properly configured.
2. Clone the repository:
   ```bash
   git clone https://github.com/itsfuad/walrus.git
   cd walrus
   ```
3. Test the project setup:
   ```bash
   go run main.go
   ```
4. Run tests to ensure everything is working:
   ```bash
   go test ./...
   ```

## Communication
- For questions and discussions, feel free to [start a discussion](https://github.com/itsfuad/walrus/discussions).
- For urgent matters, contact the maintainers via email (add contact info here).

We appreciate your contributions and look forward to building Walrus together!
