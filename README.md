# A third-party API server implementation for Rotaeno
## Implemented functions:
1. Use official cloud services for cloud synchronization and other cloud services.
2. Provide additional purchasing information for other accounts.

## How to use
1. [Download](https://go.dev/dl/) the go compiler for your platform.
2. Go to the project directory and type `go build`.
3. Run the compiled executable file.
4. Modify `config.json` in the root directory and restart program.
- You should modify `baseObjectID` to the object id of the user you want to share.
- You should modify `keys` in the format of `map[string]string` to make the key string correspond to the DLC package name.
