
# Midgarts Client

Welcome to the **Midgarts Client**, a graphical client primarily built using SDL2, OpenGL, and various custom and third-party libraries. The project is designed to create an interactive and visually appealing environment for manipulating and rendering game characters and actions. The application showcases entities, systems, and OpenGL integration for real-time character movement and rendering.

Current Screenshots:

<p align="center">
    <img src="https://user-images.githubusercontent.com/696982/117575166-fff95980-b0b6-11eb-8afa-acd7dcdd6b34.gif" width="25%" />
    <img src="https://user-images.githubusercontent.com/696982/197043590-041d711b-a5d6-4d58-bf3c-8ea98c1afdc6.gif" width="50%" />
</p>
<p align="center">
    <img src="https://user-images.githubusercontent.com/696982/116827693-c2557780-ab70-11eb-90cd-b093004361db.gif" width="34%" />
    <img src="https://user-images.githubusercontent.com/696982/115995910-96a42180-a5b3-11eb-8200-1cfae06bf5bc.gif" width="34%" />
</p>

---

## Table of Contents

- [Project Overview](#project-overview)
- [Features](#features)
- [Requirements](#requirements)
- [Setup](#setup)
- [Usage](#usage)
- [Folder Structure](#folder-structure)
- [License](#license)

---

## Project Overview

The Midgarts Client uses the **Entity-Component-System (ECS)** architecture to manage game objects and their interactions. It focuses on simulating a modeled game environment, rendering scene objects (like characters), and controlling game entities.

The main goals include:
- Configuring characters with properties like job sprites, direction, position, and states.
- Using **OpenGL** for rendering.
- SDL2 for window and event management.
- Integration with GRF file format for asset loading.

---

## Features

1. **Character Creation and Rendering**:
   - Multiple characters with configurable sprites, positions, and states.
   - Supports movement and states like "Standing" and "Walking".

2. **OpenGL Integration**:
   - Real-time rendering of characters using a perspective camera.
   - Efficient use of OpenGL viewport settings and caching.

3. **Keyboard and Mouse Controls**:
   - Move characters using `W`, `A`, `S`, `D` keys.
   - Adjust camera position using `Z`, `X`, `C`, `V` keys.
   - Mouse-click-based direction control.

4. **Game Assets from GRF Files**:
   - Reads sprite data and configuration files from `.grf` file systems.

5. **Modular Architecture with ECS**:
   - Encapsulation of rendering and action logic into systems.
   - Seamless addition/removal of entities or other systems.

6. **Logging and Debugging**:
   - Uses [zerolog](https://github.com/rs/zerolog) for structured logging.
   - Debug output including input states and errors.

---

## Requirements

To build and run the Midgarts Client, the following dependencies must be installed:

- **Go SDK 1.21 or later**
- **Libraries**:
  - [Engo ECS](https://github.com/EngoEngine/ecs): Entity-Component-System architecture.
  - [SDL2](https://github.com/veandco/go-sdl2): SDL2 bindings for Go (For windowing and events).
  - [OpenGL](https://github.com/go-gl/gl): OpenGL bindings for Go.
  - [MathGL](https://github.com/go-gl/mathgl): Vector and matrix operations.
  - [Godotenv](https://github.com/joho/godotenv): Automatic `.env` file loading.
  - [GRF](https://github.com/project-midgard/midgarts/internal/fileformat/grf): Custom library for `.grf` files.
  - [Zerolog](https://github.com/rs/zerolog): Structured and fast logging.

---

## Setup

### Step 1: Clone the Repository

```sh
git clone <repository-url>
cd midgarts-client
```

### Step 2: Install Dependencies

Use Go to download all the required modules:

```sh
go mod tidy
```

### Step 3: Set Environment Variables

The application requires the `.env` file or environmental variable `GRF_FILE_PATH` to locate required assets:

```env
GRF_FILE_PATH=/path/to/your/grf/file
```

### Step 4: Run the Application

After setting up everything, simply run:

```sh
go run main.go
```

---

## Usage

### Controls

#### **Character Movement**
| Action                      | Input                                       |
|-----------------------------|---------------------------------------------|
| Move Up                     | `W`                                         |
| Move Down                   | `S`                                         |
| Move Left                   | `A`                                         |
| Move Right                  | `D`                                         |
| Diagonal Movement           | `W+D`, `W+A`, `S+D`, `S+A`                 |

#### **Character Direction via Mouse**
| Action                               | Input                          |
|--------------------------------------|--------------------------------|
| Set Direction (Mouse Click)          | Top-left, Bottom-left, etc. in respective viewport |

#### **Camera Controls**
| Action                   | Input    |
|--------------------------|----------|
| Move Camera Backward     | `Z`      |
| Move Camera Forward      | `X`      |
| Move Camera Left         | `C`      |
| Move Camera Right        | `V`      |

---

## Folder Structure

The following are some key directories in the project:

- **`internal/camera`**: Perspective camera logic.
- **`internal/character`**: Character properties (direction, state, jobs, etc.).
- **`internal/entity`**: Definitions for character entities.
- **`internal/system`**: Systems for action handling and rendering logic.
- **`internal/window`**: SDL2-based window utilities.
- **`pkg/version`**: Application version management.

---

## License

This project is licensed under the MIT License. For more details, refer to the `LICENSE` file.

---

Enjoy building with the Midgarts Client! For contributions or bug reporting, please reach out via the project's issue tracker.
