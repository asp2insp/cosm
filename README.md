# Naked Mole Rat

Supervisor/Process runner framework for golang processes.
Can validate, run, restart, and kill jobs that are shipped
as Go binaries. Possible use cases include microservice
architecture, cluster job processing, and auto updating of
client side code

Need to think about how this is different from [https://github.com/asim/micro](https://github.com/asim/micro)

Main points that I see:

 * Easy up/down of services
 * Built-in config
 * Supervisor trees
 * Monitoring interface

Desired workflow:

 * Start service tree and see it running in the web interface
 * Write some code as a service with message handlers registered for different messages
 * Post that code to the process runner with `gotree run` which should give me a hash
 * See that code running and responding to heartbeat in the monitor web interface
 * Write another service and test it - even running locally or in go test it should talk to running services
 * Post that up and get a new hash
 * Edit the code for the first code and do an update on it. I should see all traffic routed to the new version of the service while the old one gracefully finishes its work queue and terminates.
 * My new update is crashing, but the service tree should restart it every time
 * Write a new service which depends on a non-existant service and post it up. It should crash and not start again until a new service appears (since the dependencies might now be resolved)
