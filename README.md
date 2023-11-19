# Distributed Mutual Exclusion Simulation

This repository contains a simulation of distributed mutual exclusion using the Token Ring algorithm.

## Getting Started

1. Download the repository.
2. Open three separate terminals and run the following commands in each terminal:

    ```bash
    go run . -client 1
    go run . -client 2
    go run . -client 3
    ```

   Make sure to execute all three commands in each terminal before proceeding.

3. Select "Connect to other clients" in each terminal after running the commands.

## Usage

Now that the simulation is set up, you have two options:

- **Send Messages:** Choose to send messages to other clients.
- **Enter Critical Section:** Choose to enter the critical section. Upon entering, a log will indicate that you are inside the critical section.

## Note

This simulation is based on the Token Ring algorithm. Only one client can enter the critical seciton, however, be aware that, at the moment, there is no inddicator who is the next in line in the ring, so it is just which client notice first that the critical section is empty, will enter the critical section

For a detailed explanation of how we want the implementation, please refer to the report PDF file provided in the hand-in submission.



