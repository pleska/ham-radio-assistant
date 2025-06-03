# Ham Radio Assistant

A Model Context Protocol (MCP) server providing useful tools for amateur radio operators. This application enables AI assistants to perform ham radio-related calculations and lookups through standardized tool interfaces.

## Features

- **Callsign Lookup**: Query information about amateur radio callsigns
- **Antenna Bearing Calculations**: Calculate bearing between coordinates
- **Callsign-to-Callsign Bearing**: Calculate bearing between two amateur radio operators based on their callsigns
- **POTA Park Lookup**: Query information about Parks on the Air (POTA) locations

## Model Context Protocol (MCP)

This application is built as an MCP server, allowing AI assistants to access its tools through a standardized interface. MCP is a protocol for communication between AI systems and tools, enabling AI models to extend their capabilities through external services.

## Available Tools

### 1. Callsign Lookup

Retrieves detailed information about an amateur radio callsign using the callook.info API.

**Tool ID**: `callsign-lookup`

**Inputs:**
- `callsign` (string, required): Amateur radio callsign in the format of 1-2 letters/numbers, followed by a number, followed by 1-3 letters

**Returns:**
- Formatted text with details including:
  - License class
  - Operator name
  - Address information
  - Grid square
  - Coordinates (latitude/longitude)
  - License grant date, expiry date
  - FRN (FCC Registration Number)
  - Previous callsign (if applicable)
  - Link to FCC ULS page

### 2. Antenna Bearing Calculator

Calculates bearing and distance between two geographic coordinates.

**Tool ID**: `antenna-bearing`

**Inputs:**
- `origin-latitude` (string, required): Origin station latitude in decimal degrees
- `origin-longitude` (string, required): Origin station longitude in decimal degrees
- `destination-latitude` (string, required): Destination station latitude in decimal degrees
- `destination-longitude` (string, required): Destination station longitude in decimal degrees

**Returns:**
- Distance in miles and kilometers
- Bearing in degrees from North

### 3. Callsign Bearing Calculator

Calculates bearing and distance between two amateur radio operators based on their callsigns.

**Tool ID**: `callsign-bearing`

**Inputs:**
- `origin-callsign` (string, required): Your callsign
- `destination-callsign` (string, required): Destination callsign

**Returns:**
- Origin and destination locations with coordinates and grid squares
- Distance in miles and kilometers
- Bearing in degrees from North

### 4. POTA Park Lookup

Retrieves detailed information about a Parks on the Air (POTA) location.

**Tool ID**: `pota-park-lookup`

**Inputs:**
- `reference` (string, required): POTA park reference (e.g., US-2312)

**Returns:**
  - Park name and reference
  - Park type and status (active/inactive)
  - Location information (state/country)
  - Geographic details:
    - Coordinates (latitude/longitude)
    - Grid square (4-character and 6-character)
  - Access and activation methods
  - Special comments or restrictions
  - Park website URL (when available)
  - First activation information (callsign and date)
  - Link to POTA website for the park

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Docker (optional, for containerized deployment)

### Installation

1. Clone the repository:
```
git clone https://github.com/yourusername/ham-radio-assistant.git
cd ham-radio-assistant
```
2. Build the docker Deployment
```
docker build -t ham-radio-assistant .
```

### MCP Configuration

#### VS Code with Docker

To configure the Ham Radio Assistant as an MCP server in VS Code using Docker:

1. Ensure you have built the Docker image:
   ```
   docker build -t ham-radio-assistant .
   ```

2. Create or update the `.vscode/mcp.json` file with:
   ```json
   {
       "inputs": [],
       "servers": {
           "HamRadioAssistant": {
               "command": "docker",
               "args": [
                   "run",
                   "-i",
                   "--rm",
                   "ham-radio-assistant"
               ]
           }
       }
   }
   ```

3. Configure your VS Code extension to utilize the MCP server defined in the `mcp.json` file.

#### Claude Desktop Configuration

To use Ham Radio Assistant with Claude Desktop:

1. Ensure you have built the Docker image:
   ```
   docker build -t ham-radio-assistant .
   ```

2. Create a `claude_desktop_config.json` file:
   ```json
    {
        "mcpServers": {
            "HamRadioAssistant": {
                "command": "docker",
                "args": [
                    "run",
                    "-i",
                    "--rm",
                    "ham-radio-assistant"
                ]
            }
        }
    }
   ```

3. Open Claude Desktop and load the configuration file via Settings.


## Acknowledgments

- [callook.info](https://callook.info/) for providing the callsign lookup API
- [pota.app](https://pota.app) for providing the parks on the air (pota) parks list CSV. 
- [Mark3Labs](https://github.com/mark3labs/mcp-go) for the MCP Go implementation
