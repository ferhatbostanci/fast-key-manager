# Fast Key Manager (fkm)

## 1. Introduction

Fast Key Manager (fkm) is a CLI tool designed to help users manage SSH keys on Linux and other Unix-based systems. This tool allows users to add SSH keys from their GitHub or GitLab accounts, list existing SSH keys, and remove selected keys from their system.

## 2. Purpose

The purpose of this project is to provide a convenient way for users to manage their SSH keys. By allowing users to fetch and add SSH keys directly from their GitHub or GitLab accounts, the tool aims to simplify the process of SSH key management.

## 3. Target Audience

- Developers
- System Administrators
- DevOps Engineers

## 4. Features

### 4.1 Add SSH Key

#### 4.1.1 Manual Copy-Paste
- Users can manually copy an SSH key and add it to their system using the CLI tool.

#### 4.1.2 Add from GitHub
- Users can add SSH keys from their GitHub account by providing their GitHub username.

#### 4.1.3 Add from GitLab
- Users can add SSH keys from their GitLab account by providing their GitLab username.

### 4.2 Remove SSH Key

- Users can list existing SSH keys and select the one they want to remove from their system.

### 4.3 List SSH Keys

- Users can list all existing SSH keys. The output will show the short name (comment) and the key itself.

### 4.4 Automated Releases
- GitHub Actions workflow automatically creates releases when a tag is pushed to the master branch
- Release artifacts include:
  - Binary files for supported platforms (Linux, macOS)
  - Source code archives
  - Release notes based on commit history
- Semantic versioning (vX.Y.Z) is used for release tags

## 5. Technical Requirements

### 5.1 Programming Language
- Go

### 5.2 Platforms
- Linux
- Other Unix-based distributions

### 5.3 Dependencies
- GitHub API
- GitLab API
- github.com/spf13/cobra (for CLI framework)
- github.com/manifoldco/promptui (for interactive command-line interface)
- github.com/fatih/color (for colorful CLI output)

## 6. User Interface

The CLI tool will be used as follows:

### Add SSH Key
```
# Manual Copy-Paste (Default)
fkm add --manual

# Add from GitHub
fkm add --github <username>

# Add from GitLab
fkm add --gitlab <username>
```

### Remove SSH Key
```
fkm remove
```
Users will be prompted to select the SSH key they want to remove from a list of existing keys.

### List SSH Keys
```
fkm list
```
Users can list all existing SSH keys, displaying the short name and the key itself.

## 7. Security and Privacy

- Users' SSH keys will be managed securely.
- API calls will be made with proper authentication to protect user information.

## 8. Project Timeline

### Phase 1: Requirements Gathering and Planning
- Identify requirements and create the PRD.
- Set up GitHub repository and project structure.

### Phase 2: Design and Development
- Design and develop the core functionalities of the CLI tool.

### Phase 3: Testing and Validation
- Test and validate the developed tool.

### Phase 4: Release and Maintenance
- Set up automated release process.
- Release the tool and provide updates based on user feedback.

## 9. Release Process

### 9.1 Automated Releases
- Releases are automatically created when a tag is pushed to the master branch
- Tag format must follow semantic versioning (e.g., v1.0.0)
- GitHub Actions workflow will:
  1. Build binaries for all supported platforms
  2. Create a GitHub release
  3. Attach build artifacts
  4. Generate release notes

### 9.2 Release Artifacts
- Linux amd64 binary
- macOS amd64 binary
- Source code (zip and tar.gz)
- SHA256 checksums for all artifacts
- Changelog/Release notes

## 10. Conclusion

Fast Key Manager (fkm) aims to simplify SSH key management for users by allowing them to easily add, list, and remove SSH keys from their system. With automated releases through GitHub Actions, the tool ensures consistent and reliable delivery of updates to users.
