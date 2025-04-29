
# LND-Fuzz

**LND-Fuzz** is a Go-native fuzzing tool that scans your repository for fuzz targets, runs them concurrently, and persists the generated input corpus in a Git storage repo. It helps you continuously test and improve your codebase’s resilience with minimal configuration.

---

##  Features

- **Automatic Target Detection**  
  Scans your Go modules for all functions starting with `Fuzz`.
- **Concurrent Execution**  
  Spawns multiple fuzzing processes in parallel (default: number of CPU cores).
- **Customizable Configuration**  
  Control run duration, target package, and concurrency via environment variables.
- **Corpus Persistence**  
  Pushes new inputs to a designated Git repository after each run for regression testing.

---

## Configuration

Create a `.env` file in the project root with the following variables:

| Variable              | Required | Description                                                                                       | Default                   |
|-----------------------|----------|---------------------------------------------------------------------------------------------------|---------------------------|
| `FUZZ_NUM_PROCESSES`  | No       | Number of concurrent fuzzers.                                                                      | CPU core count            |
| `PROJECT_SRC_PATH`    | Yes      | Git URL of the repository to fuzz.                                                                | —                         |
| `GIT_STORAGE_REPO`    | Yes      | Git URL of the repository where corpus will be pushed (must allow push with your credentials).    | —                         |
| `FUZZ_TIME`           | No       | Duration for each fuzz cycle (in seconds).                                                        | `120`                     |
| `FUZZ_PKG`            | Yes      | Go package path to target for fuzzing (e.g., `github.com/OWNER/REPO/pkg/target`).                 | —                         |

Example `.env`:

```env
FUZZ_NUM_PROCESSES=4
PROJECT_SRC_PATH=https://github.com/OWNER/REPO.git
GIT_STORAGE_REPO=https://oauth2:TOKEN@github.com/OWNER/STORAGEREPO.git
FUZZ_TIME=180
FUZZ_PKG=github.com/OWNER/REPO/pkg/target
```

---

## ⚙️ Usage

1. **Clone the repository**:
   ```bash
   git clone https://github.com/OWNER/LND-Fuzz.git
   cd LND-Fuzz
   ```
2. **Install dependencies** (requires Go 1.20+):
   ```bash
   go mod download
   ```
3. **Run the fuzzer**:
   ```bash
   go run main.go
   ```

Each run will:
- Detect all fuzz targets in `FUZZ_PKG`.
- Execute for `FUZZ_TIME` seconds.
- Commit and push any new corpus files to `GIT_STORAGE_REPO`.

---

##  Deployment

Deploy LND-Fuzz as a long-running service on any cloud VM (e.g., AWS EC2, GCP, DigitalOcean). No external CI required:

```bash
# On your server
git clone https://github.com/OWNER/LND-Fuzz.git
cd LND-Fuzz
# set up .env
go run main.go &
```

The tool will restart automatically after each cycle defined by `FUZZ_TIME`.

---

##  Contributing

1. Fork the repo
2. Create a branch: `git checkout -b feature/your-feature`
3. Commit your changes: `git commit -m "Add awesome feature"`
4. Push and open a Pull Request

Please follow standard Go conventions and include tests for new functionality.

---

##  License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---


