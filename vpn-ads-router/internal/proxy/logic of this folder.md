# Proxy and Router Folder Responsibilities

## Proxy Folder
The `proxy` folder is responsible for handling the initial communication layer and acts as an intermediary between incoming connections and the router logic. Its primary responsibilities include:

1. **Accepting Incoming TCP Connections**:
    - Listens for incoming client connections on a specified address and port.
    - Manages the lifecycle of these connections.

2. **Reading Raw ADS/AMS Messages**:
    - Reads raw data packets from the connected clients.
    - Extracts and parses the relevant information, such as the source NetID.

3. **Parsing Source NetID**:
    - Extracts the source NetID from the incoming ADS/AMS messages.
    - Ensures malformed packets are handled gracefully with fallback values.

4. **Pushing Messages to the Scheduler Queue**:
    - Packages the parsed data (NetID and payload) into a `ClientMessage` struct.
    - Sends the `ClientMessage` to a shared channel (`incomingChan`) for further processing by the scheduler.

### Key Exports
- **`StartListener(address string)`**: Starts the TCP listener to accept client connections.
- **`incomingChan <- ClientMessage`**: A channel where parsed messages are pushed for consumption by the scheduler.

---

## Router Folder
The `router` folder is responsible for higher-level logic, including routing, connection management, and response handling. Its primary responsibilities include:

1. **Routing Logic**:
    - Matches incoming requests to the appropriate destination based on the NetID.
    - Ensures responses are routed back to the correct client.

2. **Connection Registry**:
    - Maintains a registry of active connections and their associated NetIDs.
    - Tracks the state of each connection for efficient routing.

3. **Response Logic**:
    - Handles responses from devices or services.
    - Matches responses back to the originating connection using the NetID.

4. **Integration with Scheduler**:
    - Consumes messages from the scheduler queue.
    - Processes and forwards them to the appropriate destination.

---

## Separation of Concerns
- The **proxy folder** focuses on low-level connection handling and message parsing.
- The **router folder** handles higher-level logic, such as routing, connection management, and response processing.

This separation ensures modularity and simplifies future enhancements or debugging.  