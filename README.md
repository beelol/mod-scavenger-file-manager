# Minecraft Forge Server Mod Manager

This tool is designed to manage mods on a **Minecraft Forge** server by creating symlinks from the `server-mods` directory to the `mods` directory. It uses a lock file (`server-mods.lock`) to track the mods that have been symlinked. The lock file ensures that only the mods currently in the `server-mods` directory are linked, and any mods that have been removed from `server-mods` are unlinked from the `mods` directory.

## Features

- Automatically manages mods for **Minecraft Forge** servers.
- Creates symlinks in the `mods` directory based on the mods found in `server-mods`.
- Keeps track of mods in a lock file to ensure only active mods are linked.
- Removes outdated symlinks if a mod has been removed from `server-mods`.
- Verbose mode for detailed output.

## Directory Structure

Before running the tool, ensure your directory structure is set up as follows:

```
minecraft-server/
│
├── server-mods/        # This directory contains the mods specific to your server
│   ├── mod1.jar
│   ├── mod2.jar
│   └── ...
│
├── mods/               # This is the directory where the server looks for mods (Forge mods folder)
│   ├── <symlink to mod1.jar>
│   ├── <symlink to mod2.jar>
│   └── ...
│
├── server-mods.lock    # Lock file to track which mods are linked
│
└── server-mod_manager         # The executable for this tool
```

- **server-mods/**: This directory contains all the mods you want to load on your Forge server.
- **mods/**: This is where symlinks to mods from `server-mods/` are created. Minecraft Forge will load mods from this directory.
- **server-mods.lock**: This lock file keeps track of which mods are linked from `server-mods` to `mods`.
- **server-mod_manager**: The executable for this tool.

## Usage

### 1. Set Up Your Directories

Make sure you have the `server-mods` and `mods` directories set up as shown in the directory structure above.

### 2. Running the Tool

You can run the tool with the following command:

```bash
./server-mod_manager --verbose
```

#### Flags:

- `--verbose`: Enables detailed output, showing what mods are being linked, which symlinks are being removed, and any errors that occur.

### 3. What Happens When You Run It

- The tool will first check if a `server-mods.lock` file exists. If it doesn't, it will create one.
- It will then check the `server-mods` directory for any mods (files with `.jar` extension).
- For each mod found in `server-mods`, it will check if a symlink already exists in the `mods` directory. If not, it will create a new symlink.
- Any mods that no longer exist in `server-mods` will have their symlinks removed from `mods`.
- The `server-mods.lock` file is updated to reflect the current state of symlinked mods.

## Example Shell Script to Include Before Starting the Minecraft Server

To ensure the mod manager runs before starting the Minecraft Forge server, you can include it in a shell script:

```bash
#!/bin/bash

# Step 1: Run the mod manager to update the symlinks
echo "Updating mods with mod manager..."
/path/to/server-mod_manager --verbose

# Step 2: Start the Minecraft Forge server
echo "Starting Minecraft Forge server..."
java -Xmx1024M -Xms1024M -jar forge-server.jar nogui
```

Make this shell script executable by running:

```bash
chmod +x start_forge_server.sh
```

Then, you can use it to start your server:

```bash
./start_forge_server.sh
```

## Lock File

The lock file (`server-mods.lock`) is used to keep track of which mods are symlinked to the `mods` directory. It ensures that:

- New mods added to `server-mods` are symlinked.
- Mods removed from `server-mods` are unlinked.
- The state of the symlinked mods is preserved between runs.

## Notes

- This tool is designed specifically for **Minecraft Forge servers**.
- Ensure that the tool has proper permissions to read from `server-mods`, write to `mods`, and create/modify the `server-mods.lock` file.
- If mods are manually added or removed from the `mods` directory, it is recommended to run the tool again to ensure the state is synced properly.