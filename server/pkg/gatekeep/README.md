# Gatekeep

Gatekeep is a package to support arbitrary resource locking.
The package understands Corteza's resource structure and their hierarchy.

## Service

The service provides bare metal functionality to manage locks.

## Locker

The locker provides a state around the stateless service.
The locker manages requested locks for ease of use throughout the system.

## Inmem Bus

The inmem bus provides a simple event bus to handle listeners for gatekeep events.
At some point, this might have to become some distributed system, but for current needs and development, this is plenty.

## Inmem Store

The inmem store provides a simple storage thing for the locks.
The storage utilizes the trie data structure to be more optimal.

At a later point, when we start didling with micro services, we'll need to introduce additional drivers to support distributed databases.
