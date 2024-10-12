# Mod Scavenger: Universal Game Mod Manager

This tool is designed to manage mods for **any game** that supports local mod installations by creating symlinks from a source mod directory to a destination mods directory. While it supports any game, it currently works best with **Minecraft Forge** servers as a primary use case. The tool uses a lock file to track which mods are symlinked and to ensure that only active mods are linked, removing outdated symlinks when necessary.

## Features

- Supports **locally installed mods** for multiple games. 
- Manages mods by creating symlinks from specified source mod directories to the game's mods directory.
- Tracks mods with a lock file to ensure only active mods are linked.
- Removes outdated symlinks when mods are removed from the source directory.
- Provides verbose output for detailed operation logs.

## Supported Games
- **Minecraft Forge** server & client

## Planned Updates

- Remote mod management: Future updates will allow downloading and managing mods from remote sources.

## Directory Structure

Before running the tool, ensure your directory structure is set up as follows for **Minecraft Forge** or any other supported game:

```
game-server/
│
├── manualmods/                 # This directory contains the mods specific to your server/game. 
│   ├── server/                 # Server-only mods
│   │   ├── mod1.jar
│   │   └── ...
│   ├── client/                 # Client-only mods
│   │   ├── mod2.jar
│   │   └── ...
│   └── mod3.jar                # Mods in manualmods root are considered enviornment agnostic
│
├── mods/                       # This is the directory where the game looks for mods
│   ├── <symlink to mod1.jar>
│   ├── <symlink to mod2.jar>
│   └── ...
│
├── mods.lock                   # Lock file to track which mods are linked
│
└── mod_scavenger_file_manager  # The executable for this tool

```

- **server-mods/**: This directory contains all the mods you want to load on your Forge server.
- **mods/**: This is where symlinks to mods from `server-mods/` are created. Minecraft Forge will load mods from this directory.
- **mods.lock**: This lock file keeps track of which mods are linked from `server-mods` to `mods`.
- **server-mod_manager**: The executable for this tool.

**manualmods/**: The directory that contains the mods to be linked, separated into server/, client/, and base directory as agnostic.
**mods/**: The destination directory where symlinks are created, typically the folder where the game loads mods.
**mods.lock**: The lock file that tracks which mods are symlinked and ensures consistency between runs.
**mod_scavenger_file_manager**: The executable for this tool.


## Usage

### 1. Set Up Your Directories

Ensure you have the `manualmods` and `mods` directories set up as shown above.

### 2. Running the Tool

You can run the tool with the following command:

```bash
./mod_scavenger_file_manager --verbose
```

#### Flags:

- `--verbose`: Enables detailed output, showing which mods are being linked, unlinked, or skipped, along with any errors encountered.

### 3. What Happens When You Run It

- The tool first checks if a mods.lock file exists. If it doesn't, one will be created.
- It then scans the manualmods directory for mods (files with .jar or any other relevant extension).
- For each mod found in `server-mods`, it will check if a symlink already exists in the `mods` directory. If not, it will create a new symlink.
  - **Server mods** are found in manualmods/server/.
  - **Client mods** are found in manualmods/client/.
  - **Agnostic mods** are found in manualmods/.
- For each mod found, it checks if a symlink exists in the `mods` directory. If not, a symlink is created.
- Mods removed from manualmods will have their symlinks removed from the `mods` directory.
- The mods.lock file is updated to reflect the current state of symlinked mods.

## Example Shell Script to Automate the Process

To ensure the mod manager runs before starting your game (e.g., a Minecraft Forge server), you can include it in a shell script:
 
```bash
#!/bin/bash

# Step 1: Run the mod manager to update the symlinks
echo "Updating mods with mod manager..."
/path/to/mod_scavenger_file_manager --verbose

# Step 2: Start the game server
echo "Starting the game server..."
java -Xmx1024M -Xms1024M -jar game-server.jar nogui

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

The lock file (`mods.lock`) is used to keep track of the mods that have been symlinked. It ensures that:

- New mods added to `manualmods` are symlinked.
- Mods removed from `manualmods` are unlinked from the `mods` directory.
- The state of the symlinked mods is preserved between runs.

## Example Lock File

```yaml
mods:
- name: resurgent-essentials-0.5.0.jar
  version: 1.0.0
  environment: server
  source: local
  file_path: manualmods/server/resurgent-essentials-0.5.0.jar
- name: SuperChickenMod-2.0.jar
  version: 1.0.0
  environment: client
  source: local
  file_path: manualmods/client/SuperChickenMod-2.0.jar
- name: LaserPandas-1.3.jar
  version: 1.0.0
  environment: server
  source: local
  file_path: manualmods/server/LaserPandas-1.3.jar
- name: InfiniteBeef-4.7.2.jar
  version: 1.0.0
  environment: server
  source: local
  file_path: manualmods/server/InfiniteBeef-4.7.2.jar
- name: GiantSandwich-1.0.1.jar
  version: 1.0.0
  environment: agnostic
  source: local
  file_path: manualmods/GiantSandwich-1.0.1.jar
```

## Notes

- This tool is built to work with **locally installed mods** for any game that supports mods through a folder system, with **Minecraft Forge** as a fully supported platform.
- The tool will require read/write permissions to the `manualmods`, `mods`, and `mods.lock` files.
- Running the tool is recommended anytime mods are added, removed, or modified to ensure the mods folder is properly synced.


